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

package state

import (
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"io"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/rlp"
)

var emptyCodeHash = crypto.Keccak256(nil)

const VALIDATOR_COEFFICIENT = 70

type Code []byte

func (c Code) String() string {
	return string(c) //strings.Join(Disassemble(c), " ")
}

type Storage map[common.Hash]common.Hash

func (s Storage) String() (str string) {
	for key, value := range s {
		str += fmt.Sprintf("%X : %X\n", key, value)
	}

	return
}

func (s Storage) Copy() Storage {
	cpy := make(Storage)
	for key, value := range s {
		cpy[key] = value
	}

	return cpy
}

// stateObject represents an Ethereum account which is being modified.
//
// The usage pattern is as follows:
// First you need to obtain a state object.
// Account values can be accessed and modified through the object.
// Finally, call CommitTrie to write the modified storage trie into a database.
type stateObject struct {
	address  common.Address
	addrHash common.Hash // hash of ethereum address of the account
	data     Account
	//nft 	AccountNFT
	db *StateDB

	// DB error.
	// State objects are used by the consensus core and VM which are
	// unable to deal with database-level errors. Any error that occurs
	// during a database read is memoized here and will eventually be returned
	// by StateDB.Commit.
	dbErr error

	// Write caches.
	trie Trie // storage trie, which becomes non-nil on first access
	code Code // contract bytecode, which gets set when code is loaded

	originStorage  Storage // Storage cache of original entries to dedup rewrites, reset for every transaction
	pendingStorage Storage // Storage entries that need to be flushed to disk, at the end of an entire block
	dirtyStorage   Storage // Storage entries that have been modified in the current transaction execution
	fakeStorage    Storage // Fake storage which constructed by caller for debugging purpose.

	// Cache flags.
	// When an object is marked suicided it will be delete from the trie
	// during the "update" phase of the state transition.
	dirtyCode bool // true if the code was updated
	suicided  bool
	deleted   bool
}

// empty returns whether the account is considered empty.
func (s *stateObject) empty() bool {
	//return s.data.Nonce == 0 && s.data.Balance.Sign() == 0 && bytes.Equal(s.data.CodeHash, emptyCodeHash)
	return s.data.Nonce == 0 &&
		s.data.Balance.Sign() == 0 &&
		bytes.Equal(s.data.CodeHash, emptyCodeHash) &&
		s.data.Worm == nil &&
		s.data.Nft == nil &&
		s.data.Staker == nil &&
		s.data.Extra == nil
}

// Account is the Ethereum consensus representation of accounts.
// These objects are stored in the main account trie.
type Account struct {
	Nonce    uint64
	Balance  *big.Int
	Root     common.Hash // merkle root of the storage trie
	CodeHash []byte
	Worm     *types.WormholesExtension `rlp:"nil"`
	Nft      *types.AccountNFT         `rlp:"nil"`
	Staker   *types.AccountStaker      `rlp:"nil"`
	Extra    []byte
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

// *** modify to support nft transaction 20211215 begin ***

func (acc *Account) IsApproveAddress(address common.Address) bool {
	for _, addr := range acc.Worm.ApproveAddressList {
		if addr == address {
			return true
		}
	}
	return false
}

//type AccountNFT struct {
//	//Account
//	Name   string
//	Symbol string
//	//Price                 *big.Int
//	//Direction             uint8 // 0:no_tx,1:by,2:sell
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

// *** modify to support nft transaction 20211215 end ***

func (acc *Account) IsNFTApproveAddress(address common.Address) bool {
	//for _, addr := range accNft.NFTApproveAddressList {
	//	if addr == address {
	//		return true
	//	}
	//}

	if address == acc.Nft.NFTApproveAddressList {
		return true
	}

	return false
}

// newObject creates a state object.
func newObject(db *StateDB, address common.Address, data Account) *stateObject {
	if data.Balance == nil {
		data.Balance = new(big.Int)
	}
	if data.CodeHash == nil {
		data.CodeHash = emptyCodeHash
	}
	if data.Root == (common.Hash{}) {
		data.Root = emptyRoot
	}
	newStateObject := &stateObject{
		db:             db,
		address:        address,
		addrHash:       crypto.Keccak256Hash(address[:]),
		originStorage:  make(Storage),
		pendingStorage: make(Storage),
		dirtyStorage:   make(Storage),
	}

	newStateObject.data.Nonce = data.Nonce
	newStateObject.data.Balance = new(big.Int).Set(data.Balance)
	newStateObject.data.CodeHash = make([]byte, len(data.CodeHash))
	copy(newStateObject.data.CodeHash, data.CodeHash)
	newStateObject.data.Root = data.Root

	if data.Worm != nil {
		newStateObject.data.Worm = data.Worm.DeepCopy()
	}

	if data.Nft != nil {
		newStateObject.data.Nft = data.Nft.DeepCopy()
	}

	if data.Staker != nil {
		newStateObject.data.Staker = data.Staker.DeepCopy()
	}

	newStateObject.data.Extra = make([]byte, len(data.Extra))
	copy(newStateObject.data.Extra, data.Extra)

	return newStateObject
}

// EncodeRLP implements rlp.Encoder.
func (s *stateObject) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, s.data)
}

// setError remembers the first non-nil error it is called with.
func (s *stateObject) setError(err error) {
	if s.dbErr == nil {
		s.dbErr = err
	}
}

func (s *stateObject) markSuicided() {
	s.suicided = true
}

func (s *stateObject) touch() {
	s.db.journal.append(touchChange{
		account: &s.address,
	})
	if s.address == ripemd {
		// Explicitly put it in the dirty-cache, which is otherwise generated from
		// flattened journals.
		s.db.journal.dirty(s.address)
	}
}

func (s *stateObject) getTrie(db Database) Trie {
	if s.trie == nil {
		// Try fetching from prefetcher first
		// We don't prefetch empty tries
		if s.data.Root != emptyRoot && s.db.prefetcher != nil {
			// When the miner is creating the pending state, there is no
			// prefetcher
			s.trie = s.db.prefetcher.trie(s.data.Root)
		}
		if s.trie == nil {
			var err error
			s.trie, err = db.OpenStorageTrie(s.addrHash, s.data.Root)
			if err != nil {
				s.trie, _ = db.OpenStorageTrie(s.addrHash, common.Hash{})
				s.setError(fmt.Errorf("can't create storage trie: %v", err))
			}
		}
	}
	return s.trie
}

// GetState retrieves a value from the account storage trie.
func (s *stateObject) GetState(db Database, key common.Hash) common.Hash {
	// If the fake storage is set, only lookup the state here(in the debugging mode)
	if s.fakeStorage != nil {
		return s.fakeStorage[key]
	}
	// If we have a dirty value for this state entry, return it
	value, dirty := s.dirtyStorage[key]
	if dirty {
		return value
	}
	// Otherwise return the entry's original value
	return s.GetCommittedState(db, key)
}

