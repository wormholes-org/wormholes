package nftexchangev1

import (
	"encoding/json"
	"fmt"
	"github.com/nftexchange/nftserver/controllers"
	"github.com/nftexchange/nftserver/models"
	"io/ioutil"
)

//用户登录(不存在时创建):post
func (nft *NftExchangeControllerV1) UserLogin() {
	var httpResponseData controllers.HttpResponseData
	nd := new(models.NftDb)
	err := nd.ConnectDB(models.Sqldsndb)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	var data map[string]string
	bytes, _ := ioutil.ReadAll(nft.Ctx.Request.Body)
	defer nft.Ctx.Request.Body.Close()
	err = json.Unmarshal(bytes, &data)
	if err == nil {
		err = nd.Login(data["user_addr"], data["sig"])
		if err == nil {
			httpResponseData.Code = "200"
			httpResponseData.Data = []interface{}{}
		} else {
			httpResponseData.Code = "500"
			httpResponseData.Msg = err.Error()
		}

	} else {
		httpResponseData.Code = "500"
		httpResponseData.Msg = "输入的用户信息错误"
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	//nft.Data["json"] = responseData
	//nft.ServeJSON()
}
