package ethhelper

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/nftexchange/nftserver/ethhelper/common"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"
)

// AdminList
// @description 管理员列表
// @auth chen.gang 2021/9/1 10:57
// @param
// @return []string
func AdminList() ([]string, error) {
	var tmp big.Int
	payload := make([]byte, 4)
	tmp.SetString(common.AdminListHash, 0)
	copy(payload[:4], tmp.Bytes())

	params := common.CallParam{From: "0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900", To: common.Admin, Data: "0x" + hex.EncodeToString(payload)}
	var ret string
	var list []string
	if err := common.Client().Call(&ret, "eth_call", params, "latest"); err != nil {
		return nil, errors.New("Call failed:" + err.Error())
	} else {
		ret = ret[2:]
		if len(ret) > 128 {
			for i := 128; i < len(ret); i += 64 {
				list = append(list, ret[i+24:i+64])
			}
		}

		return list, nil
	}
}

func postAuctionTx(data interface{}) (error, string) {
	contentType := "application/json"
	client := &http.Client{Timeout: 5 * time.Second}
	jsonStr, _ := json.Marshal(data)
	resp, err := client.Post("http://localhost:1880/v1/auction", contentType, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err, ""
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	return nil, string(result)
}

// SendDealAuctionTx
// @description 发送离线签名交易
// @auth chen.gang 2021/9/1 10:57
// @param
// @return []string
func SendDealAuctionTx(from, to, nftAddr, tokenId, price, sig string) (error, string) {
	type Res struct {
		From    string `json:"from"`
		To      string `json:"to"`
		NftAddr string `json:"nftAddr"`
		TokenId string `json:"tokenId"`
		Price   string `json:"price"`
		Sig     string `json:"sig"`
	}
	data := Res{From: from, To: to, NftAddr: nftAddr, TokenId: tokenId, Price: price, Sig: sig}
	return postAuctionTx(data)
}