// GetCommittedState retrieves a value from the committed account storage trie.
func (s *stateObject) GetCommittedState(db Database, key common.Hash) common.Hash {
	// If the fake storage is set, only lookup the state here(in the debugging mode)
	if s.fakeStorage != nil {
		return s.fakeStorage[key]
	}
	// If we have a pending write or clean cached, return that
	if value, pending := s.pendingStorage[key]; pending {
		return value
	}
	if value, cached := s.originStorage[key]; cached {
		return value
	}
	// If no live objects are available, attempt to use snapshots
	var (
		enc   []byte
		err   error
		meter *time.Duration
	)
	readStart := time.Now()
	if metrics.EnabledExpensive {
		// If the snap is 'under construction', the first lookup may fail. If that
		// happens, we don't want to double-count the time elapsed. Thus this
		// dance with the metering.
		defer func() {
			if meter != nil {
				*meter += time.Since(readStart)
			}
		}()
	}
	if s.db.snap != nil {
		if metrics.EnabledExpensive {
			meter = &s.db.SnapshotStorageReads
		}
		// If the object was destructed in *this* block (and potentially resurrected),
		// the storage has been cleared out, and we should *not* consult the previous
		// snapshot about any storage values. The only possible alternatives are:
		//   1) resurrect happened, and new slot values were set -- those should
		//      have been handles via pendingStorage above.
		//   2) we don't have new values, and can deliver empty response back
		if _, destructed := s.db.snapDestructs[s.addrHash]; destructed {
			return common.Hash{}
		}
		enc, err = s.db.snap.Storage(s.addrHash, crypto.Keccak256Hash(key.Bytes()))
	}
	// If snapshot unavailable or reading from it failed, load from the database
	if s.db.snap == nil || err != nil {
		if meter != nil {
			// If we already spent time checking the snapshot, account for it
			// and reset the readStart
			*meter += time.Since(readStart)
			readStart = time.Now()
		}
		if metrics.EnabledExpensive {
			meter = &s.db.StorageReads
		}
		if enc, err = s.getTrie(db).TryGet(key.Bytes()); err != nil {
			s.setError(err)
			return common.Hash{}
		}
	}
	var value common.Hash
	if len(enc) > 0 {
		_, content, _, err := rlp.Split(enc)
		if err != nil {
			s.setError(err)
		}
		value.SetBytes(content)
	}
	s.originStorage[key] = value
	return value
}

// SetState updates a value in account storage.
func (s *stateObject) SetState(db Database, key, value common.Hash) {
	// If the fake storage is set, put the temporary state update here.
	if s.fakeStorage != nil {
		s.fakeStorage[key] = value
		return
	}
	// If the new value is the same as old, don't set
	prev := s.GetState(db, key)
	if prev == value {
		return
	}
	// New value is different, update and journal the change
	s.db.journal.append(storageChange{
		account:  &s.address,
		key:      key,
		prevalue: prev,
	})
	s.setState(key, value)
}

// SetStorage replaces the entire state storage with the given one.
//
// After this function is called, all original state will be ignored and state
// lookup only happens in the fake state storage.
//
// Note this function should only be used for debugging purpose.
func (s *stateObject) SetStorage(storage map[common.Hash]common.Hash) {
	// Allocate fake storage if it's nil.
	if s.fakeStorage == nil {
		s.fakeStorage = make(Storage)
	}
	for key, value := range storage {
		s.fakeStorage[key] = value
	}
	// Don't bother journal since this function should only be used for
	// debugging and the `fake` storage won't be committed to database.
}

func (s *stateObject) setState(key, value common.Hash) {
	s.dirtyStorage[key] = value
}

// finalise moves all dirty storage slots into the pending area to be hashed or
// committed later. It is invoked at the end of every transaction.
func (s *stateObject) finalise(prefetch bool) {
	slotsToPrefetch := make([][]byte, 0, len(s.dirtyStorage))
	for key, value := range s.dirtyStorage {
		s.pendingStorage[key] = value
		if value != s.originStorage[key] {
			slotsToPrefetch = append(slotsToPrefetch, common.CopyBytes(key[:])) // Copy needed for closure
		}
	}
	if s.db.prefetcher != nil && prefetch && len(slotsToPrefetch) > 0 && s.data.Root != emptyRoot {
		s.db.prefetcher.prefetch(s.data.Root, slotsToPrefetch)
	}
	if len(s.dirtyStorage) > 0 {
		s.dirtyStorage = make(Storage)
	}
}

// updateTrie writes cached storage modifications into the object's storage trie.
// It will return nil if the trie has not been loaded and no changes have been made
func (s *stateObject) updateTrie(db Database) Trie {
	// Make sure all dirty slots are finalized into the pending storage area
	s.finalise(false) // Don't prefetch any more, pull directly if need be
	if len(s.pendingStorage) == 0 {
		return s.trie
	}
	// Track the amount of time wasted on updating the storage trie
	if metrics.EnabledExpensive {
		defer func(start time.Time) { s.db.StorageUpdates += time.Since(start) }(time.Now())
	}
	// The snapshot storage map for the object
	var storage map[common.Hash][]byte
	// Insert all the pending updates into the trie
	tr := s.getTrie(db)
	hasher := s.db.hasher

	usedStorage := make([][]byte, 0, len(s.pendingStorage))
	for key, value := range s.pendingStorage {
		// Skip noop changes, persist actual changes
		if value == s.originStorage[key] {
			continue
		}
		s.originStorage[key] = value

		var v []byte
		if (value == common.Hash{}) {
			s.setError(tr.TryDelete(key[:]))
		} else {
			// Encoding []byte cannot fail, ok to ignore the error.
			v, _ = rlp.EncodeToBytes(common.TrimLeftZeroes(value[:]))
			s.setError(tr.TryUpdate(key[:], v))
		}
		// If state snapshotting is active, cache the data til commit
		if s.db.snap != nil {
			if storage == nil {
				// Retrieve the old storage map, if available, create a new one otherwise
				if storage = s.db.snapStorage[s.addrHash]; storage == nil {
					storage = make(map[common.Hash][]byte)
					s.db.snapStorage[s.addrHash] = storage
				}
			}
			storage[crypto.HashData(hasher, key[:])] = v // v will be nil if value is 0x00
		}
		usedStorage = append(usedStorage, common.CopyBytes(key[:])) // Copy needed for closure
	}
	if s.db.prefetcher != nil {
		s.db.prefetcher.used(s.data.Root, usedStorage)
	}
	if len(s.pendingStorage) > 0 {
		s.pendingStorage = make(Storage)
	}
	return tr
}

// UpdateRoot sets the trie root to the current root hash of
func (s *stateObject) updateRoot(db Database) {
	// If nothing changed, don't bother with hashing anything
	if s.updateTrie(db) == nil {
		return
	}
	// Track the amount of time wasted on hashing the storage trie
	if metrics.EnabledExpensive {
		defer func(start time.Time) { s.db.StorageHashes += time.Since(start) }(time.Now())
	}
	s.data.Root = s.trie.Hash()
}

// CommitTrie the storage trie of the object to db.
// This updates the trie root.
func (s *stateObject) CommitTrie(db Database) error {
	// If nothing changed, don't bother with hashing anything
	if s.updateTrie(db) == nil {
		return nil
	}
	if s.dbErr != nil {
		return s.dbErr
	}
	// Track the amount of time wasted on committing the storage trie
	if metrics.EnabledExpensive {
		defer func(start time.Time) { s.db.StorageCommits += time.Since(start) }(time.Now())
	}
	root, err := s.trie.Commit(nil)
	if err == nil {
		s.data.Root = root
	}
	return err
}

// AddBalance adds amount to s's balance.
// It is used to add funds to the destination account of a transfer.
func (s *stateObject) AddBalance(amount *big.Int) {
	// EIP161: We must check emptiness for the objects such that the account
	// clearing (0,0,0 objects) can take effect.
	if amount.Sign() == 0 {
		if s.empty() {
			s.touch()
		}
		return
	}
	s.SetBalance(new(big.Int).Add(s.Balance(), amount))
}

