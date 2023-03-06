// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package servicetest

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
	"github.com/thejerf/abtime"
	"reflect"
	"testing"
	"time"
)

// Allow the testing.T object to be injected via the context.Context
type contextKey int

const contextKeyTesting contextKey = iota

type ServiceTest struct {
	Context struct {
		Base         context.Context
		Config       *config.Config
		TokenDetails *securitytest.MockTokenDetailsProvider
		Injectors    []types.ContextInjector
	}
	Clock              *abtime.ManualTime
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

func (s *ServiceTest) WithGot(gots ...interface{}) *ServiceTest {
	s.Got = gots
	return s
}

func (s *ServiceTest) WithWant(wants ...interface{}) *ServiceTest {
	s.Want = wants
	return s
}

func (s *ServiceTest) WithWantErr(wantErr bool) *ServiceTest {
	s.WantErr = wantErr
	return s
}

func (s *ServiceTest) WithHasErr(hasErr bool) *ServiceTest {
	s.HasErr = hasErr
	return s
}

func (r *ServiceTest) WithAssertExpectations(assert bool) *ServiceTest {
	r.AssertExpectations = assert
	return r
}

func (s *ServiceTest) WithRecording(rec *logtest.Recording) *ServiceTest {
	s.Recording = rec
	return s
}

func (s *ServiceTest) WithNow(t time.Time) *ServiceTest {
	s.Clock = abtime.NewManualAtTime(t)
	return s
}

func (s *ServiceTest) WithContext(ctx context.Context) *ServiceTest {
	s.Context.Base = ctx
	return s
}

func (s *ServiceTest) WithContextInjector(i types.ContextInjector) *ServiceTest {
	s.Context.Injectors = append(s.Context.Injectors, i)
	return s
}

func (s *ServiceTest) WithConfig(cfg *config.Config) *ServiceTest {
	s.Context.Config = cfg
	return s
}

func (s *ServiceTest) WithTokenDetailsProvider(provider *securitytest.MockTokenDetailsProvider) *ServiceTest {
	s.Context.TokenDetails = provider
	return s
}

func (s *ServiceTest) WithLogCheck(l logtest.Check) *ServiceTest {
	s.Checks.Log = append(s.Checks.Log, l)
	return s
}

func (s *ServiceTest) checkWant() (results []error) {
	if len(s.Got) != len(s.Want) {
		results = append(results, errors.Errorf("Wanted %d values, %d values returned", len(s.Want), len(s.Got)))
	}

	for n, want := range s.Want {
		if len(s.Got) < n+1 {
			// Handled above
			continue
		}
		if s.HasErr && n == len(s.Want)-1 {
			// Handled by WantErr
			continue
		}
		got := s.Got[n]

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

func (s *ServiceTest) checkWantErr() error {
	return testhelpers.CheckErr(s.Got, s.HasErr, s.WantErr)
}

func (s *ServiceTest) Test(t *testing.T, fn testhelpers.ServiceTestFunc) {
	var cfg *config.Config
	if s.Context.Config != nil {
		cfg = s.Context.Config
	} else {
		cfg = configtest.NewInMemoryConfig(nil)
	}

	err := fs.ConfigureFileSystem(cfg)
	assert.NoError(t, err)

	if s.Recording == nil {
		s.Recording = logtest.RecordLogging()
	}

	ctx := s.Context.Base
	if ctx == nil {
		ctx = context.Background()
	}

	ctx = config.ContextWithConfig(ctx, cfg)

	ctx = context.WithValue(ctx, contextKeyTesting, t)
	if s.Context.TokenDetails != nil {
		ctx = s.Context.TokenDetails.Inject(ctx)
	}

	if s.Clock != nil {
		ctx = types.ContextWithClock(ctx, s.Clock)
	}

	for _, injector := range s.Context.Injectors {
		ctx = injector(ctx)
	}

	// Execute the test
	fn(t, ctx)

	// Check the logs
	for _, logCheck := range s.Checks.Log {
		errs := logCheck.Check(s.Recording)
		s.Errors.Log = append(s.Errors.Log, errs...)
	}

	// Check Got vs Want
	if len(s.Want) > 0 {
		s.Errors.Want = s.checkWant()
	}

	// Check WantErr
	wantErr := s.checkWantErr()
	if wantErr != nil {
		s.Errors.WantErr = []error{wantErr}
	}

	// Report any errors
	testhelpers.ReportErrors(t, "Want", s.Errors.Want)
	testhelpers.ReportErrors(t, "WantErr", s.Errors.WantErr)
	testhelpers.ReportErrors(t, "Log", s.Errors.Log)
}
