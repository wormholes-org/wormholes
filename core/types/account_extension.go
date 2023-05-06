package types

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type WormholesExtension struct {
	PledgedBalance     *big.Int
	PledgedBlockNumber *big.Int
	// *** modify to support nft transaction 20211215 ***
	//Owner common.Address
	// whether the account has a NFT exchanger
	ExchangerFlag    bool
	BlockNumber      *big.Int
	ExchangerBalance *big.Int
	VoteBlockNumber  *big.Int
	VoteWeight       *big.Int
	Coefficient      uint8
	// The ratio that exchanger get.
	FeeRate       uint16
	ExchangerName string
	ExchangerURL  string
	// ApproveAddress have the right to handle all nfts of the account
	ApproveAddressList []common.Address
	// NFTBalance is the nft number that the account have
	//NFTBalance uint64
	// Indicates the reward method chosen by the miner
	//RewardFlag uint8 // 0:SNFT 1:ERB default:1
}

type AccountNFT struct {
	//Account
	Name   string
	Symbol string
	//Price                 *big.Int
	//Direction             uint8 // 0:no_tx,1:by,2:sell
	Owner                 common.Address
	NFTApproveAddressList common.Address
	//Auctions map[string][]common.Address
	// MergeLevel is the level of NFT merged
	MergeLevel  uint8
	MergeNumber uint32
	//PledgedFlag           bool
	//NFTPledgedBlockNumber *big.Int

	Creator   common.Address
	Royalty   uint16
	Exchanger common.Address
	MetaURL   string
}
