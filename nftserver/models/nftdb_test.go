package models

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	_ "github.com/beego/beego/v2/server/web"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/nftexchange/nftserver/ethhelper"
	"github.com/nftexchange/nftserver/ethhelper/database"
	"math/big"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

const sqlsvrLcT = "admin:user123456@tcp(192.168.1.237:3306)/"
//const sqlsvrLcT = "demo:123456@tcp(192.168.56.128:3306)/"

//const vpnsvr = "demo:123456@tcp(192.168.1.238:3306)/"
//var SqlSvrT = "admin:user123456@tcp(192.168.1.238:3306)/"
//const dbNameT = "nftdbdemo"
const dbNameT = "nftdb"
const localtimeT = "?parseTime=true&loc=Local"

//const localtimeT = "?charset=utf8mb4&parseTime=True&loc=Local"

const sqldsnT = sqlsvrLcT + dbNameT + localtimeT

func TestCreateDb(t *testing.T) {
	nd := new(NftDb)
	err := nd.InitDb(sqlsvrLcT, dbNameT)
	if err != nil {
		fmt.Printf("InitDb() err=%s\n", err)
	}
}

func TestDbMaxConnect(t *testing.T) {
	for i := 0; i < 2000; i++ {
		_, err := NewNftDb(sqldsnT)
		if err != nil {
			fmt.Printf("connet count=%d err=%s\n", i, err)
			break
		}
	}
	fmt.Println("end.")
}

func TestLogin(t *testing.T) {
	wd := sync.WaitGroup{}
	wd.Add(1)
	go func() {
		nd, err := NewNftDb(sqldsnT)
		if err != nil {
			fmt.Printf("connect database err = %s\n", err)
		}
		defer nd.Close()
		err = nd.LoginNew("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e162", "sigdata")
		if err != nil {
			fmt.Printf("login err.\n")
		}
		wd.Done()
	}()
	wd.Add(1)
	go func() {
		nd, err := NewNftDb(sqldsnT)
		if err != nil {
			fmt.Printf("connect database err = %s\n", err)
		}
		defer nd.Close()
		err = nd.LoginNew("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e162", "sigdata")
		if err != nil {
			fmt.Printf("login err.\n")
		}
		wd.Done()
	}()
	wd.Wait()
	fmt.Println("login test end.")
}

func TestModifyUserInfo(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	_, err = nd.QueryUserInfo("bbbbbbbbbbbbbbbbbbbbb")
	if err != nil {
		fmt.Println("QueryUserInfo() err=", err)
	}
	err = nd.ModifyUserInfo("bbbbbbbbbbbbbbbbbbbbb", "renameuser",
		"portrait", "my bio.", "test@test.com", "sigdata")
	if err != nil {
		fmt.Println("ModifyUserInfo() err=", err)
	}
	_, err = nd.QueryUserInfo("bbbbbbbbbbbbbbbbbbbbb")
	if err != nil {
		fmt.Println("QueryUserInfo() err=", err)
	}
	err = nd.ModifyUserInfo("bbbbbbbbbbbbbbbbbcbbbb", "renameuser",
		"portrait", "my bio.", "test@test.com", "sigdata")
	if err != nil {
		fmt.Println("ModifyUserInfo() err=", err)
	}
}

func TestQueryUserInfo(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	_, err = nd.QueryUserInfo("0x572bcacb7ae32db658c8dee49e156d455ad59ec8")
	if err != nil {
		fmt.Println("QueryUserInfo() err=", err)
	}
}

/*func TestUpload(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	err = nd.UploadNft("0x81e4F3538eff2d3761B7637d90E8A1EaD83d44BC",
		"md5",
		"url",
		"1000",
		"signdata",
		"0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F",
		"1631679689395",
		"0x81e4F3538eff2d3761B7637d90E8A1EaD83d44BC",
		"image",
		"false")
	if err != nil {
		fmt.Printf("uploadNft err=%s.\n", err)
	}
	err = nd.UploadNft("0x81e4F3538eff2d3761B7637d90E8A1EaD83d44BC",
	"md5",
	"url",
	"2000",
	"signdata",
	"0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F",
	"1631679689395",
	"0x81e4F3538eff2d3761B7637d90E8A1EaD83d44BC",
	"image",
	"false")
	if err != nil {
		fmt.Printf("uploadNft err=%s.\n", err)
	}
	err = nd.UploadNft("useraddr",
		"md5",
		"url",
		"3000",
		"signdata",
		"contract22",
		"tokenid22",
		"ownaddr22",
		"image",
		"false")
	if err != nil {
		fmt.Printf("uploadNft err=%s.\n", err)
	}
	err = nd.UploadNft("useraddr",
		"md5",
		"url",
		"5000",
		"signdata",
		"contract55",
		"tokenid55",
		"ownaddr55",
		"image",
		"false")
	if err != nil {
		fmt.Printf("uploadNft err=%s.\n", err)
	}
}
*/

