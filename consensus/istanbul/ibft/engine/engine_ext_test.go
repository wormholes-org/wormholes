package ibftengine

import (
	core "github.com/ethereum/go-ethereum/consensus/istanbul/ibft/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"testing"
)

func TestGetValidatorList(t *testing.T) {
	rwdExtra, err := core.GenTestExtra()

	messages, err := DecodeMessages(rwdExtra)
	if err != nil {
		t.Error(err)
	}
	t.Log(messages)

	addrs, _ := MessagesToAddress(messages)
	t.Log(addrs)

	var header = types.Header{}
	istblExtra := types.IstanbulExtra{
		RewardList: rwdExtra,
	}
	header.Extra, _ = rlp.EncodeToBytes(istblExtra)
	addrs, err = GetValidatorRewardList(&header)
	if err != nil {
		t.Error(err)
	}
	t.Log(addrs)
}
