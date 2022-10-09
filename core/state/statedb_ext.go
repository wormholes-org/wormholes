package state

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
)

func (s *StateDB) CreateRewardValidators(validators []common.Address, blocknumber *big.Int) {
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

func (s *StateDB) CreateRewardExchanges(exchangers []common.Address, blocknumber *big.Int) {
	log.Info("CreateNFTByOfficial16", "exchangers len=", len(exchangers), "blocknumber=", blocknumber.Uint64())
	for _, addr := range exchangers {
		log.Info("CreateNFTByOfficial16", "exchangers=", addr.Hex(), "blocknumber=", blocknumber.Uint64())
	}
	for _, owner := range exchangers {
		nftAddr := common.Address{}
		var metaUrl string
		var royalty uint32
		var creator string
		//nftAddr, info, ok := s.SNFTExchangePool.PopAddress(blocknumber)
		//if !ok {
		nftAddr = common.BytesToAddress(s.MintDeep.OfficialMint.Bytes())
		injectedInfo := s.OfficialNFTPool.GetInjectedInfo(nftAddr)
		if injectedInfo == nil {
			return
		}
		metaUrl = injectedInfo.Dir + "/" + nftAddr.String()
		royalty = injectedInfo.Royalty
		creator = injectedInfo.Creator
		//} else {
		//	metaUrl = info.MetalUrl + "/" + nftAddr.String()
		//	royalty = info.Royalty
		//	creator = info.Creator
		//}
		log.Info("CreateNFTByOfficial16()", "--nftAddr=", nftAddr.String(), "blocknumber=", blocknumber.Uint64())

		s.CreateAccount(nftAddr)
		stateObject := s.GetOrNewStateObject(nftAddr)
		if stateObject != nil {
			stateObject.SetNFTInfo(
				"",
				"",
				//big.NewInt(0),
				//0,
				owner,
				common.Address{},
				0,
				1,
				false,
				big.NewInt(0),
				common.HexToAddress(creator),
				royalty,
				common.Address{},
				metaUrl)
			s.MergeNFT16(nftAddr)
			//if !ok {
			s.OfficialNFTPool.DeleteExpireElem(s.MintDeep.OfficialMint)
			s.MintDeep.OfficialMint.Add(s.MintDeep.OfficialMint, big.NewInt(1))
			//}
		}
	}

	if s.OfficialNFTPool.RemainderNum(s.MintDeep.OfficialMint) <= 110 {
		s.ElectNominatedOfficialNFT()
	}
}
