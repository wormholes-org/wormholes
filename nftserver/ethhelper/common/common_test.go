package common

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
	str, err := GetImgBase64("https://gateway.pinata.cloud/ipfs/QmSzyVo3sZyhrbcECLSHPTnow4tXgJwCB1d2j3iCQ4VpVY/hashersreloaded-500-999/selv2_00645.jpg")
	fmt.Println(str, err)
	//c, err := ethhelper.AdminList()
	//fmt.Println(c, err)
}
func TestGetBlock(t *testing.T) {
	var b big.Int
	b.SetString("9358322", 0)
	s := hex.EncodeToString(b.Bytes())
	GetBlock("0x" + s)
}
