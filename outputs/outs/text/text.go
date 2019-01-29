package text

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/luckywinds/rshell/types"
)

var taskresult types.Taskresult

type TEXT struct {
}

func (t TEXT) Print(intime bool, result types.Hostresult, hg types.Hostgroup) {
	taskresult.Results = append(taskresult.Results, result)
	if intime {
		printHeader(result, hg)
		printItem(result)
	}
}

func (t TEXT) Break(intime bool, hg types.Hostgroup)  {
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
			if intime {
				printItem(item)
			}
			taskresult.Results = append(taskresult.Results, item)
		}
	}
}

func (t TEXT) Finish(intime bool, hg types.Hostgroup) {
	if !intime {
		printHeader(taskresult.Results[0], hg)

		m := make(map[string]types.Hostresult)
		for _, v := range taskresult.Results {
			m[v.Hostaddr] = v
		}

		for _, h := range hg.Ips {
			printItem(m[h])
		}
	}
	taskresult = types.Taskresult{}
	header = false
}

func (t TEXT) End() {

}

var header = false
func printHeader(result types.Hostresult, hg types.Hostgroup) {
	if !header {
		color.Yellow("TASK [%-50s] *********************\n", result.Actionname + "@" + hg.Groupname)
		header = true
	}
}

func printItem(result types.Hostresult) {
	color.Green("HOST [%-16s] =======================================================\n", result.Hostaddr)
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