package core

import (
	"fmt"
	"github.com/luckywinds/rshell/modes/sftp"
	"github.com/luckywinds/rshell/modes/ssh"
	"github.com/luckywinds/rshell/options"
	"github.com/luckywinds/rshell/outputs"
	"github.com/luckywinds/rshell/pkg/checkers"
	"github.com/luckywinds/rshell/pkg/crypt"
	"github.com/luckywinds/rshell/pkg/filters"
	"github.com/luckywinds/rshell/types"
	"path"
	"strings"
)

func GetArgFields(line, keyword, sep string) ([]string) {
	as := strings.Fields(strings.TrimLeft(line, keyword))
	if sep == " " {
		return as
	} else {
		return strings.Split(strings.Join(as, " "), sep)
	}
}

func GetAuthHostgroup(o options.Options) (types.Auth, types.Hostgroup) {
	var au types.Auth
	var hg types.Hostgroup

	au = o.Authsm[o.CurrentEnv.Authname]
	if checkers.IsIpv4(o.CurrentEnv.Hostgroupname) {
		hg = types.Hostgroup{
			Groupname:  "TEMPHOST",
			Authmethod: o.CurrentEnv.Authname,
			Sshport:    o.CurrentEnv.Port,
			Hosts:      nil,
			Groups:     nil,
			Hostranges: nil,
			Ips:        []string{o.CurrentEnv.Hostgroupname},
		}
	} else {
		hg = o.Hostgroupsm[o.CurrentEnv.Hostgroupname]
	}
	hg.Sshport = o.CurrentEnv.Port

	return au, hg
}

func EnvCheck(env options.CurrentEnv) error {
	if env.Hostgroupname == "" || env.Authname == "" || env.Port == 0 {
		return fmt.Errorf("current env not found, please run step[load] first")
	}
	return nil
}

func SecurityCheck(bcmds []types.BlackCmd, cmds []string) error {
	for _, value := range cmds {
		if filters.IsBlackCmd(value, bcmds) {
			return fmt.Errorf("DANGER: [%s] in black command list.", value)
		}
	}
	return nil
}

func GetPlainPassword(c types.Cfg, au types.Auth) (nau types.Auth, err error) {
	if c.Passcrypttype != "" {
		if au.Password, err = crypt.AesDecrypt(au.Password, c); err != nil {
			return au, fmt.Errorf("decrypt password error [%v] crypt type is [%s]", err, c.Passcrypttype)
		}
		if au.Passphrase, err = crypt.AesDecrypt(au.Passphrase, c); err != nil {
			return au, fmt.Errorf("decrypt password error [%v] crypt type is [%s]", err, c.Passcrypttype)
		}
		if au.Sudopass, err = crypt.AesDecrypt(au.Sudopass, c); err != nil {
			return au, fmt.Errorf("decrypt password error [%v] crypt type is [%s]", err, c.Passcrypttype)
		}
	}
	return au,nil
}

func RunSshCommands(Concurrency int, actionname, actiontype string, au types.Auth, hg types.Hostgroup, cmds []string) {
	limit := make(chan bool, Concurrency)
	defer close(limit)

	taskchs := make(chan types.Hostresult, len(hg.Ips))
	defer close(taskchs)

	for _, ip := range hg.Ips {
		limit <- true
		go func(actionname, actiontype, groupname, host string, port int, user, pass, keyname, passphrase, sudotype, sudopass string, timeout int, ciphers, cmds []string) {
			var stdout, stderr string
			var err error
			switch actiontype {
			case "do":
				stdout, stderr, err = ssh.DO(groupname, host, port, user, pass, keyname, passphrase, timeout, ciphers, cmds)
			case "sudo":
				stdout, stderr, err = ssh.SUDO(groupname, host, port, user, pass, keyname, passphrase, sudotype, sudopass, timeout, ciphers, cmds)
			default:
				err = fmt.Errorf("action not supported")
			}
			var result types.Hostresult
			result.Actionname = actionname
			result.Actiontype = actiontype
			result.Groupname = groupname
			result.Hostaddr = host
			result.Stdout = stdout
			result.Stderr = stderr
			if err != nil {
				result.Error = err.Error()
			}
			taskchs <- result
			<-limit
		}(actionname, actiontype, hg.Groupname, ip, hg.Sshport, au.Username, au.Password, au.Privatekey, au.Passphrase, au.Sudotype, au.Sudopass,60, []string{}, cmds)
	}

	outputs.Output(taskchs, hg)
}

func RunSftpCommands(Concurrency int, actionname, actiontype string, au types.Auth, hg types.Hostgroup, srcFilePath, desDirPath string) {
	limit := make(chan bool, Concurrency)
	defer close(limit)

	taskchs := make(chan types.Hostresult, len(hg.Ips))
	defer close(taskchs)

	for _, ip := range hg.Ips {
		limit <- true
		go func(actionname, actiontype, groupname, host string, port int, user, pass, keyname, passphrase string, timeout int, ciphers []string, srcFilePath, desDirPath string) {
			if !strings.HasSuffix(desDirPath, "/") {
				desDirPath = desDirPath + "/"
			}
			var sfs []string
			var err error
			switch actiontype {
			case "download":
				sfs, err = sftp.Download(groupname, host, port, user, pass, keyname, passphrase, timeout, ciphers, srcFilePath, desDirPath)
			case "upload":
				sfs, err = sftp.Upload(groupname, host, port, user, pass, keyname, passphrase, timeout, ciphers, srcFilePath, desDirPath)
			default:
				err = fmt.Errorf("action not supported")
			}
			var result types.Hostresult
			result.Actionname = actionname
			result.Actiontype = actiontype
			result.Groupname = groupname
			result.Hostaddr = host
			if err == nil {
				result.Stdout = "DOWNLOAD Success [" + srcFilePath + " -> " + path.Join(desDirPath, hg.Groupname) + "] :\n" + strings.Join(sfs, "\n")
			} else {
				result.Stderr = "DOWNLOAD Failed [" + srcFilePath + " -> " + path.Join(desDirPath, hg.Groupname) + "] @" + err.Error()
			}
			taskchs <- result
			<-limit
		}(actionname, actiontype, hg.Groupname, ip, hg.Sshport, au.Username, au.Password, au.Privatekey, au.Passphrase, 60, []string{}, srcFilePath, desDirPath)
	}

	outputs.Output(taskchs, hg)
}