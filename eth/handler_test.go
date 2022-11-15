// Copyright 2015 The go-ethereum Authors
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

package eth

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/downloader"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/params"
)

var (
	// testKey is a private key to use for funding a tester account.
	testKey, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")

	// testAddr is the Ethereum address of the tester account.
	testAddr = crypto.PubkeyToAddress(testKey.PublicKey)
)

const SimAllocData = "" +
	"0x091DBBa95B26793515cc9aCB9bEb5124c479f27F:0xd3c21bcecceda1000000," +
	"0x107837Ea83f8f06533DDd3fC39451Cd0AA8DA8BD:0xed2b525841adfc00000," +
	"0x612DFa56DcA1F581Ed34b9c60Da86f1268Ab6349:0xd3c21bcecceda1000000," +
	"0x84d84e6073A06B6e784241a9B13aA824AB455326:0xed2b525841adfc00000," +
	"0x9e4d5C72569465270232ed7Af71981Ee82d08dBF:0xd3c21bcecceda1000000," +
	"0xa270bBDFf450EbbC2d0413026De5545864a1b6d6:0xed2b525841adfc00000," +
	"0x4110E56ED25e21267FBeEf79244f47ada4e2E963:0xd3c21bcecceda1000000," +
	"0xdb33217fE3F74bD41c550B06B624E23ab7f55d05:0xed2b525841adfc00000," +
	"0xE2FA892CC5CC268a0cC1d924EC907C796351C645:0xd3c21bcecceda1000000," +
	"0x52EAE6D396E82358D703BEDeC2ab11E723127230:0xed2b525841adfc00000," +
	"0x31534d5C7b1eabb73425c2361661b878F4322f9D:0xd3c21bcecceda1000000," +
	"0xbbaE84E9879F908700c6ef5D15e928Abfb556a21:0xed2b525841adfc00000," +
	"0x20cb28AE861c322A9A86b4F9e36Ad6977930fA05:0xd3c21bcecceda1000000," +
	"0xFfAc4cd934f026dcAF0f9d9EEDDcD9af85D8943e:0xed2b525841adfc00000," +
	"0xc067825f4B7a53Bb9f2Daf72fF22C8EE39736afF:0xd3c21bcecceda1000000," +
	"0x7bf72621Dd7C4Fe4AF77632e3177c08F53fdAF09:0xed2b525841adfc00000," +
	"0x68B14e0F18C3EE322d3e613fF63B87E56D86Df60:0xf2dc7d47f15600000," +
	"0xeEF79493F62dA884389312d16669455A7E0045c1:0xf2dc7d47f15600000," +
	"0xa5999Cc1DEC36a632dF735064Dc75eF6af0E7389:0xf2dc7d47f15600000," +
	"0x63d913dfDB75C7B09a1465Fe77B8Ec167793096b:0xf2dc7d47f15600000," +
	"0xF50f73B83721c108E8868C5A2706c5b194A0FDB1:0xf2dc7d47f15600000," +
	"0xB000811Aff6e891f8c0F1aa07f43C1976D4c3076:0xf2dc7d47f15600000," +
	"0x257F3c6749a0690d39c6FBCd2DceB3fB464f0F94:0xf2dc7d47f15600000," +
	"0x43ee5cB067F29B920CC44d5d5367BCEb162B4d9E:0xf2dc7d47f15600000," +
	"0x85D3FDA364564c365870233E5aD6B611F2227846:0xf2dc7d47f15600000," +
	"0xDc807D83d864490C6EEDAC9C9C071E9AAeD8E7d7:0xf2dc7d47f15600000," +
	"0xFA623BCC71BE5C3aBacfe875E64ef97F91B7b110:0xf2dc7d47f15600000," +
	"0xb17fAe1710f80Eb9a39732862B0058077F338B21:0xf2dc7d47f15600000," +
	"0x86FFd3e5a6D310Fcb4668582eA6d0cfC1c35da49:0xf2dc7d47f15600000," +
	"0x6bB0599bC9c5406d405a8a797F8849dB463462D0:0xf2dc7d47f15600000," +
	"0x765C83dbA2712582C5461b2145f054d4F85a3080:0xf2dc7d47f15600000," +
	"0x8520dc57A2800e417696bdF93553E63bCF31e597:0x0," +
	"0x7BE4A4a66BBf205aa55d1F18f05489f4b34c2A2D:0x0," +
	"0x66f9e46b49EDDc40F0dA18D67C07ae755b3643CE:0x0," +
	"0x4cc4F114639e22c5205a2716A673B8b625ab58fA:0x0," +
	"0x96f2A9f08c92c174700A0bdb452EA737633382A0:0x0," +
	"0xDa8dcd42B942eF5C4F701BD574AaB5eF420DD5a3:0x0," +
	"0x3E6a45b12E2A4E25fb0176c7Aa1855459E8e862b:0x0," +
	"0x84610975eBF0882600c71d594a74ebf89a972DBe:0x0," +
	"0x4d0A8127D3120684CC70eC12e6E8F44eE990b5aC:0x0," +
	"0xc2617DeFBB00aC0ACd67EA96C2D6E334cAb22ce1:0x0," +
	"0x2DbdaCc91fd967E2A5C3F04D321752d99a7741C8:0x0," +
	"0x45D706D1B80F65F262eEA09Ae9557381aFD7dfA1:0x0," +
	"0x36c1550F16c43B5Dd85f1379E708d89DA9789d9b:0x0," +
	"0x7ddcC65DABAd6b5e1c54d0784b4290f704259E34:0x0," +
	"0xbad3F0edd751B3b8DeF4AaDDbcF5533eC93452C2:0x0," +
	"0x969b308f059cACdD7a3aEb232584fb09fbCBb3A1:0x0," +
	"0x4854F8324009AFDC20C5f651D70fFA5eF6c036B8:0x2a8bf44e200bfb75600000," +
	"0xBa7B3387D88Bd7675DE8B492a9067dc6B7A59311:0x9ed194db19b238c000000"

