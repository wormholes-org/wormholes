// Copyright 2014 The go-ethereum Authors
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

// Package state provides a caching layer atop the Ethereum state trie.
package state

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	gomath "math"
	"math/big"
	"sort"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state/snapshot"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

type revision struct {
	id           int
	journalIndex int
}

var (
	// emptyRoot is the known root hash of an empty trie.
	emptyRoot      = common.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
	DREBlockReward = big.NewInt(1.1e+17)
	// reduce 12% block reward in per period
	ReduceRewardPeriod = uint64(365 * 720 * 24)
	ExchangePeriod     = uint64(6160) // 365 * 720 * 24 * 4 / 4096
	//for test
	//ReduceRewardPeriod = uint64(1024)
	//ExchangePeriod     = uint64(1) // 365 * 720 * 24 * 4 / 4096
)

type proofList [][]byte

func (n *proofList) Put(key []byte, value []byte) error {
	*n = append(*n, value)
	return nil
}

func (n *proofList) Delete(key []byte) error {
	panic("not supported")
}

// StateDB structs within the ethereum protocol are used to store anything
// within the merkle trie. StateDBs take care of caching and storing
// nested states. It's the general query interface to retrieve:
// * Contracts
// * Accounts
type StateDB struct {
	db           Database
	prefetcher   *triePrefetcher
	originalRoot common.Hash // The pre-state root, before any changes were made
	trie         Trie
	hasher       crypto.KeccakState

	snaps         *snapshot.Tree
	snap          snapshot.Snapshot
	snapDestructs map[common.Hash]struct{}
	snapAccounts  map[common.Hash][]byte
	snapStorage   map[common.Hash]map[common.Hash][]byte

	// This map holds 'live' objects, which will get modified while processing a state transition.
	stateObjects        map[common.Address]*stateObject
	stateObjectsPending map[common.Address]struct{} // State objects finalized but not yet written to the trie
	stateObjectsDirty   map[common.Address]struct{} // State objects modified in the current execution

	// DB error.
	// State objects are used by the consensus core and VM which are
	// unable to deal with database-level errors. Any error that occurs
	// during a database read is memoized here and will eventually be returned
	// by StateDB.Commit.
	dbErr error

	// The refund counter, also used by state transitioning.
	refund uint64

	thash   common.Hash
	txIndex int
	logs    map[common.Hash][]*types.Log
	logSize uint

	preimages map[common.Hash][]byte

	// Per-transaction access list
	accessList *accessList

	// Journal of state modifications. This is the backbone of
	// Snapshot and RevertToSnapshot.
	journal        *journal
	validRevisions []revision
	nextRevisionId int

	// Measurements gathered during execution for debugging purposes
	AccountReads         time.Duration
	AccountHashes        time.Duration
	AccountUpdates       time.Duration
	AccountCommits       time.Duration
	StorageReads         time.Duration
	StorageHashes        time.Duration
	StorageUpdates       time.Duration
	StorageCommits       time.Duration
	SnapshotAccountReads time.Duration
	SnapshotStorageReads time.Duration
	SnapshotCommits      time.Duration

	//deep for mint NFT
	MintDeep *types.MintDeep
	//SNFT exchange pool
	//SNFTExchangePool     *types.SNFTExchangeList
	PledgedTokenPool     []*types.PledgedToken
	ExchangerTokenPool   []*types.PledgedToken
	OfficialNFTPool      *types.InjectedOfficialNFTList
	NominatedOfficialNFT *types.NominatedOfficialNFT

	ValidatorPool []*types.Validator
}

// New creates a new state from a given trie.
func New(root common.Hash, db Database, snaps *snapshot.Tree) (*StateDB, error) {
	tr, err := db.OpenTrie(root)
	if err != nil {
		return nil, err
	}
	sdb := &StateDB{
		db:                  db,
		trie:                tr,
		originalRoot:        root,
		snaps:               snaps,
		stateObjects:        make(map[common.Address]*stateObject),
		stateObjectsPending: make(map[common.Address]struct{}),
		stateObjectsDirty:   make(map[common.Address]struct{}),
		logs:                make(map[common.Hash][]*types.Log),
		preimages:           make(map[common.Hash][]byte),
		journal:             newJournal(),
		accessList:          newAccessList(),
		hasher:              crypto.NewKeccakState(),
	}
	if sdb.snaps != nil {
		if sdb.snap = sdb.snaps.Snapshot(root); sdb.snap != nil {
			sdb.snapDestructs = make(map[common.Hash]struct{})
			sdb.snapAccounts = make(map[common.Hash][]byte)
			sdb.snapStorage = make(map[common.Hash]map[common.Hash][]byte)
		}
	}
	return sdb, nil
}

// StartPrefetcher initializes a new trie prefetcher to pull in nodes from the
// state trie concurrently while the state is mutated so that when we reach the
// commit phase, most of the needed data is already hot.
func (s *StateDB) StartPrefetcher(namespace string) {
	if s.prefetcher != nil {
		s.prefetcher.close()
		s.prefetcher = nil
	}
	if s.snap != nil {
		s.prefetcher = newTriePrefetcher(s.db, s.originalRoot, namespace)
	}
}

// StopPrefetcher terminates a running prefetcher and reports any leftover stats
// from the gathered metrics.
func (s *StateDB) StopPrefetcher() {
	if s.prefetcher != nil {
		s.prefetcher.close()
		s.prefetcher = nil
	}
}

// setError remembers the first non-nil error it is called with.
func (s *StateDB) setError(err error) {
	if s.dbErr == nil {
		s.dbErr = err
	}
}

func (s *StateDB) Error() error {
	return s.dbErr
}

func (s *StateDB) AddLog(log *types.Log) {
	s.journal.append(addLogChange{txhash: s.thash})

	log.TxHash = s.thash
	log.TxIndex = uint(s.txIndex)
	log.Index = s.logSize
	s.logs[s.thash] = append(s.logs[s.thash], log)
	s.logSize++
}

func (s *StateDB) GetLogs(hash common.Hash, blockHash common.Hash) []*types.Log {
	logs := s.logs[hash]
	for _, l := range logs {
		l.BlockHash = blockHash
	}
	return logs
}

func (s *StateDB) Logs() []*types.Log {
	var logs []*types.Log
	for _, lgs := range s.logs {
		logs = append(logs, lgs...)
	}
	return logs
}

// AddPreimage records a SHA3 preimage seen by the VM.
func (s *StateDB) AddPreimage(hash common.Hash, preimage []byte) {
	if _, ok := s.preimages[hash]; !ok {
		s.journal.append(addPreimageChange{hash: hash})
		pi := make([]byte, len(preimage))
		copy(pi, preimage)
		s.preimages[hash] = pi
	}
}

// Preimages returns a list of SHA3 preimages that have been submitted.
func (s *StateDB) Preimages() map[common.Hash][]byte {
	return s.preimages
}

// AddRefund adds gas to the refund counter
func (s *StateDB) AddRefund(gas uint64) {
	s.journal.append(refundChange{prev: s.refund})
	s.refund += gas
}

// SubRefund removes gas from the refund counter.
// This method will panic if the refund counter goes below zero
func (s *StateDB) SubRefund(gas uint64) {
	s.journal.append(refundChange{prev: s.refund})
	if gas > s.refund {
		panic(fmt.Sprintf("Refund counter below zero (gas: %d > refund: %d)", gas, s.refund))
	}
	s.refund -= gas
}

// Exist reports whether the given account address exists in the state.
// Notably this also returns true for suicided accounts.
func (s *StateDB) Exist(addr common.Address) bool {
	return s.getStateObject(addr) != nil
}

// Empty returns whether the state object is either non-existent
// or empty according to the EIP161 specification (balance = nonce = code = 0)
func (s *StateDB) Empty(addr common.Address) bool {
	so := s.getStateObject(addr)
	return so == nil || so.empty()
}

// GetBalance retrieves the balance from the given address or 0 if object not found
func (s *StateDB) GetBalance(addr common.Address) *big.Int {
	stateObject := s.getStateObject(addr)
	if stateObject != nil {
		return stateObject.Balance()
	}
	return common.Big0
}

func (s *StateDB) GetNonce(addr common.Address) uint64 {
	stateObject := s.getStateObject(addr)
	if stateObject != nil {
		return stateObject.Nonce()
	}

	return 0
}

// TxIndex returns the current transaction index set by Prepare.
func (s *StateDB) TxIndex() int {
	return s.txIndex
}

func (s *StateDB) GetCode(addr common.Address) []byte {
	stateObject := s.getStateObject(addr)
	if stateObject != nil {
		return stateObject.Code(s.db)
	}
	return nil
}

func (s *StateDB) GetCodeSize(addr common.Address) int {
	stateObject := s.getStateObject(addr)
	if stateObject != nil {
		return stateObject.CodeSize(s.db)
	}
	return 0
}

func (s *StateDB) GetCodeHash(addr common.Address) common.Hash {
	stateObject := s.getStateObject(addr)
	if stateObject == nil {
		return common.Hash{}
	}
	return common.BytesToHash(stateObject.CodeHash())
}

// GetState retrieves a value from the given account's storage trie.
func (s *StateDB) GetState(addr common.Address, hash common.Hash) common.Hash {
	stateObject := s.getStateObject(addr)
	if stateObject != nil {
		return stateObject.GetState(s.db, hash)
	}
	return common.Hash{}
}

// GetProof returns the Merkle proof for a given account.
func (s *StateDB) GetProof(addr common.Address) ([][]byte, error) {
	return s.GetProofByHash(crypto.Keccak256Hash(addr.Bytes()))
}

// GetProofByHash returns the Merkle proof for a given account.
func (s *StateDB) GetProofByHash(addrHash common.Hash) ([][]byte, error) {
	var proof proofList
	err := s.trie.Prove(addrHash[:], 0, &proof)
	return proof, err
}

// GetStorageProof returns the Merkle proof for given storage slot.
func (s *StateDB) GetStorageProof(a common.Address, key common.Hash) ([][]byte, error) {
	var proof proofList
	trie := s.StorageTrie(a)
	if trie == nil {
		return proof, errors.New("storage trie for requested address does not exist")
	}
	err := trie.Prove(crypto.Keccak256(key.Bytes()), 0, &proof)
	return proof, err
}

// GetCommittedState retrieves a value from the given account's committed storage trie.
func (s *StateDB) GetCommittedState(addr common.Address, hash common.Hash) common.Hash {
	stateObject := s.getStateObject(addr)
	if stateObject != nil {
		return stateObject.GetCommittedState(s.db, hash)
	}
	return common.Hash{}
}

// Database retrieves the low level database supporting the lower level trie ops.
func (s *StateDB) Database() Database {
	return s.db
}

// StorageTrie returns the storage trie of an account.
// The return value is a copy and is nil for non-existent accounts.
func (s *StateDB) StorageTrie(addr common.Address) Trie {
	stateObject := s.getStateObject(addr)
	if stateObject == nil {
		return nil
	}
	cpy := stateObject.deepCopy(s)
	cpy.updateTrie(s.db)
	return cpy.getTrie(s.db)
}

func (s *StateDB) HasSuicided(addr common.Address) bool {
	stateObject := s.getStateObject(addr)
	if stateObject != nil {
		return stateObject.suicided
	}
	return false
}

/*
 * SETTERS
 */

// AddBalance adds amount to the account associated with addr.
func (s *StateDB) AddBalance(addr common.Address, amount *big.Int) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.AddBalance(amount)
	}
}

// SubBalance subtracts amount from the account associated with addr.
func (s *StateDB) SubBalance(addr common.Address, amount *big.Int) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SubBalance(amount)
	}
}

func (s *StateDB) SetBalance(addr common.Address, amount *big.Int) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SetBalance(amount)
	}
}

func (s *StateDB) SetNonce(addr common.Address, nonce uint64) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SetNonce(nonce)
	}
}

func (s *StateDB) SetCode(addr common.Address, code []byte) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SetCode(crypto.Keccak256Hash(code), code)
	}
}

func (s *StateDB) SetState(addr common.Address, key, value common.Hash) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SetState(s.db, key, value)
	}
}

// SetStorage replaces the entire storage for the specified account with given
// storage. This function should only be used for debugging.
func (s *StateDB) SetStorage(addr common.Address, storage map[common.Hash]common.Hash) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SetStorage(storage)
	}
}

// Suicide marks the given account as suicided.
// This clears the account balance.
//
// The account's state object is still available until the state is committed,
// getStateObject will return a non-nil account after Suicide.
func (s *StateDB) Suicide(addr common.Address) bool {
	stateObject := s.getStateObject(addr)
	if stateObject == nil {
		return false
	}
	s.journal.append(suicideChange{
		account:     &addr,
		prev:        stateObject.suicided,
		prevbalance: new(big.Int).Set(stateObject.Balance()),
	})
	stateObject.markSuicided()
	stateObject.data.Balance = new(big.Int)

	return true
}

//
// Setting, updating & deleting state object methods.
//

// updateStateObject writes the given object to the trie.
func (s *StateDB) updateStateObject(obj *stateObject) {
	// Track the amount of time wasted on updating the account from the trie
	if metrics.EnabledExpensive {
		defer func(start time.Time) { s.AccountUpdates += time.Since(start) }(time.Now())
	}
	// Encode the account and update the account trie
	addr := obj.Address()

	data, err := rlp.EncodeToBytes(obj)

	var tempObj stateObject
	var acc Account
	rlp.DecodeBytes(data, &tempObj)
	rlp.DecodeBytes(data, &acc)

	if err != nil {
		panic(fmt.Errorf("can't encode object at %x: %v", addr[:], err))
	}
	if err = s.trie.TryUpdate(addr[:], data); err != nil {
		s.setError(fmt.Errorf("updateStateObject (%x) error: %v", addr[:], err))
	}

	// If state snapshotting is active, cache the data til commit. Note, this
	// update mechanism is not symmetric to the deletion, because whereas it is
	// enough to track account updates at commit time, deletions need tracking
	// at transaction boundary level to ensure we capture state clearing.
	if s.snap != nil {
		s.snapAccounts[obj.addrHash] = snapshot.SlimAccountRLP(obj.data.Nonce,
			obj.data.Balance,
			obj.data.Root,
			obj.data.CodeHash,
			obj.data.PledgedBalance,
			obj.data.PledgedBlockNumber,
			obj.data.ExchangerFlag,
			obj.data.BlockNumber,
			obj.data.ExchangerBalance,
			obj.data.VoteBlockNumber,
			obj.data.VoteWeight,
			obj.data.Coefficient,
			obj.data.FeeRate,
			obj.data.ExchangerName,
			obj.data.ExchangerURL,
			obj.data.ApproveAddressList,
			//obj.data.NFTBalance,
			obj.data.Extra,
			obj.data.Name,
			obj.data.Symbol,
			//obj.data.Price,
			//obj.data.Direction,
			obj.data.Owner,
			obj.data.NFTApproveAddressList,
			obj.data.MergeLevel,
			obj.data.MergeNumber,
			//obj.data.PledgedFlag,
			//obj.data.NFTPledgedBlockNumber,
			obj.data.Creator,
			obj.data.Royalty,
			obj.data.Exchanger,
			obj.data.MetaURL)
		//s.snapAccounts[obj.addrHash] = snapshot.SlimAccountRLP(obj.data.Nonce, obj.data.Balance, obj.data.Root, obj.data.CodeHash)
	}
}

