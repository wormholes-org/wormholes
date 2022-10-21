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

// This file contains some shares testing functionality, common to multiple
// different files and modules being tested. Client based network and Server
// based network can be created easily with available APIs.

package les

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/mclock"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/contracts/checkpointoracle/contract"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/forkid"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/les/checkpointoracle"
	"github.com/ethereum/go-ethereum/les/flowcontrol"
	vfs "github.com/ethereum/go-ethereum/les/vflux/server"
	"github.com/ethereum/go-ethereum/light"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/params"
)

var (
	bankKey, _ = crypto.GenerateKey()
	bankAddr   = crypto.PubkeyToAddress(bankKey.PublicKey)
	bankFunds  = big.NewInt(1_000_000_000_000_000_000)

	userKey1, _ = crypto.GenerateKey()
	userKey2, _ = crypto.GenerateKey()
	userAddr1   = crypto.PubkeyToAddress(userKey1.PublicKey)
	userAddr2   = crypto.PubkeyToAddress(userKey2.PublicKey)

	testContractAddr         common.Address
	testContractCode         = common.Hex2Bytes("606060405260cc8060106000396000f360606040526000357c01000000000000000000000000000000000000000000000000000000009004806360cd2685146041578063c16431b914606b57603f565b005b6055600480803590602001909190505060a9565b6040518082815260200191505060405180910390f35b60886004808035906020019091908035906020019091905050608a565b005b80600060005083606481101560025790900160005b50819055505b5050565b6000600060005082606481101560025790900160005b5054905060c7565b91905056")
	testContractCodeDeployed = testContractCode[16:]
	testContractDeployed     = uint64(2)

	testEventEmitterCode = common.Hex2Bytes("60606040523415600e57600080fd5b7f57050ab73f6b9ebdd9f76b8d4997793f48cf956e965ee070551b9ca0bb71584e60405160405180910390a160358060476000396000f3006060604052600080fd00a165627a7a723058203f727efcad8b5811f8cb1fc2620ce5e8c63570d697aef968172de296ea3994140029")

	// Checkpoint oracle relative fields
	oracleAddr   common.Address
	signerKey, _ = crypto.GenerateKey()
	signerAddr   = crypto.PubkeyToAddress(signerKey.PublicKey)
)

