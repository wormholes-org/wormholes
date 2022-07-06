package main

import "github.com/nftexchange/nftserver"

func main() {
	nftServerStop := make(chan struct{})
	nftserver.NftServerRun(nftServerStop)
}