// deleteStateObject removes the given object from the state trie.
func (s *StateDB) deleteStateObject(obj *stateObject) {
	// Track the amount of time wasted on deleting the account from the trie
	if metrics.EnabledExpensive {
		defer func(start time.Time) { s.AccountUpdates += time.Since(start) }(time.Now())
	}
	// Delete the account from the trie
	addr := obj.Address()
	if err := s.trie.TryDelete(addr[:]); err != nil {
		s.setError(fmt.Errorf("deleteStateObject (%x) error: %v", addr[:], err))
	}
}

// getStateObject retrieves a state object given by the address, returning nil if
// the object is not found or was deleted in this execution context. If you need
// to differentiate between non-existent/just-deleted, use getDeletedStateObject.
func (s *StateDB) getStateObject(addr common.Address) *stateObject {
	if obj := s.getDeletedStateObject(addr); obj != nil && !obj.deleted {
		return obj
	}
	return nil
}

// getDeletedStateObject is similar to getStateObject, but instead of returning
// nil for a deleted state object, it returns the actual object with the deleted
// flag set. This is needed by the state journal to revert to the correct s-
// destructed object instead of wiping all knowledge about the state object.
func (s *StateDB) getDeletedStateObject(addr common.Address) *stateObject {
	// Prefer live objects if any is available
	if obj := s.stateObjects[addr]; obj != nil {
		return obj
	}
	// If no live objects are available, attempt to use snapshots
	var (
		data *Account
		err  error
	)
	if s.snap != nil {
		if metrics.EnabledExpensive {
			defer func(start time.Time) { s.SnapshotAccountReads += time.Since(start) }(time.Now())
		}
		var acc *snapshot.Account
		if acc, err = s.snap.Account(crypto.HashData(s.hasher, addr.Bytes())); err == nil {
			if acc == nil {
				return nil
			}
			data = &Account{
				Nonce:              acc.Nonce,
				Balance:            acc.Balance,
				CodeHash:           acc.CodeHash,
				Root:               common.BytesToHash(acc.Root),
				PledgedBalance:     acc.PledgedBalance,
				PledgedBlockNumber: acc.PledgedBlockNumber,
				ExchangerFlag:      acc.ExchangerFlag,
				BlockNumber:        acc.BlockNumber,
				ExchangerBalance:   acc.ExchangerBalance,
				VoteBlockNumber:    acc.VoteBlockNumber,
				VoteWeight:         acc.VoteWeight,
				Coefficient:        acc.Coefficient,
				FeeRate:            acc.FeeRate,
				ExchangerName:      acc.ExchangerName,
				ExchangerURL:       acc.ExchangerURL,
				//NFTBalance:         acc.NFTBalance,
				Extra: acc.Extra,
				//RewardFlag:         acc.RewardFlag,
				// *** modify to support nft transaction 20211217 begin ***
				AccountNFT: AccountNFT{
					Name:   acc.Name,
					Symbol: acc.Symbol,
					//Price:      acc.Price,
					//Direction:  acc.Direction,
					Owner:       acc.Owner,
					MergeLevel:  acc.MergeLevel,
					MergeNumber: acc.MergeNumber,
					//PledgedFlag:           acc.PledgedFlag,
					//NFTPledgedBlockNumber: acc.NFTPledgedBlockNumber,
					Creator:   acc.Creator,
					Royalty:   acc.Royalty,
					Exchanger: acc.Exchanger,
					MetaURL:   acc.MetaURL,
				},
				// *** modify to support nft transaction 20211217 end ***
			}
			data.ApproveAddressList = append(data.ApproveAddressList, acc.ApproveAddressList...)
			//data.NFTApproveAddressList = append(data.NFTApproveAddressList, acc.NFTApproveAddressList...)
			data.NFTApproveAddressList = acc.NFTApproveAddressList
			if len(data.CodeHash) == 0 {
				data.CodeHash = emptyCodeHash
			}
			if data.Root == (common.Hash{}) {
				data.Root = emptyRoot
			}
		}
	}
	// If snapshot unavailable or reading from it failed, load from the database
	if s.snap == nil || err != nil {
		if metrics.EnabledExpensive {
			defer func(start time.Time) { s.AccountReads += time.Since(start) }(time.Now())
		}
		enc, err := s.trie.TryGet(addr.Bytes())
		if err != nil {
			s.setError(fmt.Errorf("getDeleteStateObject (%x) error: %v", addr.Bytes(), err))
			return nil
		}
		if len(enc) == 0 {
			return nil
		}
		data = new(Account)
		if err := rlp.DecodeBytes(enc, data); err != nil {
			log.Error("Failed to decode state object", "addr", addr, "err", err)
			return nil
		}
	}
	// Insert into the live set
	obj := newObject(s, addr, *data)
	s.setStateObject(obj)
	return obj
}

//for test
func (s *StateDB) getDeletedStateObject2(addr common.Address) *stateObject {
	// Prefer live objects if any is available
	//if obj := s.stateObjects[addr]; obj != nil {
	//	return obj
	//}
	// If no live objects are available, attempt to use snapshots
	var (
		data *Account
		err  error
	)
	//if s.snap != nil {
	//	if metrics.EnabledExpensive {
	//		defer func(start time.Time) { s.SnapshotAccountReads += time.Since(start) }(time.Now())
	//	}
	//	var acc *snapshot.Account
	//	if acc, err = s.snap.Account(crypto.HashData(s.hasher, addr.Bytes())); err == nil {
	//		if acc == nil {
	//			return nil
	//		}
	//		data = &Account{
	//			Nonce:    acc.Nonce,
	//			Balance:  acc.Balance,
	//			CodeHash: acc.CodeHash,
	//			Root:     common.BytesToHash(acc.Root),
	//			// *** modify to support nft transaction 20211217 begin ***
	//			Owner: acc.Owner,
	//			// *** modify to support nft transaction 20211217 end ***
	//		}
	//		if len(data.CodeHash) == 0 {
	//			data.CodeHash = emptyCodeHash
	//		}
	//		if data.Root == (common.Hash{}) {
	//			data.Root = emptyRoot
	//		}
	//	}
	//}
	// If snapshot unavailable or reading from it failed, load from the database
	//if s.snap == nil || err != nil {
	if metrics.EnabledExpensive {
		defer func(start time.Time) { s.AccountReads += time.Since(start) }(time.Now())
	}
	enc, err := s.trie.TryGet(addr.Bytes())
	if err != nil {
		s.setError(fmt.Errorf("getDeleteStateObject (%x) error: %v", addr.Bytes(), err))
		return nil
	}
	if len(enc) == 0 {
		return nil
	}
	data = new(Account)
	if err := rlp.DecodeBytes(enc, data); err != nil {
		log.Error("Failed to decode state object", "addr", addr, "err", err)
		return nil
	}
	//}
	// Insert into the live set
	obj := newObject(s, addr, *data)
	s.setStateObject(obj)
	return obj
}

func (s *StateDB) setStateObject(object *stateObject) {
	s.stateObjects[object.Address()] = object
}

// GetOrNewStateObject retrieves a state object or create a new state object if nil.
func (s *StateDB) GetOrNewStateObject(addr common.Address) *stateObject {
	stateObject := s.getStateObject(addr)
	if stateObject == nil {
		stateObject, _ = s.createObject(addr)
	}
	if stateObject.data.PledgedBlockNumber == nil {
		stateObject.data.PledgedBlockNumber = big.NewInt(0)
	}
	if stateObject.data.PledgedBalance == nil {
		stateObject.data.PledgedBalance = big.NewInt(0)
	}
	if stateObject.data.BlockNumber == nil {
		stateObject.data.BlockNumber = big.NewInt(0)
	}
	if stateObject.data.ExchangerBalance == nil {
		stateObject.data.ExchangerBalance = big.NewInt(0)
	}
	if stateObject.data.VoteBlockNumber == nil {
		stateObject.data.VoteBlockNumber = big.NewInt(0)
	}
	if stateObject.data.VoteWeight == nil {
		stateObject.data.VoteWeight = big.NewInt(0)
	}

	return stateObject
}

// createObject creates a new state object. If there is an existing account with
// the given address, it is overwritten and returned as the second return value.
func (s *StateDB) createObject(addr common.Address) (newobj, prev *stateObject) {
	prev = s.getDeletedStateObject(addr) // Note, prev might have been deleted, we need that!

	var prevdestruct bool
	if s.snap != nil && prev != nil {
		_, prevdestruct = s.snapDestructs[prev.addrHash]
		if !prevdestruct {
			s.snapDestructs[prev.addrHash] = struct{}{}
		}
	}
	//newobj = newObject(s, addr, Account{RewardFlag: 1})
	newobj = newObject(s, addr, Account{})
	if prev == nil {
		s.journal.append(createObjectChange{account: &addr})
	} else {
		s.journal.append(resetObjectChange{prev: prev, prevdestruct: prevdestruct})
	}
	s.setStateObject(newobj)
	if prev != nil && !prev.deleted {
		return newobj, prev
	}
	return newobj, nil
}

// CreateAccount explicitly creates a state object. If a state object with the address
// already exists the balance is carried over to the new account.
//
// CreateAccount is called during the EVM CREATE operation. The situation might arise that
// a contract does the following:
//
//   1. sends funds to sha(account ++ (nonce + 1))
//   2. tx_create(sha(account ++ nonce)) (note that this gets the address of 1)
//
// Carrying over the balance ensures that Ether doesn't disappear.
func (s *StateDB) CreateAccount(addr common.Address) {
	newObj, prev := s.createObject(addr)
	if prev != nil {
		newObj.setBalance(prev.data.Balance)
	}
}

func (db *StateDB) ForEachStorage(addr common.Address, cb func(key, value common.Hash) bool) error {
	so := db.getStateObject(addr)
	if so == nil {
		return nil
	}
	it := trie.NewIterator(so.getTrie(db.db).NodeIterator(nil))

	for it.Next() {
		key := common.BytesToHash(db.trie.GetKey(it.Key))
		if value, dirty := so.dirtyStorage[key]; dirty {
			if !cb(key, value) {
				return nil
			}
			continue
		}

		if len(it.Value) > 0 {
			_, content, _, err := rlp.Split(it.Value)
			if err != nil {
				return err
			}
			if !cb(key, common.BytesToHash(content)) {
				return nil
			}
		}
	}
	return nil
}

// Copy creates a deep, independent copy of the state.
// Snapshots of the copied state cannot be applied to the copy.
func (s *StateDB) Copy() *StateDB {
	// Copy all the basic fields, initialize the memory ones
	state := &StateDB{
		db:                  s.db,
		trie:                s.db.CopyTrie(s.trie),
		stateObjects:        make(map[common.Address]*stateObject, len(s.journal.dirties)),
		stateObjectsPending: make(map[common.Address]struct{}, len(s.stateObjectsPending)),
		stateObjectsDirty:   make(map[common.Address]struct{}, len(s.journal.dirties)),
		refund:              s.refund,
		logs:                make(map[common.Hash][]*types.Log, len(s.logs)),
		logSize:             s.logSize,
		preimages:           make(map[common.Hash][]byte, len(s.preimages)),
		journal:             newJournal(),
		hasher:              crypto.NewKeccakState(),
		MintDeep:            new(types.MintDeep),
		//SNFTExchangePool:     new(types.SNFTExchangeList),
		PledgedTokenPool:     make([]*types.PledgedToken, 0),
		ExchangerTokenPool:   make([]*types.PledgedToken, 0),
		OfficialNFTPool:      new(types.InjectedOfficialNFTList),
		NominatedOfficialNFT: new(types.NominatedOfficialNFT),
	}
	// Copy the dirty states, logs, and preimages
	for addr := range s.journal.dirties {
		// As documented [here](https://github.com/ethereum/go-ethereum/pull/16485#issuecomment-380438527),
		// and in the Finalise-method, there is a case where an object is in the journal but not
		// in the stateObjects: OOG after touch on ripeMD prior to Byzantium. Thus, we need to check for
		// nil
		if object, exist := s.stateObjects[addr]; exist {
			// Even though the original object is dirty, we are not copying the journal,
			// so we need to make sure that anyside effect the journal would have caused
			// during a commit (or similar op) is already applied to the copy.
			state.stateObjects[addr] = object.deepCopy(state)

			state.stateObjectsDirty[addr] = struct{}{}   // Mark the copy dirty to force internal (code/state) commits
			state.stateObjectsPending[addr] = struct{}{} // Mark the copy pending to force external (account) commits
		}
	}
	// Above, we don't copy the actual journal. This means that if the copy is copied, the
	// loop above will be a no-op, since the copy's journal is empty.
	// Thus, here we iterate over stateObjects, to enable copies of copies
	for addr := range s.stateObjectsPending {
		if _, exist := state.stateObjects[addr]; !exist {
			state.stateObjects[addr] = s.stateObjects[addr].deepCopy(state)
		}
		state.stateObjectsPending[addr] = struct{}{}
	}
	for addr := range s.stateObjectsDirty {
		if _, exist := state.stateObjects[addr]; !exist {
			state.stateObjects[addr] = s.stateObjects[addr].deepCopy(state)
		}
		state.stateObjectsDirty[addr] = struct{}{}
	}
	for hash, logs := range s.logs {
		cpy := make([]*types.Log, len(logs))
		for i, l := range logs {
			cpy[i] = new(types.Log)
			*cpy[i] = *l
		}
		state.logs[hash] = cpy
	}
	for hash, preimage := range s.preimages {
		state.preimages[hash] = preimage
	}
	// Do we need to copy the access list? In practice: No. At the start of a
	// transaction, the access list is empty. In practice, we only ever copy state
	// _between_ transactions/blocks, never in the middle of a transaction.
	// However, it doesn't cost us much to copy an empty list, so we do it anyway
	// to not blow up if we ever decide copy it in the middle of a transaction
	state.accessList = s.accessList.Copy()

	// If there's a prefetcher running, make an inactive copy of it that can
	// only access data but does not actively preload (since the user will not
	// know that they need to explicitly terminate an active copy).
	if s.prefetcher != nil {
		state.prefetcher = s.prefetcher.copy()
	}
	if s.snaps != nil {
		// In order for the miner to be able to use and make additions
		// to the snapshot tree, we need to copy that aswell.
		// Otherwise, any block mined by ourselves will cause gaps in the tree,
		// and force the miner to operate trie-backed only
		state.snaps = s.snaps
		state.snap = s.snap
		// deep copy needed
		state.snapDestructs = make(map[common.Hash]struct{})
		for k, v := range s.snapDestructs {
			state.snapDestructs[k] = v
		}
		state.snapAccounts = make(map[common.Hash][]byte)
		for k, v := range s.snapAccounts {
			state.snapAccounts[k] = v
		}
		state.snapStorage = make(map[common.Hash]map[common.Hash][]byte)
		for k, v := range s.snapStorage {
			temp := make(map[common.Hash][]byte)
			for kk, vv := range v {
				temp[kk] = vv
			}
			state.snapStorage[k] = temp
		}
	}

	if s.MintDeep != nil {
		state.MintDeep.UserMint = big.NewInt(0)
		state.MintDeep.OfficialMint = big.NewInt(0)
		if s.MintDeep.UserMint != nil {
			state.MintDeep.UserMint.Set(s.MintDeep.UserMint)
		}
		if s.MintDeep.OfficialMint != nil {
			state.MintDeep.OfficialMint.Set(s.MintDeep.OfficialMint)
		}
	}

	//state.SNFTExchangePool.SNFTExchanges = make([]*types.SNFTExchange, 0)
	//if s.SNFTExchangePool != nil && len(s.SNFTExchangePool.SNFTExchanges) > 0 {
	//	for _, snftExchange := range s.SNFTExchangePool.SNFTExchanges {
	//		var tempSNFTExchange types.SNFTExchange
	//		tempSNFTExchange.NFTAddress = snftExchange.NFTAddress
	//		tempSNFTExchange.MergeLevel = snftExchange.MergeLevel
	//		tempSNFTExchange.CurrentMintAddress = snftExchange.CurrentMintAddress
	//		tempSNFTExchange.BlockNumber = new(big.Int).Set(snftExchange.BlockNumber)
	//		tempSNFTExchange.MetalUrl = snftExchange.MetalUrl
	//		tempSNFTExchange.Royalty = snftExchange.Royalty
	//		tempSNFTExchange.Creator = snftExchange.Creator
	//		state.SNFTExchangePool.SNFTExchanges = append(state.SNFTExchangePool.SNFTExchanges, &tempSNFTExchange)
	//	}
	//}

	state.OfficialNFTPool.InjectedOfficialNFTs = make([]*types.InjectedOfficialNFT, 0)
	if s.OfficialNFTPool != nil && len(s.OfficialNFTPool.InjectedOfficialNFTs) > 0 {
		for _, OfficialNFT := range s.OfficialNFTPool.InjectedOfficialNFTs {
			var tempOfficialNFT types.InjectedOfficialNFT
			tempOfficialNFT.Dir = OfficialNFT.Dir
			tempOfficialNFT.StartIndex = new(big.Int).Set(OfficialNFT.StartIndex)
			tempOfficialNFT.Number = OfficialNFT.Number
			tempOfficialNFT.Royalty = OfficialNFT.Royalty
			tempOfficialNFT.Creator = OfficialNFT.Creator
			tempOfficialNFT.Address = OfficialNFT.Address
			if OfficialNFT.VoteWeight != nil {
				tempOfficialNFT.VoteWeight = new(big.Int).Set(OfficialNFT.VoteWeight)
			}
			state.OfficialNFTPool.InjectedOfficialNFTs = append(state.OfficialNFTPool.InjectedOfficialNFTs, &tempOfficialNFT)
		}
	}

	if s.PledgedTokenPool != nil && len(s.PledgedTokenPool) > 0 {
		for _, v := range s.PledgedTokenPool {
			var pledgedToken types.PledgedToken
			pledgedToken.Address = v.Address
			pledgedToken.Amount = new(big.Int).Set(v.Amount)
			pledgedToken.Flag = v.Flag
			pledgedToken.ProxyAddress = v.ProxyAddress
			state.PledgedTokenPool = append(state.PledgedTokenPool, &pledgedToken)
		}
	}

	if s.ExchangerTokenPool != nil && len(s.ExchangerTokenPool) > 0 {
		for _, v := range s.ExchangerTokenPool {
			var exchangerToken types.PledgedToken
			exchangerToken.Address = v.Address
			exchangerToken.Amount = new(big.Int).Set(v.Amount)
			exchangerToken.Flag = v.Flag
			state.ExchangerTokenPool = append(state.ExchangerTokenPool, &exchangerToken)
		}
	}
	if s.NominatedOfficialNFT != nil {
		state.NominatedOfficialNFT.Dir = s.NominatedOfficialNFT.Dir
		state.NominatedOfficialNFT.StartIndex = new(big.Int).Set(s.NominatedOfficialNFT.StartIndex)
		state.NominatedOfficialNFT.Number = s.NominatedOfficialNFT.Number
		state.NominatedOfficialNFT.Royalty = s.NominatedOfficialNFT.Royalty
		state.NominatedOfficialNFT.Creator = s.NominatedOfficialNFT.Creator
		state.NominatedOfficialNFT.Address = s.NominatedOfficialNFT.Address
	}

	state.ValidatorPool = make([]*types.Validator, 0)
	if s.ValidatorPool != nil && len(s.ValidatorPool) < 0 {
		for _, v := range s.ValidatorPool {
			a := types.Validator{
				Addr:    v.Addr,
				Proxy:   v.Proxy,
				Balance: v.Balance,
			}
			state.ValidatorPool = append(state.ValidatorPool, &a)
		}
	}

	return state
}

