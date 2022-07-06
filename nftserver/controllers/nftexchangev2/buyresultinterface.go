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

func (nft *NftExchangeControllerV2) BuyResultInterface() {
	fmt.Println("BuyResultInterface()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd := new(models.NftDb)
	err := nd.ConnectDB(models.Sqldsndb)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	var data map[string]string
	bytes, _ := ioutil.ReadAll(nft.Ctx.Request.Body)
	//fmt.Printf("receive data = %s\n", string(bytes))
	defer nft.Ctx.Request.Body.Close()
	err = json.Unmarshal(bytes, &data)
	if err == nil {
		token := nft.Ctx.Request.Header.Get("Token")
		inputDataErr := nft.verifyInputData_BuyResultInterface(data, token)
		if inputDataErr != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = inputDataErr.Error()
			httpResponseData.Data = []interface{}{}
		} else {

			rawData := signature.RemoveSignData(string(bytes))
			_, inputDataErr = signature.IsValidVerifyAddr(rawData, data["admin_sig"])

			if inputDataErr != nil {
				httpResponseData.Code = "500"
				httpResponseData.Msg = inputDataErr.Error()
				httpResponseData.Data = []interface{}{}
			} else {
				inputDataErr = nd.InsertSigData(data["admin_sig"], rawData)
				if inputDataErr != nil {
					httpResponseData.Code = "500"
					httpResponseData.Msg = inputDataErr.Error()
					httpResponseData.Data = []interface{}{}
				} else {
					inputDataErr = nd.BuyResultInterface(data["admin_addr"], data["from"], data["to"],
						data["nft_contract_addr"], data["nft_token_id"],
						data["trade_sig"], data["price"], data["royalty"],
						data["txhash"], data["sig"], data["admin_sig"])
					if inputDataErr == nil {
						httpResponseData.Code = "200"
						httpResponseData.Data = []interface{}{}
					} else {
						httpResponseData.Code = "500"
						httpResponseData.Msg = inputDataErr.Error()
						httpResponseData.Data = []interface{}{}
					}
				}
			}
		}

	} else {
		httpResponseData.Code = "500"
		httpResponseData.Data = []interface{}{}
		httpResponseData.Msg = "输入的用户信息错误"
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("BuyResultInterface()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) verifyInputData_BuyResultInterface(data map[string]string, token string) error {
	regString, _ := regexp.Compile(PattenString)
	regNumber, _ := regexp.Compile(PattenNumber)
	if data["admin_addr"] != "" {
		match := regString.MatchString(data["admin_addr"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["from"] != "" {
		match := regString.MatchString(data["from"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["to"] != "" {
		match := regString.MatchString(data["to"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["nft_contract_addr"] != "" {
		match := regString.MatchString(data["nft_contract_addr"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["nft_token_id"] != "" {
		match := regString.MatchString(data["nft_token_id"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["price"] != "" {
		match := regNumber.MatchString(data["price"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["trade_sig"] != "" {
		match := regString.MatchString(data["trade_sig"])
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
	if data["txhash"] != "" {
		match := regString.MatchString(data["txhash"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["royalty"] != "" {
		match := regNumber.MatchString(data["royalty"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["admin_sig"] != "" {
		match := regString.MatchString(data["admin_sig"])
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