func TestBuyNft(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	err = nd.BuyNft("mynft", "tradeSig", "sigdata", "contract11", "TokenId11")
	if err != nil {
		fmt.Printf("buyNft err=%s.\n", err)
	}
	err = nd.BuyNft("mynft", "tradeSig", "sigdata", "contract22", "TokenId22")
	if err != nil {
		fmt.Printf("buyNft err=%s.\n", err)
	}
}

func TestQueryNft(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	_, err = nd.QueryNft()
	if err != nil {
		fmt.Printf("uploadNft err=%s.\n", err)
	}
}

func TestNftbyUser(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	_, err = nd.QueryNftbyUser("mynft")
	if err != nil {
		fmt.Printf("uploadNft err=%s.\n", err)
	}
}

func TestRenameTab(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	nd.db.Migrator().RenameTable("users", "user_infos")
}

//func TestTimePro(t *testing.T) {
//	TimeProc(sqldsnT)
//}

func TestMash(t *testing.T) {
	type test struct {
		Num int64 `json:"num"`
	}
	price, _ := strconv.ParseUint("", 10, 64)
	fmt.Println(price)
	tt := test{98708097097987098}
	marshal, _ := json.Marshal(tt)
	t.Log(string(marshal))
}

func TestFavorited(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	err = nd.Like("useraddr", "contract", "tokenid", "sig")
	if err != nil {
		fmt.Printf("AddFavor err = %s\n", err)
	}
	err = nd.Like("useraddr", "contract11", "tokenid11", "sig")
	if err != nil {
		fmt.Printf("AddFavor err = %s\n", err)
	}
	_, err = nd.QueryNftFavorited("useraddr")
	if err != nil {
		fmt.Printf("QueryFavorited err = %s\n", err)
	}
	err = nd.DelNftFavor("useraddr", "contract11", "tokenid11")
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	_, err = nd.QueryNftFavorited("useraddr")
	if err != nil {
		fmt.Printf("QueryFavorited err = %s\n", err)
	}
}

func TestUserFavorited(t *testing.T) {

}

func TestSell(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	err = nd.Sell("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169",
		"",
		"0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F",
		"0569376186306", "HighestBid", "paychan",
		1, 1001, 2000, "royalty", "美元", "false", "sigdate", "tradedate")
	if err != nil {
		fmt.Printf("Sell() err = %s\n", err)
	}
	err = nd.MakeOffer("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169", "0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F",
		"0569376186306", "1", "1", 1100, "tradeSig", 0, "sig")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	err = nd.MakeOffer("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169", "0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F",
		"0569376186306", "1", "1", 1200, "Tradesig", 0, "sigdata")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	err = nd.MakeOffer("0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F", "0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F",
		"0569376186306", "1", "1", 1500, "TradeSig", 0, "sigdata")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	//test2
	err = nd.Sell("ownAddr11", "", "contract11", "tokenid11",
		"FixPrice", "paychan",
		1, 2001, 5000, "royalty", "use", "false", "sigdata", "tradedate")
	if err != nil {
		fmt.Printf("Sell() err = %s\n", err)
	}
	err = nd.MakeOffer("buyer1", "contract11", "tokenid11", "1", "1", 2100, "Tradesig", 0, "sigdata")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	err = nd.MakeOffer("buyer2", "contract11", "tokenid11", "1", "1", 2200, "Tradesig", 0, "sigdata")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	err = nd.MakeOffer("buyer3", "contract11", "tokenid11", "1", "1", 6300,
		"Tradesig", 0, "sigdata")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	//test3
	err = nd.Sell("ownAddr22", "", "contract22", "tokenid22", "HighestBid", "paychan",
		1, 6000, 8000, "royalty", "use", "false", "sigdata", "tradeSig")
	if err != nil {
		fmt.Printf("Sell() err = %s\n", err)
	}
	if err != nil {
		fmt.Printf("Sell() err = %s\n", err)
	}
	err = nd.MakeOffer("buyer1", "contract22", "tokenid22", "1", "1", 6100, "tradesig", 0, "sigdata")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	err = nd.MakeOffer("buyer2", "contract22", "tokenid22", "1", "1", 6200, "TradeSig", 0, "sigdata")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	err = nd.MakeOffer("buyer3", "contract22", "tokenid22", "1", "1", 6300, "tradesig", 0, "sigdata")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	nd.Close()
}

func signHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}

