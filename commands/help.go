package commands

import "fmt"

func Help() {
	fmt.Println(`Usage: KEYWORDS HOSTGROUP AGRUMENTS

do HOSTGROUP cmd1; cmd2; cmd3
    --- Run cmds on HOSTGROUP use normal user
sudo HOSTGROUP cmd1; cmd2; cmd3
    --- Run cmds on HOSTGROUP use root which auto change from normal user
download HOSTGROUP srcFile desDir
    --- Download srcFile from HOSTGROUP to local desDir
upload HOSTGROUP srcFile desDir
    --- Upload srcFile from local to HOSTGROUP desDir

encrypt_aes cleartext_password
    --- Encrypt cleartext_password with aes 256 cfb
decrypt_aes ciphertext_password
    --- Decrypt ciphertext_password with aes 256 cfb
exit
    --- Exit rshell
?
    --- Help

> Use HOSTGROUP[name@ip:port] as ip address for single host
name: must in auth.yaml
ip:   must be ipv4
port: optional, default: 22`)
}
