package main

import (
	"github.com/colek42/grpclangserv"
)

func main() {
	cfg := grpclangserv.GetCfg()
	grpclangserv.StartRPCServeri(cfg)

}
