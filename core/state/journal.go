// Copyright 2016 The go-ethereum Authors
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

package state

import (
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// journalEntry is a modification entry in the state change journal that can be
// reverted on demand.
type journalEntry interface {
	// revert undoes the changes introduced by this journal entry.
	revert(*StateDB)

	// dirtied returns the Ethereum address modified by this journal entry.
	dirtied() *common.Address
}

// journal contains the list of state modifications applied since the last state
// commit. These are tracked to be able to be reverted in case of an execution
// exception or revertal request.
type journal struct {
	entries []journalEntry         // Current changes tracked by the journal
	dirties map[common.Address]int // Dirty accounts and the number of changes
}

// newJournal create a new initialized journal.
func newJournal() *journal {
	return &journal{
		dirties: make(map[common.Address]int),
	}
}

// append inserts a new modification entry to the end of the change journal.
func (j *journal) append(entry journalEntry) {
	j.entries = append(j.entries, entry)
	if addr := entry.dirtied(); addr != nil {
		j.dirties[*addr]++
	}
}

// revert undoes a batch of journalled modifications along with any reverted
// dirty handling too.
func (j *journal) revert(statedb *StateDB, snapshot int) {
	for i := len(j.entries) - 1; i >= snapshot; i-- {
		// Undo the changes made by the operation
		j.entries[i].revert(statedb)

		// Drop any dirty tracking induced by the change
		if addr := j.entries[i].dirtied(); addr != nil {
			if j.dirties[*addr]--; j.dirties[*addr] == 0 {
				delete(j.dirties, *addr)
			}
		}
	}
	j.entries = j.entries[:snapshot]
}

// dirty explicitly sets an address to dirty, even if the change entries would
// otherwise suggest it as clean. This method is an ugly hack to handle the RIPEMD
// precompile consensus exception.
func (j *journal) dirty(addr common.Address) {
	j.dirties[addr]++
}

// length returns the current number of entries in the journal.
func (j *journal) length() int {
	return len(j.entries)
}

type (
	// Changes to the account trie.
	createObjectChange struct {
		account *common.Address
	}
	resetObjectChange struct {
		prev         *stateObject
		prevdestruct bool
	}
	suicideChange struct {
		account     *common.Address
		prev        bool // whether account had already suicided
		prevbalance *big.Int
	}

	// Changes to individual accounts.
	balanceChange struct {
		account *common.Address
		prev    *big.Int
	}
	nonceChange struct {
		account *common.Address
		prev    uint64
	}
	storageChange struct {
		account       *common.Address
		key, prevalue common.Hash
	}
	codeChange struct {
		account            *common.Address
		prevcode, prevhash []byte
	}

	// Changes to other state values.
	refundChange struct {
		prev uint64
	}
	addLogChange struct {
		txhash common.Hash
	}
	addPreimageChange struct {
		hash common.Hash
	}
	touchChange struct {
		account *common.Address
	}
	// Changes to the access list
	accessListAddAccountChange struct {
		address *common.Address
	}
	accessListAddSlotChange struct {
		address *common.Address
		slot    *common.Hash
	}

	// *** modify to support nft transaction 20211215 begin ***
	// change to the owner
	nftOwnerChange struct {
		nftAddr  *common.Address
		oldOwner common.Address
	}
	// *** modify to support nft transaction 20211215 end ***
	nftApproveAddressChange struct {
		nftAddr               *common.Address
		oldApproveAddressList []common.Address
	}

	nftApproveAddressChangeOne struct {
		nftAddr                  *common.Address
		oldNFTApproveAddressList common.Address
	}

	openExchangerChange struct {
		address          *common.Address
		oldExchangerFlag bool
		oldBlockNumber   *big.Int
		oldFeeRate       uint16
		oldExchangerName string
		oldExchangerURL  string
	}

	nftInfoChange struct {
		address                  *common.Address
		oldName                  string
		oldSymbol                string
		oldOwner                 common.Address
		oldNFTApproveAddressList common.Address
		oldMergeLevel            uint8
		oldMergeNumber           uint32
		//oldPledgedFlag           bool
		//oldNFTPledgedBlockNumber *big.Int
		oldCreator   common.Address
		oldRoyalty   uint16
		oldExchanger common.Address
		oldMetaURL   string
	}

	pledgedBalanceChange struct {
		account *common.Address
		prev    *big.Int
	}

	pledgedBlockNumberChange struct {
		account *common.Address
		prev    *big.Int
	}

	exchangerBalanceChange struct {
		account *common.Address
		prev    *big.Int
	}

	blockNumberChange struct {
		account *common.Address
		prev    *big.Int
	}

	voteBlockNumberChange struct {
		account *common.Address
		prev    *big.Int
	}

	voteWeightChange struct {
		account *common.Address
		prev    *big.Int
	}

	//pledgedNFTInfo struct {
	//	account               *common.Address
	//	pledgedFlag           bool
	//	nftPledgedBlockNumber *big.Int
	//}

	//RewardFlagChange struct {
	//	account    *common.Address
	//	rewardFlag uint8
	//}
	coefficientChange struct {
		account *common.Address
		prev    uint8
	}

	extraChange struct {
		account *common.Address
		prev    []byte
	}

	userMintChange struct {
		account *common.Address
		prev    *big.Int
	}

	officialMintChange struct {
		account *common.Address
		prev    *big.Int
	}

	validatorsChange struct {
		account       *common.Address
		oldValidators types.ValidatorList
	}

	stakersChange struct {
		account    *common.Address
		oldStakers types.StakerList
	}

	snftsChange struct {
		account  *common.Address
		oldSnfts types.InjectedOfficialNFTList
	}

	nomineeChange struct {
		account    *common.Address
		oldNominee types.NominatedOfficialNFT
	}
)

func (ch createObjectChange) revert(s *StateDB) {
	delete(s.stateObjects, *ch.account)
	delete(s.stateObjectsDirty, *ch.account)
}

func (ch createObjectChange) dirtied() *common.Address {
	return ch.account
}

func (ch resetObjectChange) revert(s *StateDB) {
	s.setStateObject(ch.prev)
	if !ch.prevdestruct && s.snap != nil {
		delete(s.snapDestructs, ch.prev.addrHash)
	}
}

func (ch resetObjectChange) dirtied() *common.Address {
	return nil
}

func (ch suicideChange) revert(s *StateDB) {
	obj := s.getStateObject(*ch.account)
	if obj != nil {
		obj.suicided = ch.prev
		obj.setBalance(ch.prevbalance)
	}
}

func (ch suicideChange) dirtied() *common.Address {
	return ch.account
}

var ripemd = common.HexToAddress("0000000000000000000000000000000000000003")

func (ch touchChange) revert(s *StateDB) {
}

func (ch touchChange) dirtied() *common.Address {
	return ch.account
}

func (ch balanceChange) revert(s *StateDB) {
	s.getStateObject(*ch.account).setBalance(ch.prev)
}

func (ch balanceChange) dirtied() *common.Address {
	return ch.account
}

func (ch nonceChange) revert(s *StateDB) {
	s.getStateObject(*ch.account).setNonce(ch.prev)
}

func (ch nonceChange) dirtied() *common.Address {
	return ch.account
}

func (ch codeChange) revert(s *StateDB) {
	s.getStateObject(*ch.account).setCode(common.BytesToHash(ch.prevhash), ch.prevcode)
}

func (ch codeChange) dirtied() *common.Address {
	return ch.account
}

func (ch storageChange) revert(s *StateDB) {
	s.getStateObject(*ch.account).setState(ch.key, ch.prevalue)
}

func (ch storageChange) dirtied() *common.Address {
	return ch.account
}

func (ch refundChange) revert(s *StateDB) {
	s.refund = ch.prev
}

func (ch refundChange) dirtied() *common.Address {
	return nil
}

func (ch addLogChange) revert(s *StateDB) {
	logs := s.logs[ch.txhash]
	if len(logs) == 1 {
		delete(s.logs, ch.txhash)
	} else {
		s.logs[ch.txhash] = logs[:len(logs)-1]
	}
	s.logSize--
}

func (ch addLogChange) dirtied() *common.Address {
	return nil
}

func (ch addPreimageChange) revert(s *StateDB) {
	delete(s.preimages, ch.hash)
}

func (ch addPreimageChange) dirtied() *common.Address {
	return nil
}

func (ch accessListAddAccountChange) revert(s *StateDB) {
	/*
		One important invariant here, is that whenever a (addr, slot) is added, if the
		addr is not already present, the add causes two journal entries:
		- one for the address,
		- one for the (address,slot)
		Therefore, when unrolling the change, we can always blindly delete the
		(addr) at this point, since no storage adds can remain when come upon
		a single (addr) change.
	*/
	s.accessList.DeleteAddress(*ch.address)
}

func (ch accessListAddAccountChange) dirtied() *common.Address {
	return nil
}

func (ch accessListAddSlotChange) revert(s *StateDB) {
	s.accessList.DeleteSlot(*ch.address, *ch.slot)
}

func (ch accessListAddSlotChange) dirtied() *common.Address {
	return nil
}

// *** modify to support nft transaction 20211215 begin ***

func (ch nftOwnerChange) revert(s *StateDB) {
	s.getStateObject(*ch.nftAddr).setOwner(ch.oldOwner)
}

func (ch nftOwnerChange) dirtied() *common.Address {
	return ch.nftAddr
}

// *** modify to support nft transaction 20211215 end ***
func (ch nftApproveAddressChange) revert(s *StateDB) {
	s.getStateObject(*ch.nftAddr).setJournalApproveAddress(ch.oldApproveAddressList)
}

func (ch nftApproveAddressChange) dirtied() *common.Address {
	return ch.nftAddr
}

// *** modify to support nft transaction 20211215 end ***
func (ch nftApproveAddressChangeOne) revert(s *StateDB) {
	s.getStateObject(*ch.nftAddr).setJournalNFTApproveAddress(ch.oldNFTApproveAddressList)
}

func (ch nftApproveAddressChangeOne) dirtied() *common.Address {
	return ch.nftAddr
}

func (ch openExchangerChange) revert(s *StateDB) {
	s.getStateObject(*ch.address).setExchangerInfo(
		ch.oldExchangerFlag,
		ch.oldBlockNumber,
		ch.oldFeeRate,
		ch.oldExchangerName,
		ch.oldExchangerURL)
}

func (ch openExchangerChange) dirtied() *common.Address {
	return ch.address
}

func (ch nftInfoChange) revert(s *StateDB) {
	s.getStateObject(*ch.address).setJournalNFTInfo(
		ch.oldName,
		ch.oldSymbol,
		nil,
		0,
		ch.oldOwner,
		ch.oldNFTApproveAddressList,
		ch.oldMergeLevel,
		ch.oldCreator,
		ch.oldRoyalty,
		ch.oldExchanger,
		ch.oldMetaURL)
}

func (ch nftInfoChange) dirtied() *common.Address {
	return ch.address
}

func (ch pledgedBalanceChange) revert(s *StateDB) {
	s.getStateObject(*ch.account).setPledgedBalance(ch.prev)
}

func (ch pledgedBalanceChange) dirtied() *common.Address {
	return ch.account
}

func (ch pledgedBlockNumberChange) revert(s *StateDB) {
	s.getStateObject(*ch.account).setPledgedBlockNumber(ch.prev)
}

func (ch pledgedBlockNumberChange) dirtied() *common.Address {
	return ch.account
}

func (ch exchangerBalanceChange) revert(s *StateDB) {
	s.getStateObject(*ch.account).setExchangerBalance(ch.prev)
}

func (ch exchangerBalanceChange) dirtied() *common.Address {
	return ch.account
}

func (ch blockNumberChange) revert(s *StateDB) {
	s.getStateObject(*ch.account).setBlockNumber(ch.prev)
}

func (ch blockNumberChange) dirtied() *common.Address {
	return ch.account
}

func (ch voteBlockNumberChange) revert(s *StateDB) {
	s.getStateObject(*ch.account).setVoteBlockNumber(ch.prev)
}

func (ch voteBlockNumberChange) dirtied() *common.Address {
	return ch.account
}

func (ch voteWeightChange) revert(s *StateDB) {
	s.getStateObject(*ch.account).setVoteWeight(ch.prev)
}

func (ch voteWeightChange) dirtied() *common.Address {
	return ch.account
}

//func (r RewardFlagChange) revert(s *StateDB) {
//	s.getStateObject(*r.account).setRewardFlag(r.rewardFlag)
//}
//
//func (r RewardFlagChange) dirtied() *common.Address {
//	return r.account
//}

//func (ch pledgedNFTInfo) revert(s *StateDB) {
//	s.getStateObject(*ch.account).setPledgedNFTInfo(ch.pledgedFlag, ch.nftPledgedBlockNumber)
//}
//
//func (ch pledgedNFTInfo) dirtied() *common.Address {
//	return ch.account
//}

func (ch coefficientChange) revert(s *StateDB) {
	s.getStateObject(*ch.account).setCoefficient(ch.prev)
}

func (ch coefficientChange) dirtied() *common.Address {
	return ch.account
}

func (ch extraChange) revert(s *StateDB) {
	s.getStateObject(*ch.account).setExtra(ch.prev)
}

func (ch extraChange) dirtied() *common.Address {
	return ch.account
}

func (ch userMintChange) revert(s *StateDB) {
	s.getStateObject(*ch.account).setUserMint(ch.prev)
}

func (ch userMintChange) dirtied() *common.Address {
	return ch.account
}

func (ch officialMintChange) revert(s *StateDB) {
	s.getStateObject(*ch.account).setOfficialMint(ch.prev)
}

func (ch officialMintChange) dirtied() *common.Address {
	return ch.account
}

func (ch validatorsChange) revert(s *StateDB) {
	s.getStateObject(*ch.account).setValidators(&ch.oldValidators)
}

func (ch validatorsChange) dirtied() *common.Address {
	return ch.account
}

func (ch stakersChange) revert(s *StateDB) {
	s.getStateObject(*ch.account).setStakers(&ch.oldStakers)
}

func (ch stakersChange) dirtied() *common.Address {
	return ch.account
}

func (ch snftsChange) revert(s *StateDB) {
	s.getStateObject(*ch.account).setSnfts(&ch.oldSnfts)
}

func (ch snftsChange) dirtied() *common.Address {
	return ch.account
}

func (ch nomineeChange) revert(s *StateDB) {
	s.getStateObject(*ch.account).setNominee(&ch.oldNominee)
}

func (ch nomineeChange) dirtied() *common.Address {
	return ch.account
}
