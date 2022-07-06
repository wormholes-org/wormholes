package main

import (
	"fmt"
	ethhelper2 "github.com/nftexchange/nftserver/ethhelper"
	"log"
	"strconv"
)

func main() {
	//ethhelper.GenCreateNftSign("0xEe501cc4eC6d23cfdEE5eD297FEB3016B5cC7E9D","0x40eF9242BaDa1A92a917E8aa8aF70e635455dd5F", "8","1","1000","abcd")
	//database.GetMaxNfts()

	//go ethhelper.NftDbProcess()
	//t1 := time.Now()
	//ethhelper.SyncNftFromChain("13447819")
	//t2 := time.Now().Sub(t1)
	//fmt.Println(t2.Seconds())

	//database.GetMaxNfts()
	//database.GetTopNft()
	//go ethhelper.NftDbProcess()

	for i := 13448167; i < 13491916; i = i + 1000000 {
		if i+10000 > 13491916 {
			fetchBlock(i, 13448167)
			break
		}
		fetchBlock(i, i+10000)
	}
	//fmt.Println(len("115792089237316195423570985008687907853269984665640564039457584007913129639935"))
	//buyResultCh := make(chan []*database2.NftTx,400)
	//ethhelper2.SyncNftFromChain(strconv.Itoa(9497539), true, buyResultCh,nil)
	//fmt.Println(buyResultCh)
	for {
		select {}
	}
}
func fetchBlock(start, end int) {
	go func() {
		log.Println("start:", start)
		for i := start; i < end; i++ {
			if i%500000 == 0 {
				fmt.Println("Index:", i)
				log.Println("Index:", i)
			}
			ethhelper2.SyncNftFromChain(strconv.Itoa(i), false, nil, nil)
		}
		log.Println("end:", end)
	}()
}
