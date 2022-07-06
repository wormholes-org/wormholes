package ethhelper

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/nftexchange/nftserver/ethhelper/common"
	"math/big"
)

// BalanceOfWeth
// @description 查询账户weth的余额
// @auth chen.gang 2021/9/1 10:57
// @param address
// @return balance  err
func BalanceOfWeth(addr string) (uint64, error) {
	var tmp big.Int
	payload := make([]byte, 36)

	tmp.SetString(common.BalanceOfHash, 0)
	copy(payload[:4], tmp.Bytes())
	tmp.SetString(addr, 0)
	copy(payload[36-len(tmp.Bytes()):], tmp.Bytes())

	params := common.CallParamTemp{To: common.WETH, Data: "0x" + hex.EncodeToString(payload)}
	jsonData, err := json.Marshal(params)
	if err != nil {
		return 0, errors.New("Umarshal failed:" + err.Error() + string(jsonData))
	}

	var ret string
	if err = common.Client().Call(&ret, "eth_call", params, "latest"); err != nil {
		return 0, errors.New("Call failed:" + err.Error())
	} else {
		tmp.SetString(ret, 0)
		return tmp.Uint64(), nil
	}
}

// AllowanceOfWeth
// @description weth的授权额度
// @auth chen.gang 2021/9/1 10:57
// @param address
// @return balance  err
func AllowanceOfWeth(addr string) (string, error) {
	var tmp big.Int
	payload := make([]byte, 68)
	tmp.SetString(common.AllowanceHash, 0)
	copy(payload[:4], tmp.Bytes())
	tmp.SetString(addr, 0)
	copy(payload[36-len(tmp.Bytes()):36], tmp.Bytes())
	tmp.SetString(common.TradeCore, 0)
	copy(payload[68-len(tmp.Bytes()):], tmp.Bytes())

	params := common.CallParamTemp{To: common.WETH, Data: "0x" + hex.EncodeToString(payload)}
	jsonData, err := json.Marshal(params)
	if err != nil {
		return "0", errors.New("Umarshal failed:" + err.Error() + string(jsonData))
	}

	var ret string
	if err = common.Client().Call(&ret, "eth_call", params, "latest"); err != nil {
		return "0", errors.New("Call failed:" + err.Error())
	} else {
		tmp.SetString(ret, 0)
		return tmp.String(), nil
	}
}