// SubBalance removes amount from s's balance.
// It is used to remove funds from the origin account of a transfer.
func (s *stateObject) SubBalance(amount *big.Int) {
	if amount.Sign() == 0 {
		return
	}
	s.SetBalance(new(big.Int).Sub(s.Balance(), amount))
}

func (s *stateObject) SetBalance(amount *big.Int) {
	s.db.journal.append(balanceChange{
		account: &s.address,
		prev:    new(big.Int).Set(s.data.Balance),
	})
	s.setBalance(amount)
}

func (s *stateObject) setBalance(amount *big.Int) {
	s.data.Balance = amount
}

func (s *stateObject) deepCopy(db *StateDB) *stateObject {
	stateObject := newObject(db, s.address, s.data)
	if s.trie != nil {
		stateObject.trie = db.db.CopyTrie(s.trie)
	}
	stateObject.code = s.code
	stateObject.dirtyStorage = s.dirtyStorage.Copy()
	stateObject.originStorage = s.originStorage.Copy()
	stateObject.pendingStorage = s.pendingStorage.Copy()
	stateObject.suicided = s.suicided
	stateObject.dirtyCode = s.dirtyCode
	stateObject.deleted = s.deleted
	return stateObject
}

//
// Attribute accessors
//

// Returns the address of the contract/account
func (s *stateObject) Address() common.Address {
	return s.address
}

// Code returns the contract code associated with this object, if any.
func (s *stateObject) Code(db Database) []byte {
	if s.code != nil {
		return s.code
	}
	if bytes.Equal(s.CodeHash(), emptyCodeHash) {
		return nil
	}
	code, err := db.ContractCode(s.addrHash, common.BytesToHash(s.CodeHash()))
	if err != nil {
		s.setError(fmt.Errorf("can't load code hash %x: %v", s.CodeHash(), err))
	}
	s.code = code
	return code
}

// CodeSize returns the size of the contract code associated with this object,
// or zero if none. This method is an almost mirror of Code, but uses a cache
// inside the database to avoid loading codes seen recently.
func (s *stateObject) CodeSize(db Database) int {
	if s.code != nil {
		return len(s.code)
	}
	if bytes.Equal(s.CodeHash(), emptyCodeHash) {
		return 0
	}
	size, err := db.ContractCodeSize(s.addrHash, common.BytesToHash(s.CodeHash()))
	if err != nil {
		s.setError(fmt.Errorf("can't load code size %x: %v", s.CodeHash(), err))
	}
	return size
}

func (s *stateObject) SetCode(codeHash common.Hash, code []byte) {
	prevcode := s.Code(s.db.db)
	s.db.journal.append(codeChange{
		account:  &s.address,
		prevhash: s.CodeHash(),
		prevcode: prevcode,
	})
	s.setCode(codeHash, code)
}

func (s *stateObject) setCode(codeHash common.Hash, code []byte) {
	s.code = code
	s.data.CodeHash = codeHash[:]
	s.dirtyCode = true
}

func (s *stateObject) SetNonce(nonce uint64) {
	s.db.journal.append(nonceChange{
		account: &s.address,
		prev:    s.data.Nonce,
	})
	s.setNonce(nonce)
}

func (s *stateObject) setNonce(nonce uint64) {
	s.data.Nonce = nonce
}

func (s *stateObject) CodeHash() []byte {
	return s.data.CodeHash
}

func (s *stateObject) Balance() *big.Int {
	return s.data.Balance
}

func (s *stateObject) Nonce() uint64 {
	return s.data.Nonce
}

// Never called, but must be present to allow stateObject to be used
// as a vm.Account interface that also satisfies the vm.ContractRef
// interface. Interfaces are awesome.
func (s *stateObject) Value() *big.Int {
	panic("Value on stateObject should never be called")
}

// *** modify to support nft transaction 20211215 begin ***

// ChangeNFTOwner change nft's owner to newOwner.
func (s *stateObject) ChangeNFTOwner(newOwner common.Address) {
	if s.data.Nft.Owner == newOwner {
		return
	}
	s.SetOwner(newOwner)
	// clear nft's approved address
	s.SetNFTApproveAddress(common.Address{})
}

func (s *stateObject) SetOwner(newOwner common.Address) {
	s.db.journal.append(nftOwnerChange{
		nftAddr:  &s.address,
		oldOwner: s.data.Nft.Owner,
	})
	s.setOwner(newOwner)
}

func (s *stateObject) setOwner(newOwner common.Address) {
	s.data.Nft.Owner = newOwner
}

func (s *stateObject) NFTOwner() common.Address {
	return s.data.Nft.Owner
}

// *** modify to support nft transaction 20211215 end ***

func (s *stateObject) GetNFTMergeLevel() uint8 {
	return s.data.Nft.MergeLevel
}

func (s *stateObject) ChangeApproveAddress(newApproveAddress common.Address) {
	if s.data.IsApproveAddress(newApproveAddress) {
		log.Info("ChangeApproveAddress()", "Is approved", true)
		return
	}
	s.SetApproveAddress(newApproveAddress)
}

func (s *stateObject) SetApproveAddress(newApproveAddress common.Address) {
	change := nftApproveAddressChange{
		nftAddr: &s.address,
	}
	change.oldApproveAddressList = append(change.oldApproveAddressList, s.data.Worm.ApproveAddressList...)
	s.db.journal.append(change)
	s.setApproveAddress(newApproveAddress)
}

func (s *stateObject) setApproveAddress(newApproveAddress common.Address) {
	s.data.Worm.ApproveAddressList = append(s.data.Worm.ApproveAddressList, newApproveAddress)
}
func (s *stateObject) setJournalApproveAddress(approveAddressList []common.Address) {
	if len(s.data.Worm.ApproveAddressList) > 0 {
		s.data.Worm.ApproveAddressList = s.data.Worm.ApproveAddressList[:0]
	}

	s.data.Worm.ApproveAddressList = append(s.data.Worm.ApproveAddressList, approveAddressList...)
}

func (s *stateObject) CancelApproveAddress(approveAddress common.Address) {
	if !s.data.IsApproveAddress(approveAddress) {
		return
	}
	s.RemoveApproveAddress(approveAddress)
}

func (s *stateObject) RemoveApproveAddress(approveAddress common.Address) {
	change := nftApproveAddressChange{
		nftAddr: &s.address,
	}
	change.oldApproveAddressList = append(change.oldApproveAddressList, s.data.Worm.ApproveAddressList...)
	s.db.journal.append(change)
	s.removeApproveAddress(approveAddress)
}

func (s *stateObject) removeApproveAddress(approveAddress common.Address) {
	var index int
	for k, addr := range s.data.Worm.ApproveAddressList {
		if addr == approveAddress {
			index = k
			break
		}
	}

	s.data.Worm.ApproveAddressList = append(s.data.Worm.ApproveAddressList[:index], s.data.Worm.ApproveAddressList[index+1:]...)
}

func (s *stateObject) ChangeNFTApproveAddress(newApproveAddress common.Address) {
	if s.data.IsNFTApproveAddress(newApproveAddress) {
		log.Info("ChangeNFTApproveAddress()", "Is approved", true)
		return
	}
	s.SetNFTApproveAddress(newApproveAddress)
}

func (s *stateObject) SetNFTApproveAddress(newApproveAddress common.Address) {
	changeOne := nftApproveAddressChangeOne{
		nftAddr: &s.address,
	}
	//changeOne.oldNFTApproveAddressList = append(changeOne.oldNFTApproveAddressList, s.data.NFTApproveAddressList...)
	changeOne.oldNFTApproveAddressList = s.data.Nft.NFTApproveAddressList
	s.db.journal.append(changeOne)
	s.setNFTApproveAddress(newApproveAddress)
}

