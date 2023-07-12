package types

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

// NFT and sNFT minting sequence storage address
var MintDeepStorageAddress = common.HexToAddress("0x0000000000000000000000000000000000000001")

// validators storage address
var ValidatorStorageAddress = common.HexToAddress("0x0000000000000000000000000000000000000002")

// stakers storage address
var StakerStorageAddress = common.HexToAddress("0x0000000000000000000000000000000000000003")

// Storage address for selected creator information
var SnftInjectedStorageAddress = common.HexToAddress("0x0000000000000000000000000000000000000004")

// Information of the creator with the highest voting weight
var NominatedStorageAddress = common.HexToAddress("0x0000000000000000000000000000000000000005")

// Next dividend amount address
var PreDividendAmountAddress = common.HexToAddress("0x0000000000000000000000000000000000000006")

// Dividend amount generated per block
var DividendAmountEachBlock, _ = new(big.Int).SetString("1800000000000000000", 10)

// Current dividend amount address
var DividendAmountAddress = common.HexToAddress("0x0000000000000000000000000000000000000007")

// Each dividend cycle
var DividendBlockInterval uint64 = 120960 // a week

// Storage address for the snft level 3 address list
var SNFTLevel3AddressList = common.HexToAddress("0x0000000000000000000000000000000000000008")

// Storage address for the information of the address list of SNFT level 3 participants in the current dividend
var DividendAddressList = common.HexToAddress("0x0000000000000000000000000000000000000009")

// Voting contract address
var VoteContractAddress = common.HexToAddress("0x0000000000000000000000000000000000000010")

// The amount of voting contract generated per block
var VoteAmountEachBlock, _ = new(big.Int).SetString("800000000000000000", 10)

// validator reward 0.16 ERB
var DREBlockReward = big.NewInt(1.6e+17)

// Deflation rate
var DeflationRate = 0.85

// Deflation time of validator's reward
// reduce 15% block reward in per period
var ReduceRewardPeriod = uint64(365 * 720 * 24)

// Deflation time of staker's reward
var ExchangePeriod = uint64(6160) // 365 * 720 * 24 * 4 / 4096

// snft exchange price
var SNFTL0 = "30000000000000000"
var SNFTL1 = "60000000000000000"
var SNFTL2 = "180000000000000000"
var SNFTL3 = "1000000000000000000"

// Redemption staking time
var CancelDayPledgedInterval int64 = 1 // blockNumber of per hour * 24h
// for test
// var CancelDayPledgedInterval int64 = 5 // blockNumber of per hour * 24h

// number of validators participating in consensus
var ConsensusValidatorsNum = 11

// Number of validators receiving rewards
var ValidatorRewardNum = 7

// The number of stakers receiving SNFT rewards
var StakerRewardNum = 4

// 一期包含的snft碎片数量
// snft版税
// 系统默认的snft的创建者
// The default location for storing metadata of SNFT in the system
var DefaultDir string = "/ipfs/Qmf3xw9rEmsjJdQTV3ZcyF4KfYGtxMkXdNQ8YkVqNmLHY8"

// The number of SNFT fragments included in the first phase
var DefaultNumber uint64 = 4096

// snft royalty
var DefaultRoyalty uint16 = 1000

// default creator
var DefaultCreator string = "0x0000000000000000000000000000000000000000"
