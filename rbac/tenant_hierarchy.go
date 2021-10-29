//go:generate mockery --inpackage --name=TenantHierarchyApi --structname=MockTenantHierarchyApi --filename=mock_tenanthierarchyapi.go

package rbac

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration/usermanagement"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
	"sync"
)

var errTenantNotFound = errors.New("tenant not found in cache")

type TenantHierarchyApi interface {
	Parent(ctx context.Context, tenantId types.UUID) (types.UUID, error)
	Ancestors(ctx context.Context, tenantId types.UUID) ([]types.UUID, error)
	Root(ctx context.Context) (types.UUID, error)
}

type TenantHierarchyCache struct {
	rootId     types.UUID              // root tenantId
	tenants    sync.Map                // maps tenantId to parentId
	loader     TenantHierarchyApi      // loads tenant hierarchy
	partitions *types.PartitionedMutex // prevents multiple identical simultaneous requests
}

type tenantHierarchyQuery struct {
	queryType string
	tenantId  string
}

type cacheRetrievalFunc func() error

func (t *TenantHierarchyCache) load(ctx context.Context, queryType string, tenantUuid types.UUID, cachedFunc cacheRetrievalFunc, loadFunc types.ActionFunc) (err error) {
	if err = cachedFunc(); err == nil {
		return nil
	}

	partitionTenantId := ""
	if tenantUuid != nil {
		partitionTenantId = tenantUuid.String()
	}

	return t.partitions.WithPartitionLock(
		tenantHierarchyQuery{
			queryType: queryType,
			tenantId:  partitionTenantId,
		},
		ctx,
		func(ctx context.Context) (err error) {
			if err := cachedFunc(); err == nil {
				return nil
			}

			if err := loadFunc(ctx); err != nil {
				return err
			}

			return cachedFunc()
		})
}

// Parent returns the parent id of the specified tenant
func (t *TenantHierarchyCache) Parent(ctx context.Context, tenantId types.UUID) (result types.UUID, err error) {
	err = t.load(ctx, "ancestors", tenantId,
		func() error {
			result, err = t.cachedParent(tenantId)
			return err
		},
		func(ctx context.Context) error {
			return t.loadAncestors(ctx, tenantId)
		})
	return
}

// Ancestors returns the consecutive parent ids of ancestors of the tenant, starting with the immediate parent and ending with the root tenant.
func (t *TenantHierarchyCache) Ancestors(ctx context.Context, tenantId types.UUID) (results []types.UUID, err error) {
	err = t.load(ctx, "ancestors", tenantId,
		func() error {
			results, err = t.cachedAncestors(tenantId)
			return err
		},
		func(ctx context.Context) error {
			return t.loadAncestors(ctx, tenantId)
		})
	return
}

// Root returns the root tenant id.
func (t *TenantHierarchyCache) Root(ctx context.Context) (result types.UUID, err error) {
	err = t.load(ctx, "root", nil,
		func() error {
			result, err = t.cachedRoot()
			return err
		},
		t.loadRoot)
	return
}

func (t *TenantHierarchyCache) loadRoot(ctx context.Context) error {
	rootId, err := t.loader.Root(ctx)
	if err != nil {
		return err
	}

	t.storeRoot(rootId)
	return nil
}

func (t *TenantHierarchyCache) loadAncestors(ctx context.Context, tenantId types.UUID) error {
	ancestors, err := t.loader.Ancestors(ctx, tenantId)
	if err != nil {
		return err
	}

	t.storeAncestors(tenantId, ancestors)
	return nil
}

func (t *TenantHierarchyCache) cachedRoot() (result types.UUID, err error) {
	if t.rootId == nil {
		return nil, errTenantNotFound
	}
	return t.rootId, nil
}

func (t *TenantHierarchyCache) cachedAncestors(tenantUuid types.UUID) (results []types.UUID, err error) {
	if t.rootId != nil && tenantUuid.Equals(t.rootId) {
		return
	}

	tenantId := tenantUuid.String()
	for !tenantUuid.Equals(t.rootId) {
		parentId, ok := t.tenants.Load(tenantId)
		if !ok {
			return nil, errors.Wrap(errTenantNotFound, tenantId)
		}
		tenantId = parentId.(string)
		tenantUuid = types.MustParseUUID(tenantId)
		results = append(results, tenantUuid)
	}

	return results, nil
}

func (t *TenantHierarchyCache) cachedParent(tenantId types.UUID) (result types.UUID, err error) {
	if t.rootId != nil && tenantId.Equals(t.rootId) {
		return nil, nil
	}

	tenantIdString := tenantId.String()
	parentId, ok := t.tenants.Load(tenantIdString)
	if ok {
		result = types.MustParseUUID(parentId.(string))
	} else {
		err = errTenantNotFound
	}
	return
}

