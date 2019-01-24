package encrypt

import (
	"fmt"
	"github.com/luckywinds/rshell/options"
	"github.com/luckywinds/rshell/pkg/crypt"
	"github.com/luckywinds/rshell/plugins/core"
)

var ACTION = "encrypt"

func Help() {
	fmt.Println(`encrypt cleartext_password
    --- Encrypt cleartext_password with aes 256 cfb

    Examples:
    - encrypt Cloud12#$`)
}

func Command(o options.Options, line string) (string, error) {
	as := core.GetArgFields(line, ACTION, " ")
	if len(as) != 1 || as[0] == "" {
		return "", fmt.Errorf("arguments illegal")
	}

	switch o.Cfg.Passcrypttype {
	case "aes":
		return crypt.AesEncrypt(as[0], o.Cfg)
	default:
		return "", fmt.Errorf("crypt type[%s] not support", o.Cfg.Passcrypttype)
	}
}
