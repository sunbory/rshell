package load

import (
	"fmt"
	"github.com/luckywinds/rshell/options"
	"github.com/luckywinds/rshell/pkg/checkers"
	"github.com/luckywinds/rshell/plugins/core"
	"strconv"
	"strings"
)

var ACTION = "load"

func Help() {
	fmt.Println(`load -A<auth> -H<host> -P<port>
    --- Load current env(auth@host:port)

    Tips:
    - host，must in hosts.yaml or ipv4
    - auth, optional if config with host in hosts.yaml, must in auth.yaml
    - port, optional if config with host in hosts.yaml, default: 22

    Examples:
    - load -Hhostgroup03
    - load -Aalpha-env-root-pass -H192.168.31.63
    - load -Aalpha-env-root-pass -H192.168.31.63 -P23`)
}

var aname, hname string
var port int
func Command(o options.Options, line string) (string, string, int, error) {
	aname = o.CurrentEnv.Authname
	hname = o.CurrentEnv.Hostgroupname
	port = o.CurrentEnv.Port

	as := core.GetArgFields(line, ACTION, " ")
	if len(as) > 3 || len(as) == 0 {
		return aname, hname, port, fmt.Errorf("arguments illegal")
	}

	for _, value := range as {
		if err := args(o, value); err != nil {
			return aname, hname, port, err
		}
	}

	if hname == "" {
		return aname, hname, port, fmt.Errorf("load arguments illegal, %s", "host empty")
	}
	if checkers.IsIpv4(hname) {
		if aname == "" {
			return aname, hname, port, fmt.Errorf("load arguments illegal, %s", "auth empty")
		}
		if port == 0 {
			port = 22
		}
	}

	return aname, hname, port, nil
}

func args(o options.Options, a string) (err error){
	switch {
	case strings.HasPrefix(a, "-A"):
		aname = strings.TrimLeft(a, "-A")
		if _, ok := o.Authsm[aname]; !ok {
			return fmt.Errorf("auth name [%s] not found", a)
		}
	case strings.HasPrefix(a, "-H"):
		hname = strings.TrimLeft(a, "-H")
		if _, ok := o.Hostgroupsm[hname]; !ok {
			if !checkers.IsIpv4(hname) {
				return fmt.Errorf("host name [%s] illegal", a)
			}
		}
	case strings.HasPrefix(a, "-P"):
		port, err = strconv.Atoi(strings.TrimLeft(a, "-P"))
		if err != nil {
			return fmt.Errorf("port number [%s] illegal", a)
		}
	default:
		return fmt.Errorf("argument [%s] illegal", a)
	}
	return nil
}