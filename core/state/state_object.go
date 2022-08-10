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
	"github.com/ethereum/go-ethereum/log"
	"io"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/rlp"
)

var emptyCodeHash = crypto.Keccak256(nil)

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
		s.data.PledgedBalance == nil &&
		!s.data.ExchangerFlag &&
		s.data.BlockNumber == nil &&
		s.data.FeeRate == 0 &&
		s.data.ExchangerName == "" &&
		s.data.ExchangerURL == "" &&
		len(s.data.ApproveAddressList) == 0 &&
		s.data.Name == "" &&
		s.data.Symbol == "" &&
		bytes.Equal(s.data.Owner.Bytes(), common.Address{}.Bytes()) &&
		len(s.data.NFTApproveAddressList) == 0 &&
		s.data.MergeLevel == 0 &&
		bytes.Equal(s.data.Creator.Bytes(), common.Address{}.Bytes()) &&
		s.data.Royalty == 0 &&
		bytes.Equal(s.data.Exchanger.Bytes(), common.Address{}.Bytes()) &&
		s.data.MetaURL == ""
}

// Account is the Ethereum consensus representation of accounts.
// These objects are stored in the main account trie.
type Account struct {
	Nonce   uint64
	Balance *big.Int
	// *** modify to support nft transaction 20211220 begin ***
	//NFTCount uint64		// number of nft who account have
	// *** modify to support nft transaction 20211220 end ***
	Root               common.Hash // merkle root of the storage trie
	CodeHash           []byte
	PledgedBalance     *big.Int
	PledgedBlockNumber *big.Int
	// *** modify to support nft transaction 20211215 ***
	//Owner common.Address
	// whether the account has a NFT exchanger
	ExchangerFlag    bool
	BlockNumber      *big.Int
	ExchangerBalance *big.Int
	VoteWeight       *big.Int
	// The ratio that exchanger get.
	FeeRate       uint32
	ExchangerName string
	ExchangerURL  string
	// ApproveAddress have the right to handle all nfts of the account
	ApproveAddressList []common.Address
	// NFTBalance is the nft number that the account have
	NFTBalance uint64
	// Indicates the reward method chosen by the miner
	RewardFlag uint8 // 0:SNFT 1:ERB default:0
	AccountNFT
}

// *** modify to support nft transaction 20211215 begin ***

func (acc *Account) IsApproveAddress(address common.Address) bool {
	for _, addr := range acc.ApproveAddressList {
		if addr == address {
			return true
		}
	}
	return false
}

type AccountNFT struct {
	//Account
	Name                  string
	Symbol                string
	Price                 *big.Int
	Direction             uint8 // 0:未交易,1:买入,2:卖出
	Owner                 common.Address
	NFTApproveAddressList common.Address
	//Auctions map[string][]common.Address
	// MergeLevel is the level of NFT merged
	MergeLevel uint8

	Creator   common.Address
	Royalty   uint32
	Exchanger common.Address
	MetaURL   string
}

// *** modify to support nft transaction 20211215 end ***

