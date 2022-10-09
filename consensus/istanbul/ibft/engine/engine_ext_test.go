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
	var header = types.Header{}
	header.Extra = istblExtra.EncodeRLP()

	messages, err := DecodeMessages(extra)
	if err != nil {
		t.Error(err)
	}
	t.Log(messages)

	addrs, err := GetValidatorRewardList(&header)
	if err != nil {
		t.Error(err)
	}
	t.Log(addrs)
}