const SimStakeData = "" +
	"0x68B14e0F18C3EE322d3e613fF63B87E56D86Df60:0xf2dc7d47f15600000:250:exchanger:www.wormholesexchanger.com," +
	"0xeEF79493F62dA884389312d16669455A7E0045c1:0xf2dc7d47f15600000:250:exchanger:www.wormholesexchanger.com," +
	"0xa5999Cc1DEC36a632dF735064Dc75eF6af0E7389:0xf2dc7d47f15600000:250:exchanger:www.wormholesexchanger.com," +
	"0x63d913dfDB75C7B09a1465Fe77B8Ec167793096b:0xf2dc7d47f15600000:250:exchanger:www.wormholesexchanger.com," +
	"0xF50f73B83721c108E8868C5A2706c5b194A0FDB1:0xf2dc7d47f15600000:250:exchanger:www.wormholesexchanger.com," +
	"0xB000811Aff6e891f8c0F1aa07f43C1976D4c3076:0xf2dc7d47f15600000:250:exchanger:www.wormholesexchanger.com," +
	"0x257F3c6749a0690d39c6FBCd2DceB3fB464f0F94:0xf2dc7d47f15600000:250:exchanger:www.wormholesexchanger.com," +
	"0x43ee5cB067F29B920CC44d5d5367BCEb162B4d9E:0xf2dc7d47f15600000:250:exchanger:www.wormholesexchanger.com," +
	"0x85D3FDA364564c365870233E5aD6B611F2227846:0xf2dc7d47f15600000:250:exchanger:www.wormholesexchanger.com," +
	"0xDc807D83d864490C6EEDAC9C9C071E9AAeD8E7d7:0xf2dc7d47f15600000:250:exchanger:www.wormholesexchanger.com," +
	"0xFA623BCC71BE5C3aBacfe875E64ef97F91B7b110:0xf2dc7d47f15600000:250:exchanger:www.wormholesexchanger.com," +
	"0xb17fAe1710f80Eb9a39732862B0058077F338B21:0xf2dc7d47f15600000:250:exchanger:www.wormholesexchanger.com," +
	"0x86FFd3e5a6D310Fcb4668582eA6d0cfC1c35da49:0xf2dc7d47f15600000:250:exchanger:www.wormholesexchanger.com," +
	"0x6bB0599bC9c5406d405a8a797F8849dB463462D0:0xf2dc7d47f15600000:250:exchanger:www.wormholesexchanger.com," +
	"0x765C83dbA2712582C5461b2145f054d4F85a3080:0xf2dc7d47f15600000:250:exchanger:www.wormholesexchanger.com"