// Snapshot returns an identifier for the current revision of the state.
func (s *StateDB) Snapshot() int {
	id := s.nextRevisionId
	s.nextRevisionId++
	s.validRevisions = append(s.validRevisions, revision{id, s.journal.length()})
	return id
}

// RevertToSnapshot reverts all state changes made since the given revision.
func (s *StateDB) RevertToSnapshot(revid int) {
	// Find the snapshot in the stack of valid snapshots.
	idx := sort.Search(len(s.validRevisions), func(i int) bool {
		return s.validRevisions[i].id >= revid
	})
	if idx == len(s.validRevisions) || s.validRevisions[idx].id != revid {
		panic(fmt.Errorf("revision id %v cannot be reverted", revid))
	}
	snapshot := s.validRevisions[idx].journalIndex

	// Replay the journal to undo changes and remove invalidated snapshots
	s.journal.revert(s, snapshot)
	s.validRevisions = s.validRevisions[:idx]
}

// GetRefund returns the current value of the refund counter.
func (s *StateDB) GetRefund() uint64 {
	return s.refund
}

// Finalise finalises the state by removing the s destructed objects and clears
// the journal as well as the refunds. Finalise, however, will not push any updates
// into the tries just yet. Only IntermediateRoot or Commit will do that.
func (s *StateDB) Finalise(deleteEmptyObjects bool) {
	addressesToPrefetch := make([][]byte, 0, len(s.journal.dirties))
	for addr := range s.journal.dirties {
		obj, exist := s.stateObjects[addr]
		if !exist {
			// ripeMD is 'touched' at block 1714175, in tx 0x1237f737031e40bcde4a8b7e717b2d15e3ecadfe49bb1bbc71ee9deb09c6fcf2
			// That tx goes out of gas, and although the notion of 'touched' does not exist there, the
			// touch-event will still be recorded in the journal. Since ripeMD is a special snowflake,
			// it will persist in the journal even though the journal is reverted. In this special circumstance,
			// it may exist in `s.journal.dirties` but not in `s.stateObjects`.
			// Thus, we can safely ignore it here
			continue
		}
		if obj.suicided || (deleteEmptyObjects && obj.empty()) {
			obj.deleted = true

			// If state snapshotting is active, also mark the destruction there.
			// Note, we can't do this only at the end of a block because multiple
			// transactions within the same block might self destruct and then
			// ressurrect an account; but the snapshotter needs both events.
			if s.snap != nil {
				s.snapDestructs[obj.addrHash] = struct{}{} // We need to maintain account deletions explicitly (will remain set indefinitely)
				delete(s.snapAccounts, obj.addrHash)       // Clear out any previously updated account data (may be recreated via a ressurrect)
				delete(s.snapStorage, obj.addrHash)        // Clear out any previously updated storage data (may be recreated via a ressurrect)
			}
		} else {
			obj.finalise(true) // Prefetch slots in the background
		}
		s.stateObjectsPending[addr] = struct{}{}
		s.stateObjectsDirty[addr] = struct{}{}

		// At this point, also ship the address off to the precacher. The precacher
		// will start loading tries, and when the change is eventually committed,
		// the commit-phase will be a lot faster
		addressesToPrefetch = append(addressesToPrefetch, common.CopyBytes(addr[:])) // Copy needed for closure
	}
	if s.prefetcher != nil && len(addressesToPrefetch) > 0 {
		s.prefetcher.prefetch(s.originalRoot, addressesToPrefetch)
	}
	// Invalidate journal because reverting across transactions is not allowed.
	s.clearJournalAndRefund()
}

// IntermediateRoot computes the current root hash of the state trie.
// It is called in between transactions to get the root hash that
// goes into transaction receipts.
func (s *StateDB) IntermediateRoot(deleteEmptyObjects bool) common.Hash {
	// Finalise all the dirty storage states and write them into the tries
	//log.Info("caver|IntermediateRoot|enter=0", "triehash", s.trie.Hash().String())
	s.Finalise(deleteEmptyObjects)
	//log.Info("caver|IntermediateRoot|enter=1", "triehash", s.trie.Hash().String())

	// If there was a trie prefetcher operating, it gets aborted and irrevocably
	// modified after we start retrieving tries. Remove it from the statedb after
	// this round of use.
	//
	// This is weird pre-byzantium since the first tx runs with a prefetcher and
	// the remainder without, but pre-byzantium even the initial prefetcher is
	// useless, so no sleep lost.
	prefetcher := s.prefetcher
	if s.prefetcher != nil {
		defer func() {
			s.prefetcher.close()
			s.prefetcher = nil
		}()
	}
	// Although naively it makes sense to retrieve the account trie and then do
	// the contract storage and account updates sequentially, that short circuits
	// the account prefetcher. Instead, let's process all the storage updates
	// first, giving the account prefeches just a few more milliseconds of time
	// to pull useful data from disk.
	for addr := range s.stateObjectsPending {
		if obj := s.stateObjects[addr]; !obj.deleted {
			obj.updateRoot(s.db)
		}
	}
	// Now we're about to start to write changes to the trie. The trie is so far
	// _untouched_. We can check with the prefetcher, if it can give us a trie
	// which has the same root, but also has some content loaded into it.
	if prefetcher != nil {
		if trie := prefetcher.trie(s.originalRoot); trie != nil {
			s.trie = trie
		}
	}
	usedAddrs := make([][]byte, 0, len(s.stateObjectsPending))
	for addr := range s.stateObjectsPending {
		if obj := s.stateObjects[addr]; obj.deleted {
			s.deleteStateObject(obj)
			//log.Info("caver|IntermediateRoot|deleteStateObject", "addr",addr.String(),"triehash", s.trie.Hash().String(), "benefiaddr", obj.data.Owner.Hex())
		} else {
			s.updateStateObject(obj)
			//log.Info("caver|IntermediateRoot|updateStateObject", "addr",addr.String(),"triehash", s.trie.Hash().String(), "benefiaddr", obj.data.Owner.Hex())
		}
		usedAddrs = append(usedAddrs, common.CopyBytes(addr[:])) // Copy needed for closure
	}
	if prefetcher != nil {
		prefetcher.used(s.originalRoot, usedAddrs)
	}
	if len(s.stateObjectsPending) > 0 {
		s.stateObjectsPending = make(map[common.Address]struct{})
	}
	// Track the amount of time wasted on hashing the account trie
	if metrics.EnabledExpensive {
		defer func(start time.Time) { s.AccountHashes += time.Since(start) }(time.Now())
	}
	//log.Info("caver|IntermediateRoot|enter=3", "triehash", s.trie.Hash().String())
	return s.trie.Hash()
}

// Prepare sets the current transaction hash and index which are
// used when the EVM emits new state logs.
func (s *StateDB) Prepare(thash common.Hash, ti int) {
	s.thash = thash
	s.txIndex = ti
	s.accessList = newAccessList()
}

func (s *StateDB) clearJournalAndRefund() {
	if len(s.journal.entries) > 0 {
		s.journal = newJournal()
		s.refund = 0
	}
	s.validRevisions = s.validRevisions[:0] // Snapshots can be created without journal entires
}

// Commit writes the state to the underlying in-memory trie database.
func (s *StateDB) Commit(deleteEmptyObjects bool) (common.Hash, error) {
	if s.dbErr != nil {
		return common.Hash{}, fmt.Errorf("commit aborted due to earlier error: %v", s.dbErr)
	}
	// Finalize any pending changes and merge everything into the tries
	s.IntermediateRoot(deleteEmptyObjects)

	// Commit objects to the trie, measuring the elapsed time
	codeWriter := s.db.TrieDB().DiskDB().NewBatch()
	for addr := range s.stateObjectsDirty {
		if obj := s.stateObjects[addr]; !obj.deleted {
			// Write any contract code associated with the state object
			if obj.code != nil && obj.dirtyCode {
				rawdb.WriteCode(codeWriter, common.BytesToHash(obj.CodeHash()), obj.code)
				obj.dirtyCode = false
			}
			// Write any storage changes in the state object to its storage trie
			if err := obj.CommitTrie(s.db); err != nil {
				return common.Hash{}, err
			}
		}
	}
	if len(s.stateObjectsDirty) > 0 {
		s.stateObjectsDirty = make(map[common.Address]struct{})
	}
	if codeWriter.ValueSize() > 0 {
		if err := codeWriter.Write(); err != nil {
			log.Crit("Failed to commit dirty codes", "error", err)
		}
	}
	// Write the account trie changes, measuing the amount of wasted time
	var start time.Time
	if metrics.EnabledExpensive {
		start = time.Now()
	}
	// The onleaf func is called _serially_, so we can reuse the same account
	// for unmarshalling every time.
	var account Account
	root, err := s.trie.Commit(func(_ [][]byte, _ []byte, leaf []byte, parent common.Hash) error {
		if err := rlp.DecodeBytes(leaf, &account); err != nil {
			return nil
		}
		if account.Root != emptyRoot {
			s.db.TrieDB().Reference(account.Root, parent)
		}
		return nil
	})
	if metrics.EnabledExpensive {
		s.AccountCommits += time.Since(start)
	}
	// If snapshotting is enabled, update the snapshot tree with this new version
	if s.snap != nil {
		if metrics.EnabledExpensive {
			defer func(start time.Time) { s.SnapshotCommits += time.Since(start) }(time.Now())
		}
		// Only update if there's a state transition (skip empty Clique blocks)
		if parent := s.snap.Root(); parent != root {
			if err := s.snaps.Update(root, parent, s.snapDestructs, s.snapAccounts, s.snapStorage); err != nil {
				log.Warn("Failed to update snapshot tree", "from", parent, "to", root, "err", err)
			}
			// Keep 128 diff layers in the memory, persistent layer is 129th.
			// - head layer is paired with HEAD state
			// - head-1 layer is paired with HEAD-1 state
			// - head-127 layer(bottom-most diff layer) is paired with HEAD-127 state
			if err := s.snaps.Cap(root, 128); err != nil {
				log.Warn("Failed to cap snapshot tree", "root", root, "layers", 128, "err", err)
			}
		}
		s.snap, s.snapDestructs, s.snapAccounts, s.snapStorage = nil, nil, nil, nil
	}
	return root, err
}

// PrepareAccessList handles the preparatory steps for executing a state transition with
// regards to both EIP-2929 and EIP-2930:
//
// - Add sender to access list (2929)
// - Add destination to access list (2929)
// - Add precompiles to access list (2929)
// - Add the contents of the optional tx access list (2930)
//
// This method should only be called if Berlin/2929+2930 is applicable at the current number.
func (s *StateDB) PrepareAccessList(sender common.Address, dst *common.Address, precompiles []common.Address, list types.AccessList) {
	s.AddAddressToAccessList(sender)
	if dst != nil {
		s.AddAddressToAccessList(*dst)
		// If it's a create-tx, the destination will be added inside evm.create
	}
	for _, addr := range precompiles {
		s.AddAddressToAccessList(addr)
	}
	for _, el := range list {
		s.AddAddressToAccessList(el.Address)
		for _, key := range el.StorageKeys {
			s.AddSlotToAccessList(el.Address, key)
		}
	}
}

