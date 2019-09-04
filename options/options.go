package options

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"

	"github.com/luckywinds/rshell/pkg/checkers"
	"github.com/luckywinds/rshell/pkg/prompt"
	. "github.com/luckywinds/rshell/types"
	"gopkg.in/yaml.v2"
)

type Options struct {
	Cfg Cfg

	Hostgroups  Hostgroups
	Hostgroupsm map[string]Hostgroup

	Auths  Auths
	Authsm map[string]Auth

	Tasks  Tasks
	Values map[string]interface{}

	IsScriptMode      bool
	IsCommandlineMode bool

	CurrentEnv CurrentEnv

	Line string
}

type CurrentEnv struct {
	Hostgroupname string
	Authname      string
	Port          int
}

var (
	cfgFile      = path.Join(".rshell", "cfg.yaml")
	hostsFile    = path.Join(".rshell", "hosts.yaml")
	authFile     = path.Join(".rshell", "auth.yaml")
	script       = flag.String("f", "", "The script yaml.")
	scriptValues = flag.String("v", "", "The script values yaml.")
	authName     = flag.String("A", "", "The auth method name.")
	hostName     = flag.String("H", "", "The host name.")
	sshport      = flag.Int("P", 22, "The ssh port.")
	cmdline      = flag.String("L", "", "The command line.")
)

func init() {
	os.Mkdir(".rshell", os.ModeDir)
	initFlag()

	initCfg()
	initHostgroups()
	initAuths()
}

func initFlag() {
	cmd := os.Args[0]
	flag.Usage = func() {
		fmt.Println(`Usage:`, cmd, `[<options>]

Options:`)
		flag.PrintDefaults()
	}
	flag.Parse()
}

var cfg Cfg

func initCfg() {
	c, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		log.Fatalf("Can not find cfg file[%s].", cfgFile)
	}
	err = yaml.Unmarshal(c, &cfg)
	if err != nil {
		log.Fatalf("YAML[%s] Unmarshal error: %v", cfgFile, err)
	}
}

var option *Options

func New() *Options {

	if option != nil {
		return option
	}

	var cfg = GetCfg()
	var auths, authsm = GetAuths()
	var hostgroups, hostgroupsm = GetHostgroups()
	var tasks = GetTasks()

	for key, _ := range authsm {
		prompt.AddAuth("-A" + key)
	}

	for key, _ := range hostgroupsm {
		prompt.AddHostgroup("-H" + key)
	}

	option := &Options{
		Cfg:               cfg,
		Hostgroups:        hostgroups,
		Hostgroupsm:       hostgroupsm,
		Auths:             auths,
		Authsm:            authsm,
		Tasks:             tasks,
		Values:            nil,
		IsScriptMode:      IsScriptMode(),
		IsCommandlineMode: IsCommandlineMode(),
		CurrentEnv: CurrentEnv{
			Hostgroupname: *hostName,
			Authname:      *authName,
			Port:          *sshport,
		},
		Line: *cmdline,
	}

	return option
}

func IsScriptMode() bool {
	return *script != ""
}

func IsCommandlineMode() bool {
	return *cmdline != ""
}

