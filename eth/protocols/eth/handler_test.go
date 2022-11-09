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
	"math"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/trie"
)

var (
	// testKey is a private key to use for funding a tester account.
	testKey, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")

	// testAddr is the Ethereum address of the tester account.
	testAddr = crypto.PubkeyToAddress(testKey.PublicKey)
)

const simAllocData = "" +
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

const simStakeData = "" +
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

const simValidatorData_v2 = "" +
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

func decodePreWormholesInfo(data string) core.GenesisAlloc {
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

func decodePreWormholesInfoV2(data string) core.GenesisAlloc {
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

func decodePreWormholesInfoV3(data string) core.GenesisAlloc {
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

// testBackend is a mock implementation of the live Ethereum message handler. Its
// purpose is to allow testing the request/reply workflows and wire serialization
// in the `eth` protocol without actually doing any data processing.
type testBackend struct {
	db     ethdb.Database
	chain  *core.BlockChain
	txpool *core.TxPool
}

func (b *testBackend) HandleWorkerMsg(msg Decoder, peer *Peer) error {
	return nil
}

// newTestBackend creates an empty chain and wraps it into a mock backend.
func newTestBackend(blocks int) *testBackend {
	return newTestBackendWithGenerator(blocks, nil)
}

// newTestBackend creates a chain with a number of explicitly defined blocks and
// wraps it into a mock backend.
func newTestBackendWithGenerator(blocks int, generator func(int, *core.BlockGen)) *testBackend {
	// Create a database pre-initialize with a genesis block
	db := rawdb.NewMemoryDatabase()
	//(&core.Genesis{
	//	Config: params.TestChainConfig,
	//	Alloc:  core.GenesisAlloc{testAddr: {Balance: big.NewInt(100_000_000_000_000_000)}},
	//}).MustCommit(db)
	genesis := core.Genesis{
		Config:       params.AllEthashProtocolChanges,
		Nonce:        0,
		ExtraData:    hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000000f90182f9013b9444d952db5dfb4cbb54443554f4bb9cbebee2194c94085abc35ed85d26c2795b64c6ffb89b68ab1c47994edfc22e9cfb4e24815c3a12e81bf10cab9ce4d26949a1711a10e3d5baa4e0ce970df6e33dc50ef099294b31b41e5ef219fb0cc9935ad914158cf8970db4494fff531a2da46d051fde4c47f042ee6322407df3f94d8861d235134ef573894529b577af28ae0e3449c949d196915f63dbdb97dea552648123655109d98a594b685eb3226d5f0d549607d2cc18672b756fd090c9483c43f6f7bb4d8e429b21ff303a16b4c99a59b059416e6ee04db765a7d3bb07966d1af025d197ac3b694033eecd45d8c8ec84516359f39b11c260a56719e9493f24e8a3162b45611ab17a62dd0c95999cda60f94f50cbaffa72cc902de3f4f1e61132d858f3361d9948b07aff2327a3b7e2876d899cafac99f7ae16b10b8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c0"),
		Difficulty:   big.NewInt(1),
		Alloc:        decodePreWormholesInfo(simAllocData),
		Stake:        decodePreWormholesInfoV3(simStakeData),
		Validator:    decodePreWormholesInfoV2(simValidatorData_v2),
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
	genesis.Alloc[testAddr] = core.GenesisAccount{Balance: big.NewInt(100_000_000_000_000_000)}
	genesis.MustCommit(db)
	chain, _ := core.NewBlockChain(db, nil, params.TestChainConfig, ethash.NewFaker(), vm.Config{}, nil, nil)

	bs, _ := core.GenerateChain(params.TestChainConfig, chain.Genesis(), ethash.NewFaker(), db, blocks, generator)
	if _, err := chain.InsertChain(bs); err != nil {
		panic(err)
	}
	txconfig := core.DefaultTxPoolConfig
	txconfig.Journal = "" // Don't litter the disk with test journals

	return &testBackend{
		db:     db,
		chain:  chain,
		txpool: core.NewTxPool(txconfig, params.TestChainConfig, chain),
	}
}

// close tears down the transaction pool and chain behind the mock backend.
func (b *testBackend) close() {
	b.txpool.Stop()
	b.chain.Stop()
}

func (b *testBackend) Chain() *core.BlockChain     { return b.chain }
func (b *testBackend) StateBloom() *trie.SyncBloom { return nil }
func (b *testBackend) TxPool() TxPool              { return b.txpool }

func (b *testBackend) RunPeer(peer *Peer, handler Handler) error {
	// Normally the backend would do peer mainentance and handshakes. All that
	// is omitted and we will just give control back to the handler.
	return handler(peer)
}
func (b *testBackend) PeerInfo(enode.ID) interface{} { panic("not implemented") }

func (b *testBackend) AcceptTxs() bool {
	panic("data processing tests should be done in the handler package")
}
func (b *testBackend) Handle(*Peer, Packet) error {
	panic("data processing tests should be done in the handler package")
}

// Tests that block headers can be retrieved from a remote chain based on user queries.
func TestGetBlockHeaders65(t *testing.T) { testGetBlockHeaders(t, ETH65) }
func TestGetBlockHeaders66(t *testing.T) { testGetBlockHeaders(t, ETH66) }

func testGetBlockHeaders(t *testing.T, protocol uint) {
	t.Parallel()

	backend := newTestBackend(maxHeadersServe + 15)
	defer backend.close()

	peer, _ := newTestPeer("peer", protocol, backend)
	defer peer.close()

	// Create a "random" unknown hash for testing
	var unknown common.Hash
	for i := range unknown {
		unknown[i] = byte(i)
	}
	// Create a batch of tests for various scenarios
	limit := uint64(maxHeadersServe)
	tests := []struct {
		query  *GetBlockHeadersPacket // The query to execute for header retrieval
		expect []common.Hash          // The hashes of the block whose headers are expected
	}{
		// A single random block should be retrievable by hash and number too
		{
			&GetBlockHeadersPacket{Origin: HashOrNumber{Hash: backend.chain.GetBlockByNumber(limit / 2).Hash()}, Amount: 1},
			[]common.Hash{backend.chain.GetBlockByNumber(limit / 2).Hash()},
		}, {
			&GetBlockHeadersPacket{Origin: HashOrNumber{Number: limit / 2}, Amount: 1},
			[]common.Hash{backend.chain.GetBlockByNumber(limit / 2).Hash()},
		},
		// Multiple headers should be retrievable in both directions
		{
			&GetBlockHeadersPacket{Origin: HashOrNumber{Number: limit / 2}, Amount: 3},
			[]common.Hash{
				backend.chain.GetBlockByNumber(limit / 2).Hash(),
				backend.chain.GetBlockByNumber(limit/2 + 1).Hash(),
				backend.chain.GetBlockByNumber(limit/2 + 2).Hash(),
			},
		}, {
			&GetBlockHeadersPacket{Origin: HashOrNumber{Number: limit / 2}, Amount: 3, Reverse: true},
			[]common.Hash{
				backend.chain.GetBlockByNumber(limit / 2).Hash(),
				backend.chain.GetBlockByNumber(limit/2 - 1).Hash(),
				backend.chain.GetBlockByNumber(limit/2 - 2).Hash(),
			},
		},
		// Multiple headers with skip lists should be retrievable
		{
			&GetBlockHeadersPacket{Origin: HashOrNumber{Number: limit / 2}, Skip: 3, Amount: 3},
			[]common.Hash{
				backend.chain.GetBlockByNumber(limit / 2).Hash(),
				backend.chain.GetBlockByNumber(limit/2 + 4).Hash(),
				backend.chain.GetBlockByNumber(limit/2 + 8).Hash(),
			},
		}, {
			&GetBlockHeadersPacket{Origin: HashOrNumber{Number: limit / 2}, Skip: 3, Amount: 3, Reverse: true},
			[]common.Hash{
				backend.chain.GetBlockByNumber(limit / 2).Hash(),
				backend.chain.GetBlockByNumber(limit/2 - 4).Hash(),
				backend.chain.GetBlockByNumber(limit/2 - 8).Hash(),
			},
		},
		// The chain endpoints should be retrievable
		{
			&GetBlockHeadersPacket{Origin: HashOrNumber{Number: 0}, Amount: 1},
			[]common.Hash{backend.chain.GetBlockByNumber(0).Hash()},
		}, {
			&GetBlockHeadersPacket{Origin: HashOrNumber{Number: backend.chain.CurrentBlock().NumberU64()}, Amount: 1},
			[]common.Hash{backend.chain.CurrentBlock().Hash()},
		},
		// Ensure protocol limits are honored
		{
			&GetBlockHeadersPacket{Origin: HashOrNumber{Number: backend.chain.CurrentBlock().NumberU64() - 1}, Amount: limit + 10, Reverse: true},
			backend.chain.GetBlockHashesFromHash(backend.chain.CurrentBlock().Hash(), limit),
		},
		// Check that requesting more than available is handled gracefully
		{
			&GetBlockHeadersPacket{Origin: HashOrNumber{Number: backend.chain.CurrentBlock().NumberU64() - 4}, Skip: 3, Amount: 3},
			[]common.Hash{
				backend.chain.GetBlockByNumber(backend.chain.CurrentBlock().NumberU64() - 4).Hash(),
				backend.chain.GetBlockByNumber(backend.chain.CurrentBlock().NumberU64()).Hash(),
			},
		}, {
			&GetBlockHeadersPacket{Origin: HashOrNumber{Number: 4}, Skip: 3, Amount: 3, Reverse: true},
			[]common.Hash{
				backend.chain.GetBlockByNumber(4).Hash(),
				backend.chain.GetBlockByNumber(0).Hash(),
			},
		},
		// Check that requesting more than available is handled gracefully, even if mid skip
		{
			&GetBlockHeadersPacket{Origin: HashOrNumber{Number: backend.chain.CurrentBlock().NumberU64() - 4}, Skip: 2, Amount: 3},
			[]common.Hash{
				backend.chain.GetBlockByNumber(backend.chain.CurrentBlock().NumberU64() - 4).Hash(),
				backend.chain.GetBlockByNumber(backend.chain.CurrentBlock().NumberU64() - 1).Hash(),
			},
		}, {
			&GetBlockHeadersPacket{Origin: HashOrNumber{Number: 4}, Skip: 2, Amount: 3, Reverse: true},
			[]common.Hash{
				backend.chain.GetBlockByNumber(4).Hash(),
				backend.chain.GetBlockByNumber(1).Hash(),
			},
		},
		// Check a corner case where requesting more can iterate past the endpoints
		{
			&GetBlockHeadersPacket{Origin: HashOrNumber{Number: 2}, Amount: 5, Reverse: true},
			[]common.Hash{
				backend.chain.GetBlockByNumber(2).Hash(),
				backend.chain.GetBlockByNumber(1).Hash(),
				backend.chain.GetBlockByNumber(0).Hash(),
			},
		},
		// Check a corner case where skipping overflow loops back into the chain start
		{
			&GetBlockHeadersPacket{Origin: HashOrNumber{Hash: backend.chain.GetBlockByNumber(3).Hash()}, Amount: 2, Reverse: false, Skip: math.MaxUint64 - 1},
			[]common.Hash{
				backend.chain.GetBlockByNumber(3).Hash(),
			},
		},
		// Check a corner case where skipping overflow loops back to the same header
		{
			&GetBlockHeadersPacket{Origin: HashOrNumber{Hash: backend.chain.GetBlockByNumber(1).Hash()}, Amount: 2, Reverse: false, Skip: math.MaxUint64},
			[]common.Hash{
				backend.chain.GetBlockByNumber(1).Hash(),
			},
		},
		// Check that non existing headers aren't returned
		{
			&GetBlockHeadersPacket{Origin: HashOrNumber{Hash: unknown}, Amount: 1},
			[]common.Hash{},
		}, {
			&GetBlockHeadersPacket{Origin: HashOrNumber{Number: backend.chain.CurrentBlock().NumberU64() + 1}, Amount: 1},
			[]common.Hash{},
		},
	}
	// Run each of the tests and verify the results against the chain
	for i, tt := range tests {
		// Collect the headers to expect in the response
		var headers []*types.Header
		for _, hash := range tt.expect {
			headers = append(headers, backend.chain.GetBlockByHash(hash).Header())
		}
		// Send the hash request and verify the response
		if protocol <= ETH65 {
			p2p.Send(peer.app, GetBlockHeadersMsg, tt.query)
			if err := p2p.ExpectMsg(peer.app, BlockHeadersMsg, headers); err != nil {
				t.Errorf("test %d: headers mismatch: %v", i, err)
			}
		} else {
			p2p.Send(peer.app, GetBlockHeadersMsg, GetBlockHeadersPacket66{
				RequestId:             123,
				GetBlockHeadersPacket: tt.query,
			})
			if err := p2p.ExpectMsg(peer.app, BlockHeadersMsg, BlockHeadersPacket66{
				RequestId:          123,
				BlockHeadersPacket: headers,
			}); err != nil {
				t.Errorf("test %d: headers mismatch: %v", i, err)
			}
		}
		// If the test used number origins, repeat with hashes as the too
		if tt.query.Origin.Hash == (common.Hash{}) {
			if origin := backend.chain.GetBlockByNumber(tt.query.Origin.Number); origin != nil {
				tt.query.Origin.Hash, tt.query.Origin.Number = origin.Hash(), 0

				if protocol <= ETH65 {
					p2p.Send(peer.app, GetBlockHeadersMsg, tt.query)
					if err := p2p.ExpectMsg(peer.app, BlockHeadersMsg, headers); err != nil {
						t.Errorf("test %d: headers mismatch: %v", i, err)
					}
				} else {
					p2p.Send(peer.app, GetBlockHeadersMsg, GetBlockHeadersPacket66{
						RequestId:             456,
						GetBlockHeadersPacket: tt.query,
					})
					if err := p2p.ExpectMsg(peer.app, BlockHeadersMsg, BlockHeadersPacket66{
						RequestId:          456,
						BlockHeadersPacket: headers,
					}); err != nil {
						t.Errorf("test %d: headers mismatch: %v", i, err)
					}
				}
			}
		}
	}
}

// Tests that block contents can be retrieved from a remote chain based on their hashes.
func TestGetBlockBodies65(t *testing.T) { testGetBlockBodies(t, ETH65) }
func TestGetBlockBodies66(t *testing.T) { testGetBlockBodies(t, ETH66) }

func testGetBlockBodies(t *testing.T, protocol uint) {
	t.Parallel()

	backend := newTestBackend(maxBodiesServe + 15)
	defer backend.close()

	peer, _ := newTestPeer("peer", protocol, backend)
	defer peer.close()

	// Create a batch of tests for various scenarios
	limit := maxBodiesServe
	tests := []struct {
		random    int           // Number of blocks to fetch randomly from the chain
		explicit  []common.Hash // Explicitly requested blocks
		available []bool        // Availability of explicitly requested blocks
		expected  int           // Total number of existing blocks to expect
	}{
		{1, nil, nil, 1},             // A single random block should be retrievable
		{10, nil, nil, 10},           // Multiple random blocks should be retrievable
		{limit, nil, nil, limit},     // The maximum possible blocks should be retrievable
		{limit + 1, nil, nil, limit}, // No more than the possible block count should be returned
		{0, []common.Hash{backend.chain.Genesis().Hash()}, []bool{true}, 1},      // The genesis block should be retrievable
		{0, []common.Hash{backend.chain.CurrentBlock().Hash()}, []bool{true}, 1}, // The chains head block should be retrievable
		{0, []common.Hash{{}}, []bool{false}, 0},                                 // A non existent block should not be returned

		// Existing and non-existing blocks interleaved should not cause problems
		{0, []common.Hash{
			{},
			backend.chain.GetBlockByNumber(1).Hash(),
			{},
			backend.chain.GetBlockByNumber(10).Hash(),
			{},
			backend.chain.GetBlockByNumber(100).Hash(),
			{},
		}, []bool{false, true, false, true, false, true, false}, 3},
	}
	// Run each of the tests and verify the results against the chain
	for i, tt := range tests {
		// Collect the hashes to request, and the response to expectva
		var (
			hashes []common.Hash
			bodies []*BlockBody
			seen   = make(map[int64]bool)
		)
		for j := 0; j < tt.random; j++ {
			for {
				num := rand.Int63n(int64(backend.chain.CurrentBlock().NumberU64()))
				if !seen[num] {
					seen[num] = true

					block := backend.chain.GetBlockByNumber(uint64(num))
					hashes = append(hashes, block.Hash())
					if len(bodies) < tt.expected {
						bodies = append(bodies, &BlockBody{Transactions: block.Transactions(), Uncles: block.Uncles()})
					}
					break
				}
			}
		}
		for j, hash := range tt.explicit {
			hashes = append(hashes, hash)
			if tt.available[j] && len(bodies) < tt.expected {
				block := backend.chain.GetBlockByHash(hash)
				bodies = append(bodies, &BlockBody{Transactions: block.Transactions(), Uncles: block.Uncles()})
			}
		}
		// Send the hash request and verify the response
		if protocol <= ETH65 {
			p2p.Send(peer.app, GetBlockBodiesMsg, hashes)
			if err := p2p.ExpectMsg(peer.app, BlockBodiesMsg, bodies); err != nil {
				t.Errorf("test %d: bodies mismatch: %v", i, err)
			}
		} else {
			p2p.Send(peer.app, GetBlockBodiesMsg, GetBlockBodiesPacket66{
				RequestId:            123,
				GetBlockBodiesPacket: hashes,
			})
			if err := p2p.ExpectMsg(peer.app, BlockBodiesMsg, BlockBodiesPacket66{
				RequestId:         123,
				BlockBodiesPacket: bodies,
			}); err != nil {
				t.Errorf("test %d: bodies mismatch: %v", i, err)
			}
		}
	}
}

// Tests that the state trie nodes can be retrieved based on hashes.
func TestGetNodeData65(t *testing.T) { testGetNodeData(t, ETH65) }
func TestGetNodeData66(t *testing.T) { testGetNodeData(t, ETH66) }

func testGetNodeData(t *testing.T, protocol uint) {
	t.Parallel()

	// Define three accounts to simulate transactions with
	acc1Key, _ := crypto.HexToECDSA("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")
	acc2Key, _ := crypto.HexToECDSA("49a7b37aa6f6645917e7b807e9d1c00d4fa71f18343b0d4122a4d2df64dd6fee")
	acc1Addr := crypto.PubkeyToAddress(acc1Key.PublicKey)
	acc2Addr := crypto.PubkeyToAddress(acc2Key.PublicKey)

	signer := types.HomesteadSigner{}
	// Create a chain generator with some simple transactions (blatantly stolen from @fjl/chain_markets_test)
	generator := func(i int, block *core.BlockGen) {
		switch i {
		case 0:
			// In block 1, the test bank sends account #1 some ether.
			tx, _ := types.SignTx(types.NewTransaction(block.TxNonce(testAddr), acc1Addr, big.NewInt(10_000_000_000_000_000), params.TxGas, block.BaseFee(), nil), signer, testKey)
			block.AddTx(tx)
		case 1:
			// In block 2, the test bank sends some more ether to account #1.
			// acc1Addr passes it on to account #2.
			tx1, _ := types.SignTx(types.NewTransaction(block.TxNonce(testAddr), acc1Addr, big.NewInt(1_000_000_000_000_000), params.TxGas, block.BaseFee(), nil), signer, testKey)
			tx2, _ := types.SignTx(types.NewTransaction(block.TxNonce(acc1Addr), acc2Addr, big.NewInt(1_000_000_000_000_000), params.TxGas, block.BaseFee(), nil), signer, acc1Key)
			block.AddTx(tx1)
			block.AddTx(tx2)
		case 2:
			// Block 3 is empty but was mined by account #2.
			block.SetCoinbase(acc2Addr)
			block.SetExtra([]byte("yeehaw"))
		case 3:
			// Block 4 includes blocks 2 and 3 as uncle headers (with modified extra data).
			b2 := block.PrevBlock(1).Header()
			b2.Extra = []byte("foo")
			block.AddUncle(b2)
			b3 := block.PrevBlock(2).Header()
			b3.Extra = []byte("foo")
			block.AddUncle(b3)
		}
	}
	// Assemble the test environment
	backend := newTestBackendWithGenerator(4, generator)
	defer backend.close()

	peer, _ := newTestPeer("peer", protocol, backend)
	defer peer.close()

	// Fetch for now the entire chain db
	var hashes []common.Hash

	it := backend.db.NewIterator(nil, nil)
	for it.Next() {
		if key := it.Key(); len(key) == common.HashLength {
			hashes = append(hashes, common.BytesToHash(key))
		}
	}
	it.Release()

	if protocol <= ETH65 {
		p2p.Send(peer.app, GetNodeDataMsg, hashes)
	} else {
		p2p.Send(peer.app, GetNodeDataMsg, GetNodeDataPacket66{
			RequestId:         123,
			GetNodeDataPacket: hashes,
		})
	}
	msg, err := peer.app.ReadMsg()
	if err != nil {
		t.Fatalf("failed to read node data response: %v", err)
	}
	if msg.Code != NodeDataMsg {
		t.Fatalf("response packet code mismatch: have %x, want %x", msg.Code, NodeDataMsg)
	}
	var data [][]byte
	if protocol <= ETH65 {
		if err := msg.Decode(&data); err != nil {
			t.Fatalf("failed to decode response node data: %v", err)
		}
	} else {
		var res NodeDataPacket66
		if err := msg.Decode(&res); err != nil {
			t.Fatalf("failed to decode response node data: %v", err)
		}
		data = res.NodeDataPacket
	}
	// Verify that all hashes correspond to the requested data, and reconstruct a state tree
	for i, want := range hashes {
		if hash := crypto.Keccak256Hash(data[i]); hash != want {
			t.Errorf("data hash mismatch: have %x, want %x", hash, want)
		}
	}
	statedb := rawdb.NewMemoryDatabase()
	for i := 0; i < len(data); i++ {
		statedb.Put(hashes[i].Bytes(), data[i])
	}
	accounts := []common.Address{testAddr, acc1Addr, acc2Addr}
	for i := uint64(0); i <= backend.chain.CurrentBlock().NumberU64(); i++ {
		trie, _ := state.New(backend.chain.GetBlockByNumber(i).Root(), state.NewDatabase(statedb), nil)

		for j, acc := range accounts {
			state, _ := backend.chain.State()
			bw := state.GetBalance(acc)
			bh := trie.GetBalance(acc)

			if (bw != nil && bh == nil) || (bw == nil && bh != nil) {
				t.Errorf("test %d, account %d: balance mismatch: have %v, want %v", i, j, bh, bw)
			}
			if bw != nil && bh != nil && bw.Cmp(bw) != 0 {
				t.Errorf("test %d, account %d: balance mismatch: have %v, want %v", i, j, bh, bw)
			}
		}
	}
}

// Tests that the transaction receipts can be retrieved based on hashes.
func TestGetBlockReceipts65(t *testing.T) { testGetBlockReceipts(t, ETH65) }
func TestGetBlockReceipts66(t *testing.T) { testGetBlockReceipts(t, ETH66) }

func testGetBlockReceipts(t *testing.T, protocol uint) {
	t.Parallel()

	// Define three accounts to simulate transactions with
	acc1Key, _ := crypto.HexToECDSA("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")
	acc2Key, _ := crypto.HexToECDSA("49a7b37aa6f6645917e7b807e9d1c00d4fa71f18343b0d4122a4d2df64dd6fee")
	acc1Addr := crypto.PubkeyToAddress(acc1Key.PublicKey)
	acc2Addr := crypto.PubkeyToAddress(acc2Key.PublicKey)

	signer := types.HomesteadSigner{}
	// Create a chain generator with some simple transactions (blatantly stolen from @fjl/chain_markets_test)
	generator := func(i int, block *core.BlockGen) {
		switch i {
		case 0:
			// In block 1, the test bank sends account #1 some ether.
			tx, _ := types.SignTx(types.NewTransaction(block.TxNonce(testAddr), acc1Addr, big.NewInt(10_000_000_000_000_000), params.TxGas, block.BaseFee(), nil), signer, testKey)
			block.AddTx(tx)
		case 1:
			// In block 2, the test bank sends some more ether to account #1.
			// acc1Addr passes it on to account #2.
			tx1, _ := types.SignTx(types.NewTransaction(block.TxNonce(testAddr), acc1Addr, big.NewInt(1_000_000_000_000_000), params.TxGas, block.BaseFee(), nil), signer, testKey)
			tx2, _ := types.SignTx(types.NewTransaction(block.TxNonce(acc1Addr), acc2Addr, big.NewInt(1_000_000_000_000_000), params.TxGas, block.BaseFee(), nil), signer, acc1Key)
			block.AddTx(tx1)
			block.AddTx(tx2)
		case 2:
			// Block 3 is empty but was mined by account #2.
			block.SetCoinbase(acc2Addr)
			block.SetExtra([]byte("yeehaw"))
		case 3:
			// Block 4 includes blocks 2 and 3 as uncle headers (with modified extra data).
			b2 := block.PrevBlock(1).Header()
			b2.Extra = []byte("foo")
			block.AddUncle(b2)
			b3 := block.PrevBlock(2).Header()
			b3.Extra = []byte("foo")
			block.AddUncle(b3)
		}
	}
	// Assemble the test environment
	backend := newTestBackendWithGenerator(4, generator)
	defer backend.close()

	peer, _ := newTestPeer("peer", protocol, backend)
	defer peer.close()

	// Collect the hashes to request, and the response to expect
	var (
		hashes   []common.Hash
		receipts [][]*types.Receipt
	)
	for i := uint64(0); i <= backend.chain.CurrentBlock().NumberU64(); i++ {
		block := backend.chain.GetBlockByNumber(i)

		hashes = append(hashes, block.Hash())
		receipts = append(receipts, backend.chain.GetReceiptsByHash(block.Hash()))
	}
	// Send the hash request and verify the response
	if protocol <= ETH65 {
		p2p.Send(peer.app, GetReceiptsMsg, hashes)
		if err := p2p.ExpectMsg(peer.app, ReceiptsMsg, receipts); err != nil {
			t.Errorf("receipts mismatch: %v", err)
		}
	} else {
		p2p.Send(peer.app, GetReceiptsMsg, GetReceiptsPacket66{
			RequestId:         123,
			GetReceiptsPacket: hashes,
		})
		if err := p2p.ExpectMsg(peer.app, ReceiptsMsg, ReceiptsPacket66{
			RequestId:      123,
			ReceiptsPacket: receipts,
		}); err != nil {
			t.Errorf("receipts mismatch: %v", err)
		}
	}
}
