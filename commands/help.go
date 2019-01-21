package commands

import "fmt"

func Help() {
	LoadHelp()
	fmt.Println()
	DoHelp()
	fmt.Println()
	SudoHelp()
	fmt.Println()
	DownloadHelp()
	fmt.Println()
	UploadHelp()
	fmt.Println()
	EncryptHelp()
	fmt.Println()
	DecryptHelp()
	fmt.Println()
	fmt.Println(`exit
    --- Exit rshell
?
    --- Help`)
}

func LoadHelp() {
	fmt.Println(`load -A<auth> -H<host> -P<port>
    --- Load current env(auth@host:port)

    Tips:
    - host，must in hosts.yaml or ipv4
    - auth, optional if config with host in hosts.yaml, must in auth.yaml
    - port, optional if config with host in hosts.yaml, default: 22

    Examples:
    - load -Hhostgroup03
    - load -Aalpha-env-root-pass -H192.168.31.63 -P22`)
}

func DoHelp() {
	fmt.Println(`do cmd1;cmd2;cmd3
    --- Run cmds on targets as normal user

    Examples:
    - do pwd
    - do pwd;whoami;date`)
}

func SudoHelp() {
	fmt.Println(`sudo cmd1;cmd2;cmd3
    --- Run cmds on targets as root which auto change from normal user

    Examples:
    - sudo pwd
    - sudo pwd;whoami;date`)
}

func DownloadHelp() {
	fmt.Println(`download srcFile desDir
    --- Download srcFile from targets to local desDir as normal user

    Examples:
    - download .bashrc .`)
}

func UploadHelp() {
	fmt.Println(`upload srcFile desDir
    --- Upload srcFile from local to targets desDir as normal user

    Examples:
    - upload README.md .`)
}

func EncryptHelp() {
	fmt.Println(`encrypt_aes cleartext_password
    --- Encrypt cleartext_password with aes 256 cfb

    Examples:
    - encrypt_aes Cloud12#$`)
}

func DecryptHelp() {
	fmt.Println(`decrypt_aes ciphertext_password
    --- Decrypt ciphertext_password with aes 256 cfb

    Examples:
    - decrypt_aes 1c15b86d686758158d5fd9551c0ccca6168a2c80f149c38bca`)
}
