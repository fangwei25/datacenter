package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	TFile struct {
		Path   string
		PreFix string
	}
	Ea struct {
		Url string
	}
}
