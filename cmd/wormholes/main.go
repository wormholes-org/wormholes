package main

import (
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/cmd/wormholes/geth"
	"github.com/ethereum/go-ethereum/sgiccommon"
	"os"
	"syscall"
)

func main() {

	//sigs := make(chan os.Signal, 1)
	stopWormhles := make(chan struct{})
	//done := make(chan bool, 1)
	//signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go geth.GethRun(stopWormhles)
	//go ipfs.IpfsRun(stopWormhles)
	//go nftserver.NftServerRun(stopWormhles)

	for {
		select {
		//case <- sigs:
		//	os.Exit(1)
		case <-stopWormhles:
			os.Exit(2)
		case <-sgiccommon.Sigc:
			utils.Sigc <- syscall.SIGTERM
		}
	}

}