var (
	// The block frequency for creating checkpoint(only used in test)
	sectionSize = big.NewInt(128)

	// The number of confirmations needed to generate a checkpoint(only used in test).
	processConfirms = big.NewInt(1)

	// The token bucket buffer limit for testing purpose.
	testBufLimit = uint64(1000000)

	// The buffer recharging speed for testing purpose.
	testBufRecharge = uint64(1000)
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
	"0x765C83dbA2712582C5461b2145f054d4F85a3080:0xf2dc7d47f15600000"

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

/*
contract test {

    uint256[100] data;

    function Put(uint256 addr, uint256 value) {
        data[addr] = value;
    }

    function Get(uint256 addr) constant returns (uint256 value) {
        return data[addr];
    }
}
*/

// prepare pre-commits specified number customized blocks into chain.
func prepare(n int, backend *backends.SimulatedBackend) {
	var (
		ctx    = context.Background()
		signer = types.HomesteadSigner{}
	)
	for i := 0; i < n; i++ {
		switch i {
		case 0:
			// Builtin-block
			//    number: 1
			//    txs:    2

			// deploy checkpoint contract
			auth, _ := bind.NewKeyedTransactorWithChainID(bankKey, big.NewInt(1337))
			oracleAddr, _, _, _ = contract.DeployCheckpointOracle(auth, backend, []common.Address{signerAddr}, sectionSize, processConfirms, big.NewInt(1))

			// bankUser transfers some ether to user1
			nonce, _ := backend.PendingNonceAt(ctx, bankAddr)
			tx, _ := types.SignTx(types.NewTransaction(nonce, userAddr1, big.NewInt(10_000_000_000_000_000), params.TxGas, big.NewInt(params.InitialBaseFee), nil), signer, bankKey)
			backend.SendTransaction(ctx, tx)
		case 1:
			// Builtin-block
			//    number: 2
			//    txs:    4

			bankNonce, _ := backend.PendingNonceAt(ctx, bankAddr)
			userNonce1, _ := backend.PendingNonceAt(ctx, userAddr1)

			// bankUser transfers more ether to user1
			tx1, _ := types.SignTx(types.NewTransaction(bankNonce, userAddr1, big.NewInt(1_000_000_000_000_000), params.TxGas, big.NewInt(params.InitialBaseFee), nil), signer, bankKey)
			backend.SendTransaction(ctx, tx1)

			// user1 relays ether to user2
			tx2, _ := types.SignTx(types.NewTransaction(userNonce1, userAddr2, big.NewInt(1_000_000_000_000_000), params.TxGas, big.NewInt(params.InitialBaseFee), nil), signer, userKey1)
			backend.SendTransaction(ctx, tx2)

			// user1 deploys a test contract
			tx3, _ := types.SignTx(types.NewContractCreation(userNonce1+1, big.NewInt(0), 200000, big.NewInt(params.InitialBaseFee), testContractCode), signer, userKey1)
			backend.SendTransaction(ctx, tx3)
			testContractAddr = crypto.CreateAddress(userAddr1, userNonce1+1)

			// user1 deploys a event contract
			tx4, _ := types.SignTx(types.NewContractCreation(userNonce1+2, big.NewInt(0), 200000, big.NewInt(params.InitialBaseFee), testEventEmitterCode), signer, userKey1)
			backend.SendTransaction(ctx, tx4)
		case 2:
			// Builtin-block
			//    number: 3
			//    txs:    2

			// bankUser transfer some ether to signer
			bankNonce, _ := backend.PendingNonceAt(ctx, bankAddr)
			tx1, _ := types.SignTx(types.NewTransaction(bankNonce, signerAddr, big.NewInt(1000000000), params.TxGas, big.NewInt(params.InitialBaseFee), nil), signer, bankKey)
			backend.SendTransaction(ctx, tx1)

			// invoke test contract
			data := common.Hex2Bytes("C16431B900000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001")
			tx2, _ := types.SignTx(types.NewTransaction(bankNonce+1, testContractAddr, big.NewInt(0), 100000, big.NewInt(params.InitialBaseFee), data), signer, bankKey)
			backend.SendTransaction(ctx, tx2)
		case 3:
			// Builtin-block
			//    number: 4
			//    txs:    1

			// invoke test contract
			bankNonce, _ := backend.PendingNonceAt(ctx, bankAddr)
			data := common.Hex2Bytes("C16431B900000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000002")
			tx, _ := types.SignTx(types.NewTransaction(bankNonce, testContractAddr, big.NewInt(0), 100000, big.NewInt(params.InitialBaseFee), data), signer, bankKey)
			backend.SendTransaction(ctx, tx)
		}
		backend.Commit()
	}
}

// testIndexers creates a set of indexers with specified params for testing purpose.
func testIndexers(db ethdb.Database, odr light.OdrBackend, config *light.IndexerConfig, disablePruning bool) []*core.ChainIndexer {
	var indexers [3]*core.ChainIndexer
	indexers[0] = light.NewChtIndexer(db, odr, config.ChtSize, config.ChtConfirms, disablePruning)
	indexers[1] = core.NewBloomIndexer(db, config.BloomSize, config.BloomConfirms)
	indexers[2] = light.NewBloomTrieIndexer(db, odr, config.BloomSize, config.BloomTrieSize, disablePruning)
	// make bloomTrieIndexer as a child indexer of bloom indexer.
	indexers[1].AddChildIndexer(indexers[2])
	return indexers[:]
}

func newTestClientHandler(backend *backends.SimulatedBackend, odr *LesOdr, indexers []*core.ChainIndexer, db ethdb.Database, peers *serverPeerSet, ulcServers []string, ulcFraction int) (*clientHandler, func()) {
	var (
		evmux  = new(event.TypeMux)
		engine = ethash.NewFaker()
		//gspec  = core.Genesis{
		//	Config:   params.AllEthashProtocolChanges,
		//	Alloc:    core.GenesisAlloc{bankAddr: {Balance: bankFunds}},
		//	GasLimit: 100000000,
		//	BaseFee:  big.NewInt(params.InitialBaseFee),
		//}
		oracle *checkpointoracle.CheckpointOracle
	)
	gspec := core.Genesis{
		Config:       params.AllEthashProtocolChanges,
		Nonce:        0,
		GasLimit:     100000000,
		BaseFee:      big.NewInt(params.InitialBaseFee),
		ExtraData:    hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000000f90182f9013b9444d952db5dfb4cbb54443554f4bb9cbebee2194c94085abc35ed85d26c2795b64c6ffb89b68ab1c47994edfc22e9cfb4e24815c3a12e81bf10cab9ce4d26949a1711a10e3d5baa4e0ce970df6e33dc50ef099294b31b41e5ef219fb0cc9935ad914158cf8970db4494fff531a2da46d051fde4c47f042ee6322407df3f94d8861d235134ef573894529b577af28ae0e3449c949d196915f63dbdb97dea552648123655109d98a594b685eb3226d5f0d549607d2cc18672b756fd090c9483c43f6f7bb4d8e429b21ff303a16b4c99a59b059416e6ee04db765a7d3bb07966d1af025d197ac3b694033eecd45d8c8ec84516359f39b11c260a56719e9493f24e8a3162b45611ab17a62dd0c95999cda60f94f50cbaffa72cc902de3f4f1e61132d858f3361d9948b07aff2327a3b7e2876d899cafac99f7ae16b10b8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c0"),
		Difficulty:   big.NewInt(1),
		Alloc:        DecodePreWormholesInfo(SimAllocData),
		Stake:        DecodePreWormholesInfo(SimStakeData),
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
	gspec.Alloc[bankAddr] = core.GenesisAccount{Balance: bankFunds}
	genesis := gspec.MustCommit(db)
	chain, _ := light.NewLightChain(odr, gspec.Config, engine, nil)
	if indexers != nil {
		checkpointConfig := &params.CheckpointOracleConfig{
			Address:   crypto.CreateAddress(bankAddr, 0),
			Signers:   []common.Address{signerAddr},
			Threshold: 1,
		}
		getLocal := func(index uint64) params.TrustedCheckpoint {
			chtIndexer := indexers[0]
			sectionHead := chtIndexer.SectionHead(index)
			return params.TrustedCheckpoint{
				SectionIndex: index,
				SectionHead:  sectionHead,
				CHTRoot:      light.GetChtRoot(db, index, sectionHead),
				BloomRoot:    light.GetBloomTrieRoot(db, index, sectionHead),
			}
		}
		oracle = checkpointoracle.New(checkpointConfig, getLocal)
	}
	client := &LightEthereum{
		lesCommons: lesCommons{
			genesis:     genesis.Hash(),
			config:      &ethconfig.Config{LightPeers: 100, NetworkId: NetworkId},
			chainConfig: params.AllEthashProtocolChanges,
			iConfig:     light.TestClientIndexerConfig,
			chainDb:     db,
			oracle:      oracle,
			chainReader: chain,
			closeCh:     make(chan struct{}),
		},
		peers:      peers,
		reqDist:    odr.retriever.dist,
		retriever:  odr.retriever,
		odr:        odr,
		engine:     engine,
		blockchain: chain,
		eventMux:   evmux,
	}
	client.handler = newClientHandler(ulcServers, ulcFraction, nil, client)

	if client.oracle != nil {
		client.oracle.Start(backend)
	}
	client.handler.start()
	return client.handler, func() {
		client.handler.stop()
	}
}

func newTestServerHandler(blocks int, indexers []*core.ChainIndexer, db ethdb.Database, clock mclock.Clock) (*serverHandler, *backends.SimulatedBackend, func()) {
	var (
		//gspec = core.Genesis{
		//	Config:   params.AllEthashProtocolChanges,
		//	Alloc:    core.GenesisAlloc{bankAddr: {Balance: bankFunds}},
		//	GasLimit: 100000000,
		//	BaseFee:  big.NewInt(params.InitialBaseFee),
		//}
		oracle *checkpointoracle.CheckpointOracle
	)

	gspec := core.Genesis{
		Config:       params.AllEthashProtocolChanges,
		Nonce:        0,
		GasLimit:     100000000,
		BaseFee:      big.NewInt(params.InitialBaseFee),
		ExtraData:    hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000000f90182f9013b9444d952db5dfb4cbb54443554f4bb9cbebee2194c94085abc35ed85d26c2795b64c6ffb89b68ab1c47994edfc22e9cfb4e24815c3a12e81bf10cab9ce4d26949a1711a10e3d5baa4e0ce970df6e33dc50ef099294b31b41e5ef219fb0cc9935ad914158cf8970db4494fff531a2da46d051fde4c47f042ee6322407df3f94d8861d235134ef573894529b577af28ae0e3449c949d196915f63dbdb97dea552648123655109d98a594b685eb3226d5f0d549607d2cc18672b756fd090c9483c43f6f7bb4d8e429b21ff303a16b4c99a59b059416e6ee04db765a7d3bb07966d1af025d197ac3b694033eecd45d8c8ec84516359f39b11c260a56719e9493f24e8a3162b45611ab17a62dd0c95999cda60f94f50cbaffa72cc902de3f4f1e61132d858f3361d9948b07aff2327a3b7e2876d899cafac99f7ae16b10b8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c0"),
		Difficulty:   big.NewInt(1),
		Alloc:        DecodePreWormholesInfo(SimAllocData),
		Stake:        DecodePreWormholesInfo(SimStakeData),
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
	gspec.Alloc[bankAddr] = core.GenesisAccount{Balance: bankFunds}
	genesis := gspec.MustCommit(db)

	// create a simulation backend and pre-commit several customized block to the database.
	simulation := backends.NewSimulatedBackendWithDatabase(db, gspec.Alloc, 100000000)
	prepare(blocks, simulation)

	txpoolConfig := core.DefaultTxPoolConfig
	txpoolConfig.Journal = ""
	txpool := core.NewTxPool(txpoolConfig, gspec.Config, simulation.Blockchain())
	if indexers != nil {
		checkpointConfig := &params.CheckpointOracleConfig{
			Address:   crypto.CreateAddress(bankAddr, 0),
			Signers:   []common.Address{signerAddr},
			Threshold: 1,
		}
		getLocal := func(index uint64) params.TrustedCheckpoint {
			chtIndexer := indexers[0]
			sectionHead := chtIndexer.SectionHead(index)
			return params.TrustedCheckpoint{
				SectionIndex: index,
				SectionHead:  sectionHead,
				CHTRoot:      light.GetChtRoot(db, index, sectionHead),
				BloomRoot:    light.GetBloomTrieRoot(db, index, sectionHead),
			}
		}
		oracle = checkpointoracle.New(checkpointConfig, getLocal)
	}
	server := &LesServer{
		lesCommons: lesCommons{
			genesis:     genesis.Hash(),
			config:      &ethconfig.Config{LightPeers: 100, NetworkId: NetworkId},
			chainConfig: params.AllEthashProtocolChanges,
			iConfig:     light.TestServerIndexerConfig,
			chainDb:     db,
			chainReader: simulation.Blockchain(),
			oracle:      oracle,
			closeCh:     make(chan struct{}),
		},
		peers:        newClientPeerSet(),
		servingQueue: newServingQueue(int64(time.Millisecond*10), 1),
		defParams: flowcontrol.ServerParams{
			BufLimit:    testBufLimit,
			MinRecharge: testBufRecharge,
		},
		fcManager: flowcontrol.NewClientManager(nil, clock),
	}
	server.costTracker, server.minCapacity = newCostTracker(db, server.config)
	server.costTracker.testCostList = testCostList(0) // Disable flow control mechanism.
	server.clientPool = vfs.NewClientPool(db, testBufRecharge, defaultConnectedBias, clock, alwaysTrueFn)
	server.clientPool.Start()
	server.clientPool.SetLimits(10000, 10000) // Assign enough capacity for clientpool
	server.handler = newServerHandler(server, simulation.Blockchain(), db, txpool, func() bool { return true })
	if server.oracle != nil {
		server.oracle.Start(simulation)
	}
	server.servingQueue.setThreads(4)
	server.handler.start()
	closer := func() { server.Stop() }
	return server.handler, simulation, closer
}

func alwaysTrueFn() bool {
	return true
}

// testPeer is a simulated peer to allow testing direct network calls.
type testPeer struct {
	cpeer *clientPeer
	speer *serverPeer

	net p2p.MsgReadWriter // Network layer reader/writer to simulate remote messaging
	app *p2p.MsgPipeRW    // Application layer reader/writer to simulate the local side
}

// handshakeWithServer executes the handshake with the remote server peer.
func (p *testPeer) handshakeWithServer(t *testing.T, td *big.Int, head common.Hash, headNum uint64, genesis common.Hash, forkID forkid.ID) {
	// It only works for the simulated client peer
	if p.cpeer == nil {
		t.Fatal("handshake for client peer only")
	}
	var sendList keyValueList
	sendList = sendList.add("protocolVersion", uint64(p.cpeer.version))
	sendList = sendList.add("networkId", uint64(NetworkId))
	sendList = sendList.add("headTd", td)
	sendList = sendList.add("headHash", head)
	sendList = sendList.add("headNum", headNum)
	sendList = sendList.add("genesisHash", genesis)
	if p.cpeer.version >= lpv4 {
		sendList = sendList.add("forkID", &forkID)
	}
	if err := p2p.ExpectMsg(p.app, StatusMsg, nil); err != nil {
		t.Fatalf("status recv: %v", err)
	}
	if err := p2p.Send(p.app, StatusMsg, sendList); err != nil {
		t.Fatalf("status send: %v", err)
	}
}

// handshakeWithClient executes the handshake with the remote client peer.
func (p *testPeer) handshakeWithClient(t *testing.T, td *big.Int, head common.Hash, headNum uint64, genesis common.Hash, forkID forkid.ID, costList RequestCostList, recentTxLookup uint64) {
	// It only works for the simulated client peer
	if p.speer == nil {
		t.Fatal("handshake for server peer only")
	}
	var sendList keyValueList
	sendList = sendList.add("protocolVersion", uint64(p.speer.version))
	sendList = sendList.add("networkId", uint64(NetworkId))
	sendList = sendList.add("headTd", td)
	sendList = sendList.add("headHash", head)
	sendList = sendList.add("headNum", headNum)
	sendList = sendList.add("genesisHash", genesis)
	sendList = sendList.add("serveHeaders", nil)
	sendList = sendList.add("serveChainSince", uint64(0))
	sendList = sendList.add("serveStateSince", uint64(0))
	sendList = sendList.add("serveRecentState", uint64(core.TriesInMemory-4))
	sendList = sendList.add("txRelay", nil)
	sendList = sendList.add("flowControl/BL", testBufLimit)
	sendList = sendList.add("flowControl/MRR", testBufRecharge)
	sendList = sendList.add("flowControl/MRC", costList)
	if p.speer.version >= lpv4 {
		sendList = sendList.add("forkID", &forkID)
		sendList = sendList.add("recentTxLookup", recentTxLookup)
	}
	if err := p2p.ExpectMsg(p.app, StatusMsg, nil); err != nil {
		t.Fatalf("status recv: %v", err)
	}
	if err := p2p.Send(p.app, StatusMsg, sendList); err != nil {
		t.Fatalf("status send: %v", err)
	}
}

// close terminates the local side of the peer, notifying the remote protocol
// manager of termination.
func (p *testPeer) close() {
	p.app.Close()
}

func newTestPeerPair(name string, version int, server *serverHandler, client *clientHandler) (*testPeer, *testPeer, error) {
	// Create a message pipe to communicate through
	app, net := p2p.MsgPipe()

	// Generate a random id and create the peer
	var id enode.ID
	rand.Read(id[:])

	peer1 := newClientPeer(version, NetworkId, p2p.NewPeer(id, name, nil), net)
	peer2 := newServerPeer(version, NetworkId, false, p2p.NewPeer(id, name, nil), app)

	// Start the peer on a new thread
	errc1 := make(chan error, 1)
	errc2 := make(chan error, 1)
	go func() {
		select {
		case <-server.closeCh:
			errc1 <- p2p.DiscQuitting
		case errc1 <- server.handle(peer1):
		}
	}()
	go func() {
		select {
		case <-client.closeCh:
			errc2 <- p2p.DiscQuitting
		case errc2 <- client.handle(peer2):
		}
	}()
	// Ensure the connection is established or exits when any error occurs
	for {
		select {
		case err := <-errc1:
			return nil, nil, fmt.Errorf("Failed to establish protocol connection %v", err)
		case err := <-errc2:
			return nil, nil, fmt.Errorf("Failed to establish protocol connection %v", err)
		default:
		}
		if atomic.LoadUint32(&peer1.serving) == 1 && atomic.LoadUint32(&peer2.serving) == 1 {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	return &testPeer{cpeer: peer1, net: net, app: app}, &testPeer{speer: peer2, net: app, app: net}, nil
}

type indexerCallback func(*core.ChainIndexer, *core.ChainIndexer, *core.ChainIndexer)

// testClient represents a client object for testing with necessary auxiliary fields.
type testClient struct {
	clock   mclock.Clock
	db      ethdb.Database
	peer    *testPeer
	handler *clientHandler

	chtIndexer       *core.ChainIndexer
	bloomIndexer     *core.ChainIndexer
	bloomTrieIndexer *core.ChainIndexer
}

// newRawPeer creates a new server peer connects to the server and do the handshake.
func (client *testClient) newRawPeer(t *testing.T, name string, version int, recentTxLookup uint64) (*testPeer, func(), <-chan error) {
	// Create a message pipe to communicate through
	app, net := p2p.MsgPipe()

	// Generate a random id and create the peer
	var id enode.ID
	rand.Read(id[:])
	peer := newServerPeer(version, NetworkId, false, p2p.NewPeer(id, name, nil), net)

	// Start the peer on a new thread
	errCh := make(chan error, 1)
	go func() {
		select {
		case <-client.handler.closeCh:
			errCh <- p2p.DiscQuitting
		case errCh <- client.handler.handle(peer):
		}
	}()
	tp := &testPeer{
		app:   app,
		net:   net,
		speer: peer,
	}
	var (
		genesis = client.handler.backend.blockchain.Genesis()
		head    = client.handler.backend.blockchain.CurrentHeader()
		td      = client.handler.backend.blockchain.GetTd(head.Hash(), head.Number.Uint64())
	)
	forkID := forkid.NewID(client.handler.backend.blockchain.Config(), genesis.Hash(), head.Number.Uint64())
	tp.handshakeWithClient(t, td, head.Hash(), head.Number.Uint64(), genesis.Hash(), forkID, testCostList(0), recentTxLookup) // disable flow control by default

	// Ensure the connection is established or exits when any error occurs
	for {
		select {
		case <-errCh:
			return nil, nil, nil
		default:
		}
		if atomic.LoadUint32(&peer.serving) == 1 {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	closePeer := func() {
		tp.speer.close()
		tp.close()
	}
	return tp, closePeer, errCh
}

// testServer represents a server object for testing with necessary auxiliary fields.
type testServer struct {
	clock   mclock.Clock
	backend *backends.SimulatedBackend
	db      ethdb.Database
	peer    *testPeer
	handler *serverHandler

	chtIndexer       *core.ChainIndexer
	bloomIndexer     *core.ChainIndexer
	bloomTrieIndexer *core.ChainIndexer
}

// newRawPeer creates a new client peer connects to the server and do the handshake.
func (server *testServer) newRawPeer(t *testing.T, name string, version int) (*testPeer, func(), <-chan error) {
	// Create a message pipe to communicate through
	app, net := p2p.MsgPipe()

	// Generate a random id and create the peer
	var id enode.ID
	rand.Read(id[:])
	peer := newClientPeer(version, NetworkId, p2p.NewPeer(id, name, nil), net)

	// Start the peer on a new thread
	errCh := make(chan error, 1)
	go func() {
		select {
		case <-server.handler.closeCh:
			errCh <- p2p.DiscQuitting
		case errCh <- server.handler.handle(peer):
		}
	}()
	tp := &testPeer{
		app:   app,
		net:   net,
		cpeer: peer,
	}
	var (
		genesis = server.handler.blockchain.Genesis()
		head    = server.handler.blockchain.CurrentHeader()
		td      = server.handler.blockchain.GetTd(head.Hash(), head.Number.Uint64())
	)
	forkID := forkid.NewID(server.handler.blockchain.Config(), genesis.Hash(), head.Number.Uint64())
	tp.handshakeWithServer(t, td, head.Hash(), head.Number.Uint64(), genesis.Hash(), forkID)

	// Ensure the connection is established or exits when any error occurs
	for {
		select {
		case <-errCh:
			return nil, nil, nil
		default:
		}
		if atomic.LoadUint32(&peer.serving) == 1 {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	closePeer := func() {
		tp.cpeer.close()
		tp.close()
	}
	return tp, closePeer, errCh
}

// testnetConfig wraps all the configurations for testing network.
type testnetConfig struct {
	blocks      int
	protocol    int
	indexFn     indexerCallback
	ulcServers  []string
	ulcFraction int
	simClock    bool
	connect     bool
	nopruning   bool
}

func newClientServerEnv(t *testing.T, config testnetConfig) (*testServer, *testClient, func()) {
	var (
		sdb    = rawdb.NewMemoryDatabase()
		cdb    = rawdb.NewMemoryDatabase()
		speers = newServerPeerSet()
	)
	var clock mclock.Clock = &mclock.System{}
	if config.simClock {
		clock = &mclock.Simulated{}
	}
	dist := newRequestDistributor(speers, clock)
	rm := newRetrieveManager(speers, dist, func() time.Duration { return time.Millisecond * 500 })
	odr := NewLesOdr(cdb, light.TestClientIndexerConfig, speers, rm)

	sindexers := testIndexers(sdb, nil, light.TestServerIndexerConfig, true)
	cIndexers := testIndexers(cdb, odr, light.TestClientIndexerConfig, config.nopruning)

	scIndexer, sbIndexer, sbtIndexer := sindexers[0], sindexers[1], sindexers[2]
	ccIndexer, cbIndexer, cbtIndexer := cIndexers[0], cIndexers[1], cIndexers[2]
	odr.SetIndexers(ccIndexer, cbIndexer, cbtIndexer)

	server, b, serverClose := newTestServerHandler(config.blocks, sindexers, sdb, clock)
	client, clientClose := newTestClientHandler(b, odr, cIndexers, cdb, speers, config.ulcServers, config.ulcFraction)

	scIndexer.Start(server.blockchain)
	sbIndexer.Start(server.blockchain)
	ccIndexer.Start(client.backend.blockchain)
	cbIndexer.Start(client.backend.blockchain)

	if config.indexFn != nil {
		config.indexFn(scIndexer, sbIndexer, sbtIndexer)
	}
	var (
		err          error
		speer, cpeer *testPeer
	)
	if config.connect {
		done := make(chan struct{})
		client.syncEnd = func(_ *types.Header) { close(done) }
		cpeer, speer, err = newTestPeerPair("peer", config.protocol, server, client)
		if err != nil {
			t.Fatalf("Failed to connect testing peers %v", err)
		}
		select {
		case <-done:
		case <-time.After(10 * time.Second):
			t.Fatal("test peer did not connect and sync within 3s")
		}
	}
	s := &testServer{
		clock:            clock,
		backend:          b,
		db:               sdb,
		peer:             cpeer,
		handler:          server,
		chtIndexer:       scIndexer,
		bloomIndexer:     sbIndexer,
		bloomTrieIndexer: sbtIndexer,
	}
	c := &testClient{
		clock:            clock,
		db:               cdb,
		peer:             speer,
		handler:          client,
		chtIndexer:       ccIndexer,
		bloomIndexer:     cbIndexer,
		bloomTrieIndexer: cbtIndexer,
	}
	teardown := func() {
		if config.connect {
			speer.close()
			cpeer.close()
			cpeer.cpeer.close()
			speer.speer.close()
		}
		ccIndexer.Close()
		cbIndexer.Close()
		scIndexer.Close()
		sbIndexer.Close()
		dist.close()
		serverClose()
		b.Close()
		clientClose()
	}
	return s, c, teardown
}

// NewFuzzerPeer creates a client peer for test purposes, and also returns
// a function to close the peer: this is needed to avoid goroutine leaks in the
// exec queue.
func NewFuzzerPeer(version int) (p *clientPeer, closer func()) {
	p = newClientPeer(version, 0, p2p.NewPeer(enode.ID{}, "", nil), nil)
	return p, func() { p.peerCommons.close() }
}