func GetCfg() Cfg {
	if cfg.Concurrency == 0 {
		cfg.Concurrency = 10
	} else if cfg.Concurrency < 0 || cfg.Concurrency > 100 {
		log.Fatalf("Config Concurrency illegal [%d] not in (0, 100].", cfg.Concurrency)
	}
	if cfg.Tasktimeout == 0 {
		cfg.Tasktimeout = 300
	} else if cfg.Tasktimeout < 0 || cfg.Tasktimeout > 86400 {
		log.Fatalf("Config Tasktimeout illegal [%d] not in (0, 86400].", cfg.Tasktimeout)
	}
	if cfg.Connecttimeout == 0 {
		cfg.Connecttimeout = 3600
	} else if cfg.Connecttimeout < 0 || cfg.Connecttimeout > 86400 {
		log.Fatalf("Config Connecttimeout illegal [%d] not in (0, 86400].", cfg.Connecttimeout)
	}
	if cfg.Sftppacketsize == 0 {
		cfg.Sftppacketsize = 32768
	} else if cfg.Sftppacketsize < 1 || cfg.Sftppacketsize > 32768 {
		log.Fatalf("Config Sftppacketsize illegal [%d] not in (0, 32768].", cfg.Sftppacketsize)
	}
	if cfg.CmdSeparator == "" {
		cfg.CmdSeparator = ";"
	} else if len(cfg.CmdSeparator) != 1 {
		log.Fatalf("Config CmdSeparator illegal [%s] not one char.", cfg.CmdSeparator)
	}
	if cfg.PromptString == "" {
		cfg.PromptString = "rshell: "
	} else if len(cfg.PromptString) > 20 {
		log.Fatalf("Config PromptString illegal [%s] length > 20.", cfg.PromptString)
	}
	if cfg.Outputtype == "" {
		cfg.Outputtype = "text"
	} else if cfg.Outputtype != "text" && cfg.Outputtype != "json" && cfg.Outputtype != "yaml" {
		log.Fatalf("Config Outputtype illegal [%s] not in [text, json, yaml].", cfg.Outputtype)
	}
	if cfg.Hostgroupsize == 0 {
		cfg.Hostgroupsize = 200
	} else if cfg.Hostgroupsize < 0 || cfg.Hostgroupsize > 1000 {
		log.Fatalf("Config Hostgroupsize illegal [%d] not in (0, 1000].", cfg.Hostgroupsize)
	}
	if cfg.Passcrypttype != "" {
		if cfg.Passcrypttype != "aes" {
			log.Fatalf("Config Passcrypttype illegal [%s] not in [aes].", cfg.Passcrypttype)
		} else if len(cfg.Passcryptkey) != 32 {
			log.Fatalf("Config Passcryptkey illegal [%s] length != 32.", cfg.Passcryptkey)
		}
	}
	if cfg.HistoryFile == "" {
		cfg.HistoryFile = ".rshell/rshell.history"
	} else if len(cfg.HistoryFile) > 200 {
		log.Fatalf("Config HistoryFile illegal [%s] length > 200.", cfg.HistoryFile)
	}

	cfg.Updateserver = append(cfg.Updateserver, "ftp://siag9x002128631")

	if len(cfg.Sshciphers) == 0 {
		cfg.Sshciphers = append(cfg.Sshciphers, "aes128-ctr")
		cfg.Sshciphers = append(cfg.Sshciphers, "aes192-ctr")
		cfg.Sshciphers = append(cfg.Sshciphers, "aes256-ctr")
		cfg.Sshciphers = append(cfg.Sshciphers, "aes128-gcm@openssh.com")
		cfg.Sshciphers = append(cfg.Sshciphers, "arcfour256")
		cfg.Sshciphers = append(cfg.Sshciphers, "arcfour128")
		cfg.Sshciphers = append(cfg.Sshciphers, "aes128-cbc")
		cfg.Sshciphers = append(cfg.Sshciphers, "3des-cbc")
		cfg.Sshciphers = append(cfg.Sshciphers, "aes192-cbc")
		cfg.Sshciphers = append(cfg.Sshciphers, "aes256-cbc")
	}
	return cfg
}

var hostgroups Hostgroups
var hostgroupsm = make(map[string]Hostgroup)

func initHostgroups() {
	h, err := ioutil.ReadFile(hostsFile)
	if err != nil {
		log.Fatalf("Can not find hosts file[%s].", hostsFile)
	}
	err = yaml.Unmarshal(h, &hostgroups)
	if err != nil {
		log.Fatalf("YAML[%s] Unmarshal error: %v", hostsFile, err)
	}
}

func incIp(s string) string {
	ip := net.ParseIP(s)
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
	return ip.String()
}

func parseHosts(hg Hostgroup) Hostgroup {
	for _, h := range hg.Hosts {
		if !checkers.ValidIP(h) {
			log.Fatalf("IP illegal [%s/%s].", hg.Groupname, h)
		}
	}
	hg.Ips = append(hg.Ips, hg.Hosts...)
	return hg
}

func parseHostrange(hg Hostgroup) Hostgroup {
	for _, hr := range hg.Hostranges {
		if !checkers.ValidIP(hr.From) || !checkers.ValidIP(hr.To) || hr.From == hr.To {
			log.Fatalf("IP Range illegal [%s/%s-%s].", hg.Groupname, hr.From, hr.To)
		}
		temp := []string{hr.From}
		found := false
		nip := hr.From
		count := 0
		for {
			count++
			nip = incIp(nip)
			if nip == hr.To {
				found = true
				temp = append(temp, nip)
				break
			}
			if count > cfg.Hostgroupsize && !found {
				log.Fatalf("Too Large Not Found. IP Range illegal [%s/%s-%s].", hg.Groupname, hr.From, hr.To)
			}
			temp = append(temp, nip)
		}
		if found {
			hg.Ips = append(hg.Ips, temp...)
		} else {
			log.Fatalf("IP Range illegal [%s/%s-%s].", hg.Groupname, hr.From, hr.To)
		}
	}
	return hg
}