func (s *stateObject) setNFTApproveAddress(newApproveAddress common.Address) {
	//s.data.NFTApproveAddressList = append(s.data.NFTApproveAddressList, newApproveAddress)
	s.data.Nft.NFTApproveAddressList = newApproveAddress
}

func (s *stateObject) setJournalNFTApproveAddress(ApproveAddressList common.Address) {
	//if len(s.data.NFTApproveAddressList) > 0 {
	//	s.data.NFTApproveAddressList = s.data.NFTApproveAddressList[:0]
	//}

	//s.data.NFTApproveAddressList = append(s.data.NFTApproveAddressList, ApproveAddressList...)
	s.data.Nft.NFTApproveAddressList = ApproveAddressList
}

func (s *stateObject) CancelNFTApproveAddress(nftApproveAddress common.Address) {
	if !s.data.IsNFTApproveAddress(nftApproveAddress) {
		return
	}

	s.RemoveNFTApproveAddress(nftApproveAddress)
}

func (s *stateObject) RemoveNFTApproveAddress(nftApproveAddress common.Address) {
	changeOne := nftApproveAddressChangeOne{
		nftAddr: &s.address,
	}
	//changeOne.oldNFTApproveAddressList = append(changeOne.oldNFTApproveAddressList, s.data.NFTApproveAddressList...)
	changeOne.oldNFTApproveAddressList = s.data.Nft.NFTApproveAddressList
	s.db.journal.append(changeOne)

	s.removeNFTApproveAddress(nftApproveAddress)
}

func (s *stateObject) removeNFTApproveAddress(nftApproveAddress common.Address) {
	//var index int
	//for k, addr := range s.data.NFTApproveAddressList {
	//	if addr == nftApproveAddress {
	//		index = k
	//		break
	//	}
	//}
	//
	//s.data.NFTApproveAddressList = append(s.data.NFTApproveAddressList[:index], s.data.NFTApproveAddressList[index+1:]...)
	s.data.Nft.NFTApproveAddressList = common.Address{}
}

func (s *stateObject) OpenExchanger(blocknumber *big.Int,
	feerate uint16,
	exchangername string,
	exchangerurl string,
	agentRecipient common.Address) {
	if s.data.Worm.ExchangerFlag {
		return
	}
	s.SetExchangerInfo(true,
		blocknumber,
		feerate,
		exchangername,
		exchangerurl,
		agentRecipient)
}

func (s *stateObject) CloseExchanger() {
	if !s.data.Worm.ExchangerFlag {
		return
	}
	s.SetExchangerInfo(false, big.NewInt(0), 0, "", "", common.Address{})
}

func (s *stateObject) SetExchangerInfo(exchangerflag bool,
	blocknumber *big.Int,
	feerate uint16,
	exchangername string,
	exchangerurl string,
	agentrecipient common.Address) {
	openExchanger := openExchangerChange{
		address:               &s.address,
		oldExchangerFlag:      s.data.Worm.ExchangerFlag,
		oldFeeRate:            s.data.Worm.FeeRate,
		oldExchangerName:      s.data.Worm.ExchangerName,
		oldExchangerURL:       s.data.Worm.ExchangerURL,
		oldSNFTAgentRecipient: s.data.Worm.SNFTAgentRecipient,
	}
	if s.data.Worm.BlockNumber == nil {
		openExchanger.oldBlockNumber = nil
	} else {
		openExchanger.oldBlockNumber = new(big.Int).Set(s.data.Worm.BlockNumber)
	}
	s.db.journal.append(openExchanger)
	s.setExchangerInfo(exchangerflag, blocknumber, feerate, exchangername, exchangerurl, agentrecipient)
}

func (s *stateObject) setExchangerInfo(exchangerflag bool,
	blocknumber *big.Int,
	feerate uint16,
	exchangername string,
	exchangerurl string,
	agentrecipient common.Address) {
	s.data.Worm.ExchangerFlag = exchangerflag
	s.data.Worm.BlockNumber = blocknumber
	s.data.Worm.FeeRate = feerate
	s.data.Worm.ExchangerName = exchangername
	s.data.Worm.ExchangerURL = exchangerurl
	s.data.Worm.SNFTAgentRecipient = agentrecipient
}

func (s *stateObject) SetExchangerInfoflag(exchangerflag bool) {
	openExchanger := openExchangerChange{
		address:               &s.address,
		oldExchangerFlag:      s.data.Worm.ExchangerFlag,
		oldFeeRate:            s.data.Worm.FeeRate,
		oldExchangerName:      s.data.Worm.ExchangerName,
		oldExchangerURL:       s.data.Worm.ExchangerURL,
		oldSNFTAgentRecipient: s.data.Worm.SNFTAgentRecipient,
	}
	if s.data.Worm.BlockNumber == nil {
		openExchanger.oldBlockNumber = nil
	} else {
		openExchanger.oldBlockNumber = new(big.Int).Set(s.data.Worm.BlockNumber)
	}
	s.db.journal.append(openExchanger)
	s.setExchangerInfoflag(exchangerflag)
}

func (s *stateObject) setExchangerInfoflag(exchangerflag bool) {
	s.data.Worm.ExchangerFlag = exchangerflag
}

func (s *stateObject) StakerPledge(addr common.Address, amount *big.Int, blocknumber *big.Int) {
	newStakers := s.data.Worm.StakerExtension.DeepCopy()
	//newStakers.AddStakerPledge(addr, amount, blocknumber)
	s.SetStakerPledge(newStakers)
}

func (s *stateObject) SetStakerPledge(newStakers *types.StakersExtensionList) {
	s.db.journal.append(stakerExtensionChange{
		account:            &s.address,
		oldStakerExtension: s.data.Worm.StakerExtension})
	s.setStakerPledge(newStakers)
}

func (s *stateObject) setStakerPledge(stakers *types.StakersExtensionList) {
	s.data.Worm.StakerExtension = *stakers
}

func (s *stateObject) CleanNFT() {
	//if s.data.NFTPledgedBlockNumber == nil {
	//	s.data.NFTPledgedBlockNumber = big.NewInt(0)
	//}
	change := nftInfoChange{
		address:        &s.address,
		oldName:        s.data.Nft.Name,
		oldSymbol:      s.data.Nft.Symbol,
		oldOwner:       s.data.Nft.Owner,
		oldMergeLevel:  s.data.Nft.MergeLevel,
		oldMergeNumber: s.data.Nft.MergeNumber,
		//oldPledgedFlag:           s.data.PledgedFlag,
		//oldNFTPledgedBlockNumber: new(big.Int).Set(s.data.NFTPledgedBlockNumber),
		oldCreator:   s.data.Nft.Creator,
		oldRoyalty:   s.data.Nft.Royalty,
		oldExchanger: s.data.Nft.Exchanger,
		oldMetaURL:   s.data.Nft.MetaURL,
	}
	//change.oldNFTApproveAddressList = append(change.oldNFTApproveAddressList, s.data.NFTApproveAddressList...)
	change.oldNFTApproveAddressList = s.data.Nft.NFTApproveAddressList
	s.db.journal.append(change)
	s.cleanNFT()
}

