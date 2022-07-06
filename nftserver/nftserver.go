package nftserver

import (
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/nftexchange/nftserver/common/contracts"
	"github.com/nftexchange/nftserver/models"
	_ "github.com/nftexchange/nftserver/routers"
)

func NftServerRun(nftServerStop chan struct{}) {
	defer close(nftServerStop)
	//8c995fd78bddf528bd548cce025f62d4c3c0658362dbfd31b23414cf7ce2e8ed
	//verify := signature.VerifyAppconf("./conf/app.conf", "0x2b0aD05ADDa21BA4E5b94C4f9aE3BCeA15A380c5")
	//if verify != true {
	//	fmt.Println("app.conf verify ")
	//	return
	//}
	err := models.InitSysParams(models.Sqldsndb)
	if err != nil {
		fmt.Println("InitSysParams err=", err)
		return
	}
	fmt.Println("NftserverRun****************************************************")
	/*err :=  os.Remove("./conf/app.conf")
	if err != nil {
		fmt.Println("delete app.conf err=", err)
		return
	}*/
	fmt.Println(models.TradeAddr)
	fmt.Println(models.NFT1155Addr)
	fmt.Println(models.AdminAddr)
	fmt.Println(models.EthersNode)
	fmt.Println(models.EthersWsNode)
	err = models.InitSyncBlockTs(models.Sqldsndb)
	if err != nil {
		fmt.Println("init err exit")
		return
	}
	go TimeProc(models.Sqldsndb)
	go contracts.EventContract(models.Sqldsndb)

	beego.Run()
}

