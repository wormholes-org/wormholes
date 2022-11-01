// Copyright 2016 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package geth

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var customGenesisTests = []struct {
	genesis string
	query   string
	result  string
}{
	// Genesis file with an empty chain configuration (ensure missing fields work)
	{
		genesis: `{
							"alloc": {
			"0x091DBBa95B26793515cc9aCB9bEb5124c479f27F": {
			  "balance": "0xd3c21bcecceda1000000"
			},
			"0x107837Ea83f8f06533DDd3fC39451Cd0AA8DA8BD": {
			  "balance": "0xed2b525841adfc00000"
			},
			"0x612DFa56DcA1F581Ed34b9c60Da86f1268Ab6349": {
			  "balance": "0xd3c21bcecceda1000000"
			},
			"0x84d84e6073A06B6e784241a9B13aA824AB455326": {
			  "balance": "0xed2b525841adfc00000"
			},
			"0x9e4d5C72569465270232ed7Af71981Ee82d08dBF": {
			  "balance": "0xd3c21bcecceda1000000"
			},
			"0xa270bBDFf450EbbC2d0413026De5545864a1b6d6": {
			  "balance": "0xed2b525841adfc00000"
			},
			"0x4110E56ED25e21267FBeEf79244f47ada4e2E963": {
			  "balance": "0xd3c21bcecceda1000000"
			},
			"0xdb33217fE3F74bD41c550B06B624E23ab7f55d05": {
			  "balance": "0xed2b525841adfc00000"
			},
			"0xE2FA892CC5CC268a0cC1d924EC907C796351C645": {
			  "balance": "0xd3c21bcecceda1000000"
			}
		  },
		  "stake": {
			"0x091DBBa95B26793515cc9aCB9bEb5124c479f27F": {
			  "balance": "0xd3c21bcecceda100000"
			},
			"0x107837Ea83f8f06533DDd3fC39451Cd0AA8DA8BD": {
			  "balance": "0xed2b525841adfc0000"
			},
			"0x612DFa56DcA1F581Ed34b9c60Da86f1268Ab6349": {
			  "balance": "0xd3c21bcecceda100000"
			},
			"0x84d84e6073A06B6e784241a9B13aA824AB455326": {
			  "balance": "0xed2b525841adfc0000"
			},
			"0x9e4d5C72569465270232ed7Af71981Ee82d08dBF": {
			  "balance": "0xd3c21bcecceda100000"
			},
			"0xa270bBDFf450EbbC2d0413026De5545864a1b6d6": {
			  "balance": "0xed2b525841adfc0000"
			},
			"0x4110E56ED25e21267FBeEf79244f47ada4e2E963": {
			  "balance": "0xd3c21bcecceda100000"
			},
			"0xdb33217fE3F74bD41c550B06B624E23ab7f55d05": {
			  "balance": "0xed2b525841adfc0000"
			},
			"0xE2FA892CC5CC268a0cC1d924EC907C796351C645": {
			  "balance": "0xd3c21bcecceda100000"
			}
		  },
		  "validator": {
			"0x091DBBa95B26793515cc9aCB9bEb5124c479f27F": {
			  "balance": "0xd3c21bcecceda100000"
			},
			"0x107837Ea83f8f06533DDd3fC39451Cd0AA8DA8BD": {
			  "balance": "0xed2b525841adfc0000"
			},
			"0x612DFa56DcA1F581Ed34b9c60Da86f1268Ab6349": {
			  "balance": "0xd3c21bcecceda100000"
			},
			"0x84d84e6073A06B6e784241a9B13aA824AB455326": {
			  "balance": "0xed2b525841adfc0000"
			},
			"0x9e4d5C72569465270232ed7Af71981Ee82d08dBF": {
			  "balance": "0xd3c21bcecceda100000"
			},
			"0xa270bBDFf450EbbC2d0413026De5545864a1b6d6": {
			  "balance": "0xed2b525841adfc0000"
			},
			"0x4110E56ED25e21267FBeEf79244f47ada4e2E963": {
			  "balance": "0xd3c21bcecceda100000"
			},
			"0xdb33217fE3F74bD41c550B06B624E23ab7f55d05": {
			  "balance": "0xed2b525841adfc0000"
			},
			"0xE2FA892CC5CC268a0cC1d924EC907C796351C645": {
			  "balance": "0xd3c21bcecceda100000"
			}
		  },
			"coinbase"   : "0x0000000000000000000000000000000000000000",
			"difficulty" : "0x20000",
			"extraData"  : "",
			"gasLimit"   : "0x2fefd8",
			"nonce"      : "0x0000000000001338",
			"mixhash"    : "0x0000000000000000000000000000000000000000000000000000000000000000",
			"parentHash" : "0x0000000000000000000000000000000000000000000000000000000000000000",
			 "timestamp": "0x00",
			  "royalty":100,
			  "creator":      "0x35636d53Ac3DfF2b2347dDfa37daD7077b3f5b6F",
			  "inject_number": 4096,
			  "start_index":   0,
			  "dir":          "/ipfs/QmS2U6Mu2X5HaUbrbVp6JoLmdcFphXiD98avZnq1My8vef",
			"config"     : {}
		}`,
		query:  "eth.getBlock(0).nonce",
		result: "0x0000000000001338",
	},
	// Genesis file with specific chain configurations
	{
		genesis: `{
							"alloc": {
			"0x091DBBa95B26793515cc9aCB9bEb5124c479f27F": {
			  "balance": "0xd3c21bcecceda1000000"
			},
			"0x107837Ea83f8f06533DDd3fC39451Cd0AA8DA8BD": {
			  "balance": "0xed2b525841adfc00000"
			},
			"0x612DFa56DcA1F581Ed34b9c60Da86f1268Ab6349": {
			  "balance": "0xd3c21bcecceda1000000"
			},
			"0x84d84e6073A06B6e784241a9B13aA824AB455326": {
			  "balance": "0xed2b525841adfc00000"
			},
			"0x9e4d5C72569465270232ed7Af71981Ee82d08dBF": {
			  "balance": "0xd3c21bcecceda1000000"
			},
			"0xa270bBDFf450EbbC2d0413026De5545864a1b6d6": {
			  "balance": "0xed2b525841adfc00000"
			},
			"0x4110E56ED25e21267FBeEf79244f47ada4e2E963": {
			  "balance": "0xd3c21bcecceda1000000"
			},
			"0xdb33217fE3F74bD41c550B06B624E23ab7f55d05": {
			  "balance": "0xed2b525841adfc00000"
			},
			"0xE2FA892CC5CC268a0cC1d924EC907C796351C645": {
			  "balance": "0xd3c21bcecceda1000000"
			}
		  },
		  "stake": {
			"0x091DBBa95B26793515cc9aCB9bEb5124c479f27F": {
			  "balance": "0xd3c21bcecceda100000"
			},
			"0x107837Ea83f8f06533DDd3fC39451Cd0AA8DA8BD": {
			  "balance": "0xed2b525841adfc0000"
			},
			"0x612DFa56DcA1F581Ed34b9c60Da86f1268Ab6349": {
			  "balance": "0xd3c21bcecceda100000"
			},
			"0x84d84e6073A06B6e784241a9B13aA824AB455326": {
			  "balance": "0xed2b525841adfc0000"
			},
			"0x9e4d5C72569465270232ed7Af71981Ee82d08dBF": {
			  "balance": "0xd3c21bcecceda100000"
			},
			"0xa270bBDFf450EbbC2d0413026De5545864a1b6d6": {
			  "balance": "0xed2b525841adfc0000"
			},
			"0x4110E56ED25e21267FBeEf79244f47ada4e2E963": {
			  "balance": "0xd3c21bcecceda100000"
			},
			"0xdb33217fE3F74bD41c550B06B624E23ab7f55d05": {
			  "balance": "0xed2b525841adfc0000"
			},
			"0xE2FA892CC5CC268a0cC1d924EC907C796351C645": {
			  "balance": "0xd3c21bcecceda100000"
			}
		  },
		  "validator": {
			"0x091DBBa95B26793515cc9aCB9bEb5124c479f27F": {
			  "balance": "0xd3c21bcecceda100000"
			},
			"0x107837Ea83f8f06533DDd3fC39451Cd0AA8DA8BD": {
			  "balance": "0xed2b525841adfc0000"
			},
			"0x612DFa56DcA1F581Ed34b9c60Da86f1268Ab6349": {
			  "balance": "0xd3c21bcecceda100000"
			},
			"0x84d84e6073A06B6e784241a9B13aA824AB455326": {
			  "balance": "0xed2b525841adfc0000"
			},
			"0x9e4d5C72569465270232ed7Af71981Ee82d08dBF": {
			  "balance": "0xd3c21bcecceda100000"
			},
			"0xa270bBDFf450EbbC2d0413026De5545864a1b6d6": {
			  "balance": "0xed2b525841adfc0000"
			},
			"0x4110E56ED25e21267FBeEf79244f47ada4e2E963": {
			  "balance": "0xd3c21bcecceda100000"
			},
			"0xdb33217fE3F74bD41c550B06B624E23ab7f55d05": {
			  "balance": "0xed2b525841adfc0000"
			},
			"0xE2FA892CC5CC268a0cC1d924EC907C796351C645": {
			  "balance": "0xd3c21bcecceda100000"
			}
		  },
			"coinbase"   : "0x0000000000000000000000000000000000000000",
			"difficulty" : "0x20000",
			"extraData"  : "",
			"gasLimit"   : "0x2fefd8",
			"nonce"      : "0x0000000000001339",
			"mixhash"    : "0x0000000000000000000000000000000000000000000000000000000000000000",
			"parentHash" : "0x0000000000000000000000000000000000000000000000000000000000000000",
			"timestamp"  : "0x00",
			  "royalty":100,
			  "creator":      "0x35636d53Ac3DfF2b2347dDfa37daD7077b3f5b6F",
			  "inject_number": 4096,
			  "start_index":   0,
			  "dir":          "/ipfs/QmS2U6Mu2X5HaUbrbVp6JoLmdcFphXiD98avZnq1My8vef",
			"config"     : {
				"homesteadBlock" : 42,
				"daoForkBlock"   : 141,
				"daoForkSupport" : true
			}
		}`,
		query:  "eth.getBlock(0).nonce",
		result: "0x0000000000001339",
	},
}

// Tests that initializing Geth with a custom genesis block and chain definitions
// work properly.
func TestCustomGenesis(t *testing.T) {
	for i, tt := range customGenesisTests {
		// Create a temporary data directory to use and inspect later
		datadir := tmpdir(t)
		defer os.RemoveAll(datadir)

		// Initialize the data directory with the custom genesis block
		json := filepath.Join(datadir, "genesis.json")
		if err := ioutil.WriteFile(json, []byte(tt.genesis), 0600); err != nil {
			t.Fatalf("test %d: failed to write genesis file: %v", i, err)
		}
		runGeth(t, "--datadir", datadir, "init", json).WaitExit()

		// Query the custom genesis block
		geth := runGeth(t, "--networkid", "1337", "--syncmode=full", "--cache", "16",
			"--datadir", datadir, "--maxpeers", "0", "--port", "0",
			"--nodiscover", "--nat", "none", "--ipcdisable",
			"--exec", tt.query, "console")
		geth.ExpectRegexp(tt.result)
		geth.ExpectExit()
	}
}