func (s *stateObject) cleanNFT() {
	s.data.Nft.Name = ""
	s.data.Nft.Symbol = ""
	s.data.Nft.Owner = common.Address{}
	//s.data.NFTApproveAddressList = s.data.NFTApproveAddressList[:0]
	s.data.Nft.NFTApproveAddressList = common.Address{}
	// Don't reset MergeLevel, because merging snft need to check this value
	// we use this value to check if snfts are in same layer
	//s.data.MergeLevel = 0
	s.data.Nft.MergeNumber = 0
	//s.data.PledgedFlag = false
	//s.data.NFTPledgedBlockNumber = big.NewInt(0)
	s.data.Nft.Creator = common.Address{}
	s.data.Nft.Royalty = 0
	s.data.Nft.Exchanger = common.Address{}
	s.data.Nft.MetaURL = ""
	s.data.Nft.SNFTRecipient = common.Address{}
}

func (s *stateObject) SetNFTInfo(
	name string,
	symbol string,
	//price *big.Int,
	//direction uint8,
	owner common.Address,
	nftApproveAddress common.Address,
	mergeLevel uint8,
	mergenumber uint32,
	//pledgedflag bool,
	//nftpledgedblocknumber *big.Int,
	creator common.Address,
	royalty uint16,
	exchanger common.Address,
	metaURL string,
	snftRecipient common.Address) {
	//if s.data.NFTPledgedBlockNumber == nil {
	//	s.data.NFTPledgedBlockNumber = big.NewInt(0)
	//}
	change := nftInfoChange{
		address:        &s.address,
		oldName:        s.data.Nft.Name,
		oldSymbol:      s.data.Nft.Symbol,
		oldOwner:       s.data.Nft.Owner,
		oldMergeLevel:  s.data.Nft.MergeLevel,
		oldMergeNumber: s.data.Nft.MergeNumber,
		//oldPledgedFlag:           s.data.PledgedFlag,
		//oldNFTPledgedBlockNumber: new(big.Int).Set(s.data.NFTPledgedBlockNumber),
		oldCreator:       s.data.Nft.Creator,
		oldRoyalty:       s.data.Nft.Royalty,
		oldExchanger:     s.data.Nft.Exchanger,
		oldMetaURL:       s.data.Nft.MetaURL,
		oldSNFTRecipient: s.data.Nft.SNFTRecipient,
	}
	//change.oldNFTApproveAddressList = append(change.oldNFTApproveAddressList, s.data.NFTApproveAddressList...)
	change.oldNFTApproveAddressList = s.data.Nft.NFTApproveAddressList
	s.db.journal.append(change)
	s.setNFTInfo(name,
		symbol,
		//price,
		//direction,
		owner,
		nftApproveAddress,
		mergeLevel,
		mergenumber,
		//pledgedflag,
		//nftpledgedblocknumber,
		creator,
		royalty,
		exchanger,
		metaURL,
		snftRecipient)
}

func (s *stateObject) setNFTInfo(
	name string,
	symbol string,
	//price *big.Int,
	//direction uint8,
	owner common.Address,
	nftApproveAddress common.Address,
	mergeLevel uint8,
	mergenumber uint32,
	//pledgedflag bool,
	//nftpledgedblocknumber *big.Int,
	creator common.Address,
	royalty uint16,
	exchanger common.Address,
	metaURL string,
	snftRecipient common.Address) {

	s.data.Nft.Name = name
	s.data.Nft.Symbol = symbol
	s.data.Nft.Owner = owner
	//s.data.NFTApproveAddressList = append(s.data.NFTApproveAddressList, nftApproveAddress)
	s.data.Nft.NFTApproveAddressList = nftApproveAddress
	s.data.Nft.MergeLevel = mergeLevel
	s.data.Nft.MergeNumber = mergenumber
	//s.data.PledgedFlag = pledgedflag
	//s.data.NFTPledgedBlockNumber = nftpledgedblocknumber
	s.data.Nft.Creator = creator
	s.data.Nft.Royalty = royalty
	s.data.Nft.Exchanger = exchanger
	s.data.Nft.MetaURL = metaURL
	s.data.Nft.SNFTRecipient = snftRecipient

}

func (s *stateObject) setJournalNFTInfo(
	name string,
	symbol string,
	price *big.Int,
	direction uint8,
	owner common.Address,
	nftApproveAddressList common.Address,
	mergeLevel uint8,
	creator common.Address,
	royalty uint16,
	exchanger common.Address,
	metaURL string) {

	s.data.Nft.Name = name
	s.data.Nft.Symbol = symbol
	s.data.Nft.Owner = owner
	//if len(s.data.NFTApproveAddressList) > 0 {
	//	s.data.NFTApproveAddressList = s.data.NFTApproveAddressList[:0]
	//}
	//s.data.NFTApproveAddressList = append(s.data.NFTApproveAddressList, nftApproveAddressList...)
	s.data.Nft.NFTApproveAddressList = nftApproveAddressList
	s.data.Nft.MergeLevel = mergeLevel
	s.data.Nft.Creator = creator
	s.data.Nft.Royalty = royalty
	s.data.Nft.Exchanger = exchanger
	s.data.Nft.MetaURL = metaURL

}

func (s *stateObject) GetNFTInfo() (
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

	return s.data.Nft.Name,
		s.data.Nft.Symbol,
		//s.data.Price,
		//s.data.Direction,
		s.data.Nft.Owner,
		s.data.Nft.NFTApproveAddressList,
		s.data.Nft.MergeLevel,
		s.data.Nft.MergeNumber,
		//s.data.PledgedFlag,
		//s.data.NFTPledgedBlockNumber,
		s.data.Nft.Creator,
		s.data.Nft.Royalty,
		s.data.Nft.Exchanger,
		s.data.Nft.MetaURL

}

func (s *stateObject) GetExchangerFlag() bool {
	return s.data.Worm.ExchangerFlag
}
func (s *stateObject) GetBlockNumber() *big.Int {
	return s.data.Worm.BlockNumber
}
func (s *stateObject) GetFeeRate() uint16 {
	return s.data.Worm.FeeRate
}
func (s *stateObject) GetExchangerName() string {
	return s.data.Worm.ExchangerName
}
func (s *stateObject) GetExchangerURL() string {
	return s.data.Worm.ExchangerURL
}
func (s *stateObject) GetApproveAddress() []common.Address {
	return s.data.Worm.ApproveAddressList
}

//func (s *stateObject) GetNFTBalance() uint64 {
//	return s.data.NFTBalance
//}

func (s *stateObject) GetName() string {
	return s.data.Nft.Name
}
func (s *stateObject) GetSymbol() string {
	return s.data.Nft.Symbol
}

//	func (s *stateObject) GetNFTApproveAddress() []common.Address {
//		return s.data.NFTApproveAddressList
//	}
func (s *stateObject) GetNFTApproveAddress() common.Address {
	return s.data.Nft.NFTApproveAddressList
}
func (s *stateObject) GetMergeLevel() uint8 {
	return s.data.Nft.MergeLevel
}

func (s *stateObject) GetMergeNumber() uint32 {
	return s.data.Nft.MergeNumber
}

//func (s *stateObject) GetPledgedFlag() bool {
//	return s.data.PledgedFlag
//}
//
//func (s *stateObject) GetNFTPledgedBlockNumber() *big.Int {
//	return s.data.NFTPledgedBlockNumber
//}

func (s *stateObject) GetCreator() common.Address {
	return s.data.Nft.Creator
}
func (s *stateObject) GetRoyalty() uint16 {
	return s.data.Nft.Royalty
}
func (s *stateObject) GetExchanger() common.Address {
	return s.data.Nft.Exchanger
}
func (s *stateObject) GetMetaURL() string {
	return s.data.Nft.MetaURL
}

func (s *stateObject) PledgedBalance() *big.Int {
	if s.data.Worm.PledgedBalance == nil {
		return big.NewInt(0)
	}
	return new(big.Int).Set(s.data.Worm.PledgedBalance)
}

