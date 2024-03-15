package main

import (
	"os"
	"path/filepath"

	"github.com/xiongwen/network-pinger/cmd/pinger"
	"github.com/xiongwen/network-pinger/pkg/util"
)

const (
	CmdPinger                = "kube-ovn-pinger"
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
