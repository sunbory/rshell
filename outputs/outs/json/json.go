package json

import (
	"encoding/json"
	"fmt"
	"github.com/luckywinds/rshell/types"
)


var taskresult types.Taskresult

type JSON struct {
}

func (j JSON) Print(intime bool, result types.Hostresult, hg types.Hostgroup) {
	taskresult.Results = append(taskresult.Results, result)
}

func (j JSON) Break(intime bool, hg types.Hostgroup)  {
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

func (j JSON) Finish(intime bool, hg types.Hostgroup) {
	taskresult.Name = taskresult.Results[0].Actionname
	d, _ := json.MarshalIndent(&taskresult, "", "  ")
	fmt.Println(string(d))
}
