package types_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestRlpEvilAction(t *testing.T) {
	h1 := &types.Header{
		ParentHash: common.HexToHash("0x1"),
		Time:       uint64(time.Now().Unix()),
	}

	h2 := &types.Header{
		ParentHash: common.HexToHash("0x2"),
		Time:       uint64(time.Now().Unix()),
	}

	evilAction := types.NewEvilAction(h1)
	evilAction.EvilHeaders = append(evilAction.EvilHeaders, h2)

	encodeFh, err := types.EncodeEvilAction(*evilAction)
	if err != nil {
		t.Error(err)
	}

	ea2, err := types.DecodeEvilAction(encodeFh)
	if err != nil {
		t.Error(err)
	}
	for _, v := range ea2.EvilHeaders {
		fmt.Printf("t: %v\\n", v)
	}
}
