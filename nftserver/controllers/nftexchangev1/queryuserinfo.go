package nftexchangev1

import (
	"encoding/json"
	"fmt"
	"github.com/nftexchange/nftserver/controllers"
	"github.com/nftexchange/nftserver/models"
	"io/ioutil"
)

//查询用户所有信息:post
func (nft *NftExchangeControllerV1) QueryUserInfo() {
	//defer nft.Ctx.Request.Body.Close()
	var httpResponseData controllers.HttpResponseData
	nd := new(models.NftDb)
	err := nd.ConnectDB(models.Sqldsndb)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()

	var data map[string]interface{}
	bytes, _ := ioutil.ReadAll(nft.Ctx.Request.Body)
	defer nft.Ctx.Request.Body.Close()
	json.Unmarshal(bytes, &data)
	s, ok := data["user_addr"].(string)
	fmt.Printf(">>>>>>>>s=%s\n", s)
	if ok {
		user, err := nd.QueryNftbyUser(s)
		if err == nil {
			httpResponseData.Code = "200"
			httpResponseData.Data = user
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
