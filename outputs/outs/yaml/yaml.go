package yaml

import (
	"fmt"
	"github.com/luckywinds/rshell/types"
	"gopkg.in/yaml.v2"
	"strings"
)

var taskresults = make(map[string][]types.Taskresult)
var taskresult types.Taskresult

type YAML struct {
}

func (y YAML) Print(intime bool, result types.Hostresult, hg types.Hostgroup) {
	taskresult.Results = append(taskresult.Results, result)
}

func (y YAML) Break(intime bool, hg types.Hostgroup)  {
	m := make(map[string]types.Hostresult)
	for _, v := range taskresult.Results {
		m[v.Hostaddr] = v
	}

	for _, h := range hg.Ips {
		if _, ok := m[h]; !ok {
			item := types.Hostresult{
				Actionname: taskresult.Results[0].Actionname,
				Actiontype: taskresult.Results[0].Actiontype,
				Groupname:  hg.Groupname,
				Hostaddr:   h,
				Error:      "TIMEOUT",
				Stdout:     "",
				Stderr:     "",
			}
			taskresult.Results = append(taskresult.Results, item)
		}
	}
}

func (y YAML) Finish(intime bool, hg types.Hostgroup) {
	var taskName, staskName string
	names := strings.Split(taskresult.Results[0].Actionname, "/")
	if len(names) == 2 {
		taskName = names[0]
		staskName = names[1]
	} else {
		taskName = names[0]
		staskName = names[0]
	}

	taskresult.Name = staskName
	taskresults[taskName] = append(taskresults[taskName], taskresult)
	taskresult = types.Taskresult{}
}

func (y YAML) End() {
	if len(taskresults) > 0 {
		d, _ := yaml.Marshal(&taskresults)
		fmt.Println(string(d))
		for key, _ := range taskresults {
			delete(taskresults, key)
		}
	}
}