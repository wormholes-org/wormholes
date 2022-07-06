package ethhelper

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
)

func TestBalanceOf(t *testing.T) {
	balance, err := BalanceOf("0x077d34394Ed01b3f31fBd9816cF35d4558146066")
	fmt.Println(balance, err)
}
func TestTxCount(t *testing.T) {
	c, err := TransactionCount("0x077d34394Ed01b3f31fBd9816cF35d4558146066")
	fmt.Println(c, err)
}
func TestTransactionReceipt(t *testing.T) {
	c, err := TransactionReceipt("0x6d1443e8b5682f94aca101a9455508ee4777c4c68b81752b9c1a04734cb0c919")
	fmt.Println(c, err)
}
func TestAdminList(t *testing.T) {
	c, err := AdminList()
	fmt.Println(c, err)
}
func TestSendDealAuctionTx(t *testing.T) {
	SendDealAuctionTx("0x10CEc672c6BB2f6782BEED65987E020902B7bD15", "0x572bcAcB7ae32Db658C8dEe49e156d455Ad59eC8", "0x58C68d71F7E8063c25097d938e7857582D5a1c70", "34", "20000000000000000", "0x844c0fe2c4183f3b73931b9f165aaff1013fbaec7fc64c4c491e52e16e1d0f12475bbb36e0028a5948e6fd59d666e6e6b918d5cb2cf309f07b9d67fd4c8586f51c")
}
func TestGetBlock(t *testing.T) {
	var b big.Int
	b.SetString("9358322", 0)
	s := hex.EncodeToString(b.Bytes())
	GetBlock("0x" + s)
}
func TestGenCreateNftSign(t *testing.T)  {
	GenCreateNftSign("0xEe501cc4eC6d23cfdEE5eD297FEB3016B5cC7E9D","0x40eF9242BaDa1A92a917E8aa8aF70e635455dd5F", "abcd","8","1","1000")
}
