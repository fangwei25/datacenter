package svc

import (
	"web_game/service/trace/rpc/internal/config"
	"web_game/service/trace/rpc/internal/tracefile"
)

type ServiceContext struct {
	Config    config.Config
	TraceFile *tracefile.TraceFile
}

func NewServiceContext(c config.Config) *ServiceContext {
	tf, err := tracefile.NewFile(c.TFile.Path, c.TFile.PreFix)
	if nil != err {
		panic(err)
	}

	return &ServiceContext{
		Config:    c,
		TraceFile: tf,
	}
}
