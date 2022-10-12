package testutils

import (
	"bytes"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	istanbulcommon "github.com/ethereum/go-ethereum/consensus/istanbul/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
)

func Genesis(validators []common.Address, isQBFT bool) *core.Genesis {
	addr1 := common.Address{}
	alloc := core.GenesisAlloc{
		addr1: {Balance: big.NewInt(1000000)},
	}

	// generate genesis block
	genesis := core.DefaultGenesisBlock()
	genesis.Config = params.TestChainConfig
	// force enable Istanbul engine
	genesis.Config.Istanbul = &params.IstanbulConfig{}
	genesis.Config.Ethash = nil
	genesis.Difficulty = istanbulcommon.DefaultDifficulty
	genesis.Nonce = istanbulcommon.EmptyBlockNonce.Uint64()
	genesis.Mixhash = types.IstanbulDigest
	genesis.StartIndex = big.NewInt(0)
	genesis.Alloc = alloc
	genesis.Validator = alloc

	if isQBFT {
		appendValidators(genesis, validators)
	} else {
		appendValidatorsIstanbulExtra(genesis, validators)
	}

	return genesis
}

func GenesisAndKeys(n int, isQBFT bool) (*core.Genesis, []*ecdsa.PrivateKey) {
	// Setup validators
	var nodeKeys []*ecdsa.PrivateKey
	var addrs []common.Address

	for i := 0; i < n; i++ {
		// nodeKeys[i], _ = crypto.GenerateKey()
		// addrs[i] = crypto.PubkeyToAddress(nodeKeys[i].PublicKey)
		prikey, _ := crypto.HexToECDSA("f616c4d20311a2e73c67ef334630f834b7fb42304a1d4448fb2058e9940ecc0a")
		nodeKeys = append(nodeKeys, prikey)

		addrs = append(addrs, crypto.PubkeyToAddress(prikey.PublicKey))
	}

	// generate genesis block
	genesis := Genesis(addrs, isQBFT)

	return genesis, nodeKeys
}

func appendValidatorsIstanbulExtra(genesis *core.Genesis, addrs []common.Address) {
	if len(genesis.ExtraData) < types.IstanbulExtraVanity {
		genesis.ExtraData = append(genesis.ExtraData, bytes.Repeat([]byte{0x00}, types.IstanbulExtraVanity)...)
	}
	genesis.ExtraData = genesis.ExtraData[:types.IstanbulExtraVanity]

	ist := &types.IstanbulExtra{
		Validators:    addrs,
		Seal:          []byte{},
		CommittedSeal: [][]byte{},
	}

	istPayload, err := rlp.EncodeToBytes(&ist)
	if err != nil {
		panic("failed to encode istanbul extra")
	}
	genesis.ExtraData = append(genesis.ExtraData, istPayload...)
}

func appendValidators(genesis *core.Genesis, addrs []common.Address) {
	vanity := append(genesis.ExtraData, bytes.Repeat([]byte{0x00}, types.IstanbulExtraVanity-len(genesis.ExtraData))...)
	ist := &types.QBFTExtra{
		VanityData:    vanity,
		Validators:    addrs,
		Vote:          nil,
		CommittedSeal: [][]byte{},
		Round:         0,
	}

	istPayload, err := rlp.EncodeToBytes(&ist)
	if err != nil {
		panic("failed to encode istanbul extra")
	}
	genesis.ExtraData = istPayload
}