const SimValidatorData_v2 = "" +
	"0x091DBBa95B26793515cc9aCB9bEb5124c479f27F:0xd3c21bcecceda1000000:0x8520dc57A2800e417696bdF93553E63bCF31e597," +
	"0x107837Ea83f8f06533DDd3fC39451Cd0AA8DA8BD:0xed2b525841adfc00000:0x0000000000000000000000000000000000000000," +
	"0x612DFa56DcA1F581Ed34b9c60Da86f1268Ab6349:0xd3c21bcecceda1000000:0x66f9e46b49EDDc40F0dA18D67C07ae755b3643CE," +
	"0x84d84e6073A06B6e784241a9B13aA824AB455326:0xed2b525841adfc00000:0x0000000000000000000000000000000000000000," +
	"0x9e4d5C72569465270232ed7Af71981Ee82d08dBF:0xd3c21bcecceda1000000:0x96f2A9f08c92c174700A0bdb452EA737633382A0," +
	"0xa270bBDFf450EbbC2d0413026De5545864a1b6d6:0xed2b525841adfc00000:0x0000000000000000000000000000000000000000," +
	"0x4110E56ED25e21267FBeEf79244f47ada4e2E963:0xd3c21bcecceda1000000:0x3E6a45b12E2A4E25fb0176c7Aa1855459E8e862b," +
	"0xdb33217fE3F74bD41c550B06B624E23ab7f55d05:0xed2b525841adfc00000:0x0000000000000000000000000000000000000000," +
	"0xE2FA892CC5CC268a0cC1d924EC907C796351C645:0xd3c21bcecceda1000000:0x4d0A8127D3120684CC70eC12e6E8F44eE990b5aC," +
	"0x52EAE6D396E82358D703BEDeC2ab11E723127230:0xed2b525841adfc00000:0x0000000000000000000000000000000000000000," +
	"0x31534d5C7b1eabb73425c2361661b878F4322f9D:0xd3c21bcecceda1000000:0x2DbdaCc91fd967E2A5C3F04D321752d99a7741C8," +
	"0xbbaE84E9879F908700c6ef5D15e928Abfb556a21:0xed2b525841adfc00000:0x0000000000000000000000000000000000000000," +
	"0x20cb28AE861c322A9A86b4F9e36Ad6977930fA05:0xd3c21bcecceda1000000:0x36c1550F16c43B5Dd85f1379E708d89DA9789d9b," +
	"0xFfAc4cd934f026dcAF0f9d9EEDDcD9af85D8943e:0xed2b525841adfc00000:0x0000000000000000000000000000000000000000," +
	"0xc067825f4B7a53Bb9f2Daf72fF22C8EE39736afF:0xd3c21bcecceda1000000:0xbad3F0edd751B3b8DeF4AaDDbcF5533eC93452C2," +
	"0x7bf72621Dd7C4Fe4AF77632e3177c08F53fdAF09:0xed2b525841adfc00000:0x0000000000000000000000000000000000000000"

func DecodePreWormholesInfo(data string) core.GenesisAlloc {
	ga := make(core.GenesisAlloc)

	accountInfos := strings.Split(data, ",")
	for _, accountInfo := range accountInfos {
		index := strings.Index(accountInfo, ":")
		if index > 0 {
			acc := string([]byte(accountInfo)[:index])
			balance := string([]byte(accountInfo)[index+1:])
			bigBalance := big.NewInt(0)
			if strings.HasPrefix(balance, "0x") ||
				strings.HasPrefix(balance, "0X") {
				balance = string([]byte(balance)[2:])
				bigBalance, _ = new(big.Int).SetString(balance, 16)
			} else {
				bigBalance, _ = new(big.Int).SetString(balance, 16)
			}

			genesisAcc := core.GenesisAccount{
				Balance: bigBalance,
			}
			ga[common.HexToAddress(acc)] = genesisAcc
		}
	}
	return ga
}

