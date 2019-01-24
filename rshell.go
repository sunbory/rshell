package main

import (
	"fmt"
	"github.com/luckywinds/rshell/options"
	"github.com/luckywinds/rshell/pkg/checkers"
	"github.com/luckywinds/rshell/pkg/prompt"
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
	go update.Update(opts.Cfg, version)

	if !opts.IsScriptMode {
		interactiveRun()
	} else {
		scriptRun()
	}
}

var version = "6.1"
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

		line = strings.TrimLeft(line, " ")
		switch {
		case strings.HasPrefix(line, "load ") || line == "load":
			a, h, p, err := load.Command(*opts, line)
			if err != nil {
				fmt.Println(err)
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
					fmt.Printf("%v\n", err.Error())
					load.Help()
				} else {
					opts.Cfg.PromptString = "[" + opts.CurrentEnv.Authname + "@" + opts.CurrentEnv.Hostgroupname + ":" + strconv.Itoa(opts.CurrentEnv.Port) + "]# "
				}
			}
		case strings.HasPrefix(line, "sudo ") || line == "sudo":
			ret, err := sudo.Command(*opts, line)
			if err != nil {
				fmt.Println(err)
				sudo.Help()
			} else {
				for _, value := range ret {
					prompt.AddCmd(strings.Trim(value, " "))
				}
			}
		case strings.HasPrefix(line, "download ") || line == "download":
			ret, err := download.Command(*opts, line)
			if err != nil {
				fmt.Println(err)
				download.Help()
			} else {
				prompt.AddSrcFile(strings.Trim(ret[0], " "))
				prompt.AddDesDir(strings.Trim(ret[1], " "))
			}
		case strings.HasPrefix(line, "upload ") || line == "upload":
			ret, err := upload.Command(*opts, line)
			if err != nil {
				fmt.Println(err)
				upload.Help()
			} else {
				prompt.AddSrcFile(strings.Trim(ret[0], " "))
				prompt.AddDesDir(strings.Trim(ret[1], " "))
			}
		case strings.HasPrefix(line, "encrypt ") || line == "encrypt":
			ret, err := encrypt.Command(*opts, line)
			if err != nil {
				fmt.Println(err)
				encrypt.Help()
			} else {
				fmt.Println(ret)
			}
		case strings.HasPrefix(line, "decrypt ") || line == "decrypt":
			ret, err := decrypt.Command(*opts, line)
			if err != nil {
				fmt.Println(err)
				decrypt.Help()
			} else {
				fmt.Println(ret)
			}
		case line == "":
		case line == "?":
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
			return
		default:
			ret, err := do.Command(*opts, line)
			if err != nil {
				fmt.Println(err)
				do.Help()
			} else {
				for _, value := range ret {
					prompt.AddCmd(strings.Trim(value, " "))
				}
			}
		}
	}
}

func scriptRun() {
	for _, task := range opts.Tasks.Ts {
		if task.Name == "" || task.Hostgroup == "" {
			log.Fatal("The task's name or hostgroup empty.")
		}

		if len(task.Subtasks) == 0 {
			log.Fatal("SSH or SFTP Tasks empty.")
		}

		opts.CurrentEnv.Hostgroupname = task.Hostgroup
		opts.CurrentEnv.Port = opts.Hostgroupsm[task.Hostgroup].Sshport
		opts.CurrentEnv.Authname = opts.Hostgroupsm[task.Hostgroup].Authmethod

		for _, stask := range task.Subtasks {
			name := task.Name + "/" + stask.Name
			if stask.Mode == SSH {
				if stask.Sudo {
					if err := sudo.Script(*opts, name, stask); err != nil {
						log.Fatalf("ERROR: %s/%s/%s/%v", name, task.Hostgroup, SUDO, err)
					}
				} else {
					if err := do.Script(*opts, name, stask); err != nil {
						log.Fatalf("ERROR: %s/%s/%s/%v", name, task.Hostgroup, DO, err)
					}
				}
			} else if stask.Mode == SFTP {
				if stask.FtpType == DOWNLOAD {
					if err := download.Script(*opts, name, stask); err != nil {
						log.Fatalf("ERROR: %s/%s/%s/%v", name, task.Hostgroup, DOWNLOAD, err)
					}
				} else if stask.FtpType == UPLOAD {
					if err := upload.Script(*opts, name, stask); err != nil {
						log.Fatalf("ERROR: %s/%s/%s/%v", name, task.Hostgroup, UPLOAD, err)
					}
				} else {
					log.Fatalf("ERROR: %s/%s/%s/%s", name, task.Hostgroup, stask.FtpType, "Not support")
				}
			} else {
				log.Fatalf("ERROR: %s/%s/%s/%s", name, task.Hostgroup, stask.Mode, "Not support")
			}
		}
	}
}
