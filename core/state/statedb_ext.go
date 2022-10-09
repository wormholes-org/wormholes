package state

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
)

func (s *StateDB) RewardOnline(validators []common.Address, blocknumber *big.Int) {
	// reward ERB to Validator
	log.Info("validators len=", len(validators), "blocknumber=", blocknumber.Uint64())
	for _, addr := range validators {
		log.Info("validators=", addr.Hex(), "blocknumber=", blocknumber.Uint64())
	}
	rewardAmount := GetRewardAmount(blocknumber.Uint64(), DREBlockReward)
	for _, owner := range validators {
		ownerObject := s.GetOrNewStateObject(owner)
		if ownerObject != nil {
			log.Info("ownerobj", "addr", ownerObject.address.Hex(), "blocknumber=", blocknumber.Uint64())
			ownerObject.AddBalance(rewardAmount)
		}
	}
}
