package ssh

import (
	"bytes"
	"fmt"
	"github.com/luckywinds/rshell/modes/client"
	"golang.org/x/crypto/ssh"
	"time"
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

	if sudotype != "" {
		fmt.Fprintf(stdin, "%s\n", sudotype)
		time.Sleep(time.Millisecond * 100)
		fmt.Fprintf(stdin, "%s\n", sudopass)
		time.Sleep(time.Millisecond * 100)
		fmt.Fprintf(stdin, "%s\n", "rrretcode=$?;[ $rrretcode -eq 0 ] || exit $rrretcode")

		for _, cmd := range cmds {
			fmt.Fprintf(stdin, "%s\n", cmd)
			fmt.Fprintf(stdin, "%s\n", "rrretcode=$?;[ $rrretcode -eq 0 ] || exit $rrretcode")
		}

		fmt.Fprintf(stdin, "%s\n", "exit")
		fmt.Fprintf(stdin, "%s\n", "exit")
	} else {
		for _, cmd := range cmds {
			fmt.Fprintf(stdin, "%s\n", cmd)
			fmt.Fprintf(stdin, "%s\n", "rrretcode=$?;[ $rrretcode -eq 0 ] || exit $rrretcode")
		}

		fmt.Fprintf(stdin, "%s\n", "exit")
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

	return DO(groupname, host, port, user, pass, keyname, passphrase, sudotype, sudopass, timeout, ciphers, cmds)
}