func TestGetSign(t *testing.T) {
	var message []byte = []byte("Hello World!")
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed GenerateKey with %s.", err)
	}
	//带有0x的私钥
	fmt.Println("private key have 0x   \n", hexutil.Encode(crypto.FromECDSA(key)))
	fmt.Println("public key have 0x   \n", hexutil.Encode(crypto.FromECDSAPub(&key.PublicKey)))
	fmt.Println("addr   \n", crypto.PubkeyToAddress(key.PublicKey).String())
	//不含0x的私钥
	fmt.Println("private key no 0x \n", hex.EncodeToString(crypto.FromECDSA(key)))
	sig, err := crypto.Sign(signHash(message), key)
	if err != nil {
		t.Errorf("signature error: %s", err)
	}
	sig[64] += 27
	sigstr := hexutil.Encode(sig)
	addr, err := NftDb{}.GetEthAddr("Hello World!", sigstr)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("addr=%x\n", addr)
}

func TestSignAppconf(t *testing.T) {
	file, err := os.OpenFile("D:\\temp\\app.conf", os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	fileinfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)
	bytesread, err := file.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(bytesread)
	msg := string(buffer)
	tm := fmt.Sprintf(time.Now().String())
	msg = msg + "[time]\n" + "date = " + tm + "\n\n"

	var message []byte = []byte(msg)
	key, err := crypto.HexToECDSA("8c995fd78bddf528bd548cce025f62d4c3c0658362dbfd31b23414cf7ce2e8ed")
	if err != nil {
		fmt.Println(err)
	}
	sig, err := crypto.Sign(signHash(message), key)
	if err != nil {
		t.Errorf("signature error: %s", err)
	}
	sig[64] += 27
	sigstr := hexutil.Encode(sig)
	msg = msg + "#签名数据\n" + "[sig]\n" + "app.conf.sig = " + sigstr
	_, err = file.WriteAt([]byte(msg), 0)
	if err != nil {
		fmt.Println(err)
	}
}

func TestVerifyAppconf(t *testing.T) {
	file, err := os.OpenFile("D:\\temp\\app.conf", os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	fileinfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)
	bytesread, err := file.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(bytesread)
	msg := string(buffer)
	index := strings.Index(msg, "app.conf.sig = ")
	sig := msg[index + len("app.conf.sig = "):]
	var message []byte = []byte(msg[:strings.Index(msg, "#签名数据")])
	addr, err := NftDb{}.GetEthAddr(string(message), sig)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(addr)
}

func TestGetEthAddr(t *testing.T) {
	/*{
	  “msg": "Hello World!"
	  "address": "0x0109cc44df1c9ae44bac132ed96f146da9a26b88",
	  "msg": "0x48656c6c6f20576f726c6421",
	  "sig": "23ad293d6976499c11905c2c811502af9c47c2a0388bec4acb7cf2005554f39226a74d6aec36cdca868dd7ecf62fdd92888e2f9f45939f7f4450362eea1cb5ad1c",
	  "version": "3",
	  "signer": "MEW"
	}*/
	nd := new(NftDb)
	addr, err := nd.GetEthAddr("Hello World!", "0x23ad293d6976499c11905c2c811502af9c47c2a0388bec4acb7cf2005554f39226a74d6aec36cdca868dd7ecf62fdd92888e2f9f45939f7f4450362eea1cb5ad1c")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(addr)
}

func TestGetAdminAddr(t *testing.T) {
	addrs, err := ethhelper.AdminList()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(addrs)
}

func TestQueryNftCurTransInfo(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	nftTranInfo, err := nd.QuerySingleNft("0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d", "9985")
	if err != nil {
		fmt.Println(err)
	}
	marshal, _ := json.Marshal(nftTranInfo)
	fmt.Printf("%s\n", string(marshal))
}

func TestDbPing(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	db, err := nd.db.DB()
	if db.Ping() != nil {

	}
}

func TestQueryMarketInfo(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	nd.QueryMarketInfo()
}

func TestGetBalance(t *testing.T) {

}

func TestBuyResult(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	nd1 := new(NftDb)
	err1 := nd1.ConnectDB(sqldsnT)
	if err1 != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd1.Close()
	/*
		 auctionRec.Ownaddr= 0x81e4F3538eff2d3761B7637d90E8A1EaD83d44BC
		5873 bidRecs.Bidaddr= 0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169
		5874 auctionRec.Contract= 0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F
		5875 auctionRec.Tokenid= 1631681392629
		5876 price= 50000000000000000
	*/
	if true {
		err = nd.BuyResult("0x86c02ffd61b0aca14ced6c3fefc4c832b58b246c",
			"0x8fbf399d77bc8c14399afb0f6d32dbe22189e169",
			"0xA1e67a33e090Afe696D7317e05c506d7687Bb2E5",
			"1062183305419",
			"tradesig",
			"200000000", "sigData", "", "txhash")
		if err != nil {
			fmt.Println(err)
		}
		err = nd.BuyResult("",
			"0x86c02Ffd61b0ACA14CED6c3feFC4C832B58b246c",
			"0xA1e67a33e090Afe696D7317e05c506d",
			"9161528579394",
			"tradesig",
			"", "sigData", "200", "txhash")
		if err != nil {
			fmt.Println(err)
		}
	} else {

		go nd.BuyResult("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169",
			"0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169",
			"0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F",
			"1631753648255",
			"tradesig",
			"20000000000", "sigData", "", "txhash")
		go nd1.BuyResult("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169",
			"0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169",
			"0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F",
			"1631753648255",
			"tradesig",
			"", "sigData", "", "txhash")
		select {}
	}
}

