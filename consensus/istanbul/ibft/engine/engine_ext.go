package ibftengine

import (
	"github.com/ethereum/go-ethereum/common"
	ibfttypes "github.com/ethereum/go-ethereum/consensus/istanbul/ibft/types"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

func DecodeMessages(data []byte) ([]*ibfttypes.Message, error) {
	msg := new(ibfttypes.Message)
	err := rlp.DecodeBytes(data, &msg)
	if err != nil {
		log.Info("rlp code:", string(data))
		return nil, err
	}

	var messages []*ibfttypes.Message
	err = msg.DecodeCommitSeals(&messages)
	if err != nil {
		log.Info("commited seal:", string(msg.CommittedSeal))
		return nil, err
	}
	return messages, nil
}

func GetValidatorRewardList(header *types.Header) ([]common.Address, error) {
	istanbulExtra, err := types.ExtractIstanbulExtra(header)
	if err != nil {
		return nil, err
	}
	messages, err := DecodeMessages(istanbulExtra.RewardList)
	//Decode Message
	if err != nil {
		return nil, err
	}
	return MessagesToAddress(messages)
}

func MessagesToAddress(messages []*ibfttypes.Message) ([]common.Address, error) {
	var rwdLst []common.Address
	for _, v := range messages {
		log.Info("rwd addr:", v.Address)
		rwdLst = append(rwdLst, v.Address)
	}
	return rwdLst, nil
}
