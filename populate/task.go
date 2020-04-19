package populate

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"sort"
)

type Task interface {
	Description() string
	During() []string
	Order() int
	Populate(context.Context) error
}

type Tasks []Task

func (t Tasks) Ordered() Tasks {
	u := t[:]
	sort.Sort(u)
	return u
}

func (t Tasks) During(jobs ...string) (results Tasks) {
	for _, task := range t {
		for _, job := range jobs {
			if types.StringStack(task.During()).Contains(job) {
				results = append(results, task)
				break
			}
		}
	}
	return
}

func (t Tasks) Len() int           { return len(t) }
func (t Tasks) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t Tasks) Less(i, j int) bool { return t[i].Order() < t[j].Order() }
