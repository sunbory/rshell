package download

import (
	"fmt"
	"github.com/luckywinds/rshell/options"
	"github.com/luckywinds/rshell/pkg/rlog"
	"github.com/luckywinds/rshell/plugins/core"
	"github.com/luckywinds/rshell/types"
)

var ACTION = "download"

func Help() {
	fmt.Println(`download srcFile desDir
    --- Download srcFile from TARGETS to LOCAL desDir as normal user

    Examples:
    - download .bashrc .`)
}

func Command(o options.Options, line string) ([]string, error) {
	rlog.Info.Printf("line: %s", line)

	if err := core.EnvCheck(o.CurrentEnv); err != nil {
		return []string{}, err
	}

	as := core.GetArgFields(line, ACTION, " ")
	if len(as) != 2 {
		return []string{}, fmt.Errorf("arguments illegal")
	}
	rlog.Debug.Printf("as: %#v", as)

	au, hg := core.GetAuthHostgroup(o)

	au, err := core.GetPlainPassword(o.Cfg, au)
	if err != nil {
		return []string{}, err
	}

	core.RunSftpCommands(o.Cfg, ACTION, ACTION, au, hg, as[0], as[1])

	rlog.Debug.Printf("ret: %s", as)
	return as, nil
}

func Script(o options.Options, name string, stask types.Subtask) error {
	rlog.Info.Printf("name: %s, stask: %+v", name, stask)

	if stask.SrcFile == "" || stask.DesDir == "" {
		return fmt.Errorf("src file or des dir empty")
	}

	au, hg := core.GetAuthHostgroup(o)
	au, err := core.GetPlainPassword(o.Cfg, au)
	if err != nil {
		return err
	}

	core.RunSftpCommands(o.Cfg, name + "/" + stask.Name, ACTION, au, hg, stask.SrcFile, stask.DesDir)

	return nil
}

