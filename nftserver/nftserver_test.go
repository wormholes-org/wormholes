package nftserver

import (
	"github.com/nftexchange/nftserver/ethhelper"
	"testing"
)

func TestMain1(t *testing.T) {
	ethhelper.GenCreateNftSign("0xEe501cc4eC6d23cfdEE5eD297FEB3016B5cC7E9D","0x40eF9242BaDa1A92a917E8aa8aF70e635455dd5F", "abcd","8","1","1000")
}