func TestQueryNftByFilter(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	filters := []StQueryField{
		{"collectcreator", "=", "0x572bcacb7ae32db658c8dee49e156d455ad59ec8"},
		{"collections", "=", "Buyer"},
	}
	sorts := []StSortField{{By: "createdate", Order: "desc"}}
	nftByFilter, count, err := nd.QueryNftByFilter(filters, sorts, "0", "10")
	if err != nil {
		t.Fatalf("err = %v\n", err)
	}
	t.Logf("nft = %v %v\n", nftByFilter, count)
}

func TestTimeStamp(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	var nftRecs []Nfts
	//errr := nd.db.Where("createaddr = ?", "useraddr").Distinct("createaddr").Find(&nftRecs)
	errr := nd.db.Where("createaddr = ?", "useraddr").Find(&nftRecs)
	//errr := nd.db.Model(&Nfts{}).Find(&nftRecs)

	if errr.Error != nil {
		fmt.Println(err.Error())
	}
	/*fmt.Println(strconv.FormatInt(time.Now().UnixNano()/1e6, 10))
	fmt.Println(strconv.FormatInt(time.Now().UnixNano()/1e6, 10))
	fmt.Println(strconv.FormatInt(time.Now().UnixNano()/1e6, 10))
	fmt.Println(strconv.FormatInt(time.Now().UnixNano()/1e6, 10))
	fmt.Println(strconv.FormatInt(time.Now().UnixNano(), 10))
	fmt.Println(strconv.FormatInt(time.Now().UnixNano(), 10))
	fmt.Println(strconv.FormatInt(time.Now().UnixNano(), 10))*/
	fmt.Println(strconv.FormatInt(time.Now().UnixNano()/1e6, 10))
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10000; i++ {
		//fmt.Println(rand.Int63())
		s := fmt.Sprintf("%d", rand.Int63())
		if len(s) > 16 {
			continue
		}
		s1 := s[len(s)-13:]
		fmt.Println(s1, "=", len(s))
		//fmt.Println(rand.New(rand.NewSource(time.Now().UnixNano())).Int63())
	}
}

func TestSysParams(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	//err  = nd.SetSysParams()

}

func TestInitSysParams(t *testing.T) {
	InitSysParams(sqldsnT)
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	err = nd.UploadNft(
		"0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169",
		"0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169",
		"0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169",
		"md5 string",
		"name string",
		"desc string",
		"meta string",
		"source_url string",
		"",
		"",
		"categories string",
		"New Art",
		"asset_sample string",
		"true",
		"2",
		"1",
		"sig string")
	if err != nil {
		fmt.Printf("uploadNft err=%s.\n", err)
	}
	err = nd.SetSysParams(SysParamsRec{NFT1155addr: "0x81e4F3538eff2d3761B7637d90E8A1EaD83d44BC", Adminaddr: "", Lowprice: 100000000})
	if err != nil {
		fmt.Printf("SetSysParams() err=%s.\n", err)
	}
	err = nd.SetSysParams(SysParamsRec{NFT1155addr: "0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F", Adminaddr: "", Lowprice: 100000000})
	if err != nil {
		fmt.Printf("SetSysParams() err=%s.\n", err)
	}
	err = nd.SetSysParams(SysParamsRec{NFT1155addr: "", Adminaddr: "", Lowprice: 100000000})
	if err != nil {
		fmt.Printf("SetSysParams() err=%s.\n", err)
	}
	nd.Close()
}

func TestBalanceOfWeth(t *testing.T) {
	c, err := ethhelper.BalanceOfWeth("0x86c02Ffd61b0ACA14CED6c3feFC4C832B58b246c")
	fmt.Println(c, err)
	fmt.Println(c > 1000, err)
}

func TestAllowanceOfWeth(t *testing.T) {
	c, err := ethhelper.AllowanceOfWeth("0x86c02Ffd61b0ACA14CED6c3feFC4C832B58b246c")
	fmt.Println(c, err)
	c = c[:len(c)-9]
	fmt.Println(c, err)
	wei := new(big.Int)
	wei.SetString(c, 10)
}

