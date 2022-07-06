package contracts

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"testing"
)
//const sqlsvrLcT = "admin:user123456@tcp(192.168.1.238:3306)/"
const sqlsvrLcT = "demo:123456@tcp(192.168.56.128:3306)/"
//const vpnsvr = "demo:123456@tcp(192.168.1.238:3306)/"
//var SqlSvrT = "admin:user123456@tcp(192.168.1.238:3306)/"
const dbNameT = "nftdbdemo"
//const dbNameT = "nftdb"
const localtimeT = "?parseTime=true&loc=Local"
//const localtimeT = "?charset=utf8mb4&parseTime=True&loc=Local"

const sqldsnT = sqlsvrLcT + dbNameT + localtimeT

func TestEventSale(t *testing.T) {
	//EventSale(sqldsnT)
}

func TestEventAuction(t *testing.T) {
	//EventAuction(sqldsnT)
}

func TestEventRoyalty(t *testing.T) {
	//EventRoyalty(sqldsnT)
}

func TestEventContract(t *testing.T) {
	EventContract(sqldsnT)
}

func TestAuction(t *testing.T) {
	from := "0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169"
	to := "0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169"
	nftAddr := "0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F"
	tokenId := "1631952185883"
	price := "10000000000000000"
	sig := "0x6ffc548def52273de0c5d27e042f44564e8ac2b23114083bc9fb027f5ec8916e0cb47240f392f2ae59d34ce1bd9c212682b76746b40c5c6e15dbb389536464001b"
	_, err := Auction(from, to, nftAddr, tokenId, price, sig)
	if err != nil {
		log.Panicln(err)
	}
}

func TestSendAddr(t *testing.T) {
	privateKey, err := crypto.HexToECDSA(TradeAuthAddrPriv)
	if err != nil {
		fmt.Println(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println("fromAddress=", fromAddress) //0x077d34394Ed01b3f31fBd9816cF35d4558146066

	privateKey, err = crypto.HexToECDSA("8c995fd78bddf528bd548cce025f62d4c3c0658362dbfd31b23414cf7ce2e8ed")
	if err != nil {
		fmt.Println(err)
	}
	publicKey = privateKey.Public()
	publicKeyECDSA, ok = publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress = crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println("fromAddress=", fromAddress)  //0x2b0aD05ADDa21BA4E5b94C4f9aE3BCeA15A380c5
}