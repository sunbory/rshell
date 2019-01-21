package utils

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/luckywinds/rshell/options"
	"github.com/luckywinds/rshell/pkg/checkers"
	"github.com/luckywinds/rshell/pkg/prompt"
	"github.com/luckywinds/rshell/types"
	"strconv"
	"strings"
)

func OutputTaskHeader(name string) {
	color.Yellow("TASK [%-30s] *****************************************\n", name)
}
func OutputHostResult(result types.Hostresult) {
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

func Output(result types.Taskresult) {
	OutputTaskHeader(result.Name)
	for _, ret := range result.Results {
		OutputHostResult(ret)
	}
}


var aname, hname string
var port int

func GetLoadArgs(o options.Options, line string) (string, string, string, int, error) {
	ks := strings.Fields(line)
	if len(ks) < 2 {
		return "", "", "", 0, fmt.Errorf("load arguments illegal")
	}
	aname = o.CurrentEnv.Authname
	hname = o.CurrentEnv.Hostgroupname
	port = o.CurrentEnv.Port

	for _, value := range ks[1:] {
		if err := getLoadArgs(o, value); err != nil {
			return "", "", "", 0, err
		}
	}

	if hname == "" {
		return "", "", "", 0, fmt.Errorf("load arguments illegal, %s", "host empty")
	}
	if checkers.IsIpv4(hname) {
		if aname == "" {
			return "", "", "", 0, fmt.Errorf("load arguments illegal, %s", "auth empty")
		}
		if port == 0 {
			port = 22
		}
		prompt.AddCmd("-H" + strings.Trim(hname, " "))
	}

	return ks[0], aname, hname, port, nil
}

func getLoadArgs(o options.Options, a string) (err error){
	switch {
	case strings.HasPrefix(a, "-A"):
		aname = strings.TrimLeft(a, "-A")
		if _, ok := o.Authsm[aname]; !ok {
			return fmt.Errorf("%s", "Auth name not found")
		}
	case strings.HasPrefix(a, "-H"):
		hname = strings.TrimLeft(a, "-H")
		if _, ok := o.Hostgroupsm[hname]; !ok {
			if !checkers.IsIpv4(hname) {
				return fmt.Errorf("%s", "Hostgroup name not found")
			}
		}
	case strings.HasPrefix(a, "-P"):
		port, err = strconv.Atoi(strings.TrimLeft(a, "-P"))
		if err != nil {
			return fmt.Errorf("%s", "Port number illegal")
		}
	}
	return nil
}

func GetSSHArgs(line string) (string, []string, error) {
	ks := strings.Fields(line)
	if len(ks) < 2 {
		return "", []string{}, fmt.Errorf("ssh arguments illegal")
	}
	cfg := options.GetCfg()
	if cfg.CmdSeparator == " " {
		return ks[0], ks[1:], nil
	} else {
		l := strings.Join(ks[1:], " ")
		return ks[0], strings.Split(l, cfg.CmdSeparator), nil
	}
}

func GetSFTPArgs(line string) (string, string, string, error) {
	ks := strings.Fields(line)
	if len(ks) != 3 {
		return "", "", "", fmt.Errorf("sftp arguments illegal")
	}
	return ks[0], ks[1], ks[2], nil
}

func GetCryptArgs(line string) (string, string, error) {
	ks := strings.Fields(line)
	if len(ks) != 2 {
		return "", "", fmt.Errorf("crypt arguments illegal")
	}
	return ks[0], ks[1], nil
}
