package json

import (
	"encoding/json"
	"fmt"
	"github.com/luckywinds/rshell/types"
	"strings"
)

type SubTaskresult struct {
	Name    string       `json:"name,omitempty"`
	Results []types.Hostresult `json:"results,omitempty"`
}
var subTaskresult = SubTaskresult{}

type Tasksresults struct {
	Name string `json:"name,omitempty"`
	Results []SubTaskresult `json:"results,omitempty"`
}
var tasksresults = Tasksresults{}

var result = []Tasksresults{}

type JSON struct {
}

func (j JSON) Print(actionname, actiontype string, result types.Hostresult, hg types.Hostgroup) {
	subTaskresult.Results = append(subTaskresult.Results, result)
}

func (j JSON) Break(actionname, actiontype string, hg types.Hostgroup)  {
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

func (j JSON) Finish(actionname, actiontype string, hg types.Hostgroup) {
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

func (j JSON) End() {
	result = append(result, tasksresults)
	tasksresults = Tasksresults{}

	d, _ := json.MarshalIndent(&result, "", "  ")
	fmt.Println(string(d))

	result = []Tasksresults{}
}
