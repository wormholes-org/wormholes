package nftserver

import (
	"fmt"
	"github.com/nftexchange/nftserver/models"
	//_ "github.com/nftexchange/nftserver/routers"
	"testing"
)

const sqlsvrLcT = "admin:user123456@tcp(192.168.1.238:3306)/"
//const sqlsvrLcT = "demo:123456@tcp(192.168.56.128:3306)/"
//const vpnsvr = "demo:123456@tcp(192.168.1.238:3306)/"
//var SqlSvrT = "admin:user123456@tcp(192.168.1.238:3306)/"
//const dbNameT = "nftdbdemo"
const dbNameT = "nftdb"
const localtimeT = "?parseTime=true&loc=Local"
//const localtimeT = "?charset=utf8mb4&parseTime=True&loc=Local"

const sqldsnT = sqlsvrLcT + dbNameT + localtimeT

func TestCallContracts(t *testing.T) {
	nd, err := models.NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	CallContracts(nd)
	nd.Close()
}

func TestForeignContract(t *testing.T) {
	nd, err := models.NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	err = nd.NewCollections("0x1840Cbd8c0f7a890c77A62594C7a3019904D4191",
		"foreign-contract-0x1840Cbd8",
		"",
		"",
		"0x91C308955eeA74f17079BE2c0dF635406Bd0Ac15",
		"foreign-test.",
		"art",
		"sigedata",
	)
	if err != nil {
		fmt.Println("NewCollections() err=", err)
	}
	models.RoyaltyLimit = 10000
	err = nd.UploadNft(
		"0x1840Cbd8c0f7a890c77A62594C7a3019904D4191",
		"0x1840Cbd8c0f7a890c77A62594C7a3019904D4191",
		"0x1840Cbd8c0f7a890c77A62594C7a3019904D4191",
		"md5 string",
		"name string",
		"desc string",
		"meta string",
		"source_url string",
		"",
		"",
		"categories string",
		"foreign-contract-0x1840Cbd8",
		"asset_sample string",
		"false",
		"2",
		"1",
		"sig string")
	if err != nil {
		fmt.Println("UploadNft err=", err)
	}
	err = nd.NewCollections("0x88c41ce51023b2891DC7b6Ae4C87c1a67163a46f",
		"foreign-contract-0x88c",
		"",
		"",
		"0x91C308955eeA74f17079BE2c0dF635406Bd0Ac15",
		"foreign-test.",
		"art",
		"sigedata",
	)
	if err != nil {
		fmt.Println("NewCollections() err=", err)
	}
	models.RoyaltyLimit = 10000
	err = nd.UploadNft(
		"0x88c41ce51023b2891DC7b6Ae4C87c1a67163a46f",
		"0x88c41ce51023b2891DC7b6Ae4C87c1a67163a46f",
		"0x88c41ce51023b2891DC7b6Ae4C87c1a67163a46f",
		"md5 string",
		"name string",
		"desc string",
		"meta string",
		"source_url string",
		"",
		"",
		"categories string",
		"foreign-contract-0x88c",
		"asset_sample string",
		"false",
		"2",
		"1",
		"sig string")
	if err != nil {
		fmt.Println("UploadNft err=", err)
	}
}