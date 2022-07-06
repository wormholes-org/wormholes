package nftexchangev2

import (
	"encoding/json"
	"fmt"
	"github.com/nftexchange/nftserver/controllers"
	"github.com/nftexchange/nftserver/models"
	"gorm.io/gorm"
	"time"
)

//查询所有nft作品:get
func (nft *NftExchangeControllerV2) QueryAllNftProducts() {
	fmt.Println("QueryAllNftProducts()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd := new(models.NftDb)
	fmt.Println("sqldsndb=", models.Sqldsndb)
	err := nd.ConnectDB(models.Sqldsndb)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()

	nftData, err := nd.QueryNft()
	if err == nil {
		httpResponseData.Code = "200"
		httpResponseData.Data = nftData
	} else {
		if err == gorm.ErrRecordNotFound || err == models.ErrNftNotExist {
			httpResponseData.Code = "200"
		} else {
			httpResponseData.Code = "500"
		}

		httpResponseData.Msg = err.Error()
		httpResponseData.Data = []interface{}{}
	}

	responseData, _ := json.Marshal(httpResponseData)
	//fmt.Printf("nftData=%s", nftData)
	//fmt.Println("responseData", responseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	//nft.Data["json"] = responseData
	//nft.ServeJSON()
	fmt.Println("QueryAllNftProducts()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}
