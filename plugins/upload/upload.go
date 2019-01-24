package upload

import (
	"fmt"
	"github.com/luckywinds/rshell/options"
	"github.com/luckywinds/rshell/plugins/core"
	"github.com/luckywinds/rshell/types"
)

var ACTION = "upload"

func Help() {
	fmt.Println(`upload srcFile desDir
    --- Upload srcFile from local to targets desDir as normal user

    Examples:
    - upload README.md .`)
}

func Command(o options.Options, line string) ([]string, error) {
	if err := core.EnvCheck(o.CurrentEnv); err != nil {
		return []string{}, err
	}

	as := core.GetArgFields(line, ACTION, " ")
	if len(as) != 2 {
		return []string{}, fmt.Errorf("arguments illegal")
	}

	au, hg := core.GetAuthHostgroup(o)

	au, err := core.GetPlainPassword(o.Cfg, au)
	if err != nil {
		return []string{}, err
	}

	core.RunSftpCommands(o.Cfg.Concurrency, ACTION, ACTION, au, hg, as[0], as[1])

	return as, nil
}

func Script(o options.Options, name string, stask types.Subtask) error {
	if stask.SrcFile == "" || stask.DesDir == "" {
		return fmt.Errorf("src file or des dir empty")
	}

	au, hg := core.GetAuthHostgroup(o)

	core.RunSftpCommands(o.Cfg.Concurrency, name + "/" + stask.Name, ACTION, au, hg, stask.SrcFile, stask.DesDir)

	return nil
}


