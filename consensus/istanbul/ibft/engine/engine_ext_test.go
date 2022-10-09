package ibftengine

import (
	core "github.com/ethereum/go-ethereum/consensus/istanbul/ibft/core"
	"github.com/ethereum/go-ethereum/core/types"
	"testing"
)

func TestGetValidatorList(t *testing.T) {
	rwdExtra, err := core.GenTestExtra()
	istblExtra := types.IstanbulExtra{}
	istblExtra.RewardList = rwdExtra

	messages, err := DecodeMessages(rwdExtra)
	if err != nil {
		t.Error(err)
	}
	t.Log(messages)

	//var header = types.Header{}
	//header.Extra = istblExtra.EncodeRLP()
	//addrs, err := GetValidatorRewardList(&header)
	//if err != nil {
	//	t.Error(err)
	//}
	//t.Log(addrs)
}
