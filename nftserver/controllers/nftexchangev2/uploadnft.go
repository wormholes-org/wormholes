package nftexchangev2

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/nftexchange/nftserver/common/signature"
	"github.com/nftexchange/nftserver/controllers"
	"github.com/nftexchange/nftserver/models"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

//上传nft作品:post
func (nft *NftExchangeControllerV2) UploadNft() {
	fmt.Println("UploadNft()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
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
		inputDataErr := nft.verifyInputData_UploadNft(data, token)
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
					}
				}
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
	fmt.Println("UploadNft()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) AddFileToIpfs(content string) (string, error) {
	serverIP, _ := beego.AppConfig.String("nftIpfsServerIP")
	serverPort, _ := beego.AppConfig.String("nftIpfsServerPort")
	url := "http://" + serverIP + ":" + serverPort + "/v1/ipfsadd"
	fmt.Println("NftExchangeController.AddFileToIpfs(), url=", url)
	var mapImage map[string]string
	mapImage = make(map[string]string, 0)
	mapImage["asset"] = content
	respData, err := nft.SendPost(url, mapImage, "application/json")
	if err != nil {
		return "", err
	} else {
		cid, ok := respData.Data.(string)
		if ok {
			return cid, nil
		} else {
			return "", errors.New("nftipfs server 返回数据格式错误！")
		}
	}
}

//发送POST请求
//url:请求地址		data:POST请求提交的数据		contentType:请求体格式，如：application/json
//content:请求返回的内容
func (nft *NftExchangeControllerV2)SendPost(url string, data interface{}, contentType string) (respData controllers.HttpResponseData,  err error) {
	jsonStr, _ := json.Marshal(data)
	//fmt.Println("SendPost,url=", url, "jsonStr=", string(jsonStr))
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonStr))
	if contentType == "" {
		contentType = "application/json"
	}
	req.Header.Add("content-type", contentType)
	if err != nil {
		return
	}
	defer req.Body.Close()
	//client := &http.Client{Timeout: 5 * time.Second}
	client := &http.Client{}
	fmt.Println("SendPost url=", url)
	resp, error := client.Do(req)
	if error != nil {
		return
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(result, &respData)
	return
}


func (nft *NftExchangeControllerV2) IpfsTest() {
	content := "test content"
	s, e := nft.AddFileToIpfs(content)
	if e != nil {
		nft.Ctx.WriteString(e.Error())
	} else {
		nft.Ctx.WriteString(s)
	}
}

func (nft *NftExchangeControllerV2) verifyInputData_UploadNft(data map[string]string, token string) error {
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
	getToken, _ := tokenMap.GetToken(data["user_addr"])
	if getToken != token {
		return ERRTOKEN
	}

	return nil
}