func (s *stateObject) PledgedBlockNumber() *big.Int {
	if s.data.Worm.PledgedBlockNumber == nil {
		return big.NewInt(0)
	}
	return new(big.Int).Set(s.data.Worm.PledgedBlockNumber)
}

func (s *stateObject) StakerPledgedBlockNumber(addr common.Address) *big.Int {
	if s.data.Worm.PledgedBlockNumber == nil {
		return big.NewInt(0)
	}
	for _, value := range s.data.Worm.StakerExtension.StakerExtensions {
		if value.Addr == addr {
			return new(big.Int).Set(value.BlockNumber)
		}
	}
	return big.NewInt(0)

}

// AddPledgedBalance adds amount to s's pledged balance.
// It is used to add funds to the destination account of a transfer.
func (s *stateObject) AddPledgedBalance(amount *big.Int) {
	// EIP161: We must check emptiness for the objects such that the account
	// clearing (0,0,0 objects) can take effect.
	if amount.Sign() == 0 {
		return
	}
	pledgedBalance := s.PledgedBalance()
	if pledgedBalance == nil {
		pledgedBalance = big.NewInt(0)
	}
	s.SetPledgedBalance(new(big.Int).Add(pledgedBalance, amount))
}

// SubPledgedBalance removes amount from s's pledged balance.
// It is used to remove funds from the origin account of a transfer.
func (s *stateObject) SubPledgedBalance(amount *big.Int) {
	if amount.Sign() == 0 {
		return
	}
	s.SetPledgedBalance(new(big.Int).Sub(s.PledgedBalance(), amount))
}

func (s *stateObject) ExchangerBalance() *big.Int {
	if s.data.Worm.ExchangerBalance == nil {
		return big.NewInt(0)
	}
	return new(big.Int).Set(s.data.Worm.ExchangerBalance)
}

// AddExchangerBalance adds amount to s's exchanger balance.
// It is used to add funds to the destination account of a transfer.
func (s *stateObject) AddExchangerBalance(amount *big.Int) {
	// EIP161: We must check emptiness for the objects such that the account
	// clearing (0,0,0 objects) can take effect.
	if amount.Sign() == 0 {
		return
	}
	s.SetExchangerBalance(new(big.Int).Add(s.ExchangerBalance(), amount))
}

// SubExchangerBalance removes amount from s's pledged balance.
// It is used to remove funds from the origin account of a transfer.
func (s *stateObject) SubExchangerBalance(amount *big.Int) {
	if amount.Sign() == 0 {
		return
	}
	s.SetExchangerBalance(new(big.Int).Sub(s.ExchangerBalance(), amount))
}

func (s *stateObject) VoteWeight() *big.Int {
	if s.data.Worm.VoteWeight == nil {
		return big.NewInt(0)
	}
	return new(big.Int).Set(s.data.Worm.VoteWeight)
}

func (s *stateObject) Coefficient() uint8 {
	return s.data.Worm.Coefficient
}

func (s *stateObject) SetCoefficient(coe uint8) {
	s.db.journal.append(coefficientChange{
		account: &s.address,
		prev:    s.data.Worm.Coefficient,
	})
	s.setCoefficient(coe)
}

func (s *stateObject) setCoefficient(coe uint8) {
	s.data.Worm.Coefficient = coe
}

//	func (s *stateObject) AddCoefficient(coe uint8) {
//		var sum uint8
//		preSum := s.Coefficient() + coe
//		if preSum <= VALIDATOR_COEFFICIENT {
//			sum = preSum
//		} else {
//			sum = VALIDATOR_COEFFICIENT
//		}
//		s.SetCoefficient(sum)
//	}
func (s *stateObject) AddCoefficient(coe uint8) {
	s.SetCoefficient(VALIDATOR_COEFFICIENT)
}

func (s *stateObject) SubCoefficient(coe uint8) {
	var result uint8

	if s.Coefficient() < coe {
		result = 1
	} else {
		preSub := s.Coefficient() - coe
		if preSub >= 1 {
			result = preSub
		} else {
			result = 1
		}
	}

	s.SetCoefficient(result)
}

// AddVoteWeight adds amount to s's vote weight.
// It is used to add funds to the destination account of a vote.
func (s *stateObject) AddVoteWeight(amount *big.Int) {
	// EIP161: We must check emptiness for the objects such that the account
	// clearing (0,0,0 objects) can take effect.
	if amount.Sign() == 0 {
		return
	}
	if s.data.Worm.VoteWeight == nil {
		s.data.Worm.VoteWeight = big.NewInt(0)
	}
	s.SetVoteWeight(new(big.Int).Add(s.VoteWeight(), amount))
}

// SubVoteWeight removes amount from s's pledged balance.
// It is used to remove funds from the origin account of a vote.
func (s *stateObject) SubVoteWeight(amount *big.Int) {
	if amount.Sign() == 0 {
		return
	}
	log.Info("stateObject.SubVoteWeight", "amount", amount, "VoteWeight", s.VoteWeight())
	s.SetVoteWeight(new(big.Int).Sub(s.VoteWeight(), amount))
}

func (s *stateObject) SetVoteBlockNumber(blocknumber *big.Int) {
	if s.data.Worm.VoteBlockNumber == nil {
		s.data.Worm.VoteBlockNumber = big.NewInt(0)
	}
	s.db.journal.append(voteBlockNumberChange{
		account: &s.address,
		prev:    new(big.Int).Set(s.data.Worm.VoteBlockNumber),
	})
	s.setVoteBlockNumber(new(big.Int).Set(blocknumber))
}

func (s *stateObject) setVoteBlockNumber(blocknumber *big.Int) {
	s.data.Worm.VoteBlockNumber = blocknumber
}

func (s *stateObject) VoteBlockNumber() *big.Int {
	if s.data.Worm.VoteBlockNumber == nil {
		return big.NewInt(0)
	}
	return new(big.Int).Set(s.data.Worm.VoteBlockNumber)
}

func (s *stateObject) SetPledgedBalance(amount *big.Int) {
	if s.data.Worm.PledgedBalance == nil {
		s.data.Worm.PledgedBalance = big.NewInt(0)
	}
	s.db.journal.append(pledgedBalanceChange{
		account: &s.address,
		prev:    new(big.Int).Set(s.data.Worm.PledgedBalance),
	})
	s.setPledgedBalance(amount)
}

func (s *stateObject) setPledgedBalance(amount *big.Int) {
	s.data.Worm.PledgedBalance = amount
}

func (s *stateObject) SetPledgedBlockNumber(blocknumber *big.Int) {
	if s.data.Worm.PledgedBlockNumber == nil {
		s.data.Worm.PledgedBlockNumber = big.NewInt(0)
	}
	s.db.journal.append(pledgedBlockNumberChange{
		account: &s.address,
		prev:    new(big.Int).Set(s.data.Worm.PledgedBlockNumber),
	})
	s.setPledgedBlockNumber(new(big.Int).Set(blocknumber))
}

func (s *stateObject) setPledgedBlockNumber(blocknumber *big.Int) {
	s.data.Worm.PledgedBlockNumber = blocknumber
}

func (s *stateObject) GetAccountInfo() Account {
	return s.data
}

func (s *stateObject) SetExchangerBalance(amount *big.Int) {
	if s.data.Worm.ExchangerBalance == nil {
		s.data.Worm.ExchangerBalance = big.NewInt(0)
	}
	s.db.journal.append(exchangerBalanceChange{
		account: &s.address,
		prev:    new(big.Int).Set(s.data.Worm.ExchangerBalance),
	})
	s.setExchangerBalance(amount)
}

