package main

import (
	"fmt"
	"github.com/luckywinds/rshell/commands"
	"github.com/luckywinds/rshell/options"
	"github.com/luckywinds/rshell/pkg/checkers"
	"github.com/luckywinds/rshell/pkg/prompt"
	"github.com/luckywinds/rshell/pkg/update"
	"github.com/luckywinds/rshell/pkg/utils"
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

var version = "6.0"
func showIntro() {
	fmt.Println(`
 ______     ______     __  __     ______     __         __
/\  == \   /\  ___\   /\ \_\ \   /\  ___\   /\ \       /\ \
\ \  __<   \ \___  \  \ \  __ \  \ \  __\   \ \ \____  \ \ \____
 \ \_\ \_\  \/\_____\  \ \_\ \_\  \ \_____\  \ \_____\  \ \_____\
  \/_/ /_/   \/_____/   \/_/\/_/   \/_____/   \/_____/   \/_____/
------ Rshell @`+version+` Type "?" or "help" for more information. -----`)
}

func interactiveRun() {
	showIntro()
	opts.CurrentEnv = options.LoadEnv()
	if opts.CurrentEnv.Authname != "" && opts.CurrentEnv.Hostgroupname != "" && opts.CurrentEnv.Port != 0 {
		opts.Cfg.PromptString = "[" + opts.CurrentEnv.Authname + "@" + opts.CurrentEnv.Hostgroupname + ":" + strconv.Itoa(opts.CurrentEnv.Port) + "]# "
	}

	l, err := prompt.New(opts.Cfg, opts.Hostgroups)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
	retry:
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
		case strings.HasPrefix(line, "load "):
			_, a, h, p, err := utils.GetLoadArgs(*opts, line)
			if err != nil {
				fmt.Printf("%v\n", err.Error())
				commands.LoadHelp()
				goto retry
			}
			if !checkers.IsIpv4(h) {
				a = opts.Hostgroupsm[h].Authmethod
				p = opts.Hostgroupsm[h].Sshport
			}
			opts.CurrentEnv.Authname = a
			opts.CurrentEnv.Hostgroupname = h
			opts.CurrentEnv.Port = p
			if err := options.SetEnv(opts.CurrentEnv); err != nil {
				fmt.Printf("%v\n", err.Error())
				commands.LoadHelp()
				goto retry
			}
			opts.Cfg.PromptString = "[" + opts.CurrentEnv.Authname + "@" + opts.CurrentEnv.Hostgroupname + ":" + strconv.Itoa(opts.CurrentEnv.Port) + "]# "
		case strings.HasPrefix(line, "do "):
			if opts.CurrentEnv.Authname == "" || opts.CurrentEnv.Hostgroupname == "" || opts.CurrentEnv.Port == 0 {
				fmt.Printf("%v\n", "Please load correct env first")
				commands.LoadHelp()
				fmt.Println()
				commands.DoHelp()
				goto retry
			}
			_, c, err := utils.GetSSHArgs(line)
			if err != nil {
				fmt.Printf("%v\n", err.Error())
				commands.DoHelp()
				goto retry
			}
			if err = commands.DO(*opts, "do", opts.CurrentEnv.Authname, opts.CurrentEnv.Hostgroupname, opts.CurrentEnv.Port, c); err != nil {
				fmt.Printf("%v\n", err.Error())
				commands.DoHelp()
				goto retry
			}
			for _, value := range c {
				prompt.AddCmd(strings.Trim(value, " "))
			}
		case strings.HasPrefix(line, "sudo "):
			if opts.CurrentEnv.Authname == "" || opts.CurrentEnv.Hostgroupname == "" || opts.CurrentEnv.Port == 0 {
				fmt.Printf("%v\n", "Please load correct env first")
				commands.LoadHelp()
				fmt.Println()
				commands.SudoHelp()
				goto retry
			}
			_, c, err := utils.GetSSHArgs(line)
			if err != nil {
				fmt.Printf("%v\n", err.Error())
				commands.SudoHelp()
				goto retry
			}
			if err = commands.SUDO(*opts, "sudo", opts.CurrentEnv.Authname, opts.CurrentEnv.Hostgroupname, opts.CurrentEnv.Port, c); err != nil {
				fmt.Printf("%v\n", err.Error())
				commands.SudoHelp()
				goto retry
			}
			for _, value := range c {
				prompt.AddCmd(strings.Trim(value, " "))
			}
		case strings.HasPrefix(line, "download "):
			if opts.CurrentEnv.Authname == "" || opts.CurrentEnv.Hostgroupname == "" || opts.CurrentEnv.Port == 0 {
				fmt.Printf("%v\n", "Please load correct env first")
				commands.LoadHelp()
				fmt.Println()
				commands.DownloadHelp()
				goto retry
			}
			_, sf, dd, err := utils.GetSFTPArgs(line)
			if err != nil {
				fmt.Printf("%v\n", err.Error())
				commands.DownloadHelp()
				goto retry
			}
			if err = commands.Download(*opts, "download", opts.CurrentEnv.Authname, opts.CurrentEnv.Hostgroupname, opts.CurrentEnv.Port, sf, dd); err != nil {
				fmt.Printf("%v\n", err.Error())
				commands.DownloadHelp()
				goto retry
			}
			prompt.AddSrcFile(strings.Trim(sf, " "))
			prompt.AddDesDir(strings.Trim(dd, " "))
		case strings.HasPrefix(line, "upload "):
			if opts.CurrentEnv.Authname == "" || opts.CurrentEnv.Hostgroupname == "" || opts.CurrentEnv.Port == 0 {
				fmt.Printf("%v\n", "Please load correct env first")
				commands.LoadHelp()
				fmt.Println()
				commands.UploadHelp()
				goto retry
			}
			_, sf, dd, err := utils.GetSFTPArgs(line)
			if err != nil {
				fmt.Printf("%v\n", err.Error())
				commands.UploadHelp()
				goto retry
			}
			if err = commands.Upload(*opts, "upload", opts.CurrentEnv.Authname, opts.CurrentEnv.Hostgroupname, opts.CurrentEnv.Port, sf, dd); err != nil {
				fmt.Printf("%v\n", err.Error())
				commands.UploadHelp()
				goto retry
			}
			prompt.AddSrcFile(strings.Trim(sf, " "))
			prompt.AddDesDir(strings.Trim(dd, " "))
		case strings.HasPrefix(line, "encrypt_aes "):
			_, t, err := utils.GetCryptArgs(line)
			if err != nil {
				fmt.Printf("%v\n", err)
				commands.EncryptHelp()
				goto retry
			}
			p, err := commands.Encrypt(t)
			if err != nil {
				fmt.Printf("%v\n", err)
				commands.EncryptHelp()
				goto retry
			}
			fmt.Println(p)
		case strings.HasPrefix(line, "decrypt_aes "):
			_, t, err := utils.GetCryptArgs(line)
			if err != nil {
				fmt.Printf("%v\n", err)
				commands.DecryptHelp()
				goto retry
			}
			p, err := commands.Decrypt(t)
			if err != nil {
				fmt.Printf("%v\n", err)
				commands.DecryptHelp()
				goto retry
			}
			fmt.Println(p)
		case line == "?":
			commands.Help()
			goto retry
		case line == "":
			goto retry
		case line == "exit":
			return
		default:
			commands.Help()
			goto retry
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

		for _, stask := range task.Subtasks {
			name := task.Name + "/" + stask.Name
			if stask.Mode == SSH {
				if stask.Sudo {
					if err := commands.SUDO(*opts, name, "", task.Hostgroup, 0, stask.Cmds); err != nil {
						log.Fatalf("ERROR: %s/%s/%s/%v", name, task.Hostgroup, SUDO, err)
					}
				} else {
					if err := commands.DO(*opts, name, "", task.Hostgroup, 0, stask.Cmds); err != nil {
						log.Fatalf("ERROR: %s/%s/%s/%v", name, task.Hostgroup, DO, err)
					}
				}
			} else if stask.Mode == SFTP {
				if stask.FtpType == DOWNLOAD {
					if err := commands.Download(*opts, name, "", task.Hostgroup, 0, stask.SrcFile, stask.DesDir); err != nil {
						log.Fatalf("ERROR: %s/%s/%s/%v", name, task.Hostgroup, DOWNLOAD, err)
					}
				} else if stask.FtpType == UPLOAD {
					if err := commands.Upload(*opts, name, "", task.Hostgroup, 0, stask.SrcFile, stask.DesDir); err != nil {
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
