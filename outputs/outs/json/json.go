package json

import (
	"encoding/json"
	"fmt"
	"github.com/luckywinds/rshell/types"
	"strings"
)

var taskresults = make(map[string][]types.Taskresult)
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

func (j JSON) End() {
	if len(taskresults) > 0 {
		d, _ := json.MarshalIndent(&taskresults, "", "  ")
		fmt.Println(string(d))
		for key, _ := range taskresults {
			delete(taskresults, key)
		}
	}
}
