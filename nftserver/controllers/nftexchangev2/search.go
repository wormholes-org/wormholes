package nftexchangev2

import (
	"encoding/json"
	"fmt"
	"github.com/nftexchange/nftserver/controllers"
	"github.com/nftexchange/nftserver/models"
	"io/ioutil"
	"time"
)

//查询用户所有信息:post
func (nft *NftExchangeControllerV2) Search() {
	fmt.Println("Search()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
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
	s, ok := data["match"].(string)
	fmt.Printf(">>>>>>>>s=%s\n", s)
	if ok {
		inputDataErr := nft.verifyInputData_Search(s)
		if inputDataErr != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = inputDataErr.Error()
			httpResponseData.Data = []interface{}{}
		} else {

			searchInfo, err := nd.Search(s)
			if err == nil {
				httpResponseData.Code = "200"
				httpResponseData.Data = searchInfo
			} else {
				httpResponseData.Code = "500"
				httpResponseData.Msg = err.Error()
				httpResponseData.Data = []interface{}{}
			}
		}

	} else {
		httpResponseData.Code = "500"
		httpResponseData.Msg = "输入的用户信息错误"
		httpResponseData.Data = []interface{}{}
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	//nft.Data["json"] = responseData
	//nft.ServeJSON()
	fmt.Println("Search()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) verifyInputData_Search(user string) error {
	//regString, _ := regexp.Compile(PattenString)
	//if user != "" {
	//	match := regString.MatchString(user)
	//	if !match {
	//		return ERRINPUTINVALID
	//	}
	//}

	return nil
}