// AddAddressToAccessList adds the given address to the access list
func (s *StateDB) AddAddressToAccessList(addr common.Address) {
	if s.accessList.AddAddress(addr) {
		s.journal.append(accessListAddAccountChange{&addr})
	}
}

// AddSlotToAccessList adds the given (address, slot)-tuple to the access list
func (s *StateDB) AddSlotToAccessList(addr common.Address, slot common.Hash) {
	addrMod, slotMod := s.accessList.AddSlot(addr, slot)
	if addrMod {
		// In practice, this should not happen, since there is no way to enter the
		// scope of 'address' without having the 'address' become already added
		// to the access list (via call-variant, create, etc).
		// Better safe than sorry, though
		s.journal.append(accessListAddAccountChange{&addr})
	}
	if slotMod {
		s.journal.append(accessListAddSlotChange{
			address: &addr,
			slot:    &slot,
		})
	}
}

// AddressInAccessList returns true if the given address is in the access list.
func (s *StateDB) AddressInAccessList(addr common.Address) bool {
	return s.accessList.ContainsAddress(addr)
}

// SlotInAccessList returns true if the given (address, slot)-tuple is in the access list.
func (s *StateDB) SlotInAccessList(addr common.Address, slot common.Hash) (addressPresent bool, slotPresent bool) {
	return s.accessList.Contains(addr, slot)
}

// *** modify to support nft transaction 20211215 begin ***

// ChangeNFTOwner change nft's owner to newOwner.
//func (s *StateDB) ChangeNFTOwner(nftAddr common.Address, newOwner common.Address) {
//	stateObject := s.GetOrNewStateObject(nftAddr)
//	if stateObject != nil {
//		s.SplitNFT(nftAddr, 0)
//		stateObject.ChangeNFTOwner(newOwner)
//		// merge nft automatically
//		s.MergeNFT(nftAddr)
//	}
//}

// GetNFTOwner retrieves the nft owner from the given nft address
func (s *StateDB) GetNFTOwner(nftAddr common.Address) common.Address {
	storeAddr, _, ok := s.GetNFTStoreAddress(nftAddr, 0)
	if ok {
		log.Info("StateDB.GetNFTOwner()", "nftAddr", nftAddr.String(), "storeAddr", storeAddr.String())
		stateObject := s.getStateObject(storeAddr)
		//stateObject := s.getDeletedStateObject2(nftAddr)
		if stateObject != nil {
			return stateObject.NFTOwner()
		}
	}

	return common.Address{}
}

// *** modify to support nft transaction 20211215 end ***

// *** modify to merge NFT 20211224 begin ***

func (s *StateDB) IsCanMergeNFT(nftAddr common.Address) bool {
	if len(nftAddr) == 0 {
		return false
	}
	nftAddrS := nftAddr.String()
	if strings.HasPrefix(nftAddrS, "0x") ||
		strings.HasPrefix(nftAddrS, "0X") {
		nftAddrS = string([]byte(nftAddrS)[2:])
	}

	// 1. get nftaddr's owner
	//nftOwner := s.GetNFTOwner(nftAddr)
	nftStateObject := s.getStateObject(nftAddr)
	validNftAddrLen := len(nftAddr) - int(nftStateObject.GetNFTMergeLevel())

	// 2. convert nft Addr to bigInt
	parentAddrS := string([]byte(nftAddrS)[:len(nftAddrS)-int(2*(nftStateObject.GetNFTMergeLevel()+1))])
	addrInt := big.NewInt(0)
	addrInt.SetString(parentAddrS, 16)
	addrInt.Lsh(addrInt, 8)

	// 3. retrieve all the sibling leaf nodes of nftAddr
	siblingInt := big.NewInt(0)
	//nftAddrSLen := len(nftAddrS)
	for i := 0; i < 256; i++ {
		// 4. convert bigInt to common.Address, and then get Account from the trie.
		siblingInt.Add(addrInt, big.NewInt(int64(i)))
		//siblingAddr := common.BigToAddress(siblingInt)
		siblingAddrS := hex.EncodeToString(siblingInt.Bytes())
		siblingAddrSLen := len(siblingAddrS)
		var prefix0 string
		for i := 0; i < 2*validNftAddrLen-siblingAddrSLen; i++ {
			prefix0 = prefix0 + "0"
		}
		siblingAddrS = prefix0 + siblingAddrS
		var suffix0 string
		for i := 0; i < int(2*nftStateObject.GetNFTMergeLevel()); i++ {
			suffix0 = suffix0 + "0"
		}
		siblingAddrS = siblingAddrS + suffix0
		siblingAddr := common.HexToAddress(siblingAddrS)
		//fmt.Println("siblingAddrS=", siblingAddrS)
		//fmt.Println("siblingAddr=", siblingAddr.String())
		//fmt.Println("nftAddr=", nftAddr.String())
		if siblingAddr == nftAddr {
			continue
		}

		siblingStateObject := s.getStateObject(siblingAddr)
		if siblingStateObject == nil ||
			siblingStateObject.NFTOwner() != nftStateObject.NFTOwner() ||
			siblingStateObject.GetNFTMergeLevel() != nftStateObject.GetNFTMergeLevel() {
			return false
		}

	}

	return true
}

func (s *StateDB) MergeNFT(nftAddr common.Address) error {
	if !s.IsCanMergeNFT(nftAddr) {
		return nil
	}
	nftAddrS := nftAddr.String()
	if strings.HasPrefix(nftAddrS, "0x") ||
		strings.HasPrefix(nftAddrS, "0X") {
		nftAddrS = string([]byte(nftAddrS)[2:])
	}

	// 1. get nftaddr's owner
	//nftOwner := s.GetNFTOwner(nftAddr)
	nftStateObject := s.getStateObject(nftAddr)
	nftStateObject = nftStateObject.deepCopy(s)
	validNftAddrLen := len(nftAddr) - int(nftStateObject.GetNFTMergeLevel())

	// 2. convert nft Addr to bigInt
	parentAddrS := string([]byte(nftAddrS)[:len(nftAddrS)-int(2*(nftStateObject.GetNFTMergeLevel()+1))])
	addrInt := big.NewInt(0)
	addrInt.SetString(parentAddrS, 16)
	addrInt.Lsh(addrInt, 8)

	// 3. retrieve all the sibling leaf nodes of nftAddr
	siblingInt := big.NewInt(0)
	//nftAddrSLen := len(nftAddrS)
	for i := 0; i < 256; i++ {
		// 4. convert bigInt to common.Address,
		// and then delete all sibling nodes and itself from the trie.
		siblingInt.Add(addrInt, big.NewInt(int64(i)))
		//siblingAddr := common.BigToAddress(siblingInt)
		siblingAddrS := hex.EncodeToString(siblingInt.Bytes())
		siblingAddrSLen := len(siblingAddrS)
		var prefix0 string
		for i := 0; i < 2*validNftAddrLen-siblingAddrSLen; i++ {
			prefix0 = prefix0 + "0"
		}
		siblingAddrS = prefix0 + siblingAddrS
		var suffix0 string
		for i := 0; i < int(2*nftStateObject.GetNFTMergeLevel()); i++ {
			suffix0 = suffix0 + "0"
		}
		siblingAddrS = siblingAddrS + suffix0
		siblingAddr := common.HexToAddress(siblingAddrS)
		//fmt.Println("siblingAddrS=", siblingAddrS)
		//fmt.Println("siblingAddr=", siblingAddr.String())
		//fmt.Println("nftAddr=", nftAddr.String())
		siblingStateObject := s.getStateObject(siblingAddr)
		//siblingStateObject.data.AccountNFT = AccountNFT{}
		siblingStateObject.CleanNFT()
		//s.deleteStateObject(siblingStateObject)
		//s.updateStateObject(siblingStateObject)

	}

	// new merged nft address
	newMergedAddrS := parentAddrS
	for i := 0; i < 2*len(nftAddr)-len(parentAddrS); i++ {
		newMergedAddrS = newMergedAddrS + "0"
	}
	newMergedAddr := common.HexToAddress(newMergedAddrS)
	index := strings.LastIndex(nftStateObject.data.MetaURL, "/")
	metaUrl := string([]byte(nftStateObject.data.MetaURL)[:index])
	metaUrl = metaUrl + "/" + newMergedAddr.String()
	var newMergeStateObject *stateObject
	if s.Exist(newMergedAddr) {
		newMergeStateObject = s.getStateObject(newMergedAddr)
		//newMergeStateObject.data.MergeLevel = nftStateObject.data.MergeLevel + 1
		//newMergeStateObject.data.Owner = nftStateObject.data.Owner
		var mergeNumber uint32 = 0
		mergeNumber = uint32(math.BigPow(16, int64(nftStateObject.data.MergeLevel+1)).Uint64())
		newMergeStateObject.SetNFTInfo(
			nftStateObject.data.Name,
			nftStateObject.data.Symbol,
			//nftStateObject.data.Price,
			//nftStateObject.data.Direction,
			nftStateObject.data.Owner,
			nftStateObject.data.NFTApproveAddressList,
			nftStateObject.data.MergeLevel+1,
			mergeNumber,
			//nftStateObject.data.PledgedFlag,
			//new(big.Int).Set(nftStateObject.data.NFTPledgedBlockNumber),
			nftStateObject.data.Creator,
			nftStateObject.data.Royalty,
			nftStateObject.data.Exchanger,
			metaUrl)
	} else {
		s.CreateAccount(newMergedAddr)
		newMergeStateObject = s.getStateObject(newMergedAddr)
		//newMergeStateObject.data.MergeLevel = nftStateObject.data.MergeLevel + 1
		//newMergeStateObject.data.Owner = nftStateObject.data.Owner
		var mergeNumber uint32 = 0
		mergeNumber = uint32(math.BigPow(16, int64(nftStateObject.data.MergeLevel+1)).Uint64())
		newMergeStateObject.SetNFTInfo(
			nftStateObject.data.Name,
			nftStateObject.data.Symbol,
			//nftStateObject.data.Price,
			//nftStateObject.data.Direction,
			nftStateObject.data.Owner,
			nftStateObject.data.NFTApproveAddressList,
			nftStateObject.data.MergeLevel+1,
			mergeNumber,
			//nftStateObject.data.PledgedFlag,
			//new(big.Int).Set(nftStateObject.data.NFTPledgedBlockNumber),
			nftStateObject.data.Creator,
			nftStateObject.data.Royalty,
			nftStateObject.data.Exchanger,
			metaUrl)
	}
	//s.updateStateObject(newMergeStateObject)
	s.MergeNFT(newMergedAddr)

	return nil
}

// *** modify to merge NFT 20211224 end ***

// Get the store address for a nft
const QUERYDEPTHLIMIT = 3

func (s *StateDB) GetNFTStoreAddress(address common.Address,
	depth int) (nftStoreAddress, owner common.Address, ok bool) {
	if depth > QUERYDEPTHLIMIT {
		return common.Address{}, common.Address{}, false
	}

	emptyNFTAddr := common.Address{}
	nftStateObj := s.getStateObject(address)
	if nftStateObj == nil {
		return common.Address{}, common.Address{}, false
	}
	if nftStateObj.data.Owner != emptyNFTAddr &&
		int(nftStateObj.GetNFTMergeLevel()) == depth {
		return address, nftStateObj.data.Owner, true
	} else {
		var parentAddrBytes []byte
		parentAddrBytes = append(parentAddrBytes, address[:len(address)-(depth+1)]...)
		for i := 0; i < (depth + 1); i++ {
			parentAddrBytes = append(parentAddrBytes, byte(0))
		}

		parentAddr := common.BytesToAddress(parentAddrBytes)
		depth = depth + 1
		return s.GetNFTStoreAddress(parentAddr, depth)
	}
}

//1. find store address according nft address which need to transfer.
//2. nft address which need to transfer is not exist if store address is empty.
//3. get stateobject of store address, then get mergeLevel
//4. return directly if mergelevel less than level
//5. if mergelevel more than level, check whether nft address is child address which's mergelevel equal level
//6. split snft
func (s *StateDB) SplitNFT(nftAddr common.Address, level int) {
	storeAddr, owner, ok := s.GetNFTStoreAddress(nftAddr, 0)
	if !ok {
		return
	}
	fmt.Println(storeAddr.String(), owner.String())

	storeStateObject := s.getStateObject(storeAddr)
	mergeLevel := int(storeStateObject.GetNFTMergeLevel())
	if mergeLevel <= level {
		return
	}

	storeAddrBytes := storeAddr.Bytes()
	nftAddrBytes := nftAddr.Bytes()
	//shouldNFTAddrBytes := storeAddrBytes[:len(storeAddrBytes)-mergeLevel]
	var shouldNFTAddrBytes []byte
	shouldNFTAddrBytes = append(shouldNFTAddrBytes, storeAddrBytes[:len(storeAddrBytes)-mergeLevel]...)
	shouldNFTAddrBytes = append(shouldNFTAddrBytes, nftAddrBytes[len(storeAddrBytes)-mergeLevel:len(storeAddrBytes)-level]...)
	shouldNFTAddrBytes = append(shouldNFTAddrBytes, storeAddrBytes[len(storeAddrBytes)-level:]...)
	if bytes.Compare(shouldNFTAddrBytes, nftAddrBytes) != 0 {
		return
	}

	storeStateObject = storeStateObject.deepCopy(s)

	var splitAddrBytes []byte
	var splitAddr common.Address
	var newSplitStateObject *stateObject
	var metaUrl string
	var index int
	for i := 0; i < mergeLevel-level; i++ {
		//if len(splitAddrBytes) > 0 {
		splitAddrBytes = splitAddrBytes[:0]
		//}
		splitAddrBytes = append(splitAddrBytes, storeAddrBytes[:len(storeAddrBytes)-mergeLevel]...)
		splitAddrBytes = append(splitAddrBytes, nftAddrBytes[len(storeAddrBytes)-mergeLevel:len(storeAddrBytes)-mergeLevel+i]...)
		splitAddrBytes = append(splitAddrBytes, storeAddrBytes[len(storeAddrBytes)-mergeLevel+i:]...)
		for j := 0; j < 256; j++ {
			splitAddrBytes[len(storeAddrBytes)-mergeLevel+i] = byte(j)
			splitAddr = common.BytesToAddress(splitAddrBytes)
			metaUrl = ""
			index = 0
			index = strings.LastIndex(storeStateObject.data.MetaURL, "/")
			metaUrl = string([]byte(storeStateObject.data.MetaURL)[:index])
			metaUrl = metaUrl + "/" + splitAddr.String()
			if s.Exist(splitAddr) {
				newSplitStateObject = s.getStateObject(splitAddr)
				//newSplitStateObject.data.MergeLevel = storeStateObject.data.MergeLevel - uint8(i + 1)
				//newSplitStateObject.data.Owner = storeStateObject.data.Owner
				var mergeNumber uint32 = 0
				mergeNumber = uint32(math.BigPow(16, int64(storeStateObject.data.MergeLevel-uint8(i+1))).Uint64())
				newSplitStateObject.SetNFTInfo(
					storeStateObject.data.Name,
					storeStateObject.data.Symbol,
					//storeStateObject.data.Price,
					//storeStateObject.data.Direction,
					storeStateObject.data.Owner,
					storeStateObject.data.NFTApproveAddressList,
					storeStateObject.data.MergeLevel-uint8(i+1),
					mergeNumber,
					//storeStateObject.data.PledgedFlag,
					//new(big.Int).Set(storeStateObject.data.NFTPledgedBlockNumber),
					storeStateObject.data.Creator,
					storeStateObject.data.Royalty,
					storeStateObject.data.Exchanger,
					metaUrl)
			} else {
				s.CreateAccount(splitAddr)
				newSplitStateObject = s.getStateObject(splitAddr)
				//newSplitStateObject.data.MergeLevel = storeStateObject.data.MergeLevel - uint8(i + 1)
				//newSplitStateObject.data.Owner = storeStateObject.data.Owner
				var mergeNumber uint32 = 0
				mergeNumber = uint32(math.BigPow(16, int64(storeStateObject.data.MergeLevel-uint8(i+1))).Uint64())
				newSplitStateObject.SetNFTInfo(
					storeStateObject.data.Name,
					storeStateObject.data.Symbol,
					//storeStateObject.data.Price,
					//storeStateObject.data.Direction,
					storeStateObject.data.Owner,
					storeStateObject.data.NFTApproveAddressList,
					storeStateObject.data.MergeLevel-uint8(i+1),
					mergeNumber,
					//storeStateObject.data.PledgedFlag,
					//new(big.Int).Set(storeStateObject.data.NFTPledgedBlockNumber),
					storeStateObject.data.Creator,
					storeStateObject.data.Royalty,
					storeStateObject.data.Exchanger,
					metaUrl)
			}
			//s.updateStateObject(newSplitStateObject)
		}
	}
}