func DecodePreWormholesInfoV2(data string) core.GenesisAlloc {
	ga := make(core.GenesisAlloc)

	accountInfos := strings.Split(data, ",")
	for _, accountInfo := range accountInfos {
		strs := strings.Split(accountInfo, ":")
		if len(strs) > 0 {
			acc := strs[0]
			balance := strs[1]
			proxy := strs[2]
			bigBalance := big.NewInt(0)
			if strings.HasPrefix(balance, "0x") ||
				strings.HasPrefix(balance, "0X") {
				balance = string([]byte(balance)[2:])
				bigBalance, _ = new(big.Int).SetString(balance, 16)
			} else {
				bigBalance, _ = new(big.Int).SetString(balance, 16)
			}

			genesisAcc := core.GenesisAccount{
				Balance: bigBalance,
				Proxy:   proxy,
			}
			ga[common.HexToAddress(acc)] = genesisAcc
		}
	}

	return ga
}

func DecodePreWormholesInfoV3(data string) core.GenesisAlloc {
	ga := make(core.GenesisAlloc)

	accountInfos := strings.Split(data, ",")
	for _, accountInfo := range accountInfos {
		accInfo := strings.Split(accountInfo, ":")
		acc := accInfo[0]
		balance := accInfo[1]
		bigBalance := big.NewInt(0)
		if strings.HasPrefix(balance, "0x") ||
			strings.HasPrefix(balance, "0X") {
			balance = string([]byte(balance)[2:])
			bigBalance, _ = new(big.Int).SetString(balance, 16)
		} else {
			bigBalance, _ = new(big.Int).SetString(balance, 16)
		}
		freerate, _ := strconv.Atoi(accInfo[2])

		genesisAcc := core.GenesisAccount{
			Balance:       bigBalance,
			FeeRate:       uint64(freerate),
			ExchangerName: accInfo[3],
			ExchangerUrl:  accInfo[4],
		}
		ga[common.HexToAddress(acc)] = genesisAcc
	}

	return ga
}

// testTxPool is a mock transaction pool that blindly accepts all transactions.
// Its goal is to get around setting up a valid statedb for the balance and nonce
// checks.
type testTxPool struct {
	pool map[common.Hash]*types.Transaction // Hash map of collected transactions

	txFeed event.Feed   // Notification feed to allow waiting for inclusion
	lock   sync.RWMutex // Protects the transaction pool
}

// newTestTxPool creates a mock transaction pool.
func newTestTxPool() *testTxPool {
	return &testTxPool{
		pool: make(map[common.Hash]*types.Transaction),
	}
}

// Has returns an indicator whether txpool has a transaction
// cached with the given hash.
func (p *testTxPool) Has(hash common.Hash) bool {
	p.lock.Lock()
	defer p.lock.Unlock()

	return p.pool[hash] != nil
}

// Get retrieves the transaction from local txpool with given
// tx hash.
func (p *testTxPool) Get(hash common.Hash) *types.Transaction {
	p.lock.Lock()
	defer p.lock.Unlock()

	return p.pool[hash]
}

// AddRemotes appends a batch of transactions to the pool, and notifies any
// listeners if the addition channel is non nil
func (p *testTxPool) AddRemotes(txs []*types.Transaction) []error {
	p.lock.Lock()
	defer p.lock.Unlock()

	for _, tx := range txs {
		p.pool[tx.Hash()] = tx
	}
	p.txFeed.Send(core.NewTxsEvent{Txs: txs})
	return make([]error, len(txs))
}

// Pending returns all the transactions known to the pool
func (p *testTxPool) Pending(enforceTips bool) (map[common.Address]types.Transactions, error) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	batches := make(map[common.Address]types.Transactions)
	for _, tx := range p.pool {
		from, _ := types.Sender(types.HomesteadSigner{}, tx)
		batches[from] = append(batches[from], tx)
	}
	for _, batch := range batches {
		sort.Sort(types.TxByNonce(batch))
	}
	return batches, nil
}

