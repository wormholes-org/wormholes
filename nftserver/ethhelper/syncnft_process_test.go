package ethhelper

import (
	"fmt"
	"github.com/nftexchange/nftserver/ethhelper/database"
	"math/big"
	"strconv"
	"testing"
	"time"
)

func TestSyncNftAddr(t *testing.T) {
	//Init()
	go NftDbProcess()
	//SyncNftFromChain("9369204")
	time.Sleep(1 * time.Second)
}
func TestUploadNft(t *testing.T) {
	var b big.Int
	b.SetString("449", 0)
	uploadNft("0xF522B448DbF8884038c684B5c3De95654007Fd2B", &b)
}
func TestGetImg(t *testing.T) {

}

func TestSyncNftFromChain(t *testing.T) {
	wethTransferCh := make(chan *WethTransfer, 100)
	wethApproveCh := make(chan *WethTransfer, 100)
	buyResultCh := make(chan []*database.NftTx)
	endCh := make(chan bool)
	go SyncNftFromChain(strconv.Itoa(9651405), true, buyResultCh, wethTransferCh, wethApproveCh, endCh)
	go func() {
		for {
			select {
			case transfer := <-wethTransferCh:
				//weth转移
				fmt.Println(transfer.From, transfer.To, transfer.Value)
			case approve := <-wethApproveCh:
				//weth授权额度    approve.Value为当前授权额度的实际值
				fmt.Println(approve.From, approve.To, approve.Value)
			}
		}
	}()
	isOver := false
	end := false
	for {
		select {
		case buyResult := <-buyResultCh:
			fmt.Println(buyResult)
		case <-endCh:
			end = true
			if len(wethTransferCh) == 0 && len(wethApproveCh) == 0 {
				isOver = true
			}
			break
		default:
		}
		if isOver || (end && len(wethTransferCh) == 0 && len(wethApproveCh) == 0) {
			break
		}
	}
	fmt.Println("end")
}