func parseHostgroups(hgs Hostgroups) Hostgroups {
	var tmphg = Hostgroups{}
	for _, hg := range hgs.Hgs {
		hg = parseHosts(hg)
		hg = parseHostrange(hg)
		tmphg.Hgs = append(tmphg.Hgs, hg)
	}

	var hgmap = make(map[string]Hostgroup)
	for _, value := range tmphg.Hgs {
		if !checkers.CheckHostgroupName(value.Groupname) {
			log.Fatalf("Hostgroup name [%s] illegal", value.Groupname)
		}
		hgmap[value.Groupname] = value
	}
	if len(tmphg.Hgs) != len(hgmap) {
		log.Fatal("There is duplicate hostgroup.")
	}

	var rethg = Hostgroups{}
	for _, hg := range tmphg.Hgs {
		for _, g := range hg.Groups {
			if hgmap[g].Groupname == "" {
				log.Fatalf("Not found. Groups illegal [%s/%s].", hg.Groupname, g)
			}
			hg.Ips = append(hg.Ips, hgmap[g].Ips...)
		}
		if checkers.IsDuplicate(hg.Ips) {
			log.Fatalf("IP Duplicate. Hostgroup illegal [%s].", hg.Groupname)
		}

		if !checkers.CheckHostgroupSize(hg, cfg.Hostgroupsize) {
			log.Fatalf("Too large. IP Range illegal [%s] > [%d].", hg.Groupname, cfg.Hostgroupsize)
		}

		if hg.Sshport < 0 {
			log.Fatalf("SSH port illegal [%d] < 0.", hg.Sshport)
		}

		if hg.Sshport == 0 {
			hg.Sshport = 22
		}
		rethg.Hgs = append(rethg.Hgs, hg)
	}

	for _, hg := range rethg.Hgs {
		if hg.Proxy != "" && hgmap[hg.Proxy].Groupname == "" {
			log.Fatalf("Proxy Not found. Groups illegal [%s/%s].", hg.Groupname, hg.Proxy)
		}
	}

	return rethg
}

func GetHostgroups() (Hostgroups, map[string]Hostgroup) {
	if len(hostgroups.Hgs) == 0 {
		log.Fatal("The hostgroups empty.")
	}
	hostgroups = parseHostgroups(hostgroups)
	for _, value := range hostgroups.Hgs {
		hostgroupsm[value.Groupname] = value
	}

	return hostgroups, hostgroupsm
}
func GetHostgroupByname(name string) (Hostgroup, error) {
	if v, ok := hostgroupsm[name]; !ok {
		return v, fmt.Errorf("hostgroup not exist")
	} else {
		return v, nil
	}
}

var auths Auths
var authsm = make(map[string]Auth)

func initAuths() {
	a, err := ioutil.ReadFile(authFile)
	if err != nil {
		log.Fatalf("Can not find auth file[%s].", authFile)
	}
	err = yaml.Unmarshal(a, &auths)
	if err != nil {
		log.Fatalf("YAML[%s] Unmarshal error: %v", authFile, err)
	}
}
func GetAuths() (Auths, map[string]Auth) {
	if len(auths.As) == 0 {
		log.Fatal("The auths empty.")
	}
	for _, value := range auths.As {
		if !checkers.CheckAuthmethodName(value.Name) {
			log.Fatalf("Authmethod name [%s] illegal", value.Name)
		}
		authsm[value.Name] = value
	}
	if len(auths.As) != len(authsm) {
		log.Fatal("There is duplicate auth.")
	}
	return auths, authsm
}
func GetAuthByname(name string) (Auth, error) {
	if v, ok := authsm[name]; !ok {
		return v, fmt.Errorf("auth not exist")
	} else {
		return v, nil
	}
}

var tasks Tasks

func initTasks(scriptFile string) {
	p, err := ioutil.ReadFile(scriptFile)
	if err != nil {
		log.Fatalf("Can not find script file[%s].", scriptFile)
	}

	err = yaml.Unmarshal(p, &tasks)
	if err != nil {
		log.Fatalf("YAML[%s] Unmarshal error: %v", scriptFile, err)
	}
}

var values map[string]interface{}

func initValues(scriptValues string) {
	if scriptValues == "" {
		return
	}
	p, err := ioutil.ReadFile(scriptValues)
	if err != nil {
		log.Fatalf("Can not find script file[%s].", scriptValues)
	}

	err = yaml.Unmarshal(p, &values)
	if err != nil {
		log.Fatalf("YAML[%s] Unmarshal error: %v", scriptValues, err)
	}
}

var tempScript = ".rshell/rshell.tmp"

func templateScript(script string) {
	t, err := template.ParseFiles(script)
	if err != nil {
		log.Fatalf("Parser template script file [%s] failed.", script)
	}

	f, err := os.Create(tempScript)
	if err != nil {
		log.Fatal("Create temp script file failed.")
	}
	defer f.Close()

	if err := t.Execute(f, values); err != nil {
		log.Fatal("Parser template script file task failed.")
	}
}

func GetTasks() Tasks {
	if *script != "" {
		initValues(*scriptValues)
		if len(values) != 0 {
			templateScript(*script)
			initTasks(tempScript)
		} else {
			initTasks(*script)
		}
		if len(tasks.Ts) == 0 {
			log.Fatal("The tasks empty.")
		}

		return tasks
	} else {
		return tasks
	}
}

func LoadEnv() CurrentEnv {
	var ce CurrentEnv
	e, err := ioutil.ReadFile(".rshell/rshell.env")
	if err != nil {
		return CurrentEnv{}
	}

	if yaml.Unmarshal(e, &ce) != nil {
		return CurrentEnv{}
	}

	return ce
}

func SetEnv(ce CurrentEnv) error {
	d, err := yaml.Marshal(&ce)
	if err != nil {
		return fmt.Errorf("Dump env error: %v", err)
	}

	if err := ioutil.WriteFile(".rshell/rshell.env", d, os.ModeAppend); err != nil {
		return fmt.Errorf("Set env error: %v", err)
	}

	return nil
}
