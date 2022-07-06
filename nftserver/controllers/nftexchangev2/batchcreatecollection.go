package nftexchangev2

import (
	"encoding/json"
	"fmt"
	"github.com/nftexchange/nftserver/controllers"
	"github.com/nftexchange/nftserver/models"
	"io/ioutil"
	"regexp"
	"time"
)

//创建集合:post
func (nft *NftExchangeControllerV2) BatchCreateCollection() {
	fmt.Println("BatchCreateCollection()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
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
		//token := nft.Ctx.Request.Header.Get("Token")
		//fmt.Println("create new collection, token=", token)
		inputDataErr := nft.verifyInputData_BatchCreateCollection(data)
		if inputDataErr != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = inputDataErr.Error()
			httpResponseData.Data = []interface{}{}
		} else {

			//rawData := signature.RemoveSignData(string(bytes))
			//approveAddr, _ := approveAddrsMap.GetApproveAddr(data["user_addr"])
			//_, err := signature.IsValidAddr(rawData, data["sig"], approveAddr)
			//if err != nil {
			//	httpResponseData.Code = "500"
			//	httpResponseData.Msg = err.Error()
			//	httpResponseData.Data = []interface{}{}
			//} else {
			//	err = nd.InsertSigData(data["sig"], rawData)
			//	if err != nil {
			//		httpResponseData.Code = "500"
			//		httpResponseData.Msg = err.Error()
			//		httpResponseData.Data = []interface{}{}
			//	} else {
					err = nd.NewCollections(data["user_addr"], data["name"],
						data["img"], data["contract_type"], data["contract_addr"],
						data["desc"], data["categories"], data["sig"])
					if err == nil {
						httpResponseData.Code = "200"
						httpResponseData.Data = []interface{}{}
					} else {
						httpResponseData.Code = "500"
						httpResponseData.Msg = err.Error()
						httpResponseData.Data = []interface{}{}
					}
				//}
			//}
		}
	} else {
		httpResponseData.Code = "500"
		httpResponseData.Msg = "输入的用户信息错误"
		httpResponseData.Data = []interface{}{}
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("BatchCreateCollection()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) verifyInputData_BatchCreateCollection(data map[string]string) error {
	regString, _ := regexp.Compile(PattenString)
	regAddr, _ := regexp.Compile(PattenAddr)
	regImage, _ := regexp.Compile(PattenImageBase64)
	//regNumber, _ := regexp.Compile(PattenNumber)

	if data["user_addr"] != "" {
		match := regString.MatchString(data["user_addr"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	//if data["name"] != "" {
	//	match := regString.MatchString(data["name"])
	//	if !match {
	//		return ERRINPUTINVALID
	//	}
	//}

	match := regImage.MatchString(data["img"])
	if !match {
		return ERRINPUTINVALID
	}

	if data["contract_type"] != "" {
		match := regString.MatchString(data["contract_type"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["contract_addr"] != "" {
		match := regString.MatchString(data["contract_addr"])
		if !match {
			return ERRINPUTINVALID
		}
		match = regAddr.MatchString(data["contract_addr"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	//if data["desc"] != "" {
	//	match := regString.MatchString(data["desc"])
	//	if !match {
	//		return ERRINPUTINVALID
	//	}
	//}
	/*if data["royalty"] != "" {
		match := regNumber.MatchString(data["royalty"])
		if !match {
			return ERRINPUTINVALID
		}
	}*/
	if data["categories"] != "" {
		match := regString.MatchString(data["categories"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["sig"] != "" {
		match := regString.MatchString(data["sig"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	//getToken, _ := tokenMap.GetToken(data["user_addr"])
	//if getToken != token {
	//	return ERRTOKEN
	//}

	return nil
}