func (s *stateObject) setExchangerBalance(amount *big.Int) {
	s.data.Worm.ExchangerBalance = amount
}

func (s *stateObject) SetBlockNumber(blocknumber *big.Int) {
	if s.data.Worm.BlockNumber == nil {
		s.data.Worm.BlockNumber = big.NewInt(0)
	}
	s.db.journal.append(blockNumberChange{
		account: &s.address,
		prev:    new(big.Int).Set(s.data.Worm.BlockNumber),
	})
	s.setBlockNumber(new(big.Int).Set(blocknumber))
}

func (s *stateObject) setBlockNumber(blocknumber *big.Int) {
	s.data.Worm.BlockNumber = blocknumber
}

func (s *stateObject) SetVoteWeight(amount *big.Int) {
	s.db.journal.append(voteWeightChange{
		account: &s.address,
		prev:    new(big.Int).Set(s.data.Worm.VoteWeight),
	})

	// for test
	if amount.Cmp(big.NewInt(0)) < 0 {
		log.Error("stateObject.SetVoteWeight", "negative amount", amount)
		amount = big.NewInt(0)
	}

	s.setVoteWeight(amount)
}

func (s *stateObject) setVoteWeight(amount *big.Int) {
	s.data.Worm.VoteWeight = amount
}

//func (s *stateObject) ChangeRewardFlag(flag uint8) {
//	s.db.journal.append(RewardFlagChange{
//		account:    &s.address,
//		rewardFlag: flag,
//	})
//	s.setRewardFlag(flag)
//}

//func (s *stateObject) RewardFlag() uint8 {
//	return s.data.RewardFlag
//}
//
//func (s *stateObject) setRewardFlag(flag uint8) {
//	s.data.RewardFlag = flag
//}

//func (s *stateObject) PledgeNFT(blocknumber *big.Int) {
//	s.SetPledgedNFTInfo(true, blocknumber)
//}
//
//func (s *stateObject) CancelPledgedNFT() {
//	s.SetPledgedNFTInfo(false, big.NewInt(0))
//}

//func (s *stateObject) SetPledgedNFTInfo(pledgedflag bool, blocknumber *big.Int) {
//	if s.data.NFTPledgedBlockNumber == nil {
//		s.data.NFTPledgedBlockNumber = big.NewInt(0)
//	}
//	s.db.journal.append(pledgedNFTInfo{
//		account:               &s.address,
//		pledgedFlag:           s.data.PledgedFlag,
//		nftPledgedBlockNumber: new(big.Int).Set(s.data.NFTPledgedBlockNumber),
//	})
//	s.setPledgedNFTInfo(pledgedflag, blocknumber)
//}

//func (s *stateObject) setPledgedNFTInfo(pledgedflag bool, blocknumber *big.Int) {
//	s.data.PledgedFlag = pledgedflag
//	s.data.NFTPledgedBlockNumber = blocknumber
//}

func (s *stateObject) SetExtra(extra []byte) {
	oldExtra := extraChange{
		account: &s.address,
		//prev:    s.data.Extra,
	}
	oldExtra.prev = make([]byte, 0)
	oldExtra.prev = append(oldExtra.prev, s.data.Extra...)
	s.db.journal.append(oldExtra)

	s.setExtra(extra)
}

func (s *stateObject) setExtra(extra []byte) {
	s.data.Extra = extra
}

func (s *stateObject) UserMint() *big.Int {
	if s.data.Staker != nil {
		return s.data.Staker.Mint.UserMint
	}
	return big.NewInt(0)
}

func (s *stateObject) OfficialMint() *big.Int {
	if s.data.Staker != nil {
		return s.data.Staker.Mint.OfficialMint
	}
	return big.NewInt(0)
}

func (s *stateObject) AddUserMint(amount *big.Int) {
	if amount.Sign() == 0 {
		return
	}
	s.SetUserMint(new(big.Int).Add(s.UserMint(), amount))
}

func (s *stateObject) SetUserMint(amount *big.Int) {
	s.db.journal.append(userMintChange{
		account: &s.address,
		prev:    new(big.Int).Set(s.data.Staker.Mint.UserMint),
	})

	s.setUserMint(amount)
}

func (s *stateObject) setUserMint(amount *big.Int) {
	s.data.Staker.Mint.UserMint = amount
}

func (s *stateObject) AddOfficialMint(amount *big.Int) {
	if amount.Sign() == 0 {
		return
	}
	s.SetOfficialMint(new(big.Int).Add(s.OfficialMint(), amount))
}

func (s *stateObject) SetOfficialMint(amount *big.Int) {
	s.db.journal.append(officialMintChange{
		account: &s.address,
		prev:    new(big.Int).Set(s.data.Staker.Mint.OfficialMint),
	})

	s.setOfficialMint(amount)
}

func (s *stateObject) setOfficialMint(amount *big.Int) {
	s.data.Staker.Mint.OfficialMint = amount
}

func (s *stateObject) AddValidator(addr common.Address, balance *big.Int, proxy common.Address) bool {
	newValidators := s.data.Staker.Validators.DeepCopy()
	ok := newValidators.AddValidator(addr, balance, proxy)
	if !ok {
		return false
	}

	s.SetValidators(newValidators)
	return true
}

func (s *stateObject) RemoveValidator(addr common.Address, balance *big.Int) bool {
	newValidators := s.data.Staker.Validators.DeepCopy()
	ok := newValidators.RemoveValidator(addr, balance)
	if !ok {
		return false
	}

	s.SetValidators(newValidators)
	return true
}

func (s *stateObject) SetValidators(varlidators *types.ValidatorList) {
	s.db.journal.append(validatorsChange{
		account:       &s.address,
		oldValidators: s.data.Staker.Validators,
	})

	s.setValidators(varlidators)
}

func (s *stateObject) setValidators(varlidators *types.ValidatorList) {
	s.data.Staker.Validators = *varlidators
}

func (s *stateObject) GetValidators() *types.ValidatorList {
	if s.data.Staker != nil {
		return &s.data.Staker.Validators
	}

	return nil
}

func (s *stateObject) AddStaker(addr common.Address, balance *big.Int) {
	newStakers := s.data.Staker.Stakers.DeepCopy()
	newStakers.AddStaker(addr, balance)
	s.SetStakers(newStakers)
}

func (s *stateObject) RemoveStaker(addr common.Address, balance *big.Int) {
	newStakers := s.data.Staker.Stakers.DeepCopy()
	newStakers.RemoveStaker(addr, balance)
	s.SetStakers(newStakers)
}

func (s *stateObject) SetStakers(stakers *types.StakerList) {
	s.db.journal.append(stakersChange{
		account:    &s.address,
		oldStakers: s.data.Staker.Stakers,
	})

	s.setStakers(stakers)
}

func (s *stateObject) setStakers(stakers *types.StakerList) {
	s.data.Staker.Stakers = *stakers
}

func (s *stateObject) GetStakers() *types.StakerList {
	if s.data.Staker != nil {
		return &s.data.Staker.Stakers
	}

	return nil
}

func (s *stateObject) AddInjectedSnfts(InjectedSnft *types.InjectedOfficialNFT) {
	newSnfts := s.data.Staker.Snfts.DeepCopy()
	newSnfts.InjectedOfficialNFTs = append(newSnfts.InjectedOfficialNFTs, InjectedSnft)
	s.SetSnfts(newSnfts)
}