func (t *TenantHierarchyCache) storeRoot(result types.UUID) {
	t.rootId = result
}

func (t *TenantHierarchyCache) storeAncestors(id types.UUID, ancestors []types.UUID) {
	tenantId := id.String()
	for _, parentUuid := range ancestors {
		parentId := parentUuid.String()
		t.tenants.Store(tenantId, parentId)
		tenantId = parentId
	}

	if len(ancestors) > 0 {
		t.storeRoot(ancestors[len(ancestors)-1])
	} else if len(ancestors) == 0 && t.rootId == nil {
		t.storeRoot(id)
	}
}

func newTenantHierarchyCache(ctx context.Context) (*TenantHierarchyCache, error) {
	loader, err := NewTenantHierarchyLoader(ctx)
	if err != nil {
		return nil, err
	}

	return &TenantHierarchyCache{
		tenants:    sync.Map{},
		loader:     loader,
		partitions: types.NewPartitionedMutex(),
	}, nil
}

type TenantHierarchyLoader struct{}

func (t TenantHierarchyLoader) Parent(ctx context.Context, tenantId types.UUID) (types.UUID, error) {
	api, err := usermanagement.NewIntegration(ctx)
	if err != nil {
		return nil, err
	}

	response, err := api.GetTenantHierarchyParent(tenantId)
	if err != nil {
		return nil, err
	}

	if len(response.BodyString) == 0 {
		return nil, nil
	}

	return types.ParseUUID(response.BodyString)
}

func (t TenantHierarchyLoader) Ancestors(ctx context.Context, tenantId types.UUID) ([]types.UUID, error) {
	api, err := usermanagement.NewIntegration(ctx)
	if err != nil {
		return nil, err
	}

	_, ancestors, err := api.GetTenantHierarchyAncestors(tenantId)
	if err != nil {
		return nil, err
	}

	return ancestors, nil
}

func (t TenantHierarchyLoader) Root(ctx context.Context) (types.UUID, error) {
	api, err := usermanagement.NewIntegration(ctx)
	if err != nil {
		return nil, err
	}

	response, err := api.GetTenantHierarchyRoot()
	if err != nil {
		return nil, err
	}

	tenantId, err := types.ParseUUID(response.BodyString)
	if err != nil {
		return nil, err
	}

	return tenantId, nil
}

type TenantHierarchyPoolLoader struct {
	pool   *types.WorkerPool 		// Loads tenant hierarchies in parallel
	loader TenantHierarchyLoader	// Does the actual work
}

func (l *TenantHierarchyPoolLoader) Root(ctx context.Context) (result types.UUID, err error) {
	err = l.pool.Run(
		func(ctx context.Context) error {
			result, err = l.loader.Root(ctx)
			return err
		},
		types.JobDecorator(log.RecoverLogDecorator(logger)),
		types.JobContext(ctx),
	)

	return
}

func (l *TenantHierarchyPoolLoader) Parent(ctx context.Context, tenantId types.UUID) (result types.UUID, err error) {
	err = l.pool.Run(
		func(ctx context.Context) error {
			result, err = l.loader.Parent(ctx, tenantId)
			return err
		},
		types.JobDecorator(log.RecoverLogDecorator(logger)),
		types.JobContext(ctx),
	)

	return
}

func (l *TenantHierarchyPoolLoader) Ancestors(ctx context.Context, tenantId types.UUID) (results []types.UUID, err error) {
	err = l.pool.Run(
		func(ctx context.Context) error {
			results, err = l.loader.Ancestors(ctx, tenantId)
			return err
		},
		types.JobDecorator(log.RecoverLogDecorator(logger)),
		types.JobContext(ctx),
	)

	return
}

func NewTenantHierarchyLoader(ctx context.Context) (TenantHierarchyApi, error) {
	api := TenantHierarchyLoaderFromContext(ctx)
	if api == nil {
		workerPool, err := types.NewWorkerPool(ctx, 4)
		if err != nil {
			return nil, err
		}

		api = &TenantHierarchyPoolLoader{
			pool: workerPool,
		}
	}

	return api, nil
}

var tenantHierarchy *TenantHierarchyCache

func GetTenantHierarchyApi(ctx context.Context) (service TenantHierarchyApi, err error) {
	service = TenantHierarchyFromContext(ctx)
	if service != nil {
		return service, nil
	}

	if tenantHierarchy == nil {
		tenantHierarchy, err = newTenantHierarchyCache(ctx)
		if err != nil {
			return nil, err
		}
	}

	return tenantHierarchy, nil
}
