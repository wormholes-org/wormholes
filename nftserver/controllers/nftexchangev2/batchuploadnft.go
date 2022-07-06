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

//上传nft作品:post
func (nft *NftExchangeControllerV2) BatchUploadNft() {
	fmt.Println("BatchUploadNft()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
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
		//token := nft.Ctx.Request.Header.Get("Token")
		inputDataErr := nft.verifyInputData_BatchUploadNft(data)
		if inputDataErr != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = inputDataErr.Error()
			httpResponseData.Data = []interface{}{}
		} else {

			//调用nftipfs,将图片上传到ipfs文件服务器
			//cid, err := nft.AddFileToIpfs(data["asset_sample"])
			//if err != nil {
			//	httpResponseData.Code = "400"
			//	httpResponseData.Msg = err.Error()
			//	httpResponseData.Data = []interface{}{}
			//} else {
			//	fmt.Printf(">>>>>>>cid=%s\n", cid)
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
					err = nd.UploadNft(data["user_addr"], data["creator_addr"], data["owner_addr"],
						data["md5"], data["name"], data["desc"],
						data["meta"], data["source_url"],
						data["nft_contract_addr"], data["nft_token_id"],
						data["categories"], data["collections"],
						data["asset_sample"], data["hide"], data["royalty"], data["count"],data["sig"])
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
			//}
		}
	} else {
		httpResponseData.Code = "500"
		httpResponseData.Msg = "输入的用户信息错误"
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	//nft.Data["json"] = responseData
	//nft.ServeJSON()
	fmt.Println("BatchUploadNft()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) verifyInputData_BatchUploadNft(data map[string]string) error {
	regString, _ := regexp.Compile(PattenString)
	regNumber, _ := regexp.Compile(PattenNumber)
	regImage, _ := regexp.Compile(PattenImageBase64)

	if data["user_addr"] != "" {
		match := regString.MatchString(data["user_addr"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["creator_addr"] != "" {
		match := regString.MatchString(data["creator_addr"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["owner_addr"] != "" {
		match := regString.MatchString(data["owner_addr"])
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
	//if data["desc"] != "" {
	//	match := regString.MatchString(data["desc"])
	//	if !match {
	//		return ERRINPUTINVALID
	//	}
	//}
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
	if data["md5"] != "" {
		match := regString.MatchString(data["md5"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	if data["categories"] != "" {
		match := regString.MatchString(data["categories"])
		if !match {
			return ERRINPUTINVALID
		}
	}
	//if data["collections"] != "" {
	//	match := regString.MatchString(data["collections"])
	//	if !match {
	//		return ERRINPUTINVALID
	//	}
	//}
	match := regImage.MatchString(data["asset_sample"])
	if !match {
		return ERRINPUTINVALID
	}
	if data["hide"] != "" {
		match := regString.MatchString(data["hide"])
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
	if data["count"] != "" {
		match := regNumber.MatchString(data["count"])
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