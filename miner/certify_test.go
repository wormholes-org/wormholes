package miner

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/rlp"
)

func TestDecodeMsg(t *testing.T) {
	// assemble onlinezkQuestion
	q := &OnlineZkQuestion{Height: big.NewInt(100)}
	enq, _ := Encode(q)
	msg := &Msg{
		//Code: MsgOnlineQuestion,
		Msg: enq,
	}
	enqMsg, _ := Encode(msg)
	var msg2 Msg
	err := rlp.DecodeBytes(enqMsg, &msg2)
	if err != nil {
		fmt.Println("err1", err)
	}

	var q2 OnlineZkQuestion

	err = rlp.DecodeBytes(msg2.Msg, &q2)
	if err != nil {
		fmt.Println("err2", err)
	}
	fmt.Println("height is :", q2.Height)
}
