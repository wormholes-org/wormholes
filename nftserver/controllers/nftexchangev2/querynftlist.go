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

func (nft *NftExchangeControllerV2) QueryNftList() {
	fmt.Println("QueryNftList()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData controllers.HttpResponseData
	nd := new(models.NftDb)
	err := nd.ConnectDB(models.Sqldsndb)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	var data controllers.HttpRequestFilter
	bytes, _ := ioutil.ReadAll(nft.Ctx.Request.Body)
	defer nft.Ctx.Request.Body.Close()
	err = json.Unmarshal(bytes, &data)
	if err == nil {
		inputDataErr := nft.verifyInputData_QueryNftList(data)
		if inputDataErr != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = inputDataErr.Error()
			httpResponseData.Data = []interface{}{}
		} else {

			nfts, totalCount, err := nd.QueryNftByFilter(data.Filter, data.Sort, data.StartIndex, data.Count)
			if err == nil {
				httpResponseData.Code = "200"
				httpResponseData.Data = nfts
				httpResponseData.TotalCount = totalCount
			} else {
				if err == gorm.ErrRecordNotFound || err == models.ErrNftNotExist {
					httpResponseData.Code = "200"
				} else {
					httpResponseData.Code = "500"
				}
				httpResponseData.Msg = err.Error()
				httpResponseData.Data = []interface{}{}
			}
		}

	} else {
		httpResponseData.Code = "500"
		httpResponseData.Data = []interface{}{}
		httpResponseData.Msg = "输入的用户信息错误"
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("QueryNftList()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) verifyInputData_QueryNftList(data controllers.HttpRequestFilter) error {
	regString, _ := regexp.Compile(PattenString)
	regNumber, _ := regexp.Compile(PattenNumber)
	regOperator, _ := regexp.Compile(PattenOperator)

	if len(data.Filter) > 0 {
		for _, v := range data.Filter {
			match := regString.MatchString(v.Field)
			if !match {
				fmt.Println("v.Field=", v.Field)
				return ERRINPUTINVALID
			}
			match = regOperator.MatchString(v.Operation)
			if !match {
				fmt.Println("v.Operation=", v.Operation)
				return ERRINPUTINVALID
			}

			if v.Field == "collections" ||
				strings.Contains(v.Field, "name") ||
				strings.Contains(v.Field, "desc") {
				continue
			}

			match = regString.MatchString(v.Value)
			if !match {
				fmt.Println("v.Value=", v.Value)
				return ERRINPUTINVALID
			}
		}
	}
	if len(data.Sort) > 0 {
		for _, v := range data.Sort {
			match := regString.MatchString(v.By)
			if !match {
				fmt.Println("v.By=", v.By)
				return ERRINPUTINVALID
			}
			match = regString.MatchString(v.Order)
			if !match {
				fmt.Println("v.Order=", v.Order)
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

