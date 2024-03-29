package main

import (
	"os"
	"path/filepath"

	"github.com/wenwenxiong/network-pinger/cmd/pinger"
	"github.com/wenwenxiong/network-pinger/pkg/util"
)

const (
	CmdPinger                = "network-pinger"
)

func main() {
	cmd := filepath.Base(os.Args[0])
	switch cmd {
	case CmdPinger:
		pinger.CmdMain()
	default:
		util.LogFatalAndExit(nil, "%s is an unknown command", cmd)
	}
}
