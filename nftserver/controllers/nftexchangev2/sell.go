package nftexchangev2

import (
	"encoding/json"
	"fmt"
	"github.com/nftexchange/nftserver/common/signature"
	"github.com/nftexchange/nftserver/controllers"
	"github.com/nftexchange/nftserver/models"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//售卖(上架)
func (nft *NftExchangeControllerV2) Sell() {
	fmt.Println("Sell()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
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
		inputDataErr := nft.verifyInputData_Sell(data, token)
		if inputDataErr != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = inputDataErr.Error()
			httpResponseData.Data = []interface{}{}
		} else {

			rawData := signature.RemoveSignData(string(bytes))
			approveAddr, _ := approveAddrsMap.GetApproveAddr(data["user_addr"])
			fmt.Println("sell user_addr:", data["user_addr"], ", approve_addr:", approveAddr)
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
					//modify the database value if the valification address is valid.
					price1, _ := strconv.ParseUint(data["price1"], 10, 64)
					price2, _ := strconv.ParseUint(data["price2"], 10, 64)
					//startTime, _ := time.Parse("2006-01-02 15:04:05", data["time1"])
					//closeTime, _ := time.Parse("2006-01-02 15:04:05", data["time2"])
					//startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", data["time1"], time.Local)
					//closeTime, _ := time.ParseInLocation("2006-01-02 15:04:05", data["time2"], time.Local)
					days, _ := strconv.Atoi(strings.TrimSpace(data["day"]))
					err := nd.Sell(data["user_addr"], "", data["nft_contract_addr"],
						data["nft_token_id"], data["selltype"], data["pay_channel"], days,
						price1, price2, data["royalty"], data["currency"], data["hide"], data["sig"], data["trade_sig"])
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
	fmt.Println("Sell()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())

}

func (nft *NftExchangeControllerV2) verifyInputData_Sell(data map[string]string, token string) error {
	regString, _ := regexp.Compile(PattenString)
	regNumber, _ := regexp.Compile(PattenNumber)

	if data["user_addr"] != "" {
		match := regString.MatchString(data["user_addr"])
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
	if data["selltype"] == "" {
		return ERRINPUTINVALID
	} else {
		match := regString.MatchString(data["selltype"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["pay_channel"] != "" {
		match := regString.MatchString(data["pay_channel"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["days"] != "" {
		match := regNumber.MatchString(data["days"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["price1"] != "" {
		match := regNumber.MatchString(data["price1"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["price2"] != "" {
		match := regNumber.MatchString(data["price2"])
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
	if data["currency"] != "" {
		match := regString.MatchString(data["currency"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["hide"] != "" {
		match := regString.MatchString(data["hide"])
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
	if data["trade_sig"] != "" {
		match := regString.MatchString(data["trade_sig"])
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


