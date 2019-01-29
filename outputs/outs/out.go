package outs

import (
	"github.com/luckywinds/rshell/types"
)

var taskresult types.Taskresult

type OUT interface {
	Print(intime bool, result types.Hostresult, hg types.Hostgroup)
	Break(intime bool, hg types.Hostgroup)
	Finish(intime bool, hg types.Hostgroup)
}
