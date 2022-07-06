package nftexchangev2

import (
	"encoding/json"
	"fmt"
	"github.com/nftexchange/nftserver/controllers"
	"github.com/nftexchange/nftserver/models"
	"gorm.io/gorm"
	"io/ioutil"
	"regexp"
	"strings"
	"time"
)

//获取用户交易历史
func (nft *NftExchangeControllerV2) QueryMarketTradingHistory() {
	fmt.Println("QueryMarketTradingHistory()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd := new(models.NftDb)
	err := nd.ConnectDB(models.Sqldsndb)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()

	var data controllers.HttpRequestFilter
	defer nft.Ctx.Request.Body.Close()
	bytes, _ := ioutil.ReadAll(nft.Ctx.Request.Body)
	//fmt.Printf("receive data = %s\n", string(bytes))
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		httpResponseData.Code = "500"
		httpResponseData.Msg = err.Error()
		httpResponseData.Data = []interface{}{}
	} else {

		inputDataErr := nft.verifyInputData_QueryMarketTradingHistory(data)
		if inputDataErr != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = inputDataErr.Error()
			httpResponseData.Data = []interface{}{}
		} else {
			userTradings, totalCount, err := nd.QueryMarketTradingHistory(data.Filter, data.Sort, data.StartIndex, data.Count)
			if err != nil {
				if err == gorm.ErrRecordNotFound || err == models.ErrNftNotExist {
					httpResponseData.Code = "200"
				} else {
					httpResponseData.Code = "500"
				}
				httpResponseData.Msg = err.Error()
				httpResponseData.Data = []interface{}{}
			} else {
				httpResponseData.Code = "200"
				httpResponseData.Data = userTradings
				httpResponseData.TotalCount = uint64(totalCount)
			}
		}
	}

	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("QueryUserTradingHistory()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) verifyInputData_QueryMarketTradingHistory(data controllers.HttpRequestFilter) error {
	regString, _ := regexp.Compile(PattenString)
	regNumber, _ := regexp.Compile(PattenNumber)
	regOperator, _ := regexp.Compile(PattenOperator)

	if len(data.Filter) > 0 {
		for _, v := range data.Filter {
			match := regString.MatchString(v.Field)
			if !match {
				return ERRINPUTINVALID
			}
			match = regOperator.MatchString(v.Operation)
			if !match {
				return ERRINPUTINVALID
			}

			if v.Field == "collections" ||
				strings.Contains(v.Field, "name") ||
				strings.Contains(v.Field, "desc") {
				continue
			}

			match = regString.MatchString(v.Value)
			if !match {
				return ERRINPUTINVALID
			}

		}
	}
	if len(data.Sort) > 0 {
		for _, v := range data.Sort {
			match := regString.MatchString(v.By)
			if !match {
				return ERRINPUTINVALID
			}
			match = regString.MatchString(v.Order)
			if !match {
				return ERRINPUTINVALID
			}
		}
	}
	if data.StartIndex != "" {
		match := regNumber.MatchString(data.StartIndex)
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data.Count != "" {
		match := regNumber.MatchString(data.Count)
		if !match {
			return ERRINPUTINVALID
		}
	}

	return nil
}

