package nftexchangev1

import (
	beego "github.com/beego/beego/v2/server/web"
)

type NftExchangeControllerV1 struct {
	beego.Controller
}

//func (nft *NftExchangeControllerV1) GetImageFromIPFS() {
//	var data map[string]interface{}
//	bytes, _ := ioutil.ReadAll(nft.Ctx.Request.Body)
//	defer nft.Ctx.Request.Body.Close()
//	json.Unmarshal(bytes, &data)
//	s, ok := data["hash"].(string)
//	fmt.Printf(">>>>>>>>s=%s\n", s)
//	if ok {
//
//	}
//}

