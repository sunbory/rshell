package do

import (
	"fmt"
	"github.com/luckywinds/rshell/options"
	"github.com/luckywinds/rshell/pkg/rlog"
	"github.com/luckywinds/rshell/plugins/core"
	"github.com/luckywinds/rshell/types"
)

var ACTION = "do"

func Help() {
	fmt.Println(`cmd1;cmd2;cmd3
    --- Run cmds on TARGETS as normal user

    Examples:
    - pwd
    - pwd;whoami;date`)
}

func Command(o options.Options, line string) ([]string, error) {
	rlog.Info.Printf("line: %s", line)

	if err := core.EnvCheck(o.CurrentEnv); err != nil {
		return []string{}, err
	}

	as := core.GetArgFields(line, "", o.Cfg.CmdSeparator)
	if len(as) == 0 || as[0] == "" {
		return []string{}, fmt.Errorf("arguments empty")
	}
	rlog.Debug.Printf("as: %#v", as)

	if err := core.SecurityCheck(o.Cfg.BlackCmdList, as); err != nil {
		return []string{}, err
	}

	au, hg := core.GetAuthHostgroup(o)

	au, err := core.GetPlainPassword(o.Cfg, au)
	if err != nil {
		return []string{}, err
	}

	core.RunSshCommands(o.Cfg.Concurrency, ACTION, ACTION, au, hg, as)

	rlog.Debug.Printf("ret: %s", as)
	return as, nil
}

func Script(o options.Options, name string, stask types.Subtask) error {
	rlog.Info.Printf("name: %s, stask: %+v", name, stask)

	if len(stask.Cmds) == 0 {
		return fmt.Errorf("commands empty")
	}

	au, hg := core.GetAuthHostgroup(o)
	au, err := core.GetPlainPassword(o.Cfg, au)
	if err != nil {
		return err
	}

	core.RunSshCommands(o.Cfg.Concurrency, name + "/" + stask.Name, ACTION, au, hg, stask.Cmds)

	return nil
}