func TestCollections(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	err = nd.NewCollections("0x86c02ffd61b0aca14ced6c3fefc4c832b58b246c",
		"test",
		"data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAAEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQH/2wBDAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQH/wAARCAK8ArwDAREAAhEBAxEB/8QAHwAAAQUBAQEBAQEAAAAAAAAAAAECAwQFBgcICQoL/8QAtRAAAgEDAwIEAwUFBAQAAAF9AQIDAAQRBRIhMUEGE1FhByJxFDKBkaEII0KxwRVS0fAkM2JyggkKFhcYGRolJicoKSo0NTY3ODk6Q0RFRkdISUpTVFVWV1hZWmNkZWZnaGlqc3R1dnd4eXqDhIWGh4iJipKTlJWWl5iZmqKjpKWmp6ipqrKztLW2t7i5usLDxMXGx8jJytLT1NXW19jZ2uHi4+Tl5ufo6erx8vP09fb3+Pn6/8QAHwEAAwEBAQEBAQEBAQAAAAAAAAECAwQFBgcICQoL/8QAtREAAgECBAQDBAcFBAQAAQJ3AAECAxEEBSExBhJBUQdhcRMiMoEIFEKRobHBCSMzUvAVYnLRChYkNOEl8RcYGRomJygpKjU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6goOEhYaHiImKkpOUlZaXmJmaoqOkpaanqKmqsrO0tba3uLm6wsPExcbHyMnK0tPU1dbX2Nna4uPk5ebn6Onq8vP09fb3+Pn6/9oADAMBAAIRAxEAPwD+/igAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAo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gAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoA//9k=",
		"",
		"",
		"test.",
		"art",
		"sigedata",
	)
	if err != nil {
		fmt.Println("NewCollections() err=", err)
	}
	err = nd.ModifyCollections("0x86c02ffd61b0aca14ced6c3fefc4c832b58b246c",
		"test", "img", "contract_type", "contract_addr",
		"test desc.", "art", "sig string")
	if err != nil {
		fmt.Println("NewCollections() err=", err)
	}
	err = nd.ModifyCollectionsImage("test",
		"0x86c02ffd61b0aca14ced6c3fefc4c832b58b246c",
		"data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAAEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQH/2wBDAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQH/wAARCAK8ArwDAREAAhEBAxEB/8QAHwAAAQUBAQEBAQEAAAAAAAAAAAECAwQFBgcICQoL/8QAtRAAAgEDAwIEAwUFBAQAAAF9AQIDAAQRBRIhMUEGE1FhByJxFDKBkaEII0KxwRVS0fAkM2JyggkKFhcYGRolJicoKSo0NTY3ODk6Q0RFRkdISUpTVFVWV1hZWmNkZWZnaGlqc3R1dnd4eXqDhIWGh4iJipKTlJWWl5iZmqKjpKWmp6ipqrKztLW2t7i5usLDxMXGx8jJytLT1NXW19jZ2uHi4+Tl5ufo6erx8vP09fb3+Pn6/8QAHwEAAwEBAQEBAQEBAQAAAAAAAAECAwQFBgcICQoL/8QAtREAAgECBAQDBAcFBAQAAQJ3AAECAxEEBSExBhJBUQdhcRMiMoEIFEKRobHBCSMzUvAVYnLRChYkNOEl8RcYGRomJygpKjU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6goOEhYaHiImKkpOUlZaXmJmaoqOkpaanqKmqsrO0tba3uLm6wsPExcbHyMnK0tPU1dbX2Nna4uPk5ebn6Onq8vP09fb3+Pn6/9oADAMBAAIRAxEAPwD+/igAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAo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gAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoA//9k=",
		"sig string")
	if err != nil {
		fmt.Println("NewCollections() err=", err)
	}
}

func TestUploadNftNew(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	RoyaltyLimit = 10000
	err = nd.UploadNft(
		"0x86c02ffd61b0aca14ced6c3fefc4c832b58b246c",
		"0x86c02ffd61b0aca14ced6c3fefc4c832b58b246c",
		"0x86c02ffd61b0aca14ced6c3fefc4c832b58b246c",
		"md5 string",
		"name string",
		"desc string",
		"meta string",
		"source_url string",
		"",
		"",
		"categories string",
		"test",
		"asset_sample string",
		"true",
		"2",
		"1",
		"sig string")
	if err != nil {
		fmt.Println("UploadNft err=", err)
	}
	for i := 0; i < 20; i++ {
		if i%2 == 0 {
			NFT1155Addr = "0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169"
		} else {
			NFT1155Addr = "0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F"
		}
		err = nd.UploadNft(
			"0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169",
			"0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169",
			"0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169",
			"md5 string",
			"name string",
			"desc string",
			"meta string",
			"source_url string",
			"",
			"",
			"categories string",
			"",
			"asset_sample string",
			"true",
			"2",
			"1",
			"sig string")
		if err != nil {
			fmt.Println("UploadNft err=", err)
		}
	}
}