// SubscribeNewTxsEvent should return an event subscription of NewTxsEvent and
// send events to the given channel.
func (p *testTxPool) SubscribeNewTxsEvent(ch chan<- core.NewTxsEvent) event.Subscription {
	return p.txFeed.Subscribe(ch)
}

// testHandler is a live implementation of the Ethereum protocol handler, just
// preinitialized with some sane testing defaults and the transaction pool mocked
// out.
type testHandler struct {
	db      ethdb.Database
	chain   *core.BlockChain
	txpool  *testTxPool
	handler *handler
}

// newTestHandler creates a new handler for testing purposes with no blocks.
func newTestHandler() *testHandler {
	return newTestHandlerWithBlocks(0)
}

// newTestHandlerWithBlocks creates a new handler for testing purposes, with a
// given number of initial blocks.
func newTestHandlerWithBlocks(blocks int) *testHandler {
	// Create a database pre-initialize with a genesis block
	db := rawdb.NewMemoryDatabase()
	genesis := core.Genesis{
		Config:       params.AllEthashProtocolChanges,
		Nonce:        0,
		ExtraData:    hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000000f90182f9013b9444d952db5dfb4cbb54443554f4bb9cbebee2194c94085abc35ed85d26c2795b64c6ffb89b68ab1c47994edfc22e9cfb4e24815c3a12e81bf10cab9ce4d26949a1711a10e3d5baa4e0ce970df6e33dc50ef099294b31b41e5ef219fb0cc9935ad914158cf8970db4494fff531a2da46d051fde4c47f042ee6322407df3f94d8861d235134ef573894529b577af28ae0e3449c949d196915f63dbdb97dea552648123655109d98a594b685eb3226d5f0d549607d2cc18672b756fd090c9483c43f6f7bb4d8e429b21ff303a16b4c99a59b059416e6ee04db765a7d3bb07966d1af025d197ac3b694033eecd45d8c8ec84516359f39b11c260a56719e9493f24e8a3162b45611ab17a62dd0c95999cda60f94f50cbaffa72cc902de3f4f1e61132d858f3361d9948b07aff2327a3b7e2876d899cafac99f7ae16b10b8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c0"),
		Difficulty:   big.NewInt(1),
		Alloc:        DecodePreWormholesInfo(SimAllocData),
		Stake:        DecodePreWormholesInfoV3(SimStakeData),
		Validator:    DecodePreWormholesInfoV2(SimValidatorData_v2),
		Coinbase:     common.HexToAddress("0x0000000000000000000000000000000000000000"),
		Mixhash:      common.HexToHash("0x63746963616c2062797a616e74696e65206661756c7420746f6c6572616e6365"),
		ParentHash:   common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
		Timestamp:    0,
		Dir:          "/ipfs/QmS2U6Mu2X5HaUbrbVp6JoLmdcFphXiD98avZnq1My8vef",
		InjectNumber: 4096,
		StartIndex:   big.NewInt(0),
		Royalty:      100,
		Creator:      "0x35636d53Ac3DfF2b2347dDfa37daD7077b3f5b6F",
	}
	genesis.Alloc[testAddr] = core.GenesisAccount{Balance: big.NewInt(1000000)}
	genesis.MustCommit(db)

	chain, _ := core.NewBlockChain(db, nil, params.TestChainConfig, ethash.NewFaker(), vm.Config{}, nil, nil)

	bs, _ := core.GenerateChain(params.TestChainConfig, chain.Genesis(), ethash.NewFaker(), db, blocks, nil)
	if _, err := chain.InsertChain(bs); err != nil {
		panic(err)
	}
	txpool := newTestTxPool()

	handler, _ := newHandler(&handlerConfig{
		Database:   db,
		Chain:      chain,
		TxPool:     txpool,
		Network:    1,
		Sync:       downloader.FastSync,
		BloomCache: 1,
	})
	handler.Start(1000)

	return &testHandler{
		db:      db,
		chain:   chain,
		txpool:  txpool,
		handler: handler,
	}
}

// close tears down the handler and all its internal constructs.
func (b *testHandler) close() {
	b.handler.Stop()
	b.chain.Stop()
}