func (s *StateDB) ConstructLog(mergedNFTAddress common.Address,
	owner common.Address,
	mergedNFTLevel uint8,
	mergedNFTNumber uint32,
	blockNumber *big.Int,
	mergedNFTs []*MergedNFT) *types.Log {
	var temp string = ""
	//event hash: MergeSNFT(address indexed snft,address indexed owner,uint256 pieces)
	hash1 := common.HexToHash("2b2711f6ad8adbb2fc8751c8400b9c6ebdaf9ea371995641808a7c692d89d46a")
	nftAddrString := mergedNFTAddress.Hex()
	nftAddrString = string([]byte(nftAddrString)[2 : len(nftAddrString)-int(mergedNFTLevel)])
	for i := 0; i < 64-len(nftAddrString); i++ {
		temp = temp + "0"
	}
	hash2 := common.HexToHash(temp + nftAddrString)
	ownerString := owner.Hex()
	ownerString = string([]byte(ownerString)[2:])
	hash3 := common.HexToHash("000000000000000000000000" + ownerString)

	log := &types.Log{
		Address: common.Address{},
		Topics: []common.Hash{
			hash1,
			hash2,
			hash3,
		},
		Data:        big.NewInt(int64(mergedNFTNumber)).FillBytes(make([]byte, 32)),
		BlockNumber: blockNumber.Uint64(),
	}

	nftsData, err := json.Marshal(mergedNFTs)
	if err == nil {
		toString := hex.EncodeToString(nftsData)
		log.Data = append(log.Data, []byte(toString)...)
	}

	return log
}

// ChangeNFTOwner change nft's owner to newOwner.
func (s *StateDB) ChangeNFTOwner(nftAddr common.Address,
	newOwner common.Address,
	level int,
	blocknumber *big.Int) {
	stateObject := s.GetOrNewStateObject(nftAddr)
	if stateObject != nil {
		if s.IsOfficialNFT(nftAddr) {
			//s.SplitNFT16(nftAddr, level)
			// subtract old Owner's voteweight
			initAmount := s.calculateExchangeAmount(stateObject.GetNFTMergeLevel(), stateObject.GetMergeNumber())
			amount := s.GetExchangAmount(nftAddr, initAmount)
			oldOwnerStateObject := s.GetOrNewStateObject(stateObject.NFTOwner())
			if oldOwnerStateObject.VoteWeight().Cmp(amount) < 0 {
				log.Error("StateDB.ChangeNFTOwner()", "old owner's voteweight less nft's value")
				amount.Set(oldOwnerStateObject.VoteWeight())
			}

			stateObject.ChangeNFTOwner(newOwner)
			oldOwnerStateObject.SubVoteWeight(amount)
			// merge nft automatically
			increaseValue, _ := s.MergeNFT16(nftAddr, blocknumber)

			// add new Owner's voteweight
			totalValue := new(big.Int).Add(increaseValue, amount)
			newOwnerStateObject := s.GetOrNewStateObject(newOwner)
			newOwnerStateObject.AddVoteWeight(totalValue)

		} else {
			stateObject.ChangeNFTOwner(newOwner)
		}
	}
}

// GetNFTOwner16 retrieves the nft owner from the given nft address
//func (s *StateDB) GetNFTOwner16(nftAddr common.Address) common.Address {
//	storeAddr, _, ok := s.GetNFTStoreAddress16(nftAddr, 0)
//	if ok {
//		log.Info("StateDB.GetNFTOwner16()", "nftAddr", nftAddr.String(), "storeAddr", storeAddr.String())
//		stateObject := s.getStateObject(storeAddr)
//		//stateObject := s.getDeletedStateObject2(nftAddr)
//		if stateObject != nil {
//			return stateObject.NFTOwner()
//		}
//	}
//
//	return common.Address{}
//}
// if snfts have been merged, original snfts are not exist, they become a new merged snft
func (s *StateDB) GetNFTOwner16(nftAddr common.Address) common.Address {
	stateObject := s.getStateObject(nftAddr)
	if stateObject != nil {
		return stateObject.NFTOwner()
	}

	return common.Address{}
}

func (s *StateDB) IsCanMergeNFT16(nftAddr common.Address) bool {
	if len(nftAddr) == 0 {
		return false
	}
	emptyAddress := common.Address{}

	nftAddrS := nftAddr.String()
	if strings.HasPrefix(nftAddrS, "0x") ||
		strings.HasPrefix(nftAddrS, "0X") {
		nftAddrS = string([]byte(nftAddrS)[2:])
	}

	// 1. get nftaddr's owner
	//nftOwner := s.GetNFTOwner(nftAddr)
	nftStateObject := s.getStateObject(nftAddr)
	if nftStateObject == nil {
		return false
	}
	mergeLevel := nftStateObject.GetNFTMergeLevel()
	if mergeLevel >= QUERYDEPTHLIMIT16 {
		return false
	}
	validNftAddrLen := len(nftAddrS) - int(mergeLevel)

	// 2. convert nft Addr to bigInt
	parentAddrS := string([]byte(nftAddrS)[:len(nftAddrS)-int((mergeLevel+1))])
	addrInt := big.NewInt(0)
	addrInt.SetString(parentAddrS, 16)
	addrInt.Lsh(addrInt, 4)

	// 3. retrieve all the sibling leaf nodes of nftAddr
	siblingInt := big.NewInt(0)
	//nftAddrSLen := len(nftAddrS)
	for i := 0; i < 16; i++ {
		// 4. convert bigInt to common.Address, and then get Account from the trie.
		siblingInt.Add(addrInt, big.NewInt(int64(i)))
		//siblingAddr := common.BigToAddress(siblingInt)
		siblingAddrS := hex.EncodeToString(siblingInt.Bytes())
		siblingAddrSLen := len(siblingAddrS)
		var prefix0 string
		for i := 0; i < validNftAddrLen-siblingAddrSLen; i++ {
			prefix0 = prefix0 + "0"
		}
		siblingAddrS = prefix0 + siblingAddrS
		var suffix0 string
		for i := 0; i < int(mergeLevel); i++ {
			suffix0 = suffix0 + "0"
		}
		siblingAddrS = siblingAddrS + suffix0
		siblingAddr := common.HexToAddress(siblingAddrS)
		//fmt.Println("siblingAddrS=", siblingAddrS)
		//fmt.Println("siblingAddr=", siblingAddr.String())
		//fmt.Println("nftAddr=", nftAddr.String())
		if siblingAddr == nftAddr {
			continue
		}

		siblingStateObject := s.getStateObject(siblingAddr)
		if siblingStateObject != nil &&
			siblingStateObject.NFTOwner() == emptyAddress {
			if siblingStateObject.GetNFTMergeLevel() != mergeLevel {
				return false
			} else {
				continue
			}
		}
		//if siblingStateObject == nil ||
		//	siblingStateObject.NFTOwner() != nftStateObject.NFTOwner() ||
		//	siblingStateObject.GetNFTMergeLevel() != mergeLevel ||
		//	siblingStateObject.GetPledgedFlag() != false {
		if siblingStateObject == nil ||
			siblingStateObject.NFTOwner() != nftStateObject.NFTOwner() ||
			siblingStateObject.GetNFTMergeLevel() != mergeLevel {
			return false
		}

		// add new limit by metaurl 2022-06-20
		nftIndex := strings.LastIndex(nftStateObject.data.MetaURL, "/")
		nftMetaUrl := string([]byte(nftStateObject.data.MetaURL)[:nftIndex])
		siblingIndex := strings.LastIndex(siblingStateObject.data.MetaURL, "/")
		siblingMetaUrl := string([]byte(siblingStateObject.data.MetaURL)[:siblingIndex])
		if nftMetaUrl != siblingMetaUrl {
			return false
		}

	}

	return true
}

func (s *StateDB) GetExistAddress(nftAddr common.Address, mergeLevel uint8) common.Address {
	if len(nftAddr) == 0 {
		return common.Address{}
	}
	emptyAddress := common.Address{}

	nftAddrS := nftAddr.String()
	if strings.HasPrefix(nftAddrS, "0x") ||
		strings.HasPrefix(nftAddrS, "0X") {
		nftAddrS = string([]byte(nftAddrS)[2:])
	}

	if mergeLevel >= QUERYDEPTHLIMIT16 {
		return common.Address{}
	}

	validNftAddrLen := len(nftAddrS) - int(mergeLevel)

	// 2. convert nft Addr to bigInt
	parentAddrS := string([]byte(nftAddrS)[:len(nftAddrS)-int((mergeLevel+1))])
	addrInt := big.NewInt(0)
	addrInt.SetString(parentAddrS, 16)
	addrInt.Lsh(addrInt, 4)

	// 3. retrieve all the sibling leaf nodes of nftAddr
	siblingInt := big.NewInt(0)
	//nftAddrSLen := len(nftAddrS)
	for i := 0; i < 16; i++ {
		// 4. convert bigInt to common.Address, and then get Account from the trie.
		siblingInt.Add(addrInt, big.NewInt(int64(i)))
		//siblingAddr := common.BigToAddress(siblingInt)
		siblingAddrS := hex.EncodeToString(siblingInt.Bytes())
		siblingAddrSLen := len(siblingAddrS)
		var prefix0 string
		for i := 0; i < validNftAddrLen-siblingAddrSLen; i++ {
			prefix0 = prefix0 + "0"
		}
		siblingAddrS = prefix0 + siblingAddrS
		var suffix0 string
		for i := 0; i < int(mergeLevel); i++ {
			suffix0 = suffix0 + "0"
		}
		siblingAddrS = siblingAddrS + suffix0
		siblingAddr := common.HexToAddress(siblingAddrS)
		//fmt.Println("siblingAddrS=", siblingAddrS)
		//fmt.Println("siblingAddr=", siblingAddr.String())
		//fmt.Println("nftAddr=", nftAddr.String())
		if siblingAddr == nftAddr {
			continue
		}

		siblingStateObject := s.getStateObject(siblingAddr)
		if siblingStateObject != nil &&
			siblingStateObject.NFTOwner() != emptyAddress {
			if siblingStateObject.GetNFTMergeLevel() == mergeLevel {
				return siblingAddr
			} else {
				return common.Address{}
			}
		}
	}

	for i := 0; i < int(mergeLevel+1); i++ {
		parentAddrS = parentAddrS + "0"
	}
	parentAddrS = "0x" + parentAddrS
	parentAddr := common.HexToAddress(parentAddrS)
	return s.GetExistAddress(parentAddr, mergeLevel+1)
}

type MergedNFT struct {
	Address common.Address `json:"address"`
	Number  uint32         `json:"number"`
}

