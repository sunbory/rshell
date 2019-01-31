package main

import (
	"fmt"
	"github.com/luckywinds/rshell/modes/client"
	"github.com/luckywinds/rshell/options"
	"github.com/luckywinds/rshell/outputs"
	"github.com/luckywinds/rshell/pkg/checkers"
	"github.com/luckywinds/rshell/pkg/prompt"
	"github.com/luckywinds/rshell/pkg/rlog"
	"github.com/luckywinds/rshell/pkg/update"
	"github.com/luckywinds/rshell/plugins/decrypt"
	"github.com/luckywinds/rshell/plugins/do"
	"github.com/luckywinds/rshell/plugins/download"
	"github.com/luckywinds/rshell/plugins/encrypt"
	"github.com/luckywinds/rshell/plugins/load"
	"github.com/luckywinds/rshell/plugins/sudo"
	"github.com/luckywinds/rshell/plugins/upload"
	"log"
	"strconv"
	"strings"
)

const (
	Interactive string = "interactive"
	SCRIPT      string = "script"
	DOWNLOAD    string = "download"
	UPLOAD      string = "upload"
	SSH         string = "ssh"
	SFTP        string = "sftp"
	AES         string = "aes"
	DO          string = "do"
	SUDO        string = "sudo"
)

var opts = options.New()

func main() {
	setup()

	if opts.IsCommandlineMode {
		commandRun()
	} else if opts.IsScriptMode {
		scriptRun()
	} else {
		interactiveRun()
	}
}

func setup() {
	rlog.Info.Printf("Init options : %+v", opts)

	go update.Update(opts.Cfg, version)

	client.SetupCache(opts.Cfg.Connecttimeout)
}

var version = "8.0"
func showIntro() {
	fmt.Println(`
 ______     ______     __  __     ______     __         __
/\  == \   /\  ___\   /\ \_\ \   /\  ___\   /\ \       /\ \
\ \  __<   \ \___  \  \ \  __ \  \ \  __\   \ \ \____  \ \ \____
 \ \_\ \_\  \/\_____\  \ \_\ \_\  \ \_____\  \ \_____\  \ \_____\
  \/_/ /_/   \/_____/   \/_/\/_/   \/_____/   \/_____/   \/_____/
------ Rshell @`+version+` Type "?" or "help" for more information. -----
{The Correct Step: 1.help -> 2.load -> 3.*/sudo/download/upload}`)
}

func commandRun() {
	rlog.Info.Printf("current env: %+v", opts.CurrentEnv)

	if strings.HasPrefix(opts.Line, "encrypt ") || strings.HasPrefix(opts.Line, "decrypt ") {
		runCommand(opts.Line)
	} else {
		if checkers.IsIpv4(opts.CurrentEnv.Hostgroupname) {
			if _, ok := opts.Authsm[opts.CurrentEnv.Authname]; !ok {
				rlog.Error.Fatalf("auth name [%s] not found", opts.CurrentEnv.Authname)
			}
			if opts.CurrentEnv.Port <= 0 || opts.CurrentEnv.Port >= 65535 {
				rlog.Error.Fatalf("ssh port [%d] illegal", opts.CurrentEnv.Port)
			}
		} else {
			if _, ok := opts.Hostgroupsm[opts.CurrentEnv.Hostgroupname]; !ok {
				rlog.Error.Fatalf("host group [%s] not found", opts.CurrentEnv.Hostgroupname)
			} else {
				opts.CurrentEnv.Authname = opts.Hostgroupsm[opts.CurrentEnv.Hostgroupname].Authmethod
				opts.CurrentEnv.Port = opts.Hostgroupsm[opts.CurrentEnv.Hostgroupname].Sshport
			}
		}

		runCommand(opts.Line)
	}
}
func interactiveRun() {
	showIntro()
	opts.CurrentEnv = options.LoadEnv()
	if opts.CurrentEnv.Authname != "" && opts.CurrentEnv.Hostgroupname != "" && opts.CurrentEnv.Port != 0 {
		opts.Cfg.PromptString = "[" + opts.CurrentEnv.Authname + "@" + opts.CurrentEnv.Hostgroupname + ":" + strconv.Itoa(opts.CurrentEnv.Port) + "]# "
		prompt.AddHostgroup("-H" + opts.CurrentEnv.Hostgroupname)
	}

	l, err := prompt.New(opts.Cfg, opts.Hostgroups)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		line, err := prompt.Prompt(l, opts.Cfg)
		if err == prompt.ErrPromptAborted {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == prompt.ErrPromptError {
			break
		}

		rlog.Debug.Printf("line: %s", line)
		if runCommand(line) {
			break
		}
	}
}

