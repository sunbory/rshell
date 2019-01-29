package yaml

import (
	"fmt"
	"github.com/luckywinds/rshell/types"
	"gopkg.in/yaml.v2"
)


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
	taskresult.Name = taskresult.Results[0].Actionname
	d, _ := yaml.Marshal(&taskresult)
	fmt.Println(string(d))
}