// MergeNFT16 merge snfts and return the increase of value because of merging.
func (s *StateDB) MergeNFT16(nftAddr common.Address, blocknumber *big.Int) (*big.Int, error) {
	if !s.IsCanMergeNFT16(nftAddr) {
		return big.NewInt(0), nil
	}
	emptyAddress := common.Address{}

	nftAddrS := nftAddr.String()
	if strings.HasPrefix(nftAddrS, "0x") ||
		strings.HasPrefix(nftAddrS, "0X") {
		nftAddrS = string([]byte(nftAddrS)[2:])
	}

	// 1. get nftaddr's owner
	//nftOwner := s.GetNFTOwner(nftAddr)
	nftStateObject := s.getStateObject(nftAddr)
	nftStateObject = nftStateObject.deepCopy(s)
	validNftAddrLen := len(nftAddrS) - int(nftStateObject.GetNFTMergeLevel())

	// 2. convert nft Addr to bigInt
	parentAddrS := string([]byte(nftAddrS)[:len(nftAddrS)-int((nftStateObject.GetNFTMergeLevel()+1))])
	addrInt := big.NewInt(0)
	addrInt.SetString(parentAddrS, 16)
	addrInt.Lsh(addrInt, 4)

	// 3. retrieve all the sibling leaf nodes of nftAddr
	siblingInt := big.NewInt(0)
	var mergedNFTs []*MergedNFT
	//nftAddrSLen := len(nftAddrS)
	var mergeNumber uint32 = 0
	for i := 0; i < 16; i++ {
		// 4. convert bigInt to common.Address,
		// and then delete all sibling nodes and itself from the trie.
		siblingInt.Add(addrInt, big.NewInt(int64(i)))
		//siblingAddr := common.BigToAddress(siblingInt)
		siblingAddrS := hex.EncodeToString(siblingInt.Bytes())
		siblingAddrSLen := len(siblingAddrS)
		var prefix0 string
		for i := 0; i < validNftAddrLen-siblingAddrSLen; i++ {
			prefix0 = prefix0 + "0"
		}
		siblingAddrS = prefix0 + siblingAddrS
		var suffix0 string
		for i := 0; i < int(nftStateObject.GetNFTMergeLevel()); i++ {
			suffix0 = suffix0 + "0"
		}
		siblingAddrS = siblingAddrS + suffix0
		siblingAddr := common.HexToAddress(siblingAddrS)
		//fmt.Println("siblingAddrS=", siblingAddrS)
		//fmt.Println("siblingAddr=", siblingAddr.String())
		//fmt.Println("nftAddr=", nftAddr.String())
		siblingStateObject := s.getStateObject(siblingAddr)
		//siblingStateObject.data.AccountNFT = AccountNFT{}
		if siblingStateObject.NFTOwner() != emptyAddress {
			mergedNFT := &MergedNFT{
				Address: siblingAddr,
				Number:  siblingStateObject.GetMergeNumber(),
			}
			mergedNFTs = append(mergedNFTs, mergedNFT)
			mergeNumber = mergeNumber + siblingStateObject.GetMergeNumber()
			siblingStateObject.CleanNFT()
		}
		//s.deleteStateObject(siblingStateObject)
		//s.updateStateObject(siblingStateObject)

	}

	// new merged nft address
	newMergedAddrS := parentAddrS
	for i := 0; i < 2*len(nftAddr)-len(parentAddrS); i++ {
		newMergedAddrS = newMergedAddrS + "0"
	}
	newMergedAddr := common.HexToAddress(newMergedAddrS)
	index := strings.LastIndex(nftStateObject.data.MetaURL, "/")
	metaUrl := string([]byte(nftStateObject.data.MetaURL)[:index])
	metaUrl = metaUrl + "/" + newMergedAddr.String()
	var newMergeStateObject *stateObject
	if s.Exist(newMergedAddr) {
		newMergeStateObject = s.getStateObject(newMergedAddr)
		//newMergeStateObject.data.MergeLevel = nftStateObject.data.MergeLevel + 1
		//newMergeStateObject.data.Owner = nftStateObject.data.Owner
		//var mergeNumber uint32 = 0
		//mergeNumber = uint32(math.BigPow(16, int64(nftStateObject.data.MergeLevel+1)).Uint64())
		newMergeStateObject.SetNFTInfo(
			nftStateObject.data.Name,
			nftStateObject.data.Symbol,
			//nftStateObject.data.Price,
			//nftStateObject.data.Direction,
			nftStateObject.data.Owner,
			nftStateObject.data.NFTApproveAddressList,
			nftStateObject.data.MergeLevel+1,
			mergeNumber,
			//nftStateObject.data.PledgedFlag,
			//new(big.Int).Set(nftStateObject.data.NFTPledgedBlockNumber),
			nftStateObject.data.Creator,
			nftStateObject.data.Royalty,
			nftStateObject.data.Exchanger,
			metaUrl)
	} else {
		s.CreateAccount(newMergedAddr)
		newMergeStateObject = s.getStateObject(newMergedAddr)
		//newMergeStateObject.data.MergeLevel = nftStateObject.data.MergeLevel + 1
		//newMergeStateObject.data.Owner = nftStateObject.data.Owner
		//var mergeNumber uint32 = 0
		//mergeNumber = uint32(math.BigPow(16, int64(nftStateObject.data.MergeLevel+1)).Uint64())
		newMergeStateObject.SetNFTInfo(
			nftStateObject.data.Name,
			nftStateObject.data.Symbol,
			//nftStateObject.data.Price,
			//nftStateObject.data.Direction,
			nftStateObject.data.Owner,
			nftStateObject.data.NFTApproveAddressList,
			nftStateObject.data.MergeLevel+1,
			mergeNumber,
			//nftStateObject.data.PledgedFlag,
			//new(big.Int).Set(nftStateObject.data.NFTPledgedBlockNumber),
			nftStateObject.data.Creator,
			nftStateObject.data.Royalty,
			nftStateObject.data.Exchanger,
			metaUrl)
	}
	//s.updateStateObject(newMergeStateObject)

	// calculate the increase of value
	mergedInitAmount := s.calculateExchangeAmount(newMergeStateObject.GetNFTMergeLevel(), mergeNumber)
	mergedAmount := s.GetExchangAmount(newMergedAddr, mergedInitAmount)
	noMergedInitAmount := s.calculateExchangeAmount(newMergeStateObject.GetNFTMergeLevel()-1, mergeNumber)
	noMergedAmount := s.GetExchangAmount(newMergedAddr, noMergedInitAmount)
	increaseValue := new(big.Int).Sub(mergedAmount, noMergedAmount)

	// add merge snft log
	log := s.ConstructLog(newMergedAddr,
		newMergeStateObject.data.Owner,
		newMergeStateObject.data.MergeLevel,
		mergeNumber,
		blocknumber,
		mergedNFTs)
	s.AddLog(log)

	tempValue, _ := s.MergeNFT16(newMergedAddr, blocknumber)

	totalIncreaseValue := new(big.Int).Add(increaseValue, tempValue)

	return totalIncreaseValue, nil
}

// Get the store address for a nft
const QUERYDEPTHLIMIT16 = 3

func (s *StateDB) GetNFTStoreAddress16(address common.Address,
	depth int) (nftStoreAddress, owner common.Address, ok bool) {
	if depth > QUERYDEPTHLIMIT16 {
		return common.Address{}, common.Address{}, false
	}

	emptyNFTAddr := common.Address{}
	nftStateObj := s.getStateObject(address)
	if nftStateObj == nil {
		return common.Address{}, common.Address{}, false
	}
	if nftStateObj.data.Owner != emptyNFTAddr &&
		int(nftStateObj.GetNFTMergeLevel()) == depth {
		return address, nftStateObj.data.Owner, true
	} else {
		var parentAddrBytes []byte
		addressBytes16 := []byte(address.String())
		parentAddrBytes = append(parentAddrBytes, addressBytes16[:len(addressBytes16)-(depth+1)]...)
		for i := 0; i < (depth + 1); i++ {
			parentAddrBytes = append(parentAddrBytes, byte(0+48))
		}
		parentAddr := common.HexToAddress(string(parentAddrBytes))
		depth = depth + 1
		return s.GetNFTStoreAddress16(parentAddr, depth)
	}
}

func (s *StateDB) SplitNFT16(nftAddr common.Address, level int) {
	storeAddr, owner, ok := s.GetNFTStoreAddress16(nftAddr, 0)
	if !ok {
		return
	}
	fmt.Println(storeAddr.String(), owner.String())

	storeStateObject := s.getStateObject(storeAddr)
	mergeLevel := int(storeStateObject.GetNFTMergeLevel())
	if mergeLevel <= level {
		return
	}

	storeAddrBytes := []byte(storeAddr.String())
	nftAddrBytes := []byte(nftAddr.String())
	//shouldNFTAddrBytes := storeAddrBytes[:len(storeAddrBytes)-mergeLevel]
	var shouldNFTAddrBytes []byte
	shouldNFTAddrBytes = append(shouldNFTAddrBytes, storeAddrBytes[:len(storeAddrBytes)-mergeLevel]...)
	shouldNFTAddrBytes = append(shouldNFTAddrBytes, nftAddrBytes[len(storeAddrBytes)-mergeLevel:len(storeAddrBytes)-level]...)
	shouldNFTAddrBytes = append(shouldNFTAddrBytes, storeAddrBytes[len(storeAddrBytes)-level:]...)
	if bytes.Compare(shouldNFTAddrBytes, nftAddrBytes) != 0 {
		return
	}

	storeStateObject = storeStateObject.deepCopy(s)

	var splitAddrBytes []byte
	var splitAddr common.Address
	var newSplitStateObject *stateObject
	var metaUrl string
	var index int
	for i := 0; i < mergeLevel-level; i++ {
		//if len(splitAddrBytes) > 0 {
		splitAddrBytes = splitAddrBytes[:0]
		//}
		splitAddrBytes = append(splitAddrBytes, storeAddrBytes[:len(storeAddrBytes)-mergeLevel]...)
		splitAddrBytes = append(splitAddrBytes, nftAddrBytes[len(storeAddrBytes)-mergeLevel:len(storeAddrBytes)-mergeLevel+i]...)
		splitAddrBytes = append(splitAddrBytes, storeAddrBytes[len(storeAddrBytes)-mergeLevel+i:]...)
		for j := 0; j < 16; j++ {
			if j < 10 {
				splitAddrBytes[len(storeAddrBytes)-mergeLevel+i] = byte(j + 48)
			} else {
				splitAddrBytes[len(storeAddrBytes)-mergeLevel+i] = byte(j + 55)
			}
			splitAddr = common.HexToAddress(string(splitAddrBytes))
			metaUrl = ""
			index = 0
			index = strings.LastIndex(storeStateObject.data.MetaURL, "/")
			metaUrl = string([]byte(storeStateObject.data.MetaURL)[:index])
			metaUrl = metaUrl + "/" + splitAddr.String()
			if s.Exist(splitAddr) {
				newSplitStateObject = s.getStateObject(splitAddr)
				//newSplitStateObject.data.MergeLevel = storeStateObject.data.MergeLevel - uint8(i + 1)
				//newSplitStateObject.data.Owner = storeStateObject.data.Owner
				var mergeNumber uint32 = 0
				mergeNumber = uint32(math.BigPow(16, int64(storeStateObject.data.MergeLevel-uint8(i+1))).Uint64())
				newSplitStateObject.SetNFTInfo(
					storeStateObject.data.Name,
					storeStateObject.data.Symbol,
					//storeStateObject.data.Price,
					//storeStateObject.data.Direction,
					storeStateObject.data.Owner,
					storeStateObject.data.NFTApproveAddressList,
					storeStateObject.data.MergeLevel-uint8(i+1),
					mergeNumber,
					//storeStateObject.data.PledgedFlag,
					//new(big.Int).Set(storeStateObject.data.NFTPledgedBlockNumber),
					storeStateObject.data.Creator,
					storeStateObject.data.Royalty,
					storeStateObject.data.Exchanger,
					metaUrl)
			} else {
				s.CreateAccount(splitAddr)
				newSplitStateObject = s.getStateObject(splitAddr)
				//newSplitStateObject.data.MergeLevel = storeStateObject.data.MergeLevel - uint8(i + 1)
				//newSplitStateObject.data.Owner = storeStateObject.data.Owner
				var mergeNumber uint32 = 0
				mergeNumber = uint32(math.BigPow(16, int64(storeStateObject.data.MergeLevel-uint8(i+1))).Uint64())
				newSplitStateObject.SetNFTInfo(
					storeStateObject.data.Name,
					storeStateObject.data.Symbol,
					//storeStateObject.data.Price,
					//storeStateObject.data.Direction,
					storeStateObject.data.Owner,
					storeStateObject.data.NFTApproveAddressList,
					storeStateObject.data.MergeLevel-uint8(i+1),
					mergeNumber,
					//storeStateObject.data.PledgedFlag,
					//new(big.Int).Set(storeStateObject.data.NFTPledgedBlockNumber),
					storeStateObject.data.Creator,
					storeStateObject.data.Royalty,
					storeStateObject.data.Exchanger,
					metaUrl)
			}
			//s.updateStateObject(newSplitStateObject)
		}
	}
}

// IsOfficialNFT return true if nft address is created by official
func (s *StateDB) IsOfficialNFT(nftAddress common.Address) bool {
	maskByte := byte(128)
	nftByte := nftAddress[0]
	result := maskByte & nftByte
	if result == 128 {
		return true
	}
	return false
}

func (s *StateDB) InjectOfficialNFT(dir string,
	startIndex *big.Int,
	number uint64,
	royalty uint16,
	creator string) {
	injectNFT := &types.InjectedOfficialNFT{
		Dir:        dir,
		StartIndex: new(big.Int).Set(startIndex),
		Number:     number,
		Royalty:    royalty,
		Creator:    creator,
	}
	s.OfficialNFTPool.InjectedOfficialNFTs = append(s.OfficialNFTPool.InjectedOfficialNFTs, injectNFT)
}

/*
Owner common.Address
ApproveAddress common.Address
//Auctions map[string][]common.Address
// MergeLevel is the level of NFT merged
MergeLevel uint8

Creator common.Address
Royalty uint32
Exchanger common.Address
MetaURL string
*/
//- [X]NFT official mint
//
//func (s *StateDB) CreateNFTByOfficial(owners []common.Address, blocknumber *big.Int) {
//	for _, owner := range owners {
//		nftAddr := common.Address{}
//		var metaUrl string
//		var royalty uint32
//		var creator string
//		nftAddr, info, ok := s.SNFTExchangePool.PopAddress(blocknumber)
//		if !ok {
//			nftAddr = common.BytesToAddress(s.MintDeep.OfficialMint.Bytes())
//			injectedInfo := s.OfficialNFTPool.GetInjectedInfo(nftAddr)
//			if injectedInfo == nil {
//				return
//			}
//			metaUrl = injectedInfo.Dir + "/" + nftAddr.String()
//			royalty = injectedInfo.Royalty
//			creator = injectedInfo.Creator
//		} else {
//			metaUrl = info.MetalUrl
//			royalty = info.Royalty
//			creator = info.Creator
//		}
//		log.Info("CreateNFTByOfficial()", "--nftAddr=", nftAddr.String())
//
//		s.CreateAccount(nftAddr)
//		stateObject := s.GetOrNewStateObject(nftAddr)
//		if stateObject != nil {
//			stateObject.SetNFTInfo(
//				"",
//				"",
//				//big.NewInt(0),
//				//0,
//				owner,
//				common.Address{},
//				0,
//				1,
//				false,
//				big.NewInt(0),
//				common.HexToAddress(creator),
//				royalty,
//				common.Address{},
//				metaUrl)
//			s.MergeNFT(nftAddr)
//			if !ok {
//				s.OfficialNFTPool.DeleteExpireElem(s.MintDeep.OfficialMint)
//				s.MintDeep.OfficialMint.Add(s.MintDeep.OfficialMint, big.NewInt(1))
//			}
//		}
//	}
//
//	if s.OfficialNFTPool.RemainderNum(s.MintDeep.OfficialMint) <= 110 {
//		s.ElectNominatedOfficialNFT()
//	}
//}

func GetRewardAmount(blocknumber uint64, initamount *big.Int) *big.Int {
	times := blocknumber / ReduceRewardPeriod
	rewardratio := gomath.Pow(0.88, float64(times))
	u, _ := new(big.Float).Mul(big.NewFloat(rewardratio), new(big.Float).SetInt(initamount)).Uint64()

	return new(big.Int).SetUint64(u)
}

func (s *StateDB) CreateNFTByOfficial16(validators, exchangers []common.Address, blocknumber *big.Int) {

	// reward ERB or SNFT to validators
	log.Info("CreateNFTByOfficial16", "validators len=", len(validators), "blocknumber=", blocknumber.Uint64())
	for _, addr := range validators {
		log.Info("CreateNFTByOfficial16", "validators=", addr.Hex(), "blocknumber=", blocknumber.Uint64())
	}
	rewardAmount := GetRewardAmount(blocknumber.Uint64(), DREBlockReward)
	for _, owner := range validators {
		ownerObject := s.GetOrNewStateObject(owner)
		if ownerObject != nil {
			log.Info("ownerobj", "addr", ownerObject.address.Hex(), "blocknumber=", blocknumber.Uint64())
			ownerObject.AddBalance(rewardAmount)
		}
	}

	// reward SNFT to exchangers
	log.Info("CreateNFTByOfficial16", "exchangers len=", len(exchangers), "blocknumber=", blocknumber.Uint64())
	for _, addr := range exchangers {
		log.Info("CreateNFTByOfficial16", "exchangers=", addr.Hex(), "blocknumber=", blocknumber.Uint64())
	}
	for _, owner := range exchangers {
		nftAddr := common.Address{}
		var metaUrl string
		var royalty uint16
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
				//false,
				//big.NewInt(0),
				common.HexToAddress(creator),
				royalty,
				common.Address{},
				metaUrl)

			initAmount := s.calculateExchangeAmount(0, 1)
			amount := s.GetExchangAmount(nftAddr, initAmount)
			//increaseValue, mergedNFTAddress, NFTOwner, mergedNFTLevel, mergedNFTNumber, _ := s.MergeNFT16(nftAddr)
			//emptyAddress := common.Address{}
			//if mergedNFTAddress != emptyAddress {
			//	log := s.ConstructLog(mergedNFTAddress, NFTOwner, mergedNFTLevel, mergedNFTNumber, blocknumber)
			//	s.AddLog(log)
			//}
			increaseValue, _ := s.MergeNFT16(nftAddr, blocknumber)
			totalIncreaseValue := new(big.Int).Add(increaseValue, amount)
			ownerStateObject := s.GetOrNewStateObject(owner)
			if ownerStateObject != nil {
				ownerStateObject.AddVoteWeight(totalIncreaseValue)
			}

			//if !ok {
			s.OfficialNFTPool.DeleteExpireElem(s.MintDeep.OfficialMint)
			s.MintDeep.OfficialMint.Add(s.MintDeep.OfficialMint, big.NewInt(1))
			//}
		}
	}

	if s.OfficialNFTPool.RemainderNum(s.MintDeep.OfficialMint) <= 110 {
		s.ElectNominatedOfficialNFT(blocknumber)
	}
}

