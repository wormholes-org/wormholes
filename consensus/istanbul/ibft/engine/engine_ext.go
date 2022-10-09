package ibftengine

import (
	ibfttypes "github.com/ethereum/go-ethereum/consensus/istanbul/ibft/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

func DecodeMessages(data []byte) ([]*ibfttypes.Message, error) {
	var commitseals []*ibfttypes.Message
	msg := new(ibfttypes.Message)
	err := rlp.DecodeBytes(data, &msg)
	if err != nil {
		return nil, err
	}
	err = msg.DecodeCommitSeals(&commitseals)
	if err != nil {
		log.Error("Decode Msg Err")
	}
	return commitseals, nil
}