func (s *stateObject) RemoveInjectSnfts(num *big.Int) {
	newSnfts := s.data.Staker.Snfts.DeepCopy()
	newSnfts.DeleteExpireElem(num)
	s.SetSnfts(newSnfts)
}

func (s *stateObject) SetSnfts(snfts *types.InjectedOfficialNFTList) {
	s.db.journal.append(snftsChange{
		account:  &s.address,
		oldSnfts: s.data.Staker.Snfts,
	})

	s.setSnfts(snfts)
}

func (s *stateObject) setSnfts(snfts *types.InjectedOfficialNFTList) {
	s.data.Staker.Snfts = *snfts
}

func (s *stateObject) GetSnfts() *types.InjectedOfficialNFTList {
	if s.data.Staker != nil {
		return &s.data.Staker.Snfts
	}

	return nil
}

func (s *stateObject) SetNominee(nominee *types.NominatedOfficialNFT) {
	oldNomineeChange := nomineeChange{
		account: &s.address,
	}
	if s.data.Staker.Nominee != nil {
		oldNomineeChange.oldNominee = *s.data.Staker.Nominee
	} else {
		oldNomineeChange.oldNominee = types.NominatedOfficialNFT{}
	}
	s.db.journal.append(oldNomineeChange)

	s.setNominee(nominee)
}

func (s *stateObject) setNominee(nominee *types.NominatedOfficialNFT) {
	s.data.Staker.Nominee = nominee
}

func (s *stateObject) GetNominee() *types.NominatedOfficialNFT {
	if s.data.Staker != nil {
		return s.data.Staker.Nominee
	}

	return nil
}

func (s *stateObject) GetSNFTAgentRecipient() common.Address {
	return s.data.Worm.SNFTAgentRecipient
}

func (s *stateObject) SetSNFTAgentRecipient(recipient common.Address) {
	s.db.journal.append(sNFTAgentRecipientChange{
		account:               &s.address,
		oldSNFTAgentRecipient: s.data.Worm.SNFTAgentRecipient,
	})

	s.setSNFTAgentRecipient(recipient)
}

func (s *stateObject) setSNFTAgentRecipient(recipient common.Address) {
	s.data.Worm.SNFTAgentRecipient = recipient
}

func (s *stateObject) GetSNFTNoMerge() bool {
	return s.data.Worm.SNFTNoMerge
}

func (s *stateObject) SetSNFTNoMerge(flag bool) {
	s.db.journal.append(sNFTNoMergeChange{
		account:        &s.address,
		oldSNFTNoMerge: s.data.Worm.SNFTNoMerge,
	})

	s.setSNFTNoMerge(flag)
}

func (s *stateObject) setSNFTNoMerge(flag bool) {
	s.data.Worm.SNFTNoMerge = flag
}

func (s *stateObject) GetSNFTL3Addrs() []common.Address {
	newSNFTL3Addrs := make([]common.Address, 0)
	newSNFTL3Addrs = append(newSNFTL3Addrs, s.data.Staker.SNFTL3Addrs...)
	return newSNFTL3Addrs
}

func (s *stateObject) AddSNFTL3Addrs(snftAddr common.Address) {
	newSNFTL3Addrs := make([]common.Address, 0)
	newSNFTL3Addrs = append(newSNFTL3Addrs, s.data.Staker.SNFTL3Addrs...)
	newSNFTL3Addrs = append(newSNFTL3Addrs, snftAddr)
	s.SetSNFTL3Addrs(newSNFTL3Addrs)
}

func (s *stateObject) RemoveSNFTL3Addrs(snftAddr common.Address) {
	var index int
	newSNFTL3Addrs := make([]common.Address, 0)
	newSNFTL3Addrs = append(newSNFTL3Addrs, s.data.Staker.SNFTL3Addrs...)
	for i, addr := range newSNFTL3Addrs {
		if addr == snftAddr {
			index = i
			break
		}
	}
	newSNFTL3Addrs = append(newSNFTL3Addrs[:index], newSNFTL3Addrs[index+1:]...)
	s.setSNFTL3Addrs(newSNFTL3Addrs)
}

func (s *stateObject) SetSNFTL3Addrs(snftAddrs []common.Address) {
	oldSNFTL3AddrsChange := sNFTL3AddrsChange{
		account: &s.address,
	}
	oldSNFTL3AddrsChange.oldSNFTL3Addrs = append(oldSNFTL3AddrsChange.oldSNFTL3Addrs, s.data.Staker.SNFTL3Addrs...)
	s.db.journal.append(oldSNFTL3AddrsChange)

	s.setSNFTL3Addrs(snftAddrs)
}

func (s *stateObject) setSNFTL3Addrs(snftAddrs []common.Address) {
	s.data.Staker.SNFTL3Addrs = snftAddrs[:]
}

func (s *stateObject) GetDividendAddrs() []common.Address {
	newDividendAddrs := make([]common.Address, 0)
	newDividendAddrs = append(newDividendAddrs, s.data.Staker.DividendAddrs...)
	return newDividendAddrs
}

func (s *stateObject) AddDividendAddrsOne(snftAddr common.Address) {
	newDividendAddrs := make([]common.Address, 0)
	newDividendAddrs = append(newDividendAddrs, s.data.Staker.DividendAddrs...)
	newDividendAddrs = append(newDividendAddrs, snftAddr)
	s.SetDividendAddrs(newDividendAddrs)
}

func (s *stateObject) AddDividendAddrs(snftAddrs []common.Address) {
	newDividendAddrs := make([]common.Address, 0)
	newDividendAddrs = append(newDividendAddrs, s.data.Staker.DividendAddrs...)
	newDividendAddrs = append(newDividendAddrs, snftAddrs...)
	s.SetDividendAddrs(newDividendAddrs)
}

func (s *stateObject) RemoveDividendAddrsOne(snftAddr common.Address) {
	var index int
	newDividendAddrs := make([]common.Address, 0)
	newDividendAddrs = append(newDividendAddrs, s.data.Staker.DividendAddrs...)
	for i, addr := range newDividendAddrs {
		if addr == snftAddr {
			index = i
			break
		}
	}
	newDividendAddrs = append(newDividendAddrs[:index], newDividendAddrs[index+1:]...)
	s.SetDividendAddrs(newDividendAddrs)
}

func (s *stateObject) RemoveDividendAddrsAll() {
	newDividendAddrs := make([]common.Address, 0)
	s.SetDividendAddrs(newDividendAddrs)
}

func (s *stateObject) SetDividendAddrs(snftAddrs []common.Address) {
	oldDividendAddrsChange := dividendAddrsChange{
		account: &s.address,
	}
	oldDividendAddrsChange.oldDividendAddrs = append(oldDividendAddrsChange.oldDividendAddrs, s.data.Staker.DividendAddrs...)
	s.db.journal.append(oldDividendAddrsChange)

	s.setDividendAddrs(snftAddrs)
}

func (s *stateObject) setDividendAddrs(snftAddrs []common.Address) {
	s.data.Staker.DividendAddrs = snftAddrs[:]
}

func (s *stateObject) GetLockSNFTFlag() bool {
	return s.data.Worm.LockSNFTFlag
}

func (s *stateObject) SetLockSNFTFlag(flag bool) {
	s.db.journal.append(lockSNFTFlagChange{
		account:         &s.address,
		oldLockSNFTFlag: s.data.Worm.LockSNFTFlag,
	})

	s.setLockSNFTFlag(flag)
}

func (s *stateObject) setLockSNFTFlag(flag bool) {
	s.data.Worm.LockSNFTFlag = flag
}
