package decrypt

import (
	"fmt"
	"github.com/luckywinds/rshell/options"
	"github.com/luckywinds/rshell/pkg/crypt"
	"github.com/luckywinds/rshell/plugins/core"
)

var ACTION = "decrypt"

func Help() {
	fmt.Println(`decrypt ciphertext_password
    --- Decrypt ciphertext_password with aes 256 cfb

    Examples:
    - decrypt 1c15b86d686758158d5fd9551c0ccca6168a2c80f149c38bca`)
}

func Command(o options.Options, line string) (string, error) {
	as := core.GetArgFields(line, ACTION, " ")
	if len(as) != 1 || as[0] == "" {
		return "", fmt.Errorf("arguments illegal")
	}

	switch o.Cfg.Passcrypttype {
	case "aes":
		return crypt.AesDecrypt(as[0], o.Cfg)
	default:
		return "", fmt.Errorf("crypt type[%s] not support", o.Cfg.Passcrypttype)
	}
}

