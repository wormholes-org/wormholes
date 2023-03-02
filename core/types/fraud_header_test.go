package types

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"testing"
	"time"
)

func TestRlpFraudHeader(t *testing.T) {
	h1 := &Header{
		ParentHash: common.HexToHash("0x1"),
		Time:       uint64(time.Now().Unix()),
	}

	h2 := &Header{
		ParentHash: common.HexToHash("0x2"),
		Time:       uint64(time.Now().Unix()),
	}

	fraudheader := NewFraudHeader(h1, h2)

	encodeFh, err := EncodeFraudHeader(*fraudheader)
	if err != nil {
		t.Error(err)
	}

	var fraudheader2 FraudHeader
	fh, err := DecodeFraudHeader(encodeFh, fraudheader2)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("t: %v\n%v\n", fh.RemoteHeader,
		fh.LocalHeader)
}
