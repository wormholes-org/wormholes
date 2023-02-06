package miner

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
	"testing"
)

func TestDecodeMsg(t *testing.T) {
	// assemble onlinezkQuestion
	q := &types.OnlineZkQuestion{Height: big.NewInt(100)}
	enq, _ := Encode(q)
	msg := &types.EmptyMsg{
		//Code: MsgOnlineQuestion,
		Msg: enq,
	}
	enqMsg, _ := Encode(msg)
	var msg2 types.EmptyMsg
	err := rlp.DecodeBytes(enqMsg, &msg2)
	if err != nil {
		fmt.Println("err1", err)
	}

	var q2 types.OnlineZkQuestion

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
	msg := &types.EmptyMsg{
		Code: SendSignMsg,
		Msg:  encQues,
	}

	payload, err := msg.Payload()
	msg2 := new(types.EmptyMsg)
	if err := msg2.FromPayload(payload); err != nil {
		t.Log("msg.FromPayload  error", err)
	}
	var signature *types.SignatureData
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
	ques := &types.SignatureData{
		Address: common.HexToAddress("0x2000000000000000000000000000000000000002"),
		Height:  big.NewInt(2),
		//Timestamp: uint64(time.Now().Unix()),
	}
	encQues, err := Encode(ques)
	if err != nil {
		t.Log("err", err)
	}
	msg := &types.EmptyMsg{
		Code: SendSignMsg,
		Msg:  encQues,
	}

	payload, err := msg.Payload()
	msg2 := new(types.EmptyMsg)
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

func TestMsgSignature(t *testing.T) {
	ques := &types.SignatureData{
		Address: common.HexToAddress("0x2000000000000000000000000000000000000002"),
		Height:  big.NewInt(2),
		//Timestamp: uint64(time.Now().Unix()),
	}
	encQues, err := Encode(ques)
	if err != nil {
		t.Log("err", err)
	}
	msg := &types.EmptyMsg{
		Code: SendSignMsg,
		Msg:  encQues,
	}

	payload, err := msg.PayloadNoSig()
	//msg2 := new(Msg)
	//if err := msg2.FromPayload(payload); err != nil {
	//	t.Log("msg.FromPayload  error", err)
	//}
	for i := 0; i < 2; i++ {
		hashData := crypto.Keccak256(payload)
		prv, _ := crypto.HexToECDSA("501bbf00179b7e626d8983b7d7c9e1b040c8a5d9a0f5da649bf38e10b2dbfb8d")
		msg.Signature, _ = crypto.Sign(hashData, prv)
		t.Log(hex.EncodeToString(msg.Signature))
		bytes, _ := msg.Payload()

		//recover address
		msg := &types.EmptyMsg{}
		address, err := msg.RecoverAddress(bytes)
		if err != nil {
			t.Error(err)
		} else {
			t.Log("address= ", address)
		}

		toString := hex.EncodeToString(bytes)
		t.Log(toString)
	}
}
