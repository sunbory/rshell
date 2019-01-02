package commands

import (
	"fmt"
	"github.com/luckywinds/rshell/options"
	"github.com/luckywinds/rshell/pkg/checkers"
	"github.com/luckywinds/rshell/pkg/crypt"
	"github.com/luckywinds/rshell/types"
	"strconv"
	"strings"
)

func getGroupAuthbyHostinfo(hostinfo string) (types.Hostgroup, types.Auth, error) {
	if hostinfo == "" {
		return types.Hostgroup{}, types.Auth{}, fmt.Errorf("hostinfo[%s] empty", hostinfo)
	}

	authhost := strings.SplitN(hostinfo, "@", 2)
	if len(authhost) != 2 {
		return types.Hostgroup{}, types.Auth{}, fmt.Errorf("hostinfo[%s] auth wrong", hostinfo)
	}
	authname := authhost[0]

	hostport := strings.SplitN(authhost[1], ":", 2)
	host := hostport[0]
	port := 22
	if len(hostport) == 2 {
		var err error
		port, err = strconv.Atoi(hostport[1])
		if err != nil {
			return types.Hostgroup{}, types.Auth{}, fmt.Errorf("hostinfo[%s] port wrong", hostinfo)
		}
	}

	au, err := options.GetAuthByname(authname)
	if err != nil {
		return types.Hostgroup{}, types.Auth{}, fmt.Errorf("auth[%s] not found", authname)
	}

	if !checkers.IsIpv4(host) || port < 0 || port > 65535 {
		return types.Hostgroup{}, types.Auth{}, fmt.Errorf("hostinfo[%s] host or port wrong", hostinfo)
	}

	hg := types.Hostgroup{
		Groupname:  "TEMPHOST",
		Authmethod: authname,
		Sshport:    port,
		Hosts:      nil,
		Groups:     nil,
		Hostranges: nil,
		Ips:        []string{host},
	}

	return hg, au, nil
}

func getGroupAuthbyGroupname(groupname string) (types.Hostgroup, types.Auth, error) {
	if groupname == "" {
		return types.Hostgroup{}, types.Auth{}, fmt.Errorf("groupname[%s] empty", groupname)
	}

	hg, err := options.GetHostgroupByname(groupname)
	if err != nil {
		return types.Hostgroup{}, types.Auth{}, fmt.Errorf("group[%s] not found", groupname)
	}

	au, err := options.GetAuthByname(hg.Authmethod)
	if err != nil {
		return types.Hostgroup{}, types.Auth{}, fmt.Errorf("auth[%s] not found", hg.Authmethod)
	}

	if len(hg.Ips) == 0 {
		return types.Hostgroup{}, types.Auth{}, fmt.Errorf("hostgroup[%s] hosts empty", groupname)
	}

	return hg, au, nil
}

func getPlainPass(pass string, cfg types.Cfg) (string, error) {
	text, err := crypt.AesDecrypt(pass, cfg)
	if err != nil {
		return "", err
	}
	return text, nil
}