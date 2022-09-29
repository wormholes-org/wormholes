package common

import (
	"math/big"
	"math/rand"
	"time"
)

// Local random address
func LocalRandomBytes() []byte {
	rndSeed := time.Now().UnixNano()
	rand.Seed(rndSeed)
	bigSeed := big.NewInt(rand.Int63())
	//rndDat = common.BigToHash(bigSeed)
	return bigSeed.Bytes()
}