//- create a nft by user :creator can get a nft , include exchanger, royalty and meta.
//wormholes chain will assign a nft address to the nft.
//````
//{
//from:exchanger addr
//to:creator
//data:{
//version:0
//type:0
//royalty:
//metaUrl:
//}
//}
//````
//
func (s *StateDB) CreateNFTByUser(exchanger common.Address,
	owner common.Address,
	royalty uint16,
	metaurl string) (common.Address, bool) {
	nftAddr := common.BytesToAddress(s.MintDeep.UserMint.Bytes())
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
			//false,
			//big.NewInt(0),
			owner,
			royalty,
			exchanger,
			metaurl)
		s.MintDeep.UserMint.Add(s.MintDeep.UserMint, big.NewInt(1))
		return nftAddr, true
	}

	return common.Address{}, false
}

//- approve nft's authority:[?]
//````
//{
//from: owner
//to:address who get the authority
//data:{
//version:0
//type:2
//nftAddress:NFT address
//}
//}
//````
//
// ChangeApproveAddress is to approve all nfts
func (s *StateDB) ChangeApproveAddress(addr common.Address, approveAddr common.Address) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.ChangeApproveAddress(approveAddr)
	}
}

func (s *StateDB) CancelApproveAddress(addr common.Address, approveAddr common.Address) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.CancelApproveAddress(approveAddr)
	}
}

// ChangeNFTApproveAddress is to approve a nft
func (s *StateDB) ChangeNFTApproveAddress(nftAddr common.Address, approveAddr common.Address) {
	stateObject := s.GetOrNewStateObject(nftAddr)
	if stateObject != nil {
		stateObject.ChangeNFTApproveAddress(approveAddr)
	}
}

func (s *StateDB) CancelNFTApproveAddress(nftAddr common.Address, approveAddr common.Address) {
	stateObject := s.GetOrNewStateObject(nftAddr)
	if stateObject != nil {
		stateObject.CancelNFTApproveAddress(approveAddr)
	}
}

//- exchange nft:the nft which is exchanged will be add to exchange pool.
//user can get erb according to snft level
//````
//{
//from:owner
//to:0xffff....ffff
//data:{
//version:0
//type:3
//nftAddress:NFT address
//}
//}
//````
//
func (s *StateDB) ExchangeNFTToCurrency(address common.Address,
	nftaddress common.Address,
	blocknumber *big.Int,
	level int) {
	//s.SplitNFT16(nftaddress, level)
	nftStateObject := s.GetOrNewStateObject(nftaddress)
	stateObject := s.GetOrNewStateObject(address)
	if nftStateObject != nil && stateObject != nil {
		//LastIndex := strings.LastIndex(nftStateObject.data.MetaURL, "/")
		//Url := string([]byte(nftStateObject.data.MetaURL)[:LastIndex])
		//nftExchange := types.SNFTExchange{
		//	NFTAddress:         nftStateObject.address,
		//	MergeLevel:         nftStateObject.data.MergeLevel,
		//	CurrentMintAddress: nftStateObject.address,
		//	BlockNumber:        new(big.Int).Set(blocknumber),
		//	InjectedInfo: types.InjectedInfo{
		//		//MetalUrl: nftStateObject.data.MetaURL,
		//		MetalUrl: Url,
		//		Royalty:  nftStateObject.data.Royalty,
		//		Creator:  nftStateObject.data.Creator.String(),
		//	},
		//}
		//s.SNFTExchangePool.SNFTExchanges = append(s.SNFTExchangePool.SNFTExchanges, &nftExchange)

		emptyAddress := common.Address{}
		//creator := nftStateObject.GetCreator()
		//creatorObj := s.GetOrNewStateObject(creator)
		initAmount := s.calculateExchangeAmount(nftStateObject.GetNFTMergeLevel(), nftStateObject.GetMergeNumber())
		amount := s.GetExchangAmount(nftaddress, initAmount)

		//if creator != emptyAddress && creatorObj != nil {
		//	creatorObj.AddBalance(big.NewInt(0).Div(amount, big.NewInt(10)))
		//}
		mergeLevel := nftStateObject.GetNFTMergeLevel()
		nftStateObject.CleanNFT()
		stateObject.AddBalance(amount)
		stateObject.SubVoteWeight(amount)

		//Merge SNFT
		existNftAddress := s.GetExistAddress(nftaddress, mergeLevel)
		if existNftAddress != emptyAddress {
			existNftStateObject := s.GetOrNewStateObject(existNftAddress)
			nftOwner := existNftStateObject.NFTOwner()
			increaseValue, _ := s.MergeNFT16(existNftAddress, blocknumber)
			existOwnerStateObject := s.GetOrNewStateObject(nftOwner)
			if existOwnerStateObject != nil {
				existOwnerStateObject.AddVoteWeight(increaseValue)
			}
		}
	}
}

func (s *StateDB) GetExchangAmount(nftaddress common.Address, initamount *big.Int) *big.Int {
	nftInt := new(big.Int).SetBytes(nftaddress.Bytes())
	baseInt, _ := big.NewInt(0).SetString("8000000000000000000000000000000000000000", 16)
	nftInt.Sub(nftInt, baseInt)
	//nftInt.Add(nftInt, big.NewInt(1))
	nftInt.Div(nftInt, big.NewInt(4096))
	times := nftInt.Uint64() / ExchangePeriod
	rewardratio := gomath.Pow(0.88, float64(times))
	result := big.NewInt(0)
	new(big.Float).Mul(big.NewFloat(rewardratio), new(big.Float).SetInt(initamount)).Int(result)

	return result
}

func (s *StateDB) calculateExchangeAmount(level uint8, mergenumber uint32) *big.Int {
	//nftNumber := math.BigPow(16, int64(level))
	nftNumber := big.NewInt(int64(mergenumber))
	switch {
	case level == 0:
		radix, _ := big.NewInt(0).SetString("30000000000000000", 10)
		return big.NewInt(0).Mul(nftNumber, radix)
	case level == 1:
		radix, _ := big.NewInt(0).SetString("143000000000000000", 10)
		return big.NewInt(0).Mul(nftNumber, radix)
	case level == 2:
		radix, _ := big.NewInt(0).SetString("271000000000000000", 10)
		return big.NewInt(0).Mul(nftNumber, radix)
	default:
		radix, _ := big.NewInt(0).SetString("650000000000000000", 10)
		return big.NewInt(0).Mul(nftNumber, radix)
	}
}

func (s *StateDB) CalculateExchangeAmount(level uint8, mergenumber uint32) *big.Int {
	return s.calculateExchangeAmount(level, mergenumber)
}

//- pledge nft :NFT is pledged.
// the owner of the nft can get gasfee discount according to nft's level.
// a address can only pledge one nft.

//````
//{
//from:NFT owner
//to:0xffff...ffff
//data:{
//version:0
//type:4
//nftAddress:NFT address
//}
//}
//````
//
//func (s *StateDB) PledgeNFT(nftaddress common.Address, blocknumber *big.Int) {
//	nftStateObject := s.GetOrNewStateObject(nftaddress)
//	if nftStateObject != nil {
//		amount := s.calculateExchangeAmount(nftStateObject.GetNFTMergeLevel(), nftStateObject.GetMergeNumber())
//		addr := nftStateObject.NFTOwner()
//		nftToken := &types.PledgedToken{
//			Address: addr,
//			Amount:  amount,
//			Flag:    true,
//		}
//		s.ExchangerTokenPool = append(s.ExchangerTokenPool, nftToken)
//		nftStateObject.PledgeNFT(blocknumber)
//	}
//}

//- cancel nft pledge
//````
//{
//from:NFT owner
//to:0xffff...ffff
//data:{
//version:0
//type:5
//nftAddress:NFT address
//}
//}
//````
//
//func (s *StateDB) CancelPledgedNFT(nftaddress common.Address) {
//	nftStateObject := s.GetOrNewStateObject(nftaddress)
//	if nftStateObject != nil {
//		amount := s.calculateExchangeAmount(nftStateObject.GetNFTMergeLevel(), nftStateObject.GetMergeNumber())
//		addr := nftStateObject.NFTOwner()
//		nftToken := &types.PledgedToken{
//			Address: addr,
//			Amount:  amount,
//			Flag:    false,
//		}
//		s.ExchangerTokenPool = append(s.ExchangerTokenPool, nftToken)
//		nftStateObject.CancelPledgedNFT()
//	}
//}

//-  pledge token: a user who want to be a miner need to pledge token, must more than 100000 erb
//````
//{
//from:owner
//to:0xffff...ffff
//balance:????
//data:{
//version:0
//type:6
//}
//}
//````
//
func (s *StateDB) PledgeToken(address common.Address,
	amount *big.Int,
	proxy common.Address,
	blocknumber *big.Int) error {

	if amount == nil {
		amount = big.NewInt(0)
	}

	stateObject := s.GetOrNewStateObject(address)
	log.Info("PledgeToken", "address", address.Hex(), "proxy", proxy.Hex(), "amount", amount, "ValidatorPool", len(s.ValidatorPool))

	//Resolving duplicates is delegated
	empty := common.Address{}
	for _, v := range s.ValidatorPool {
		if v.Proxy != empty && v.Addr != address && v.Proxy == proxy {
			log.Info("PledgeToken|break", "address", address, "proxy", proxy)
			return errors.New("cannot delegate repeatedly")
		}
	}

	if stateObject != nil {
		pledgeToken := types.PledgedToken{
			Address:      address,
			Amount:       amount,
			Flag:         true,
			ProxyAddress: proxy,
		}
		log.Info("caver|PledgeToken", "s.PledgedTokenPool", s.PledgedTokenPool == nil)
		s.PledgedTokenPool = append(s.PledgedTokenPool, &pledgeToken)
		stateObject.SubBalance(amount)
		stateObject.AddPledgedBalance(amount)
		stateObject.SetPledgedBlockNumber(blocknumber)
	}
	return nil
}

func (s *StateDB) MinerConsign(address common.Address, proxy common.Address) error {
	stateObject := s.GetOrNewStateObject(address)
	log.Info("MinerConsign", "address", address.Hex(), "proxy", proxy.Hex(), "ValidatorPool", len(s.ValidatorPool))
	empty := common.Address{}

	//Only pledged account can consign to an another account
	existAddress := false
	for _, v := range s.ValidatorPool {
		if address.Hex() == v.Addr.Hex() {
			existAddress = true
		}
	}
	if !existAddress {
		log.Info("MinerConsign", "err", "no repeated pledge")
		return errors.New("no repeated pledge")
	}

	//Resolving duplicates is delegated
	for _, v := range s.ValidatorPool {
		if v.Proxy.Hex() != empty.Hex() && v.Proxy.Hex() == proxy.Hex() {
			log.Info("PledgeToken|break", "address", address, "proxy", proxy)
			return errors.New("cannot delegate repeatedly")
		}
	}
	if stateObject != nil {
		pledgeToken := types.PledgedToken{
			Address:      address,
			Amount:       big.NewInt(0),
			Flag:         true,
			ProxyAddress: proxy,
		}
		log.Info("caver|PledgeToken", "s.PledgedTokenPool", s.PledgedTokenPool == nil)
		s.PledgedTokenPool = append(s.PledgedTokenPool, &pledgeToken)
	}
	return nil
}

//- cancel pledged token
//````
//{
//from: holder
//to:0xffff...ffff
//balance:???? amount of recall ERB
//data:{
//version:0
//type:7
//}
//}
//````
//
func (s *StateDB) CancelPledgedToken(address common.Address, amount *big.Int) {
	stateObject := s.GetOrNewStateObject(address)
	if stateObject != nil {
		pledgeToken := types.PledgedToken{
			Address: address,
			Amount:  amount,
			Flag:    false,
		}
		s.PledgedTokenPool = append(s.PledgedTokenPool, &pledgeToken)
		stateObject.SubPledgedBalance(amount)
		stateObject.AddBalance(amount)
	}
}

func (s *StateDB) PunishEvilValidator(address common.Address) {
	amount := s.GetPledgedBalance(address)
	evilObj := s.GetOrNewStateObject(address)
	zeroObj := s.GetOrNewStateObject(common.HexToAddress("0x0000000000000000000000000000000000000000"))
	if evilObj != nil && zeroObj != nil {
		pledgeToken := types.PledgedToken{
			Address: address,
			Amount:  amount,
			Flag:    false,
		}
		s.PledgedTokenPool = append(s.PledgedTokenPool, &pledgeToken)
		evilObj.SubPledgedBalance(amount)
		zeroObj.AddBalance(amount)
	}
}

//- open exchanger:
//````
//{
//from: address which owns the new exchanger
//to:0xffff...ffff
//balance:50ERB
//data:{
//version:0
//type:8
//feeRate:the ratio that exchanger can get
//name:exchanger name
//url:exchanger url
//}
//}
//````
func (s *StateDB) OpenExchanger(addr common.Address,
	amount *big.Int,
	blocknumber *big.Int,
	feerate uint16,
	exchangername string,
	exchangerurl string) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		exchangerToken := &types.PledgedToken{
			Address: addr,
			Amount:  amount,
			Flag:    true,
		}
		s.ExchangerTokenPool = append(s.ExchangerTokenPool, exchangerToken)
		stateObject.SubBalance(amount)
		stateObject.SetExchangerBalance(amount)
		stateObject.OpenExchanger(blocknumber, feerate, exchangername, exchangerurl)
	}
}

func (s *StateDB) ModifyOpenExchangerTime(addr common.Address, blocknumber *big.Int) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SetBlockNumber(blocknumber)
	}
}

func (s *StateDB) CloseExchanger(addr common.Address,
	blocknumber *big.Int) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		if blocknumber.Cmp(stateObject.GetBlockNumber()) > 0 {
			amount := stateObject.ExchangerBalance()
			exchangerToken := &types.PledgedToken{
				Address: addr,
				Amount:  amount,
				Flag:    false,
			}
			s.ExchangerTokenPool = append(s.ExchangerTokenPool, exchangerToken)
			stateObject.AddBalance(amount)
			stateObject.SetExchangerBalance(new(big.Int).SetInt64(0))
			stateObject.CloseExchanger()
		}
	}
}

