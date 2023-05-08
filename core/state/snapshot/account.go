// Copyright 2019 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package snapshot

import (
	"bytes"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

// *** modify to support nft transaction 20211217 begin ***

// Account is a modified version of a state.Account, where the root is replaced
// with a byte slice. This format can be used to represent full-consensus format
// or slim-snapshot format which replaces the empty root and code hash as nil
// byte slice.
type Account struct {
	Nonce    uint64
	Balance  *big.Int
	Root     []byte
	CodeHash []byte

	Worm   *types.WormholesExtension `rlp:"nil"`
	Nft    *types.AccountNFT         `rlp:"nil"`
	Staker *types.AccountStaker      `rlp:"nil"`
	Extra  []byte
}

//type WormholesExtension struct {
//	PledgedBalance     *big.Int
//	PledgedBlockNumber *big.Int
//	// *** modify to support nft transaction 20211215 ***
//	//Owner common.Address
//	// whether the account has a NFT exchanger
//	ExchangerFlag    bool
//	BlockNumber      *big.Int
//	ExchangerBalance *big.Int
//	VoteBlockNumber  *big.Int
//	VoteWeight       *big.Int
//	Coefficient      uint8
//	// The ratio that exchanger get.
//	FeeRate       uint16
//	ExchangerName string
//	ExchangerURL  string
//	// ApproveAddress have the right to handle all nfts of the account
//	ApproveAddressList []common.Address
//	// NFTBalance is the nft number that the account have
//	//NFTBalance uint64
//	// Indicates the reward method chosen by the miner
//	//RewardFlag uint8 // 0:SNFT 1:ERB default:1
//}
//
//type AccountNFT struct {
//	//Account
//	Name   string
//	Symbol string
//	//Price                 *big.Int
//	//Direction             uint8 // 0:un_tx,1:buy,2:sell
//	Owner                 common.Address
//	NFTApproveAddressList common.Address
//	//Auctions map[string][]common.Address
//	// MergeLevel is the level of NFT merged
//	MergeLevel  uint8
//	MergeNumber uint32
//	//PledgedFlag           bool
//	//NFTPledgedBlockNumber *big.Int
//
//	Creator   common.Address
//	Royalty   uint16
//	Exchanger common.Address
//	MetaURL   string
//}

// SlimAccount converts a state.Account content into a slim snapshot account
func SlimAccount(nonce uint64,
	balance *big.Int,
	root common.Hash,
	codehash []byte,
	worm *types.WormholesExtension,
	nft *types.AccountNFT,
	staker *types.AccountStaker,
	extra []byte) Account {
	//func SlimAccount(nonce uint64, balance *big.Int, root common.Hash, codehash []byte) Account {
	slim := Account{
		Nonce:   nonce,
		Balance: balance,
	}
	if root != emptyRoot {
		slim.Root = root[:]
	}
	if !bytes.Equal(codehash, emptyCode[:]) {
		slim.CodeHash = codehash
	}

	if worm != nil {
		slim.Worm = worm
	}
	if nft != nil {
		slim.Nft = nft
	}

	if staker != nil {
		slim.Staker = staker
	}

	slim.Extra = extra[:]

	return slim
}

// *** modify to support nft transaction 20211217 begin ***

// SlimAccountRLP converts a state.Account content into a slim snapshot
// version RLP encoded.
func SlimAccountRLP(nonce uint64,
	balance *big.Int,
	root common.Hash,
	codehash []byte,
	worm *types.WormholesExtension,
	nft *types.AccountNFT,
	staker *types.AccountStaker,
	extra []byte) []byte {
	data, err := rlp.EncodeToBytes(SlimAccount(nonce,
		balance,
		root,
		codehash,
		worm,
		nft,
		staker,
		extra))
	//func SlimAccountRLP(nonce uint64, balance *big.Int, root common.Hash, codehash []byte) []byte {
	//	data, err := rlp.EncodeToBytes(SlimAccount(nonce, balance, root, codehash))
	// *** modify to support nft transaction 20211217 end ***
	if err != nil {
		panic(err)
	}
	return data
}

// FullAccount decodes the data on the 'slim RLP' format and return
// the consensus format account.
func FullAccount(data []byte) (Account, error) {
	var account Account
	if err := rlp.DecodeBytes(data, &account); err != nil {
		return Account{}, err
	}
	if len(account.Root) == 0 {
		account.Root = emptyRoot[:]
	}
	if len(account.CodeHash) == 0 {
		account.CodeHash = emptyCode[:]
	}
	return account, nil
}

// FullAccountRLP converts data on the 'slim RLP' format into the full RLP-format.
func FullAccountRLP(data []byte) ([]byte, error) {
	account, err := FullAccount(data)
	if err != nil {
		return nil, err
	}
	return rlp.EncodeToBytes(account)
}
