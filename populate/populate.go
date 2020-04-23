package populate

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"strings"
)

var logger = log.NewLogger("msx.populate")

var tasks Tasks

func RegisterPopulationTask(task Task) {
	tasks = append(tasks, task)
}

func getRegisteredJobs() types.StringSet {
	var jobSet = make(types.StringSet)
	for _, task := range tasks {
		jobSet.AddAll(task.During()...)
	}
	return jobSet
}

func CustomizeCommand(cmd *cobra.Command) {
	cmd.Args = func(cmd *cobra.Command, args []string) error {
		jobSet := getRegisteredJobs()
		for _, arg := range args {
			if !jobSet.Contains(arg) {
				return errors.Errorf(
					"Unknown populate job %q.  Must be one of: %s",
					arg,
					strings.Join(jobSet.Values(), ", "))
			}
		}
		return nil
	}

	cmd.Use += " [job [...]]"

	cmd.Flags().Bool("dry-run", false, "List the populate tasks that would be executed")
	cmd.Flags().Bool("list", false, "List the populate jobs and tasks that are registered")
}

func Populate(ctx context.Context, args []string) (err error) {
	list, _ := config.FromContext(ctx).Bool("cli.flag.list")
	if list {
		listPopulateJobsAndTasks(ctx)
		return nil
	}

	dryRun, _ := config.FromContext(ctx).Bool("cli.flag.dryrun")

	if len(args) == 0 {
		args = []string{"all"}
	}

	logger.WithContext(ctx).WithField("dry-run", dryRun).Infof("Executing populate jobs: '%s'", strings.Join(args, "', '"))

	for _, task := range tasks.During(args...).Ordered() {
		logger.WithContext(ctx).WithField("dry-run", dryRun).Infof("Executing task: %s", task.Description())
		if dryRun {
			continue
		}
		err = task.Populate(ctx)
		if err != nil {
			logger.WithContext(ctx).WithField("dry-run", dryRun).WithError(err).Infof("Executing task failed: %q", task.Description())
			return
		}
	}

	logger.WithContext(ctx).WithField("dry-run", dryRun).Infof("Finished populate jobs: '%s'", strings.Join(args, "', '"))

	return nil
}

func listPopulateJobsAndTasks(ctx context.Context) {
	jobs := getRegisteredJobs()
	for job := range jobs {
		for _, task := range tasks.Ordered() {
			if types.StringStack(task.During()).Contains(job) {
				logger.WithContext(ctx).Infof("Job %q: %s", job, task.Description())
			}
		}
	}
}
