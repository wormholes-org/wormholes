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
	"time"
)
type ResponseLogin struct {
	Hash string `json:"hash"`
	Secret string `json:"secret"`
	TimeStamp int64 `json:"time_stamp"`
	Token string `json:"token"`
}

//用户登录(不存在时创建):post
func (nft *NftExchangeControllerV2) UserLogin() {
	fmt.Println("UserLogin()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
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
		inputDataErr := nft.verifyInputData_UserLogin(data)
		if inputDataErr != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = inputDataErr.Error()
			httpResponseData.Data = []interface{}{}
		} else {
			ResLogin := ResponseLogin{}
			if data["result"] == "" {
				hashString, secret, result, err := answerMapList.GetHashAndSecret()
				if err != nil {
					httpResponseData.Code = "500"
					httpResponseData.Msg = err.Error()
					httpResponseData.Data = []interface{}{}
				} else {
					ResLogin.Hash = hashString
					ResLogin.Secret = secret
					fmt.Println("resutl = ", result)
					timeKey := answerMapList.AddResult(data["user_addr"], result)
					ResLogin.TimeStamp = timeKey
					httpResponseData.Code = "200"
					httpResponseData.Data = []ResponseLogin{ResLogin}
				}
			} else {
				timeStamp, _ := strconv.ParseInt(data["time_stamp"], 10, 64)
				result, _ := strconv.ParseInt(data["result"], 10, 64)
				valid := answerMapList.IsValid(timeStamp, data["user_addr"], result)
				if valid == false {
					httpResponseData.Code = "500"
					httpResponseData.Msg = "answer is wrong!"
					httpResponseData.Data = []interface{}{}
				} else {
					rawData := signature.RemoveSignData(string(bytes))
					_, err := signature.IsValidAddr(rawData, data["sig"], data["user_addr"])
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
							err = nd.Login(data["user_addr"], data["sig"])
							if err == nil {
								token, err := GenerateToken(data["user_addr"])
								if err != nil {
									httpResponseData.Code = "500"
									httpResponseData.Data = []interface{}{}
									httpResponseData.Msg = "授权失败!"
								} else {
									httpResponseData.Code = "200"

									ResLogin.Token = token
									httpResponseData.Data = []ResponseLogin{ResLogin}
									tokenMap.AddToken(data["user_addr"], token)
									fmt.Println("login user_addr:", data["user_addr"], ", approve_addr:", data["approve_addr"])
									approveAddrsMap.AddApproveAddr(data["user_addr"], data["approve_addr"])
									//nft.Ctx.ResponseWriter.Header().Set("Token", token)
								}
							} else {
								httpResponseData.Code = "500"
								httpResponseData.Msg = err.Error()
								httpResponseData.Data = []interface{}{}
							}
						}
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
	fmt.Println("UserLogin()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}

func (nft *NftExchangeControllerV2) verifyInputData_UserLogin(data map[string]string) error {
	regString, _ := regexp.Compile(PattenString)

	if data["user_addr"] != "" {
		match := regString.MatchString(data["user_addr"])
		if !match {
			return ERRINPUTINVALID
		}
	}

	_, ok := data["result"]
	if !ok {
		return ERRINPUTINVALID
	}

	if data["sig"] != "" {
		match := regString.MatchString(data["sig"])
		if !match {
			return ERRINPUTINVALID
		}
	}

	return nil
}

