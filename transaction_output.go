package main

import (
	"bytes"
	"encoding/gob"
	"log"
)

type TXOutput struct {
	Value int		//  交易的数值
	PubKeyHash []byte   //加密value的公玥
}

func (out *TXOutput)Lock(address []byte)  {
	pubKeyHash := Base58Decode(address)
	pubKeyHash = pubKeyHash[1:len(pubKeyHash)-addressChecksumLen]
	out.PubKeyHash = pubKeyHash
}

func (out *TXOutput)IsLockedWithKey(pubKeyHash []byte) bool  {
	return bytes.Compare(out.PubKeyHash,pubKeyHash) == 0
}

func NewTXOutput(value int, address string)*TXOutput  {
	txo := &TXOutput{value,nil} //
	txo.Lock([]byte(address))
	return txo
}

// TXOutputs collects TXOutput
type TXOutputs struct {
	Outputs []TXOutput
}

// Serialize serializes TXOutputs
func (outs TXOutputs) Serialize() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(outs)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// DeserializeOutputs deserializes TXOutputs
func DeserializeOutputs(data []byte) TXOutputs {
	var outputs TXOutputs

	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&outputs)
	if err != nil {
		log.Panic(err)
	}

	return outputs
}
