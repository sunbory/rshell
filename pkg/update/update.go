package update

import (
	"github.com/luckywinds/rshell/pkg/rlog"
	"github.com/luckywinds/rshell/types"
	"github.com/secsy/goftp"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

var fileroot = ".rshell"
var filename = "rshell.latest"
var version types.Latestversion

func isNeedUpdate(cVersion string) int {
	return strings.Compare(cVersion, version.Version)
}

func getLatestVersion() {
	a, err := ioutil.ReadFile(fileroot + "/" + filename)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(a, &version)
	if err != nil {
		return
	}
}

func fromHttp(outpath string, file string, url string) error {
	out, err := os.Create(outpath + "/" + file)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url + "/" + file)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func fromFtp(outpath string, file string, url string) error {
	config := goftp.Config{
		User:               "rshell",
		Password:           "",
		Timeout:            10 * time.Second,
	}
	client, err := goftp.DialConfig(config, strings.TrimLeft(url, "ftp://"))
	if err != nil {
		return err
	}
	defer client.Close()

	out, err := os.Create(outpath + "/" + file)
	if err != nil {
		return err
	}
	defer out.Close()

	err = client.Retrieve(file, out)
	if err != nil {
		return err
	}

	return nil
}

func downloadFile(outpath string, file string, url string) error {
	if strings.HasPrefix(url, "http") {
		if err := fromHttp(outpath, file, url); err != nil {
			rlog.Error.Printf("from http error : %v", err)
			return err
		}
	}
	if strings.HasPrefix(url, "ftp") {
		if err := fromFtp(outpath, file, url); err != nil {
			rlog.Error.Printf("from ftp error : %v", err)
			return err
		}
	}
	return nil
}

func Update(c types.Cfg, cVersion string) {
	if len(c.Updateserver) == 0 {
		rlog.Warn.Print("updateserver empty")
		return
	}

	var server = ""
	for _, s := range c.Updateserver {
		rlog.Info.Printf("check updateserver : %s", s)
		downloadFile(fileroot, filename, s)
		getLatestVersion()
		if version.Version != "" && version.Release != "" {
			server = s
			break
		}
	}
	rlog.Info.Printf("updateserver = %s", server)

	if server != "" && isNeedUpdate(cVersion) != 0 {
		downloadFile(".", version.Release, server)
	}
}