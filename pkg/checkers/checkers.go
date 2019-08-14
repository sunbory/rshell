package checkers

import (
	"github.com/luckywinds/rshell/types"
	"github.com/scylladb/go-set/strset"
	"net"
	"regexp"
	"strings"
)

func IsDuplicate(ss []string) bool {
	tempset := strset.New()
	for _, value := range ss {
		tempset.Add(value)
	}
	if tempset.Size() != len(ss) {
		return true
	}
	return false
}

func ValidIP(s string) bool {
	temp := net.ParseIP(s)
	if temp == nil {
		return false
	}
// 	if temp.To4() == nil {
// 		return false
// 	}
	return true
}

func CheckHostgroupSize(h types.Hostgroup, max int) bool {
	if len(h.Ips) > max {
		return false
	}
	return true
}

func CheckTaskName(name string) bool {
	if ok, _ := regexp.MatchString("^[a-z A-Z0-9_-]+$", name); !ok {
		return false
	}
	return true
}

func CheckHostgroupName(name string) bool {
	if ok, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", name); !ok {
		return false
	}
	return true
}

func CheckAuthmethodName(name string) bool {
	if ok, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", name); !ok {
		return false
	}
	return true
}

func IsBlackCmd(cmd string, bcl []types.BlackCmd) bool {
	cstr := strings.TrimSpace(cmd)
	for _, bstr := range bcl {
		if bstr.Cmd != "" {
			if cstr == bstr.Cmd {
				return true
			}
		}
		if bstr.CmdPrefix != "" {
			if strings.HasPrefix(cstr, bstr.CmdPrefix) {
				return true
			}
		}
	}
	return false
}
