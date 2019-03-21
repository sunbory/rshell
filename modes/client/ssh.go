package client

import (
	"fmt"
	"github.com/luckywinds/rshell/pkg/checkers"
	"github.com/patrickmn/go-cache"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"time"
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

func New(groupname, host string, port int, user, pass, keyname, passphrase string, ciphers []string) (*ssh.Client, error) {
	if groupname == "" {
		groupname = "DEFAULT"
	}
	if !checkers.IsIpv4(host) || port <= 0 || port > 65535 || user == "" {
		return nil, fmt.Errorf("host[%s] or port[%d] or user[%s] illegal", host, port, user)
	}
	if pass == "" && keyname == "" {
		return nil, fmt.Errorf("pass and keyname can not be empty")
	}

	cachekey := user + "+" + keyname + "@" + groupname + "/" + host + ":" + fmt.Sprintf("%d", port)
	if v, ok := dialcache.Get(cachekey); ok {
		return v.(*ssh.Client), nil
	}

	var err error
	auth := make([]ssh.AuthMethod, 0)
	if pass != "" {
		auth = append(auth, ssh.Password(pass))
	}
	if keyname != "" {
		var (
			pemBytes []byte
			signer   ssh.Signer
		)
		pemBytes, err = ioutil.ReadFile(keyname)
		if err != nil {
			return nil, err
		}
		if passphrase == "" {
			signer, err = ssh.ParsePrivateKey(pemBytes)
		} else {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(passphrase))
		}
		if err != nil {
			return nil, err
		}
		auth = append(auth, ssh.PublicKeys(signer))
	}

	clientConfig := &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 60 * time.Second,
		Config: ssh.Config{
			Ciphers: ciphers,
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), clientConfig)
	if err != nil {
		return nil, err
	}

	dialcache.Set(cachekey, client, cache.DefaultExpiration)
	return client, nil
}