func runCommand(line string) bool {
	line = strings.TrimLeft(line, " ")
	switch {
	case strings.HasPrefix(line, "load ") || line == "load":
		rlog.Info.Printf("current env: %+v", opts.CurrentEnv)
		a, h, p, err := load.Command(*opts, line)
		if err != nil {
			rlog.Error.Printf("load: %v", err)
			load.Help()
		} else {
			if !checkers.IsIpv4(h) {
				a = opts.Hostgroupsm[h].Authmethod
				p = opts.Hostgroupsm[h].Sshport
			}
			opts.CurrentEnv.Authname = a
			opts.CurrentEnv.Hostgroupname = h
			opts.CurrentEnv.Port = p
			if err := options.SetEnv(opts.CurrentEnv); err != nil {
				rlog.Error.Printf("set env: %v", err)
				load.Help()
			} else {
				opts.Cfg.PromptString = "[" + opts.CurrentEnv.Authname + "@" + opts.CurrentEnv.Hostgroupname + ":" + strconv.Itoa(opts.CurrentEnv.Port) + "]# "
			}
		}
	case strings.HasPrefix(line, "sudo ") || line == "sudo":
		rlog.Info.Printf("current env: %+v", opts.CurrentEnv)
		ret, err := sudo.Command(*opts, line)
		if err != nil {
			rlog.Error.Printf("sudo: %v", err)
			sudo.Help()
		} else {
			for _, value := range ret {
				prompt.AddCmd(strings.Trim(value, " "))
			}
		}
		outputs.End()
	case strings.HasPrefix(line, "download ") || line == "download":
		rlog.Info.Printf("current env: %+v", opts.CurrentEnv)
		ret, err := download.Command(*opts, line)
		if err != nil {
			rlog.Error.Printf("download: %v", err)
			download.Help()
		} else {
			prompt.AddSrcFile(strings.Trim(ret[0], " "))
			prompt.AddDesDir(strings.Trim(ret[1], " "))
		}
		outputs.End()
	case strings.HasPrefix(line, "upload ") || line == "upload":
		rlog.Info.Printf("current env: %+v", opts.CurrentEnv)
		ret, err := upload.Command(*opts, line)
		if err != nil {
			rlog.Error.Printf("upload: %v", err)
			upload.Help()
		} else {
			prompt.AddSrcFile(strings.Trim(ret[0], " "))
			prompt.AddDesDir(strings.Trim(ret[1], " "))
		}
		outputs.End()
	case strings.HasPrefix(line, "encrypt ") || line == "encrypt":
		rlog.Info.Printf("current env: %+v", opts.CurrentEnv)
		ret, err := encrypt.Command(*opts, line)
		if err != nil {
			rlog.Error.Printf("encrypt: %v", err)
			encrypt.Help()
		} else {
			fmt.Println(ret)
		}
	case strings.HasPrefix(line, "decrypt ") || line == "decrypt":
		rlog.Info.Printf("current env: %+v", opts.CurrentEnv)
		ret, err := decrypt.Command(*opts, line)
		if err != nil {
			rlog.Error.Printf("decrypt: %v", err)
			decrypt.Help()
		} else {
			fmt.Println(ret)
		}
	case line == "":
	case line == "?" || line == "help":
		load.Help()
		fmt.Println()
		do.Help()
		fmt.Println()
		sudo.Help()
		fmt.Println()
		download.Help()
		fmt.Println()
		upload.Help()
		fmt.Println()
		encrypt.Help()
		fmt.Println()
		decrypt.Help()
		fmt.Println()
		fmt.Println(`exit
    --- Exit rshell
?
    --- Help`)
	case line == "exit":
		return true
	default:
		rlog.Info.Printf("current env: %+v", opts.CurrentEnv)
		ret, err := do.Command(*opts, line)
		if err != nil {
			rlog.Error.Printf("do: %v", err)
			do.Help()
		} else {
			for _, value := range ret {
				prompt.AddCmd(strings.Trim(value, " "))
			}
		}
		outputs.End()
	}
	return false
}

func scriptRun() {
	for _, task := range opts.Tasks.Ts {
		rlog.Debug.Printf("task: %+v", task)
		if !checkers.CheckTaskName(task.Name) || !checkers.CheckHostgroupName(task.Hostgroup) {
			log.Fatal("The task's name or hostgroup illegal.")
		}

		if len(task.Subtasks) == 0 {
			log.Fatal("SSH or SFTP Tasks empty.")
		}

		if _, ok := opts.Hostgroupsm[task.Hostgroup]; !ok {
			rlog.Error.Fatalf("host group [%s] not found", task.Hostgroup)
		} else {
			opts.CurrentEnv.Hostgroupname = task.Hostgroup
			opts.CurrentEnv.Authname = opts.Hostgroupsm[opts.CurrentEnv.Hostgroupname].Authmethod
			opts.CurrentEnv.Port = opts.Hostgroupsm[opts.CurrentEnv.Hostgroupname].Sshport
		}

		rlog.Info.Printf("current env: %+v", opts.CurrentEnv)
		for _, stask := range task.Subtasks {
			rlog.Debug.Printf("stask: %+v", stask)
			if !checkers.CheckTaskName(stask.Name) {
				log.Fatal("The stask's name illegal.")
			}
			if stask.Mode == SSH {
				if stask.Sudo {
					if err := sudo.Script(*opts, task.Name, stask); err != nil {
						rlog.Error.Fatalf("%s/%s/%s/%v", task.Name, task.Hostgroup, SUDO, err)
					}
				} else {
					if err := do.Script(*opts, task.Name, stask); err != nil {
						rlog.Error.Fatalf("%s/%s/%s/%v", task.Name, task.Hostgroup, DO, err)
					}
				}
			} else if stask.Mode == SFTP {
				if stask.FtpType == DOWNLOAD {
					if err := download.Script(*opts, task.Name, stask); err != nil {
						rlog.Error.Fatalf("%s/%s/%s/%v", task.Name, task.Hostgroup, DOWNLOAD, err)
					}
				} else if stask.FtpType == UPLOAD {
					if err := upload.Script(*opts, task.Name, stask); err != nil {
						rlog.Error.Fatalf("%s/%s/%s/%v", task.Name, task.Hostgroup, UPLOAD, err)
					}
				} else {
					rlog.Error.Fatalf("%s/%s/%s/%s", task.Name, task.Hostgroup, stask.FtpType, "Not support")
				}
			} else {
				rlog.Error.Fatalf("%s/%s/%s/%s", task.Name, task.Hostgroup, stask.Mode, "Not support")
			}
		}
	}
	outputs.End()
}
