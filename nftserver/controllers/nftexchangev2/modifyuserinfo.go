package nftexchangev2

import (
	"encoding/json"
	"fmt"
	"github.com/nftexchange/nftserver/common/signature"
	"github.com/nftexchange/nftserver/controllers"
	"github.com/nftexchange/nftserver/models"
	"io/ioutil"
	"regexp"
	"time"
)

//修改用户信息
func (nft *NftExchangeControllerV2) ModifyUserInfo() {
	fmt.Println("ModifyUserInfo()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd := new(models.NftDb)
	err := nd.ConnectDB(models.Sqldsndb)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()

	var data map[string]string
	defer nft.Ctx.Request.Body.Close()
	bytes, _ := ioutil.ReadAll(nft.Ctx.Request.Body)
	//fmt.Printf("receive data = %s\n", string(bytes))
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		httpResponseData.Code = "500"
		httpResponseData.Msg = err.Error()
		httpResponseData.Data = []interface{}{}
	} else {
		token := nft.Ctx.Request.Header.Get("Token")
		inputDataErr := nft.verifyInputData_ModifyUserInfo(data, token)
		if inputDataErr != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = inputDataErr.Error()
			httpResponseData.Data = []interface{}{}
		} else {

			rawData := signature.RemoveSignData(string(bytes))
			approveAddr, _ := approveAddrsMap.GetApproveAddr(data["user_addr"])
			_, err := signature.IsValidAddr(rawData, data["sig"], approveAddr)
			if err != nil {
				httpResponseData.Code = "500"
				httpResponseData.Msg = err.Error()
				httpResponseData.Data = []interface{}{}
			} else {
				err = nd.InsertSigData(data["sig"], rawData)
				if err != nil {
					httpResponseData.Code = "500"
					httpResponseData.Msg = err.Error()
					httpResponseData.Data = []interface{}{}
				} else {
					// modify the information of user if
					err:= nd.ModifyUserInfo(data["user_addr"], data["user_name"],
						data["portrait"], data["user_mail"], data["user_info"], data["sig"])
					if err != nil {
						httpResponseData.Code = "500"
						httpResponseData.Msg = err.Error()
						httpResponseData.Data = []interface{}{}
					} else {
						httpResponseData.Code = "200"
						httpResponseData.Data = []interface{}{}
					}
				}
			}
		}

	}

	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("ModifyUserInfo()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) verifyInputData_ModifyUserInfo(data map[string]string, token string) error {
	regString, _ := regexp.Compile(PattenString)
	regEmail, _ := regexp.Compile(PattenEmail)
	regImage, _ := regexp.Compile(PattenImageBase64)

	if data["user_addr"] != "" {
		match := regString.MatchString(data["user_addr"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	//if data["user_name"] !=  "" {
	//	match := regString.MatchString(data["user_name"])
	//	if !match {
	//		return ERRINPUTINVALID
	//	}
	//}
	//if data["user_info"] != "" {
	//	match := regString.MatchString(data["user_info"])
	//	if !match {
	//		return ERRINPUTINVALID
	//	}
	//}
	if data["portrait"] != "" {
		match := regImage.MatchString(data["portrait"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["user_mail"] != "" {
		match := regEmail.MatchString(data["user_mail"])
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
	getToken, _ := tokenMap.GetToken(data["user_addr"])
	if getToken != token {
		return ERRTOKEN
	}

	return nil
}

