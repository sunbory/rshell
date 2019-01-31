package yaml

import (
	"fmt"
	"github.com/luckywinds/rshell/types"
	"gopkg.in/yaml.v2"
	"strings"
)

type SubTaskresult struct {
	Name    string       `yaml:"name,omitempty"`
	Results []types.Hostresult `yaml:"results,omitempty"`
}
var subTaskresult = SubTaskresult{}

type Tasksresults struct {
	Name string `yaml:"name,omitempty"`
	Results []SubTaskresult `yaml:"results,omitempty"`
}
var tasksresults = Tasksresults{}

var result = []Tasksresults{}

type YAML struct {
}

func (y YAML) Print(actionname, actiontype string, result types.Hostresult, hg types.Hostgroup) {
	subTaskresult.Results = append(subTaskresult.Results, result)
}

func (y YAML) Break(actionname, actiontype string, hg types.Hostgroup)  {
	m := make(map[string]types.Hostresult)
	for _, v := range subTaskresult.Results {
		m[v.Hostaddr] = v
	}

	for _, h := range hg.Ips {
		if _, ok := m[h]; !ok {
			item := types.Hostresult{
				Actiontype: actiontype,
				Groupname:  hg.Groupname,
				Hostaddr:   h,
				Error:      "TIMEOUT",
				Stdout:     "",
				Stderr:     "",
			}
			subTaskresult.Results = append(subTaskresult.Results, item)
		}
	}
}

func (y YAML) Finish(actionname, actiontype string, hg types.Hostgroup) {
	var taskName, staskName string
	names := strings.Split(actionname, "/")
	if len(names) == 2 {
		taskName = names[0]
		staskName = names[1]
	} else {
		taskName = names[0]
		staskName = names[0]
	}

	if tasksresults.Name != "" && tasksresults.Name != taskName {
		result = append(result, tasksresults)
		tasksresults = Tasksresults{}
	}

	subTaskresult.Name = staskName
	tasksresults.Name = taskName
	tasksresults.Results = append(tasksresults.Results, subTaskresult)
	subTaskresult = SubTaskresult{}
}

func (y YAML) End() {
	result = append(result, tasksresults)
	tasksresults = Tasksresults{}

	d, _ := yaml.Marshal(&result)
	fmt.Println(string(d))

	result = []Tasksresults{}
}