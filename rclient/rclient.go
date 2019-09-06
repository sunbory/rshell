package rclient

import (
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"time"
        "path/filepath"
        "io"
        "bytes"
        "os"
        "runtime"
        "path"
        "strings"

	"github.com/luckywinds/rshell/pkg/checkers"
	"github.com/patrickmn/go-cache"
	"golang.org/x/crypto/ssh"
	"github.com/pkg/sftp"
)

type (
	RClient struct {
		ssh.Client
		RConfig
	}

	RConfig struct {
		Groupname string
		Host      string
		Port      int
		User      string
		Key       string
		Password  string
		Passphrase string
		Ciphers   []string
		Sudotype  string
		Sudopass  string
		Timeout   time.Duration
		Proxy     *RConfig
	}
)

var dialcache *cache.Cache

func SetupDialCache(ttl int) {
	if ttl != 0 {
		dialcache = cache.New(time.Duration(ttl)*time.Second, time.Duration(10)*time.Second)
	} else {
		dialcache = cache.New(3600*time.Second, 10*time.Second)
	}
	go func(c *cache.Cache) {
		t := time.NewTicker(10 * time.Second)
		defer t.Stop()
		for range t.C {
			for key, value := range c.Items() {
				_, _, err := value.Object.(*ssh.Client).Conn.SendRequest("keepalive@rshell", true, nil)
				if err != nil {
					dialcache.Delete(key)
				}
			}
		}
	}(dialcache)
}

func New(config *RConfig) (*RClient, error) {
	
	err := checkRConfig(config)
	if err != nil {
		return nil, err
	}

	cachekey := config.User + "+" + config.Key + "@" + config.Groupname + "/" + net.JoinHostPort(config.Host, strconv.Itoa(config.Port))
	if v, ok := dialcache.Get(cachekey); ok {
		return v.(*RClient), nil
	}

	client := &RClient{
		RConfig: *config,
	}

	rst_client, err := client.connect()
	if err != nil {
		return nil, err
	}
        client.Client = *rst_client

	dialcache.Set(cachekey, client, cache.DefaultExpiration)
	return client, nil
}

func checkRConfig (config *RConfig) (error) {

	if config.Groupname == "" {
		config.Groupname = "DEFAULT"
	}
	if !checkers.ValidIP(config.Host) || config.Port <= 0 || config.Port > 65535 || config.User == "" {
		return fmt.Errorf("host[%s] or port[%d] or user[%s] illegal", config.Host, config.Port, config.User)
	}
	if config.Password == "" && config.Key == "" {
		return fmt.Errorf("pass and keyname can not be empty")
	}

	return nil
	
}

