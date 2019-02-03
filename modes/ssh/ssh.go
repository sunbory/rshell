package ssh

import (
	"bytes"
	"fmt"
	"github.com/luckywinds/rshell/modes/client"
	"golang.org/x/crypto/ssh"
	"strings"
)

func DO(groupname, host string, port int, user, pass, keyname, passphrase, sudotype, sudopass string, timeout int, ciphers, cmds []string) (string, string, error) {
	var (
		session *ssh.Session
		stderr  bytes.Buffer
		stdout  bytes.Buffer
		err     error
	)
	c, err := client.New(groupname, host, port, user, pass, keyname, passphrase, timeout, ciphers)
	if err != nil {
		return "", "", err
	}

	session, err = c.NewSession()
	if err != nil {
		return "", "", err
	}

	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		return "", "", err
	}

	session.Stderr = &stderr
	session.Stdout = &stdout

	err = session.Shell()
	if err != nil {
		return "", "", err
	}

	if cmds[len(cmds)-1] != "exit" && !strings.HasPrefix(cmds[len(cmds)-1], "exit ") {
		cmds = append(cmds, "exit")
	}

	var newcmds = []string{}
	if sudotype != "" {
		newcmds = append(newcmds, cmds[:2]...)
		for _, value := range cmds[2 : len(cmds)-2] {
			if value != "" {
				newcmds = append(newcmds, value)
				newcmds = append(newcmds, "rrretcode=$?;[ $rrretcode -eq 0 ] || exit $rrretcode")
			}
		}
		newcmds = append(newcmds, cmds[len(cmds)-2:]...)
	} else {
		for _, value := range cmds[:len(cmds)-1] {
			if value != "" {
				newcmds = append(newcmds, value)
				newcmds = append(newcmds, "rrretcode=$?;[ $rrretcode -eq 0 ] || exit $rrretcode")
			}
		}
		newcmds = append(newcmds, cmds[len(cmds)-1])
	}

	for _, cmd := range newcmds {
		if _, e := fmt.Fprintf(stdin, "%s\n", cmd); e != nil {
			break
		}
	}

	err = session.Wait()

	return stdout.String(), stderr.String(), err
}

func SUDO(groupname, host string, port int, user, pass, keyname, passphrase, sudotype, sudopass string, timeout int, ciphers, cmds []string) (string, string, error) {
	if len(cmds) == 0 {
		return "", "", fmt.Errorf("cmds[%v] empty", cmds)
	}
	if sudotype == "" {
		sudotype = "su"
	}
	if cmds[len(cmds)-1] != "exit" && !strings.HasPrefix(cmds[len(cmds)-1], "exit ") {
		cmds = append(cmds, "exit")
	}
	cmds = append([]string{sudotype, sudopass}, cmds...)
	cmds = append(cmds, "exit")

	return DO(groupname, host, port, user, pass, keyname, passphrase, sudotype, sudopass, timeout, ciphers, cmds)
}