func TestForeignContract(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	err = nd.NewCollections("0x86c02ffd61b0aca14ced6c3fefc4c832b58b246c",
		"foreign-contract-test",
		"",
		"",
		"0x9e2576747C2525062a77667E4E88A97b6034C461",
		"foreign-test.",
		"art",
		"sigedata",
	)
	if err != nil {
		fmt.Println("NewCollections() err=", err)
	}
	RoyaltyLimit = 10000
	err = nd.UploadNft(
		"0x86c02ffd61b0aca14ced6c3fefc4c832b58b246c",
		"0x86c02ffd61b0aca14ced6c3fefc4c832b58b246c",
		"0x86c02ffd61b0aca14ced6c3fefc4c832b58b246c",
		"md5 string",
		"name string",
		"desc string",
		"meta string",
		"source_url string",
		"",
		"",
		"categories string",
		"foreign-contract-test",
		"asset_sample string",
		"false",
		"2",
		"1",
		"sig string")
	if err != nil {
		fmt.Println("UploadNft err=", err)
	}
}

func TestDbQuery(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	var count int64
	dberr := nd.db.Model(Nfts{}).Where("contract = ? ", "0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169").Count(&count)
	dberr = nd.db.Model(Nfts{}).Where("contract = ? ", "0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F").Count(&count)

	var nfttab []Nfts
	dberr = nd.db.Model(Nfts{}).Where("contract = ? ", "0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169").Limit(2).Offset(2).Find(&nfttab)
	dberr = nd.db.Model(Nfts{}).Where("contract = ? ", "0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169").Limit(2).Offset(5).Find(&nfttab)
	dberr = nd.db.Where("contract = ? ", "").First(&nfttab)
	if dberr.Error == nil {
		fmt.Println("UploadNft() err=nft already exist.")
	}
}

func TestQueryUserCollectionList(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	_, _, err = nd.QueryUserCollectionList("0x8fBC8ad616177c6519228FCa4a7D9EC7d1804900",
		"0", "5")
	if err != nil {
		fmt.Println("QueryUserCollectionList() err=", err)
	}
}

func TestQueryUserFavorited(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	_, _, err = nd.QueryUserFavoriteList("0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900",
		"0", "5")
	if err != nil {
		fmt.Println("QueryUserCollectionList() err=", err)
	}
}

func TestModifyCollectionsImage(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	nn := time.Now().Unix()
	fmt.Println(nn)
	err = nd.ModifyCollectionsImage("test", "0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900",
		"modify", "0x4a71940655b075316ae19b02457201ed0f719d14f2d20c986b8c16113233e047535d5d1cc4eb293609e79bc60daf622216b190d50a16519d6f826bee05e548051b")
	if err != nil {
		fmt.Println("QueryUserCollectionList() err=", err)
	}
}

func TestQueryUserTradingHistory(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	_, _, err = nd.QueryUserTradingHistory("0x572bcAcB7ae32Db658C8dEe49e156d455Ad59eC8",
		"0", "5")
	if err != nil {
		fmt.Println("QueryUserCollectionList() err=", err)
	}
}

func TestLike(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	err = nd.Like("0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900",
		"0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F", "1632799124069", "sig")
	if err != nil {
		fmt.Println("QueryUserCollectionList() err=", err)
	}
}

func TestSearch(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()

	//testCond1 := "name"
	testCond2 := "0x"

	//searchData1, _ := nd.Search(testCond1)
	searchdata2, _ := nd.Search(testCond2)
	//for _, data := range searchData1 {
	//	for _, data1 := range data.CollectsRecords {
	//		t.Log(data1)
	//	}
	//	for _, data1 := range data.UserAddrs {
	//		t.Log(data1)
	//	}
	//	for _, data1 := range data.NftsRecords {
	//		t.Log(data1)
	//	}
	//}
	for _, data := range searchdata2 {
		for _, data1 := range data.CollectsRecords {
			t.Log(data1)
		}
		for _, data1 := range data.UserAddrs {
			t.Log(data1)
		}
		for _, data1 := range data.NftsRecords {
			t.Log(data1)
		}
	}
}

func TestBidPriceWithBuy(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	/*err = nd.Sell("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e169",
		"",
		"0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F",
		"0569376186306", "HighestBid", "paychan",
		1, 1001, 2000, "royalty","美元", "false", "sigdate", "tradedate")
	if err != nil {
		fmt.Printf("Sell() err = %s\n", err)
	}*/
	err = nd.MakeOffer("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e160", "0xA1e67a33e090Afe696D7317e05c506d7687Bb2E5",
		"7070595686952", "1", "1", 1100, "tradeSig", 0, "sig")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	err = nd.MakeOffer("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e161", "0x53d76f1988B50674089e489B5ad1217AaC08CC85",
		"2530439535350", "1", "1", 1100, "tradeSig", 0, "sig")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	err = nd.MakeOffer("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e162", "0x53d76f1988B50674089e489B5ad1217AaC08CC85",
		"2530439535350", "1", "1", 1100, "tradeSig", 0, "sig")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	err = nd.BuyResult("0x86c02Ffd61b0ACA14CED6c3feFC4C832B58b246c",
		"0x8fBf399D77BC8C14399AFB0F6d32DBe22189e162",
		"0x53d76f1988B50674089e489B5ad1217AaC08CC85",
		"2530439535350",
		"tradesig",
		"200000000", "sigData", "", "txhash")
	if err != nil {
		fmt.Println(err)
	}
}

