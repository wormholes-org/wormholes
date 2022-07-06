// Copyright 2018 The go-ethereum Authors
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

package rawdb

import (
	"bytes"
	"encoding/json"
	"github.com/ethereum/go-ethereum/core/types"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
)

// ReadDatabaseVersion retrieves the version number of the database.
func ReadDatabaseVersion(db ethdb.KeyValueReader) *uint64 {
	var version uint64

	enc, _ := db.Get(databaseVersionKey)
	if len(enc) == 0 {
		return nil
	}
	if err := rlp.DecodeBytes(enc, &version); err != nil {
		return nil
	}

	return &version
}

// WriteDatabaseVersion stores the version number of the database
func WriteDatabaseVersion(db ethdb.KeyValueWriter, version uint64) {
	enc, err := rlp.EncodeToBytes(version)
	if err != nil {
		log.Crit("Failed to encode database version", "err", err)
	}
	if err = db.Put(databaseVersionKey, enc); err != nil {
		log.Crit("Failed to store the database version", "err", err)
	}
}

// ReadChainConfig retrieves the consensus settings based on the given genesis hash.
func ReadChainConfig(db ethdb.KeyValueReader, hash common.Hash) *params.ChainConfig {
	data, _ := db.Get(configKey(hash))
	if len(data) == 0 {
		return nil
	}
	var config params.ChainConfig
	if err := json.Unmarshal(data, &config); err != nil {
		log.Error("Invalid chain config JSON", "hash", hash, "err", err)
		return nil
	}
	return &config
}

// WriteChainConfig writes the chain config settings to the database.
func WriteChainConfig(db ethdb.KeyValueWriter, hash common.Hash, cfg *params.ChainConfig) {
	if cfg == nil {
		return
	}
	data, err := json.Marshal(cfg)
	if err != nil {
		log.Crit("Failed to JSON encode chain config", "err", err)
	}
	if err := db.Put(configKey(hash), data); err != nil {
		log.Crit("Failed to store chain config", "err", err)
	}
}

// crashList is a list of unclean-shutdown-markers, for rlp-encoding to the
// database
type crashList struct {
	Discarded uint64   // how many ucs have we deleted
	Recent    []uint64 // unix timestamps of 10 latest unclean shutdowns
}

const crashesToKeep = 10

// PushUncleanShutdownMarker appends a new unclean shutdown marker and returns
// the previous data
// - a list of timestamps
// - a count of how many old unclean-shutdowns have been discarded
func PushUncleanShutdownMarker(db ethdb.KeyValueStore) ([]uint64, uint64, error) {
	var uncleanShutdowns crashList
	// Read old data
	if data, err := db.Get(uncleanShutdownKey); err != nil {
		log.Warn("Error reading unclean shutdown markers", "error", err)
	} else if err := rlp.DecodeBytes(data, &uncleanShutdowns); err != nil {
		return nil, 0, err
	}
	var discarded = uncleanShutdowns.Discarded
	var previous = make([]uint64, len(uncleanShutdowns.Recent))
	copy(previous, uncleanShutdowns.Recent)
	// Add a new (but cap it)
	uncleanShutdowns.Recent = append(uncleanShutdowns.Recent, uint64(time.Now().Unix()))
	if count := len(uncleanShutdowns.Recent); count > crashesToKeep+1 {
		numDel := count - (crashesToKeep + 1)
		uncleanShutdowns.Recent = uncleanShutdowns.Recent[numDel:]
		uncleanShutdowns.Discarded += uint64(numDel)
	}
	// And save it again
	data, _ := rlp.EncodeToBytes(uncleanShutdowns)
	if err := db.Put(uncleanShutdownKey, data); err != nil {
		log.Warn("Failed to write unclean-shutdown marker", "err", err)
		return nil, 0, err
	}
	return previous, discarded, nil
}

// PopUncleanShutdownMarker removes the last unclean shutdown marker
func PopUncleanShutdownMarker(db ethdb.KeyValueStore) {
	var uncleanShutdowns crashList
	// Read old data
	if data, err := db.Get(uncleanShutdownKey); err != nil {
		log.Warn("Error reading unclean shutdown markers", "error", err)
	} else if err := rlp.DecodeBytes(data, &uncleanShutdowns); err != nil {
		log.Error("Error decoding unclean shutdown markers", "error", err) // Should mos def _not_ happen
	}
	if l := len(uncleanShutdowns.Recent); l > 0 {
		uncleanShutdowns.Recent = uncleanShutdowns.Recent[:l-1]
	}
	data, _ := rlp.EncodeToBytes(uncleanShutdowns)
	if err := db.Put(uncleanShutdownKey, data); err != nil {
		log.Warn("Failed to clear unclean-shutdown marker", "err", err)
	}
}

func WriteStakePool(db ethdb.KeyValueWriter, hash common.Hash, number uint64, stakerList *types.StakerList) {
	data, err := rlp.EncodeToBytes(stakerList)
	if err != nil {
		log.Crit("Failed to RLP stakePool", "err", err)
	}

	if err := db.Put(StakePoolKey(number, hash), data); err != nil {
		log.Crit("Failed to store stakePool", "err", err)
	}
}

func ReadStakePool(db ethdb.Reader, hash common.Hash, number uint64) (*types.StakerList, error) {
	data, err := db.Get(StakePoolKey(number, hash))
	if err != nil {
		return nil, err
	}

	stakeList := new(types.StakerList)
	if err := rlp.Decode(bytes.NewReader(data), stakeList); err != nil {
		log.Error("Invalid stakeAddr RLP", "hash", hash, "err", err)
		return nil, err
	}
	return stakeList, nil
}

func WriteValidatorPool(db ethdb.KeyValueWriter, hash common.Hash, number uint64, validatorList *types.ValidatorList) {
	data, err := rlp.EncodeToBytes(validatorList)
	if err != nil {
		log.Crit("Failed to RLP validatorPool", "err", err)
	}

	if err := db.Put(ValidatorPoolKey(number, hash), data); err != nil {
		log.Crit("Failed to store validatorPool", "err", err)
	}
}

func ReadValidatorPool(db ethdb.Reader, hash common.Hash, number uint64) (*types.ValidatorList, error) {
	data, err := db.Get(ValidatorPoolKey(number, hash))
	if err != nil {
		return nil, err
	}

	validatorList := new(types.ValidatorList)
	if err := rlp.Decode(bytes.NewReader(data), validatorList); err != nil {
		log.Error("Invalid validatorAddr RLP", "hash", hash, "err", err)
		return nil, err
	}
	return validatorList, nil
}

func WriteActiveMinersPool(db ethdb.KeyValueWriter, hash common.Hash, number uint64, activeMiners *types.ActiveMinerList) {
	data, err := rlp.EncodeToBytes(activeMiners)
	if err != nil {
		log.Crit("Failed to RLP activeMinersPool", "err", err)
	}

	if err := db.Put(ActiveMinersPoolKey(number, hash), data); err != nil {
		log.Crit("Failed to store activeMinersPool", "err", err)
	}
}

func ReadActiveMinersPool(db ethdb.Reader, hash common.Hash, number uint64) (*types.ActiveMinerList, error) {
	data, err := db.Get(ActiveMinersPoolKey(number, hash))
	if err != nil {
		return nil, err
	}

	activeMiners := new(types.ActiveMinerList)
	if err := rlp.Decode(bytes.NewReader(data), activeMiners); err != nil {
		log.Error("Invalid activeMiners RLP", "hash", hash, "err", err)
		return nil, err
	}
	return activeMiners, nil
}