func (accNft *AccountNFT) IsNFTApproveAddress(address common.Address) bool {
	//for _, addr := range accNft.NFTApproveAddressList {
	//	if addr == address {
	//		return true
	//	}
	//}

	if address == accNft.NFTApproveAddressList {
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
	return &stateObject{
		db:             db,
		address:        address,
		addrHash:       crypto.Keccak256Hash(address[:]),
		data:           data,
		originStorage:  make(Storage),
		pendingStorage: make(Storage),
		dirtyStorage:   make(Storage),
	}
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
	if s.data.Owner == newOwner {
		return
	}
	s.SetOwner(newOwner)
	// clear nft's approved address
	s.SetNFTApproveAddress(common.Address{})
}

func (s *stateObject) SetOwner(newOwner common.Address) {
	s.db.journal.append(nftOwnerChange{
		nftAddr:  &s.address,
		oldOwner: s.data.Owner,
	})
	s.setOwner(newOwner)
}

func (s *stateObject) setOwner(newOwner common.Address) {
	s.data.Owner = newOwner
}

func (s *stateObject) NFTOwner() common.Address {
	return s.data.Owner
}

// *** modify to support nft transaction 20211215 end ***

func (s *stateObject) GetNFTMergeLevel() uint8 {
	return s.data.MergeLevel
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
	change.oldApproveAddressList = append(change.oldApproveAddressList, s.data.ApproveAddressList...)
	s.db.journal.append(change)
	s.setApproveAddress(newApproveAddress)
}

func (s *stateObject) setApproveAddress(newApproveAddress common.Address) {
	s.data.ApproveAddressList = append(s.data.ApproveAddressList, newApproveAddress)
}
func (s *stateObject) setJournalApproveAddress(approveAddressList []common.Address) {
	if len(s.data.ApproveAddressList) > 0 {
		s.data.ApproveAddressList = s.data.ApproveAddressList[:0]
	}

	s.data.ApproveAddressList = append(s.data.ApproveAddressList, approveAddressList...)
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
	change.oldApproveAddressList = append(change.oldApproveAddressList, s.data.ApproveAddressList...)
	s.db.journal.append(change)
	s.removeApproveAddress(approveAddress)
}

func (s *stateObject) removeApproveAddress(approveAddress common.Address) {
	var index int
	for k, addr := range s.data.ApproveAddressList {
		if addr == approveAddress {
			index = k
			break
		}
	}

	s.data.ApproveAddressList = append(s.data.ApproveAddressList[:index], s.data.ApproveAddressList[index+1:]...)
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
	changeOne.oldNFTApproveAddressList = s.data.NFTApproveAddressList
	s.db.journal.append(changeOne)
	s.setNFTApproveAddress(newApproveAddress)
}

func (s *stateObject) setNFTApproveAddress(newApproveAddress common.Address) {
	//s.data.NFTApproveAddressList = append(s.data.NFTApproveAddressList, newApproveAddress)
	s.data.NFTApproveAddressList = newApproveAddress
}

func (s *stateObject) setJournalNFTApproveAddress(ApproveAddressList common.Address) {
	//if len(s.data.NFTApproveAddressList) > 0 {
	//	s.data.NFTApproveAddressList = s.data.NFTApproveAddressList[:0]
	//}

	//s.data.NFTApproveAddressList = append(s.data.NFTApproveAddressList, ApproveAddressList...)
	s.data.NFTApproveAddressList = ApproveAddressList
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
	changeOne.oldNFTApproveAddressList = s.data.NFTApproveAddressList
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
	s.data.NFTApproveAddressList = common.Address{}
}

func (s *stateObject) OpenExchanger(blocknumber *big.Int, feerate uint32, exchangername string, exchangerurl string) {
	if s.data.ExchangerFlag {
		return
	}
	s.SetExchangerInfo(true, blocknumber, feerate, exchangername, exchangerurl)
}

func (s *stateObject) CloseExchanger() {
	if !s.data.ExchangerFlag {
		return
	}
	s.SetExchangerInfo(false, big.NewInt(0), 0, "", "")
}

func (s *stateObject) SetExchangerInfo(exchangerflag bool, blocknumber *big.Int, feerate uint32, exchangername string, exchangerurl string) {
	openExchanger := openExchangerChange{
		address:          &s.address,
		oldExchangerFlag: s.data.ExchangerFlag,
		oldFeeRate:       s.data.FeeRate,
		oldExchangerName: s.data.ExchangerName,
		oldExchangerURL:  s.data.ExchangerURL,
	}
	if s.data.BlockNumber == nil {
		openExchanger.oldBlockNumber = nil
	} else {
		openExchanger.oldBlockNumber = new(big.Int).Set(s.data.BlockNumber)
	}
	s.db.journal.append(openExchanger)
	s.setExchangerInfo(exchangerflag, blocknumber, feerate, exchangername, exchangerurl)
}

func (s *stateObject) setExchangerInfo(exchangerflag bool, blocknumber *big.Int, feerate uint32, exchangername string, exchangerurl string) {
	s.data.ExchangerFlag = exchangerflag
	s.data.BlockNumber = blocknumber
	s.data.FeeRate = feerate
	s.data.ExchangerName = exchangername
	s.data.ExchangerURL = exchangerurl
}

func (s *stateObject) CleanNFT() {
	change := nftInfoChange{
		address:       &s.address,
		oldName:       s.data.Name,
		oldSymbol:     s.data.Symbol,
		oldOwner:      s.data.Owner,
		oldMergeLevel: s.data.MergeLevel,
		oldCreator:    s.data.Creator,
		oldRoyalty:    s.data.Royalty,
		oldExchanger:  s.data.Exchanger,
		oldMetaURL:    s.data.MetaURL,
	}
	//change.oldNFTApproveAddressList = append(change.oldNFTApproveAddressList, s.data.NFTApproveAddressList...)
	change.oldNFTApproveAddressList = s.data.NFTApproveAddressList
	s.db.journal.append(change)
	s.cleanNFT()
}

func (s *stateObject) cleanNFT() {
	s.data.Name = ""
	s.data.Symbol = ""
	s.data.Owner = common.Address{}
	//s.data.NFTApproveAddressList = s.data.NFTApproveAddressList[:0]
	s.data.NFTApproveAddressList = common.Address{}
	s.data.MergeLevel = 0
	s.data.Creator = common.Address{}
	s.data.Royalty = 0
	s.data.Exchanger = common.Address{}
	s.data.MetaURL = ""
}

func (s *stateObject) SetNFTInfo(
	name string,
	symbol string,
	price *big.Int,
	direction uint8,
	owner common.Address,
	nftApproveAddress common.Address,
	mergeLevel uint8,
	creator common.Address,
	royalty uint32,
	exchanger common.Address,
	metaURL string) {
	change := nftInfoChange{
		address:       &s.address,
		oldName:       s.data.Name,
		oldSymbol:     s.data.Symbol,
		oldOwner:      s.data.Owner,
		oldMergeLevel: s.data.MergeLevel,
		oldCreator:    s.data.Creator,
		oldRoyalty:    s.data.Royalty,
		oldExchanger:  s.data.Exchanger,
		oldMetaURL:    s.data.MetaURL,
	}
	//change.oldNFTApproveAddressList = append(change.oldNFTApproveAddressList, s.data.NFTApproveAddressList...)
	change.oldNFTApproveAddressList = s.data.NFTApproveAddressList
	s.db.journal.append(change)
	s.setNFTInfo(name,
		symbol,
		price,
		direction,
		owner,
		nftApproveAddress,
		mergeLevel,
		creator,
		royalty,
		exchanger,
		metaURL)
}

func (s *stateObject) setNFTInfo(
	name string,
	symbol string,
	price *big.Int,
	direction uint8,
	owner common.Address,
	nftApproveAddress common.Address,
	mergeLevel uint8,
	creator common.Address,
	royalty uint32,
	exchanger common.Address,
	metaURL string) {

	s.data.Name = name
	s.data.Symbol = symbol
	s.data.Owner = owner
	//s.data.NFTApproveAddressList = append(s.data.NFTApproveAddressList, nftApproveAddress)
	s.data.NFTApproveAddressList = nftApproveAddress
	s.data.MergeLevel = mergeLevel
	s.data.Creator = creator
	s.data.Royalty = royalty
	s.data.Exchanger = exchanger
	s.data.MetaURL = metaURL

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
	royalty uint32,
	exchanger common.Address,
	metaURL string) {

	s.data.Name = name
	s.data.Symbol = symbol
	s.data.Owner = owner
	//if len(s.data.NFTApproveAddressList) > 0 {
	//	s.data.NFTApproveAddressList = s.data.NFTApproveAddressList[:0]
	//}
	//s.data.NFTApproveAddressList = append(s.data.NFTApproveAddressList, nftApproveAddressList...)
	s.data.NFTApproveAddressList = nftApproveAddressList
	s.data.MergeLevel = mergeLevel
	s.data.Creator = creator
	s.data.Royalty = royalty
	s.data.Exchanger = exchanger
	s.data.MetaURL = metaURL

}

func (s *stateObject) GetNFTInfo() (
	string,
	string,
	*big.Int,
	uint8,
	common.Address,
	common.Address,
	uint8,
	common.Address,
	uint32,
	common.Address,
	string) {

	return s.data.Name,
		s.data.Symbol,
		s.data.Price,
		s.data.Direction,
		s.data.Owner,
		s.data.NFTApproveAddressList,
		s.data.MergeLevel,
		s.data.Creator,
		s.data.Royalty,
		s.data.Exchanger,
		s.data.MetaURL

}

func (s *stateObject) GetExchangerFlag() bool {
	return s.data.ExchangerFlag
}
func (s *stateObject) GetBlockNumber() *big.Int {
	return s.data.BlockNumber
}
func (s *stateObject) GetFeeRate() uint32 {
	return s.data.FeeRate
}
func (s *stateObject) GetExchangerName() string {
	return s.data.ExchangerName
}
func (s *stateObject) GetExchangerURL() string {
	return s.data.ExchangerURL
}
func (s *stateObject) GetApproveAddress() []common.Address {
	return s.data.ApproveAddressList
}
func (s *stateObject) GetNFTBalance() uint64 {
	return s.data.NFTBalance
}

func (s *stateObject) GetName() string {
	return s.data.Name
}
func (s *stateObject) GetSymbol() string {
	return s.data.Symbol
}

//func (s *stateObject) GetNFTApproveAddress() []common.Address {
//	return s.data.NFTApproveAddressList
//}
func (s *stateObject) GetNFTApproveAddress() common.Address {
	return s.data.NFTApproveAddressList
}
func (s *stateObject) GetMergeLevel() uint8 {
	return s.data.MergeLevel
}
func (s *stateObject) GetCreator() common.Address {
	return s.data.Creator
}
func (s *stateObject) GetRoyalty() uint32 {
	return s.data.Royalty
}
func (s *stateObject) GetExchanger() common.Address {
	return s.data.Exchanger
}
func (s *stateObject) GetMetaURL() string {
	return s.data.MetaURL
}

func (s *stateObject) PledgedBalance() *big.Int {
	return s.data.PledgedBalance
}

func (s *stateObject) PledgedBlockNumber() *big.Int {
	return s.data.PledgedBlockNumber
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
	return s.data.ExchangerBalance
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
	return s.data.VoteWeight
}

// AddVoteWeight adds amount to s's vote weight.
// It is used to add funds to the destination account of a vote.
func (s *stateObject) AddVoteWeight(amount *big.Int) {
	// EIP161: We must check emptiness for the objects such that the account
	// clearing (0,0,0 objects) can take effect.
	if amount.Sign() == 0 {
		return
	}
	if s.data.VoteWeight == nil {
		s.data.VoteWeight = big.NewInt(0)
	}
	s.SetVoteWeight(new(big.Int).Add(s.VoteWeight(), amount))
}

// SubVoteWeight removes amount from s's pledged balance.
// It is used to remove funds from the origin account of a vote.
func (s *stateObject) SubVoteWeight(amount *big.Int) {
	if amount.Sign() == 0 {
		return
	}
	s.SetVoteWeight(new(big.Int).Sub(s.VoteWeight(), amount))
}

func (s *stateObject) SetPledgedBalance(amount *big.Int) {
	if s.data.PledgedBalance == nil {
		s.data.PledgedBalance = big.NewInt(0)
	}
	s.db.journal.append(pledgedBalanceChange{
		account: &s.address,
		prev:    new(big.Int).Set(s.data.PledgedBalance),
	})
	s.setPledgedBalance(amount)
}

func (s *stateObject) setPledgedBalance(amount *big.Int) {
	s.data.PledgedBalance = amount
}

func (s *stateObject) SetPledgedBlockNumber(blocknumber *big.Int) {
	if s.data.PledgedBlockNumber == nil {
		s.data.PledgedBlockNumber = big.NewInt(0)
	}
	s.db.journal.append(pledgedBlockNumberChange{
		account: &s.address,
		prev:    new(big.Int).Set(s.data.PledgedBlockNumber),
	})
	s.setPledgedBlockNumber(blocknumber)
}

func (s *stateObject) setPledgedBlockNumber(blocknumber *big.Int) {
	s.data.PledgedBlockNumber = blocknumber
}

func (s *stateObject) GetAccountInfo() Account {
	return s.data
}

func (s *stateObject) SetExchangerBalance(amount *big.Int) {
	if s.data.ExchangerBalance == nil {
		s.data.ExchangerBalance = big.NewInt(0)
	}
	s.db.journal.append(exchangerBalanceChange{
		account: &s.address,
		prev:    new(big.Int).Set(s.data.ExchangerBalance),
	})
	s.setExchangerBalance(amount)
}

func (s *stateObject) setExchangerBalance(amount *big.Int) {
	s.data.ExchangerBalance = amount
}

func (s *stateObject) SetVoteWeight(amount *big.Int) {
	s.db.journal.append(voteWeightChange{
		account: &s.address,
		prev:    new(big.Int).Set(s.data.VoteWeight),
	})
	s.setVoteWeight(amount)
}

func (s *stateObject) setVoteWeight(amount *big.Int) {
	s.data.VoteWeight = amount
}

func (s *stateObject) ChangeRewardFlag(flag uint8) {
	s.db.journal.append(RewardFlagChange{
		account:    &s.address,
		rewardFlag: flag,
	})
	s.setRewardFlag(flag)
}

func (s *stateObject) RewardFlag() uint8 {
	return s.data.RewardFlag
}

func (s *stateObject) setRewardFlag(flag uint8) {
	s.data.RewardFlag = flag
}