func TestSignal(t *testing.T) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	s := <-c
	fmt.Println("Got signal:", s)
}

func TestBidPriceWithSell(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	err = nd.MakeOffer("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e160", "0x53d76f1988B50674089e489B5ad1217AaC08CC85",
		"2530439535350", "1", "1", 1100, "tradeSig", 0, "sig")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	err = nd.MakeOffer("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e161", "0x53d76f1988B50674089e489B5ad1217AaC08CC85",
		"2530439535350", "1", "1", 1100, "tradeSig", 0, "sig")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	err = nd.MakeOffer("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e162", "0x53d76f1988B50674089e489B5ad1217AaC08CC85",
		"2530439535350", "1", "1", 1100, "tradeSig", 0, "sig")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	err = nd.Sell("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e162",
		"",
		"0x53d76f1988B50674089e489B5ad1217AaC08CC85",
		"2530439535350", "HighestBid", "paychan",
		1, 1001, 2000, "royalty", "美元", "false", "sigdate", "tradedate")
	if err != nil {
		fmt.Printf("Sell() err = %s\n", err)
	}
}

func TestBidPriceWithTime(t *testing.T) {
	nd := new(NftDb)
	err := nd.ConnectDB(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	err = nd.MakeOffer("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e160", "0x53d76f1988B50674089e489B5ad1217AaC08CC85",
		"2530439535350", "1", "1", 1100, "tradeSig", time.Now().Unix()+10, "sig")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	err = nd.MakeOffer("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e161", "0x53d76f1988B50674089e489B5ad1217AaC08CC85",
		"2530439535350", "1", "1", 1100, "tradeSig", time.Now().Unix()+1000, "sig")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	err = nd.MakeOffer("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e162", "0x53d76f1988B50674089e489B5ad1217AaC08CC85",
		"2530439535350", "1", "1", 1100, "tradeSig", 0, "sig")
	if err != nil {
		fmt.Printf("MakeOffer() err = %s\n", err)
	}
	err = nd.Sell("0x8fBf399D77BC8C14399AFB0F6d32DBe22189e162",
		"",
		"0x53d76f1988B50674089e489B5ad1217AaC08CC85",
		"2530439535350", "HighestBid", "paychan",
		1, 1001, 2000, "royalty", "美元", "false", "sigdate", "tradedate")
	if err != nil {
		fmt.Printf("Sell() err = %s\n", err)
	}
}

func TestQueryMarketTradingHistory(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()

	sorts := []StSortField{{By: "transtime", Order: "desc"}}

	history, i, err := nd.QueryMarketTradingHistory(nil, sorts, "0", "2")

	t.Log(history)
	t.Log(i)
	t.Log(err)
}

func TestAnnouncements(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	err = nd.SetAnnouncement("title one", "content one")
	if err != nil {
		fmt.Println("insert announcement.")
	}
	err = nd.SetAnnouncement("title two", "content two")
	if err != nil {
		fmt.Println("insert announcement.")
	}
	err = nd.SetAnnouncement("title three", "content three")
	if err != nil {
		fmt.Println("insert announcement.")
	}
	err = nd.SetAnnouncement("title four", "content three")
	if err != nil {
		fmt.Println("insert announcement.")
	}
	err = nd.SetAnnouncement("title five", "content five")
	if err != nil {
		fmt.Println("insert announcement.")
	}
	err = nd.SetAnnouncement("title six", "content six")
	if err != nil {
		fmt.Println("insert announcement.")
	}
	err = nd.SetAnnouncement("title seven", "content seven")
	if err != nil {
		fmt.Println("insert announcement.")
	}
	err = nd.SetAnnouncement("title eight", "content eight")
	if err != nil {
		fmt.Println("insert announcement.")
	}
	_, err = nd.QueryAnnouncement()
	if err != nil {
		fmt.Println("insert announcement.")
	}
}

func TestSearchSql(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	var useroffer []UserOffer
	sql := "SELECT biddings.contract as Contract, biddings.tokenid as Tokenid, biddings.price as Price, " +
		"biddings.count as Count, biddings.bidtime as Bidtime FROM biddings LEFT JOIN nfts ON biddings.contract = nfts.contract AND biddings.tokenid = nfts.tokenid WHERE ownaddr = ? AND biddings.deleted_at is null"
	sql = sql + " limit 1, 2"
	db := nd.db.Raw(sql, "0x2b0aD05ADDa21BA4E5b94C4f9aE3BCeA15A380c5").Scan(&useroffer)
	if db.Error != nil {
		fmt.Println("QueryUserInfo() query Sum err=", err)
	}
	var count int64
	sql = "SELECT biddings.contract as Contract, biddings.tokenid as Tokenid, biddings.price as Price, " +
		"biddings.count as Count, biddings.bidtime as Bidtime FROM biddings LEFT JOIN nfts ON biddings.contract = nfts.contract AND biddings.tokenid = nfts.tokenid WHERE ownaddr = ? AND biddings.deleted_at is null"
	db = nd.db.Raw(sql, "0x2b0aD05ADDa21BA4E5b94C4f9aE3BCeA15A380c5").Count(&count)
	if db.Error != nil {
		fmt.Println("QueryUserInfo() query Sum err=", err)
	}
}

func TestIsValidCategory(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()

	var category1 = "virtual_worlds"
	validCategory1 := nd.IsValidCategory(category1)
	var category2 = "virtual"
	validCategory2 := nd.IsValidCategory(category2)
	t.Log("validCategory1=", validCategory1, "validCategory2=", validCategory2)
}

func TestQueryCollectionInfo(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()

	info, _ := nd.QueryCollectionInfo("0x86c02ffd61b0aca14ced6c3fefc4c832b58b246c",
		"实用合集")
	t.Log(info)
}

func TestQueryHomePage(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()

	page, err := nd.QueryHomePage()

	t.Log(page)
	t.Log(err)
}

func TestConvert(t *testing.T) {
	m := -19
	mstr := strconv.Itoa(m)
	u64, err := strconv.ParseUint(mstr, 10, 64)
	fmt.Println(u64, err)
	mstr = "ffffabdcdef"
	u64, err = strconv.ParseUint(mstr, 16, 64)
	fmt.Println(u64, err)

	mstr = ""
	u64, err = strconv.ParseUint(mstr, 10, 64)
	fmt.Println(u64, err)
	data, err := strconv.Atoi(mstr)
	fmt.Println(data, err)
	mstr = "ffffabdcdef"
	u64, err = strconv.ParseUint(mstr, 16, 64)
	fmt.Println(u64, err)
}

func TestConvertValid(t *testing.T) {
	err := IsIntDataValid("")
	if err != true {
		fmt.Println("datat err")
	}
	err = IsUint64DataValid("")
	if err != true {
		fmt.Println("datat err")
	}
}

func TestQueryUserOfferList(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	_, _, _ = nd.QueryUserOfferList("0x572bcacb7ae32db658c8dee49e156d455ad59ec8",
		"0", "10")
}

func TestQueryNftCollectionList(t *testing.T) {
	nd, err := NewNftDb(sqldsnT)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	u1, _, _ := nd.QueryNFTCollectionList("0", "25")
	u2, _, _ := nd.QueryNFTCollectionList("0", "10")
	u3, _, _ := nd.QueryNFTCollectionList("10", "10")
	u4, _, _ := nd.QueryNFTCollectionList("20", "10")
	fmt.Println(u1, u2, u3, u4)
}

func TestUpdateBlockNumber(t *testing.T) {
	UpdateBlockNumber(sqldsnT)
}

func TestSyncProc(t *testing.T) {
	//9532550
	InitSyncBlockTs(sqldsnT)
	//syncFlag := make(chan struct{})
	//SyncProc("", syncFlag)
	//<-syncFlag
}

func TestSyncNftFromChain(t *testing.T) {
	buyResultCh := make(chan []*database.NftTx)
	wethTransferCh := make(chan *ethhelper.WethTransfer)
	wethApproveCh := make(chan *ethhelper.WethTransfer)
	var BlockTxs []*database.NftTx
	var wethTransfers []*ethhelper.WethTransfer
	var wethApproves []*ethhelper.WethTransfer
	endCh := make(chan bool)
	go ethhelper.SyncNftFromChain(strconv.Itoa(9651405/*9570987*/), true, buyResultCh, wethTransferCh, wethApproveCh, endCh)
	isOver := false
	for {
		select {
		case buyResult := <-buyResultCh:
			fmt.Println(buyResult)
			BlockTxs = append(BlockTxs, buyResult...)
		case wethTransfer := <-wethTransferCh:
			fmt.Println(wethTransfers)
			wethTransfers = append(wethTransfers, wethTransfer)
		case wethApprove := <-wethApproveCh:
			fmt.Println(wethTransfers)
			wethApproves = append(wethApproves, wethApprove)
		case <-endCh:
			isOver = true
			break
		default:
		}
		if isOver {
			break
		}
	}
	fmt.Println("end")
}

func TestGetBlockTxs(t *testing.T) {
	//txs := GetBlockTxs(9508909)
	txs, _, _ := GetBlockTxs(9508910)
	fmt.Println(len(txs))
}
