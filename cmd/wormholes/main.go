package main

import (
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/cmd/wormholes/geth"
	"github.com/ethereum/go-ethereum/sgiccommon"
	"os"
	"syscall"
)

func main() {

	stopWormhles := make(chan struct{})
	go geth.GethRun(stopWormhles)

	for {
		select {
		case <-stopWormhles:
			os.Exit(2)
		case <-sgiccommon.Sigc:
			utils.Sigc <- syscall.SIGTERM
		}
	}

}
