package text

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/luckywinds/rshell/types"
	"strings"
)

type SubTaskresult struct {
	Name    string
	Results []types.Hostresult
}
var subTaskresult = SubTaskresult{}

type Tasksresults struct {
	Name string
	Results []SubTaskresult
}

var hgtmp types.Hostgroup
var lastTaskName string

type TEXT struct {
}

func (t TEXT) Print(actionname, actiontype string, result types.Hostresult, hg types.Hostgroup) {
	subTaskresult.Results = append(subTaskresult.Results, result)
}

func (t TEXT) Break(actionname, actiontype string, hg types.Hostgroup)  {
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

func (t TEXT) Finish(actionname, actiontype string, hg types.Hostgroup) {
	hgtmp = hg

	var taskName, staskName string
	names := strings.Split(actionname, "/")
	if len(names) == 2 {
		taskName = names[0]
		staskName = names[1]
	} else {
		taskName = names[0]
		staskName = names[0]
	}

	subTaskresult.Name = staskName
	printSTask(taskName, subTaskresult, hgtmp)
	subTaskresult = SubTaskresult{}
	lastTaskName = taskName
}

func (t TEXT) End() {
}

func printSTask(taskName string, st SubTaskresult, hg types.Hostgroup) {
	pirntTaskHeader(taskName)

	pirntStaskHeader(st)
	m := make(map[string]types.Hostresult)
	for _, v := range st.Results {
		m[v.Hostaddr] = v
	}
	for _, h := range hg.Ips {
		printHost(m[h])
	}
}

func pirntTaskHeader(taskName string) {
	if taskName != lastTaskName {
		color.Yellow("TASK  [%-20s] ++++++++++++++++++++++++++++++++++++++++++++++++++\n", taskName)
	}
}

func pirntStaskHeader(st SubTaskresult) {
	color.Yellow("STASK [%-20s] ==================================================\n", st.Name)
}

func printHost(result types.Hostresult) {
	color.Green("HOST  [%-20s] --------------------------------------------------\n", result.Hostaddr)
	if result.Stdout != "" {
		fmt.Printf("%s\n", result.Stdout)
	}
	if result.Stderr != "" {
		color.Red("%s\n", "STDERR =>")
		fmt.Printf("%s\n", result.Stderr)
	}
	if result.Error != "" {
		color.Red("%s\n", "SYSERR =>")
		fmt.Printf("%s\n", result.Error)
	}
	if result.Stdout == "" && result.Stderr == "" && result.Error == "" {
		fmt.Println()
	}
}
