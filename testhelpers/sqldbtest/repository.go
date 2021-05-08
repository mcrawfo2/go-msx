package sqldbtest

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/fs"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/logtest"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/securitytest"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

// Allow the testing.T object to be injected via the context.Context
type contextKey int

const contextKeyTesting contextKey = iota

type RepositoryTest struct {
	Context struct {
		Base         context.Context
		Config       *config.Config
		TokenDetails *securitytest.MockTokenDetailsProvider
		Injectors    []types.ContextInjector
	}
	HasErr             bool
	Got                []interface{}
	Want               []interface{}
	WantErr            bool
	AssertExpectations bool
	Checks             struct {
		Log []logtest.Check
	}
	Errors struct {
		Want    []error
		WantErr []error
		Log     []error
	}
	Recording *logtest.Recording
}

func (r *RepositoryTest) WithGot(gots ...interface{}) *RepositoryTest {
	r.Got = gots
	return r
}

func (r *RepositoryTest) WithWant(wants ...interface{}) *RepositoryTest {
	r.Want = wants
	return r
}

func (r *RepositoryTest) WithWantErr(wantErr bool) *RepositoryTest {
	r.WantErr = wantErr
	return r
}

func (r *RepositoryTest) WithHasErr(hasErr bool) *RepositoryTest {
	r.HasErr = hasErr
	return r
}

func (r *RepositoryTest) WithAssertExpectations(assert bool) *RepositoryTest {
	r.AssertExpectations = assert
	return r
}

func (r *RepositoryTest) WithRecording(rec *logtest.Recording) *RepositoryTest {
	r.Recording = rec
	return r
}

func (r *RepositoryTest) WithContext(ctx context.Context) *RepositoryTest {
	r.Context.Base = ctx
	return r
}

func (r *RepositoryTest) WithContextInjector(i types.ContextInjector) *RepositoryTest {
	r.Context.Injectors = append(r.Context.Injectors, i)
	return r
}

func (r *RepositoryTest) WithConfig(cfg *config.Config) *RepositoryTest {
	r.Context.Config = cfg
	return r
}

func (r *RepositoryTest) WithTokenDetailsProvider(provider *securitytest.MockTokenDetailsProvider) *RepositoryTest {
	r.Context.TokenDetails = provider
	return r
}

func (r *RepositoryTest) WithLogCheck(l logtest.Check) *RepositoryTest {
	r.Checks.Log = append(r.Checks.Log, l)
	return r
}

func (r *RepositoryTest) checkWant() (results []error) {
	if len(r.Got) != len(r.Want) {
		results = append(results, errors.Errorf("Wanted %d values, %d values returned", len(r.Want), len(r.Got)))
	}

	for n, want := range r.Want {
		if len(r.Got) < n+1 {
			// Handled above
			continue
		}
		if r.HasErr && n == len(r.Want)-1 {
			// Handled by WantErr
			continue
		}
		got := r.Got[n]

		zero := false
		if want == nil {
			gotValue := reflect.ValueOf(got)
			if gotValue.Kind() == reflect.Interface && got != nil {
				gotValue = gotValue.Elem()
			}
			switch gotValue.Kind() {
			case reflect.Ptr:
				zero = gotValue.IsNil()
			case reflect.Slice, reflect.Map:
				zero = gotValue.Len() == 0
			case reflect.Interface:
				zero = got != nil
			}
		}

		if !zero && !reflect.DeepEqual(want, got) {
			results = append(results, errors.Errorf("Returned Value %d mismatch:\n%s", n, testhelpers.Diff(want, got)))
		}
	}

	return
}

func (r *RepositoryTest) checkWantErr() error {
	if len(r.Got) == 0 {
		return errors.Errorf("Wanted Error, No values returned")
	} else if !r.HasErr {
		return errors.Errorf("Wanted Error, Method not flagged as returning error")
	} else if err, ok := r.Got[len(r.Got)-1].(error); ok && err == nil && r.WantErr {
		return errors.Errorf("Wanted Error, No error returned")
	} else if ok && !r.WantErr && err != nil {
		return errors.Errorf("Unwanted Error, Error returned:\n%s", testhelpers.Dump(err))
	}
	return nil
}

func (r *RepositoryTest) Test(t *testing.T, fn func(t *testing.T, ctx context.Context)) {
	var cfg *config.Config
	if r.Context.Config != nil {
		cfg = r.Context.Config
	} else {
		cfg = configtest.NewInMemoryConfig(nil)
	}

	err := fs.ConfigureFileSystem(cfg)
	assert.NoError(t, err)

	if r.Recording == nil {
		r.Recording = logtest.RecordLogging()
	}

	ctx := r.Context.Base
	if ctx == nil {
		ctx = context.Background()
	}

	ctx = config.ContextWithConfig(ctx, cfg)

	ctx = context.WithValue(ctx, contextKeyTesting, t)
	if r.Context.TokenDetails != nil {
		ctx = r.Context.TokenDetails.Inject(ctx)
	}

	for _, injector := range r.Context.Injectors {
		ctx = injector(ctx)
	}

	// Execute the test
	fn(t, ctx)

	// Check the logs
	for _, logCheck := range r.Checks.Log {
		errs := logCheck.Check(r.Recording)
		r.Errors.Log = append(r.Errors.Log, errs...)
	}

	// Check Got vs Want
	if len(r.Want) > 0 {
		r.Errors.Want = r.checkWant()
	}

	// Check WantErr
	wantErr := r.checkWantErr()
	if wantErr != nil {
		r.Errors.WantErr = []error{wantErr}
	}

	// Report any errors
	testhelpers.ReportErrors(t, "Want", r.Errors.Want)
	testhelpers.ReportErrors(t, "WantErr", r.Errors.WantErr)
	testhelpers.ReportErrors(t, "Log", r.Errors.Log)
}
