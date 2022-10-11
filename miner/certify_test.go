package miner

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
	"testing"
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

func TestRlpDecode(t *testing.T) {
	type SignatureData2 struct {
		Address common.Address
		Height  *big.Int
	}
	ques := &SignatureData2{
		Address: common.HexToAddress("0x2000000000000000000000000000000000000002"),
		Height:  big.NewInt(2),
	}
	encQues, err := Encode(ques)
	if err != nil {
		t.Log("err", err)
	}
	msg := &Msg{
		Code: SendSignMsg,
		Msg:  encQues,
	}

	payload, err := msg.Payload()
	msg2 := new(Msg)
	if err := msg2.FromPayload(payload); err != nil {
		t.Log("msg.FromPayload  error", err)
	}
	var signature *SignatureData
	err = msg2.Decode(&signature)
	if err != nil {
		t.Log("err 2", err)
	}
	t.Log("signature", signature)
}

func TestRlpDecode2(t *testing.T) {
	type SignatureData2 struct {
		Address common.Address
		Height  *big.Int
	}
	ques := &SignatureData{
		Address: common.HexToAddress("0x2000000000000000000000000000000000000002"),
		Height:  big.NewInt(2),
		//Timestamp: uint64(time.Now().Unix()),
	}
	encQues, err := Encode(ques)
	if err != nil {
		t.Log("err", err)
	}
	msg := &Msg{
		Code: SendSignMsg,
		Msg:  encQues,
	}

	payload, err := msg.Payload()
	msg2 := new(Msg)
	if err := msg2.FromPayload(payload); err != nil {
		t.Log("msg.FromPayload  error", err)
	}
	var signature *SignatureData2
	err = msg2.Decode(&signature)
	if err != nil {
		t.Log("err 2", err)
	}
	t.Log("signature", signature)
}