func (s *StateDB) AddExchangerToken(address common.Address, amount *big.Int) {
	stateObject := s.GetOrNewStateObject(address)
	if stateObject != nil {
		exchangerToken := types.PledgedToken{
			Address: address,
			Amount:  amount,
			Flag:    true,
		}
		s.ExchangerTokenPool = append(s.ExchangerTokenPool, &exchangerToken)
		stateObject.SubBalance(amount)
		stateObject.AddExchangerBalance(amount)
	}
}

func (s *StateDB) SubExchangerToken(address common.Address, amount *big.Int) {
	stateObject := s.GetOrNewStateObject(address)
	if stateObject != nil {
		exchangerToken := types.PledgedToken{
			Address: address,
			Amount:  amount,
			Flag:    false,
		}
		s.ExchangerTokenPool = append(s.ExchangerTokenPool, &exchangerToken)
		stateObject.SubExchangerBalance(amount)
		stateObject.AddBalance(amount)
	}
}

func (s *StateDB) SubExchangerBalance(address common.Address, amount *big.Int) {
	stateObject := s.GetOrNewStateObject(address)
	if stateObject != nil {
		exchangerToken := types.PledgedToken{
			Address: address,
			Amount:  amount,
			Flag:    false,
		}
		s.ExchangerTokenPool = append(s.ExchangerTokenPool, &exchangerToken)
		stateObject.SubExchangerBalance(amount)
	}
}

func (s *StateDB) GetNFTInfo(nftAddr common.Address) (
	string,
	string,
	//*big.Int,
	//uint8,
	common.Address,
	common.Address,
	uint8,
	uint32,
	//bool,
	//*big.Int,
	common.Address,
	uint16,
	common.Address,
	string) {
	stateObject := s.GetOrNewStateObject(nftAddr)
	if stateObject != nil {
		return stateObject.GetNFTInfo()
	}
	return "",
		"",
		//big.NewInt(0),
		//0,
		common.Address{},
		common.Address{},
		0,
		0,
		//false,
		//big.NewInt(0),
		common.Address{},
		0,
		common.Address{},
		""
}

func (s *StateDB) GetPledgedTime(addr common.Address) *big.Int {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		return new(big.Int).Set(stateObject.PledgedBlockNumber())
	}
	return common.Big0
}

func (s *StateDB) GetExchangerFlag(addr common.Address) bool {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		return stateObject.GetExchangerFlag()
	}
	return false
}
func (s *StateDB) GetOpenExchangerTime(addr common.Address) *big.Int {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		return new(big.Int).Set(stateObject.GetBlockNumber())
	}
	return common.Big0
}
func (s *StateDB) GetFeeRate(addr common.Address) uint16 {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		return stateObject.GetFeeRate()
	}
	return 0
}
func (s *StateDB) GetExchangerName(addr common.Address) string {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		return stateObject.GetExchangerName()
	}
	return ""
}
func (s *StateDB) GetExchangerURL(addr common.Address) string {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		return stateObject.GetExchangerURL()
	}
	return ""
}
func (s *StateDB) GetApproveAddress(addr common.Address) []common.Address {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		return stateObject.GetApproveAddress()
	}
	return []common.Address{}
}

//func (s *StateDB) GetNFTBalance(addr common.Address) uint64 {
//	stateObject := s.GetOrNewStateObject(addr)
//	if stateObject != nil {
//		return stateObject.GetNFTBalance()
//	}
//	return 0
//}

func (s *StateDB) GetNFTName(addr common.Address) string {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		return stateObject.GetName()
	}
	return ""
}
func (s *StateDB) GetNFTSymbol(addr common.Address) string {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		return stateObject.GetSymbol()
	}
	return ""
}

//func (s *StateDB) GetNFTApproveAddress(addr common.Address) []common.Address {
//	stateObject := s.getStateObject(addr)
//	if stateObject != nil {
//		return stateObject.GetNFTApproveAddress()
//	}
//	return []common.Address{}
//}
func (s *StateDB) GetNFTApproveAddress(addr common.Address) common.Address {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		return stateObject.GetNFTApproveAddress()
	}
	return common.Address{}
}
func (s *StateDB) GetNFTMergeLevel(addr common.Address) uint8 {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		return stateObject.GetMergeLevel()
	}
	return 0
}
func (s *StateDB) GetNFTCreator(addr common.Address) common.Address {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		return stateObject.GetCreator()
	}
	return common.Address{}
}
func (s *StateDB) GetNFTRoyalty(addr common.Address) uint16 {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		return stateObject.GetRoyalty()
	}
	return 0
}
func (s *StateDB) GetNFTExchanger(addr common.Address) common.Address {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		return stateObject.GetExchanger()
	}
	return common.Address{}
}
func (s *StateDB) GetNFTMetaURL(addr common.Address) string {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		return stateObject.GetMetaURL()
	}
	return ""
}

func (s *StateDB) GetMergeNumber(addr common.Address) uint32 {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		return stateObject.GetMergeNumber()
	}
	return 0
}

//func (s *StateDB) GetPledgedFlag(addr common.Address) bool {
//	stateObject := s.GetOrNewStateObject(addr)
//	if stateObject != nil {
//		return stateObject.GetPledgedFlag()
//	}
//	return false
//}
//
//func (s *StateDB) GetNFTPledgedBlockNumber(addr common.Address) *big.Int {
//	stateObject := s.GetOrNewStateObject(addr)
//	if stateObject != nil {
//		return stateObject.GetNFTPledgedBlockNumber()
//	}
//	return big.NewInt(0)
//}

func (s *StateDB) IsExistNFT(addr common.Address) bool {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		return stateObject.NFTOwner() != common.Address{}
	}
	return false
}

func (s *StateDB) IsApprovedOne(nftAddr common.Address, addr common.Address) bool {
	stateObject := s.GetOrNewStateObject(nftAddr)
	if stateObject != nil {
		return stateObject.data.IsNFTApproveAddress(addr)
	}
	return false
}

func (s *StateDB) IsApprovedForAll(ownerAddr common.Address, addr common.Address) bool {
	stateObject := s.GetOrNewStateObject(ownerAddr)
	if stateObject != nil {
		return stateObject.data.IsApproveAddress(addr)
	}
	return false
}

func (s *StateDB) IsApprovedForAllByNFT(nftAddr common.Address, addr common.Address) bool {
	owner := s.GetNFTOwner16(nftAddr)
	stateObject := s.GetOrNewStateObject(owner)
	if stateObject != nil {
		return stateObject.data.IsApproveAddress(addr)
	}
	return false
}

func (s *StateDB) IsApproved(nftAddr common.Address, addr common.Address) bool {
	if s.IsApprovedOne(nftAddr, addr) || s.IsApprovedForAllByNFT(nftAddr, addr) {
		return true
	}
	return false
}

// GetPledgedBalance retrieves the pledged balance from the given address or 0 if object not found
func (s *StateDB) GetPledgedBalance(addr common.Address) *big.Int {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		pledgedBalance := stateObject.PledgedBalance()
		if pledgedBalance != nil {
			return pledgedBalance
		} else {
			return common.Big0
		}
	}
	return common.Big0
}

func (s *StateDB) GetAccountInfo(addr common.Address) Account {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		return stateObject.GetAccountInfo()
	}
	return Account{}
}

// GetCoefficient retrieves the coefficient from the given address or 0 if object not found
func (s *StateDB) GetCoefficient(addr common.Address) uint8 {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		return stateObject.Coefficient()
	}
	return 0
}

// GetExchangerBalance retrieves the exchanger balance from the given address or 0 if object not found
func (s *StateDB) GetExchangerBalance(addr common.Address) *big.Int {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		exchangerBalnace := stateObject.ExchangerBalance()
		if exchangerBalnace != nil {
			return exchangerBalnace
		} else {
			return common.Big0
		}
	}
	return common.Big0
}

func (s *StateDB) VoteOfficialNFT(nominatedOfficialNFT *types.NominatedOfficialNFT, blocknumber *big.Int) error {
	voteWeight := big.NewInt(0)
	nominatedWeight := big.NewInt(0)
	voteBlockNumber := big.NewInt(0)
	nominatedVoteBlockNumber := big.NewInt(0)
	stateObject := s.GetOrNewStateObject(nominatedOfficialNFT.Address)
	if stateObject != nil {
		voteWeight = stateObject.VoteWeight()
		voteBlockNumber = stateObject.VoteBlockNumber()

	}
	emptyAddress := common.Address{}
	if s.NominatedOfficialNFT != nil && s.NominatedOfficialNFT.Address != emptyAddress {
		nominatedObject := s.GetOrNewStateObject(s.NominatedOfficialNFT.Address)
		if nominatedObject != nil {
			nominatedWeight = nominatedObject.VoteWeight()
			nominatedVoteBlockNumber = nominatedObject.VoteBlockNumber()
		}
	}

	if voteWeight == nil {
		voteWeight = big.NewInt(0)
	}
	if nominatedWeight == nil {
		nominatedWeight = big.NewInt(0)
	}
	if voteBlockNumber == nil {
		voteBlockNumber = big.NewInt(0)
	}
	if nominatedVoteBlockNumber == nil {
		nominatedVoteBlockNumber = big.NewInt(0)
	}

	voteSubNumber := new(big.Int).Sub(blocknumber, voteBlockNumber)
	nominatedSubNumber := new(big.Int).Sub(blocknumber, nominatedVoteBlockNumber)
	voteWeight.Mul(voteWeight, voteSubNumber)
	nominatedWeight.Mul(nominatedWeight, nominatedSubNumber)

	if voteWeight.Cmp(nominatedWeight) > 0 {
		tempNominatedNFT := types.NominatedOfficialNFT{}
		tempNominatedNFT.Address = nominatedOfficialNFT.Address
		tempNominatedNFT.Dir = nominatedOfficialNFT.Dir
		tempNominatedNFT.StartIndex = new(big.Int).Set(nominatedOfficialNFT.StartIndex)
		tempNominatedNFT.Number = nominatedOfficialNFT.Number
		tempNominatedNFT.Royalty = nominatedOfficialNFT.Royalty
		tempNominatedNFT.Creator = nominatedOfficialNFT.Creator
		s.NominatedOfficialNFT = &tempNominatedNFT
		return nil
	}

	return errors.New("voteweight less than previous one")
}

func (s *StateDB) ElectNominatedOfficialNFT(blocknumber *big.Int) {
	emptyAddress := common.Address{}
	if s.NominatedOfficialNFT != nil &&
		s.NominatedOfficialNFT.Address != emptyAddress {
		injectNFT := &types.InjectedOfficialNFT{
			Dir:        s.NominatedOfficialNFT.Dir,
			StartIndex: new(big.Int).Set(s.NominatedOfficialNFT.StartIndex),
			Number:     s.NominatedOfficialNFT.Number,
			Royalty:    s.NominatedOfficialNFT.Royalty,
			Creator:    s.NominatedOfficialNFT.Creator,
			Address:    s.NominatedOfficialNFT.Address,
		}
		voteWeight := s.GetVoteWeight(s.NominatedOfficialNFT.Address)
		voteBlockNumber := s.GetVoteBlockNumber(s.NominatedOfficialNFT.Address)
		subNumber := new(big.Int).Sub(blocknumber, voteBlockNumber)
		injectNFT.VoteWeight = new(big.Int).Mul(voteWeight, subNumber)
		s.OfficialNFTPool.InjectedOfficialNFTs = append(s.OfficialNFTPool.InjectedOfficialNFTs, injectNFT)
		//s.SubVoteWeight(s.NominatedOfficialNFT.Address, voteWeight)
		s.SetVoteBlockNumber(s.NominatedOfficialNFT.Address, blocknumber)

		InjectRewardAddress := common.HexToAddress("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
		injectRewardBalance := s.GetBalance(InjectRewardAddress)
		s.SubBalance(InjectRewardAddress, injectRewardBalance)
		s.AddBalance(s.NominatedOfficialNFT.Address, injectRewardBalance)

		////s.NominatedOfficialNFT = nil
		//s.NominatedOfficialNFT.Dir = types.DefaultDir
		//s.NominatedOfficialNFT.StartIndex = new(big.Int).Set(s.OfficialNFTPool.MaxIndex())
		//s.NominatedOfficialNFT.Number = types.DefaultNumber
		//s.NominatedOfficialNFT.Royalty = types.DefaultRoyalty
		//s.NominatedOfficialNFT.Creator = types.DefaultCreator
		//s.NominatedOfficialNFT.Address = common.Address{}
	} else {
		injectNFT := &types.InjectedOfficialNFT{
			Dir:        types.DefaultDir,
			StartIndex: new(big.Int).Set(s.OfficialNFTPool.MaxIndex()),
			Number:     types.DefaultNumber,
			Royalty:    types.DefaultRoyalty,
			Creator:    types.DefaultCreator,
		}
		s.OfficialNFTPool.InjectedOfficialNFTs = append(s.OfficialNFTPool.InjectedOfficialNFTs, injectNFT)
	}

	s.NominatedOfficialNFT.Dir = types.DefaultDir
	s.NominatedOfficialNFT.StartIndex = new(big.Int).Set(s.OfficialNFTPool.MaxIndex())
	s.NominatedOfficialNFT.Number = types.DefaultNumber
	s.NominatedOfficialNFT.Royalty = types.DefaultRoyalty
	s.NominatedOfficialNFT.Creator = types.DefaultCreator
	s.NominatedOfficialNFT.Address = common.Address{}
}

// AddVoteWeight adds amount to the VoteWeight associated with addr.
func (s *StateDB) AddVoteWeight(addr common.Address, amount *big.Int) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.AddVoteWeight(amount)
	}
}

// SubVoteWeight subtracts amount from the VoteWeight associated with addr.
func (s *StateDB) SubVoteWeight(addr common.Address, amount *big.Int) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SubVoteWeight(amount)
	}
}

// GetVoteWeight retrieves the VoteWeight from the given address or 0 if object not found
func (s *StateDB) GetVoteWeight(addr common.Address) *big.Int {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		return stateObject.VoteWeight()
	}
	return common.Big0
}

func (s *StateDB) SetVoteBlockNumber(addr common.Address, blocknumber *big.Int) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SetVoteBlockNumber(blocknumber)
	}
}

func (s *StateDB) GetVoteBlockNumber(addr common.Address) *big.Int {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		return stateObject.VoteBlockNumber()
	}
	return common.Big0
}

func (s *StateDB) NextIndex() *big.Int {
	return s.OfficialNFTPool.MaxIndex()
}

//func (s *StateDB) ChangeRewardFlag(addr common.Address, flag uint8) {
//	stateObject := s.GetOrNewStateObject(addr)
//	if stateObject != nil {
//		stateObject.ChangeRewardFlag(flag)
//	}
//}

// AddValidatorCoefficient adds amount to the ValidatorCoefficient associated with addr.
func (s *StateDB) AddValidatorCoefficient(addr common.Address, coe uint8) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.AddCoefficient(coe)
	}
}

// SubValidatorCoefficient subtracts amount from the ValidatorCoefficient associated with addr.
func (s *StateDB) SubValidatorCoefficient(addr common.Address, coe uint8) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SubCoefficient(coe)
	}
}

// GetValidatorCoefficient retrieves the ValidatorCoefficient from the given address or 0 if object not found
func (s *StateDB) GetValidatorCoefficient(addr common.Address) uint8 {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		coe := stateObject.Coefficient()
		return coe
	}
	return 0
}
