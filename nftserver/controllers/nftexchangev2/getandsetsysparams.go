package nftexchangev2

import (
	"encoding/json"
	"fmt"
	"github.com/nftexchange/nftserver/models"
	"io/ioutil"
)

func (nft *NftExchangeControllerV2) GetSysParams() {
	sysParams := models.SysParamsRec{}

	nd := new(models.NftDb)
	err := nd.ConnectDB(models.Sqldsndb)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()
	param, err := nd.QuerySysParams()

	if err == nil {
		sysParams = *param
	}

	responseData, _ := json.Marshal(sysParams)
	nft.Ctx.ResponseWriter.Write(responseData)
}

func (nft *NftExchangeControllerV2) SetSysParams() {
	sysParams := models.SysParamsRec{}
	nd := new(models.NftDb)
	err := nd.ConnectDB(models.Sqldsndb)
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
	}
	defer nd.Close()

	requestData, _ := ioutil.ReadAll(nft.Ctx.Request.Body)
	defer nft.Ctx.Request.Body.Close()
	err = json.Unmarshal(requestData, &sysParams)
	if err != nil {
		nft.Ctx.ResponseWriter.Write([]byte("更新数据失败！"))
	} else {
		err := nd.SetSysParams(sysParams)
		if err != nil {
			nft.Ctx.ResponseWriter.Write([]byte("更新数据失败！"))
		} else {
			nft.Ctx.ResponseWriter.Write([]byte("更新数据成功！"))
		}
	}
}