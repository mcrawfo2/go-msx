package populate

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	api "cto-github.cisco.com/NFV-BU/go-msx/integration/usermanagement"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/paging"
	"cto-github.cisco.com/NFV-BU/go-msx/populate"
	"cto-github.cisco.com/NFV-BU/go-msx/resource"
	"cto-github.cisco.com/NFV-BU/go-msx/security/service"
	"github.com/pkg/errors"
	"net/http"
	"path"
)

const (
	permissionsPopulatorConfigRoot = "populate.usermanagement.permission"
	manifestFile                   = "manifest.json"
)

var logger = log.NewLogger("msx.integration.usermanagement.populate")

type manifest struct {
	Owner               string               `json:"owner"`
	DeletedCapabilities []string             `json:"deletedCapabilities"`
	Capabilities        []capabilityArtifact `json:"capabilities"`
	DeletedRoles        []string             `json:"deletedRoles"`
	Roles               []roleArtifact       `json:"roles"`
}

type capabilityArtifact struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Operation   string `json:"operation"`
	ObjectName  string `json:"objectName"`
	Owner       string `json:"owner"`
	IsDefault   bool   `json:"isDefault,omitempty"`
}

type roleArtifact struct {
	RoleName       string   `json:"roleName"`
	CapabilityList []string `json:"capabilitylist"`
}

type roleMap map[string]api.RoleResponse

type PermissionPopulatorConfig struct {
	Enabled bool   `config:"default=false"`
	Root    string `config:"default=${populate.root}/usermanagement"`
}

type PermissionPopulator struct {
	cfg PermissionPopulatorConfig
}

func (p PermissionPopulator) depopulateCapability(ctx context.Context, idm api.Api, owner string, name string) error {
	logger.WithContext(ctx).Infof("Removing capability %q", name)

	response, err := idm.DeleteCapability(true, owner, name)
	if err != nil {
		if response != nil && response.StatusCode == http.StatusNotFound {
			logger.WithContext(ctx).Warnf("Capability %q not found.", name)
			return nil
		}

		return errors.Wrapf(err, "Failed to remove capability %q", name)
	}

	return nil
}

func (p PermissionPopulator) populateCapability(ctx context.Context, idm api.Api, owner string, capability capabilityArtifact) error {
	logger.WithContext(ctx).Infof("Removing capability %q", capability.Name)

	_, err := idm.BatchUpdateCapabilities(true, owner, []api.CapabilityUpdateRequest{
		{
			Category:    capability.Category,
			Description: capability.Description,
			DisplayName: capability.DisplayName,
			Name:        capability.Name,
			ObjectName:  capability.ObjectName,
			Operation:   capability.Operation,
		},
	})

	if err != nil {
		return errors.Wrapf(err, "Failed to update capability %q", capability.Name)
	}

	return nil
}

func (p PermissionPopulator) depopulateRole(ctx context.Context, idm api.Api, name string, roleMap roleMap) error {
	logger.WithContext(ctx).Infof("Removing role %q", name)

	if _, ok := roleMap[name]; !ok {
		logger.WithContext(ctx).Warnf("Role %s already deleted", name)
		return nil
	}

	response, err := idm.DeleteRole(name)
	if err != nil {
		if response != nil && response.StatusCode == http.StatusNotFound {
			logger.WithContext(ctx).Warnf("Role %q not found.", name)
			return nil
		}

		return errors.Wrapf(err, "Failed to remove role %q", name)
	}

	return nil
}

func (p PermissionPopulator) populateRole(ctx context.Context, idm api.Api, owner string, role roleArtifact, roleMap roleMap) (err error) {
	if _, ok := roleMap[role.RoleName]; ok {
		logger.WithContext(ctx).Infof("Updating existing role %q", role.RoleName)

		_, err = idm.UpdateRole(true, api.RoleUpdateRequest{
			CapabilityList: role.CapabilityList,
			Owner:          roleMap[role.RoleName].Owner,
			RoleName:       role.RoleName,
		})
	} else {
		logger.WithContext(ctx).Infof("Creating new role %q", role.RoleName)

		_, err = idm.CreateRole(true, api.RoleCreateRequest{
			CapabilityList: role.CapabilityList,
			Owner:          owner,
			RoleName:       role.RoleName,
		})
	}

	if err != nil {
		return errors.Wrapf(err, "Failed top populate role %q", role.RoleName)
	}

	return nil
}

func (p PermissionPopulator) getRoles(ctx context.Context, idm api.Api) (roleMap, error) {
	response, err := idm.GetRoles(true, paging.Request{Page: 0, Size: 1000})
	if err != nil {
		return nil, err
	}

	roleListResponse, ok := response.Payload.(*api.RoleListResponse)
	if !ok {
		return nil, errors.New("Unable to parse IDM role list response")
	}

	var result = make(map[string]api.RoleResponse)
	for _, role := range roleListResponse.Roles {
		result[role.RoleName] = role
	}
	return result, nil
}

func (p PermissionPopulator) Populate(ctx context.Context) error {
	if !p.cfg.Enabled {
		logger.WithContext(ctx).Warn("Role/Capability populator disabled.")
		return nil
	}

	return service.WithDefaultServiceAccount(ctx, func(ctx context.Context) error {

		var manifest manifest
		err := resource.
			Reference(path.Join(p.cfg.Root, manifestFile)).
			Unmarshal(&manifest)
		if err != nil {
			return err
		}

		idm, _ := api.NewIntegration(ctx)
		roleMap, err := p.getRoles(ctx, idm)
		if err != nil {
			return errors.Wrap(err, "Failed to retrieve role list")
		}

		logger.WithContext(ctx).Info("Populating capabilities")

		for _, name := range manifest.DeletedCapabilities {
			err = p.depopulateCapability(ctx, idm, manifest.Owner, name)
			if err != nil {
				return err
			}
		}

		for _, capability := range manifest.Capabilities {
			err = p.populateCapability(ctx, idm, manifest.Owner, capability)
			if err != nil {
				return err
			}
		}

		logger.WithContext(ctx).Info("Populating roles")

		for _, role := range manifest.Roles {
			err = p.populateRole(ctx, idm, manifest.Owner, role, roleMap)
			if err != nil {
				return err
			}
		}

		for _, name := range manifest.DeletedRoles {
			err = p.depopulateRole(ctx, idm, name, roleMap)
		}

		return nil
	})
}

func init() {
	populate.RegisterPopulationTask(
		populate.NewPopulatorTask(
			"Populate roles and capabilities",
			1000,
			[]string{"all", "customRolesAndCapabilities", "serviceMetadata"},
			func(ctx context.Context) (populate.Populator, error) {
				var cfg PermissionPopulatorConfig
				err := config.MustFromContext(ctx).Populate(&cfg, permissionsPopulatorConfigRoot)
				if err != nil {
					return nil, err
				}
				return &PermissionPopulator{
					cfg: cfg,
				}, nil
			}))
}
