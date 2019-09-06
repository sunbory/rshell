package core

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/luckywinds/rshell/types"
    "github.com/luckywinds/rshell/options"
	"github.com/luckywinds/rshell/outputs"
	"github.com/luckywinds/rshell/pkg/checkers"
	"github.com/luckywinds/rshell/pkg/crypt"
	"github.com/luckywinds/rshell/pkg/rlog"
	"github.com/luckywinds/rshell/rclient"
)

func GetArgFields(line, keyword, sep string) []string {
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
	if checkers.ValidIP(o.CurrentEnv.Hostgroupname) {
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
		if checkers.IsBlackCmd(value, bcmds) {
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
	return au, nil
}

func GetProxyConfig (hostname string) (*rclient.RConfig, error) {

	o := options.New()
	
	hg, err := options.GetHostgroupByname(hostname)
	if err != nil {
		return nil, err
	}

	au, err := options.GetAuthByname(hg.Authmethod)
	if err != nil {
		return nil, err
	}

	au, err = GetPlainPassword(options.GetCfg(), au)
	if err != nil {
		return nil, err
	}

	return &rclient.RConfig {
		Groupname  : hostname,
		Host       : hg.Ips[0],
		Port       : hg.Sshport,  
		User       : au.Username,
		Key        : au.Privatekey,
		Password   : au.Password,
		Passphrase : au.Passphrase,
		Ciphers    : o.Cfg.Sshciphers,
		Sudotype   : au.Sudotype,
		Sudopass   : au.Sudopass,
	}, nil

}

func RunSshCommands(cfg types.Cfg, actionname, actiontype string, au types.Auth, hg types.Hostgroup, cmds []string) {
	rlog.Info.Printf("concurrency: %d, actionname: %s, actiontype: %s, au: %+v, hg: %+v, cmds: %#v", cfg.Concurrency, actionname, actiontype, au, hg, cmds)

	limit := make(chan bool, cfg.Concurrency)
	defer close(limit)

	taskchs := make(chan types.Hostresult, len(hg.Ips))
	defer close(taskchs)

	ctx, _ := context.WithTimeout(context.Background(), time.Duration(cfg.Tasktimeout)*time.Second)

	for _, ip := range hg.Ips {
		limit <- true
		
		rConfig := &rclient.RConfig {
			Groupname  : hg.Groupname,
			Host       : ip,
			Port       : hg.Sshport ,  
			User       : au.Username,
			Key        : au.Privatekey,
			Password   : au.Password,
			Passphrase : au.Passphrase,
			Ciphers    : cfg.Sshciphers,
			Sudotype   : au.Sudotype,
			Sudopass   : au.Sudopass,
		} 

		if hg.Proxy != "" {
			rConfig.Proxy, _ = GetProxyConfig(hg.Proxy)
		}

		go func(ctx context.Context, actionname, actiontype string, rconfig *rclient.RConfig, cmds []string) {
			var stdout, stderr string

			var result = types.Hostresult {
				Actiontype : actiontype,
				Groupname  : rconfig.Groupname,
				Hostaddr   : rconfig.Host,
			}
			
			ssh, err := rclient.New(rconfig)
			if err != nil {
				result.Error = err.Error()
				taskchs <- result
				return 
			}

			switch actiontype {
			case "do":
				stdout, stderr, err = ssh.DO(cmds)
			case "sudo":
				stdout, stderr, err = ssh.SUDO(cmds)
				stderr = strings.Replace(stderr, rconfig.Sudopass, "******", -1)
			default:
				err = fmt.Errorf("action not supported")
			}
			result.Stdout = stdout
			result.Stderr = stderr
			if err != nil {
				result.Error = err.Error()
			}

			select {
			case <-ctx.Done():
				rlog.Warn.Printf("ACTION TIMEOUT [%v:%v:%v:%v:%v]", actionname, actiontype, rconfig.Groupname, rconfig.Host, rconfig.Port)
				return
			default:
				taskchs <- result
				<-limit
			}
		}(ctx, actionname, actiontype, rConfig, cmds)
	}

	outputs.Output(actionname, actiontype, taskchs, hg)
}

func RunSftpCommands(cfg types.Cfg, actionname, actiontype string, au types.Auth, hg types.Hostgroup, srcFilePath, desDirPath string) {
	rlog.Info.Printf("concurrency: %d, actionname: %s, actiontype: %s, au: %+v, hg: %+v, srcFilePath: %s, desDirPath: %s", cfg.Concurrency, actionname, actiontype, au, hg, srcFilePath, desDirPath)

	limit := make(chan bool, cfg.Concurrency)
	defer close(limit)

	taskchs := make(chan types.Hostresult, len(hg.Ips))
	defer close(taskchs)

	ctx, _ := context.WithTimeout(context.Background(), time.Duration(cfg.Tasktimeout)*time.Second)

	for _, ip := range hg.Ips {
		limit <- true

		rConfig := &rclient.RConfig {
			Groupname  : hg.Groupname,
			Host       : ip,
			Port       : hg.Sshport,   
			User       : au.Username,
			Key        : au.Privatekey,
			Password   : au.Password,
			Passphrase : au.Passphrase,
			Ciphers    : cfg.Sshciphers,
			Sudotype   : au.Sudotype,
			Sudopass   : au.Sudopass,
		} 

		if hg.Proxy != ""{
			rConfig.Proxy, _ = GetProxyConfig(hg.Proxy)
		}
		
		go func(ctx context.Context, actionname, actiontype string, rconfig *rclient.RConfig, maxPacketSize int, srcFilePath, desDirPath string) {
			// 			if !strings.HasSuffix(desDirPath, "/") {
			// 				desDirPath = desDirPath + "/"
			// 			}
			var sfs []string

			var result = types.Hostresult {
				Actiontype : actiontype,
				Groupname  : rconfig.Groupname,
				Hostaddr   : rconfig.Host,
			}
			
			sftp, err := rclient.New(rconfig)
			if err != nil {
				result.Error = err.Error()
				taskchs <- result
				return 
			}

			switch actiontype {
			case "download":
				sfs, err = sftp.Download(srcFilePath, desDirPath, maxPacketSize)
			case "upload":
				sfs, err = sftp.Upload(srcFilePath, desDirPath, maxPacketSize)
			default:
				err = fmt.Errorf("action not supported")
			}
			if err == nil {
				result.Stdout = "SUCCESS [" + srcFilePath + " -> " + desDirPath + "] :\n" + strings.Join(sfs, "\n")
			} else {
				result.Stderr = "FAILED [" + srcFilePath + " -> " + desDirPath + "] @" + err.Error()
			}

			select {
			case <-ctx.Done():
				rlog.Warn.Printf("ACTION TIMEOUT [%v:%v:%v:%v:%v]", actionname, actiontype, rconfig.Groupname, rconfig.Host, rconfig.Port)
				return
			default:
				taskchs <- result
				<-limit
			}
		}(ctx, actionname, actiontype, rConfig, cfg.Sftppacketsize, srcFilePath, desDirPath)
	}

	outputs.Output(actionname, actiontype, taskchs, hg)
}
