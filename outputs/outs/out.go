package outs

import (
	"github.com/luckywinds/rshell/types"
)

type OUT interface {
	Print(actionname, actiontype string, result types.Hostresult, hg types.Hostgroup)
	Break(actionname, actiontype string, hg types.Hostgroup)
	Finish(actionname, actiontype string, hg types.Hostgroup)
	End()
}