func (client *RClient) genSSHConfig () (*ssh.ClientConfig, error) {
	
	var err error
	auth := make([]ssh.AuthMethod, 0)
	if client.Password != "" {
		auth = append(auth, ssh.Password(client.Password))

		keyboardInteractiveChallenge := func(
			user,
			instruction string,
			questions []string,
			echos []bool,
		) (answers []string, err error) {
			if len(questions) == 0 {
				return []string{}, nil
			}
			return []string{client.Password}, nil
		}
		auth = append(auth, ssh.KeyboardInteractive(keyboardInteractiveChallenge))
	}
	if client.Key != "" {
		var (
			pemBytes []byte
			signer   ssh.Signer
		)
		pemBytes, err = ioutil.ReadFile(client.Key)
		if err != nil {
			return nil, err
		}
		if client.Passphrase == "" {
			signer, err = ssh.ParsePrivateKey(pemBytes)
		} else {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(client.Passphrase))
		}
		if err != nil {
			return nil, err
		}
		auth = append(auth, ssh.PublicKeys(signer))
	}

	sshConfig := &ssh.ClientConfig{
		User:    client.User,
		Auth:    auth,
		Timeout: 60 * time.Second,
		Config: ssh.Config{
			Ciphers: client.Ciphers,
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	return sshConfig, nil
}

func (client *RClient) connect() (*ssh.Client, error) {
	
    var sshClient *ssh.Client
    
    if client.Proxy != nil {

		proxyClient, err := New(client.Proxy)
		if err != nil {
			return nil, err
		}

		conn, err := proxyClient.Dial("tcp", net.JoinHostPort(client.Host,strconv.Itoa(client.Port)))
		if err != nil {
			return nil, err
		}

		targetConfig, err := client.genSSHConfig()
		if err != nil {
			return nil, err
		}

		ncc, chans, reqs, err := ssh.NewClientConn(conn, net.JoinHostPort(client.Host, strconv.Itoa(client.Port)), targetConfig)
		if err != nil {
			return nil, err
		}

		sshClient = ssh.NewClient(ncc, chans, reqs)

	} else {

		clientConfig, err := client.genSSHConfig()
		if err != nil {
			return nil, err
		}

		sshClient, err = ssh.Dial("tcp", net.JoinHostPort(client.Host, strconv.Itoa(client.Port)), clientConfig)
		if err != nil {
			return nil, err
		}
	}

	return sshClient, nil
}

func (client *RClient)  DO (cmds []string, sudo bool) (string, string, error) {
	var (
		session *ssh.Session
		stderr  bytes.Buffer
		stdout  bytes.Buffer
		err     error
	)

	session, err = client.NewSession()
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

	if sudo == true {
		fmt.Fprintf(stdin, "%s || exit 1\n", client.Sudotype)
		time.Sleep(time.Millisecond * 100)
		fmt.Fprintf(stdin, "%s\n", client.Sudopass)
		time.Sleep(time.Millisecond * 100)
		fmt.Fprintf(stdin, "%s\n", "rrretcode=$?;[ $rrretcode -eq 0 ] || exit $rrretcode")

		fmt.Fprintf(stdin, "%s\n", "echo > .rshell.sh")
		for _, cmd := range cmds {
			if cmd != "" {
				fmt.Fprintf(stdin, "%s\n", "echo '"+cmd+"' >> .rshell.sh")
				fmt.Fprintf(stdin, "%s\n", "echo 'rrretcode=$?;[ $rrretcode -eq 0 ] || exit $rrretcode' >> .rshell.sh")
			}
		}
		fmt.Fprintf(stdin, "%s\n", "sh .rshell.sh")
		fmt.Fprintf(stdin, "%s\n", "rm -f .rshell.sh")

		fmt.Fprintf(stdin, "%s\n", "exit")
		fmt.Fprintf(stdin, "%s\n", "exit")
	} else {
		for _, cmd := range cmds {
			if cmd != "" {
				fmt.Fprintf(stdin, "%s\n", cmd)
				fmt.Fprintf(stdin, "%s\n", "rrretcode=$?;[ $rrretcode -eq 0 ] || exit $rrretcode")
			}
		}

		fmt.Fprintf(stdin, "%s\n", "exit")
	}

	err = session.Wait()

	return stdout.String(), stderr.String(), err
}

func (client *RClient) SUDO(cmds []string) (string, string, error) {
	if len(cmds) == 0 {
		return "", "", fmt.Errorf("cmds[%v] empty", cmds)
	}
	if client.Sudotype == "" {
		client.Sudotype = "su"
	}

	return client.DO(cmds, true)
}

func (client *RClient) Upload(srcFilePath, desDirPath string, maxPacketSize int) ([]string, error) {
	var (
		session *sftp.Client
		err     error
	)

	session, err = sftp.NewClient(&client.Client, sftp.MaxPacket(maxPacketSize))
	if err != nil {
		return nil, err
	}
	defer session.Close()

	srcFiles, err := filepath.Glob(srcFilePath)
	if err != nil {
		return nil, err
	}
	if srcFiles != nil {
		for _, sf := range srcFiles {
			srcFile, err := os.Open(sf)
			if err != nil {
				return nil, err
			}
			defer srcFile.Close()

			var desFileName string
			if runtime.GOOS == "windows" {
				desFileName = path.Base(strings.Replace(srcFile.Name(), "\\", "/", -1))
			} else {
				desFileName = path.Base(srcFile.Name())
			}
			desFile, err := session.Create(path.Join(desDirPath, desFileName))
			if err != nil {
				return nil, err
			}
			defer desFile.Close()

			_, err = io.Copy(desFile, srcFile)
			if err != nil {
				return nil, err
			}
		}
		return srcFiles, nil
	} else {
		return nil, fmt.Errorf("files not found")
	}
}

func (client *RClient) Download(srcFilePath, desDirPath string, maxPacketSize int) ([]string, error) {
	var (
		session *sftp.Client
		err     error
	)

	session, err = sftp.NewClient(&client.Client, sftp.MaxPacket(maxPacketSize))
	if err != nil {
		return nil, err
	}
	defer session.Close()

	if err = os.Mkdir(desDirPath, os.ModeDir|os.ModePerm); err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
	}
	if err = os.Mkdir(path.Join(desDirPath, client.Groupname), os.ModeDir|os.ModePerm); err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
	}
	if err = os.Mkdir(path.Join(path.Join(desDirPath, client.Groupname), client.Host), os.ModeDir|os.ModePerm); err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
	}

	srcFiles, err := session.Glob(srcFilePath)
	if err != nil {
		return nil, err
	}
	if srcFiles != nil {
		for _, sf := range srcFiles {
			var desFileName = path.Base(sf)
			srcFile, err := session.Open(sf)
			if err != nil {
				return nil, err
			} else {
				desFile, err := os.Create(path.Join(path.Join(path.Join(desDirPath, client.Groupname), client.Host), desFileName))
				if err != nil {
					return nil, err
				}
				defer desFile.Close()

				_, err = io.Copy(desFile, srcFile)
				if err != nil {
					return nil, err
				}
			}
			defer srcFile.Close()
		}
		return srcFiles, nil
	} else {
		return nil, fmt.Errorf("files not found")
	}
}
