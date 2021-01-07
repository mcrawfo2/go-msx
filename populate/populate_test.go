package populate

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestMockImplementations(t *testing.T) {
	var _ Task = new(MockTask)
}

func TestCustomizeCommand(t *testing.T) {
	cmd := &cobra.Command{}
	CustomizeCommand(cmd)
	assert.NotNil(t, cmd.Args)
	assert.NotNil(t, cmd.Flags().Lookup("dry-run"))
	assert.NotNil(t, cmd.Flags().Lookup("list"))
}

func TestPopulate(t *testing.T) {
	var wasRun bool

	ctx := configtest.ContextWithNewInMemoryConfig(
		context.Background(),
		map[string]string{
			"cli.flag.list":   "false",
			"cli.flag.dryrun": "false",
		})

	tasks = []Task{}
	mockTask1 := new(MockTask)
	mockTask1.On("Description").Return("Populate some entities")
	mockTask1.On("During").Return([]string{"mach"})
	mockTask1.On("Order").Return(100)
	mockTask1.On("Populate", mock.AnythingOfType("*context.valueCtx")).
		Run(func(args mock.Arguments) {
			wasRun = true
		}).
		Return(nil)
	RegisterPopulationTask(mockTask1)

	type args struct {
		ctx  context.Context
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				ctx:  ctx,
				args: []string{"mach"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wasRun = false
			if err := Populate(tt.args.ctx, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("Populate() error = %v, wantErr %v", err, tt.wantErr)
			} else if err == nil {
				assert.True(t, wasRun)
			}
		})
	}
}

func TestRegisterPopulationTask(t *testing.T) {
	tasks = []Task{}
	mockTask1 := new(MockTask)
	mockTask1.On("Description").Return("Populate some entities")
	mockTask1.On("During").Return([]string{"mach"})
	RegisterPopulationTask(mockTask1)

	assert.Len(t, tasks, 1)
	assert.Equal(t, mockTask1, tasks[0])
}

func Test_getRegisteredJobs(t *testing.T) {
	tasks = []Task{}

	mockTask1 := new(MockTask)
	mockTask1.On("Description").Return("Populate some entities")
	mockTask1.On("During").Return([]string{"mach"})
	RegisterPopulationTask(mockTask1)

	mockTask2 := new(MockTask)
	mockTask2.On("Description").Return("Populate some entities")
	mockTask2.On("During").Return([]string{"mach2"})
	RegisterPopulationTask(mockTask2)

	got := getRegisteredJobs()
	for _, v := range []string{"mach", "mach2"} {
		if !got.Contains(v) {
			t.Errorf("getRegisteredJobs() does not contain %q", v)
		}
	}
}

func Test_listPopulateJobsAndTasks(t *testing.T) {
	tasks = []Task{}
	mockTask1 := new(MockTask)
	mockTask1.On("Description").Return("Populate some entities")
	mockTask1.On("During").Return([]string{"mach"})
	mockTask1.On("Order").Return(1000)
	RegisterPopulationTask(mockTask1)

	listPopulateJobsAndTasks(context.Background())
}

func Test_validateJobs(t *testing.T) {
	tasks = []Task{}

	mockTask1 := new(MockTask)
	mockTask1.On("Description").Return("Populate some entities")
	mockTask1.On("During").Return([]string{"mach"})
	RegisterPopulationTask(mockTask1)

	mockTask2 := new(MockTask)
	mockTask2.On("Description").Return("Populate some entities")
	mockTask2.On("During").Return([]string{"mach2"})
	RegisterPopulationTask(mockTask2)

	type args struct {
		jobs []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				jobs: []string{"mach", "mach2"},
			},
			wantErr: false,
		},
		{
			name: "Failure",
			args: args{
				jobs: []string{"mach3"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateJobs(tt.args.jobs); (err != nil) != tt.wantErr {
				t.Errorf("validateJobs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
