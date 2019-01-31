package outputs

import (
	"github.com/luckywinds/rshell/options"
	"github.com/luckywinds/rshell/outputs/outs"
	"github.com/luckywinds/rshell/outputs/outs/json"
	"github.com/luckywinds/rshell/outputs/outs/text"
	"github.com/luckywinds/rshell/outputs/outs/yaml"
	"github.com/luckywinds/rshell/types"
	"time"
)

var cfg = options.GetCfg()
var O outs.OUT
func init() {
	O = outFactory(cfg.Outputtype)
}

func Output(actionname, actiontype string, result chan types.Hostresult, hg types.Hostgroup) {
	for i := 0; i < len(hg.Ips); i++ {
		select {
		case res := <-result:
			O.Print(actionname, actiontype, res, hg)
		case <-time.After(time.Duration(cfg.Tasktimeout) * time.Second):
			O.Break(actionname, actiontype, hg)
		}
	}
	O.Finish(actionname, actiontype, hg)
}

func outFactory(t string) outs.OUT {
	switch t {
	case "text":
		return text.TEXT{}
	case "json":
		return json.JSON{}
	case "yaml":
		return yaml.YAML{}
	default:
		return text.TEXT{}
	}
}

func End() {
	O.End()
}
