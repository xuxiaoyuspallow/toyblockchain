package main

import "bytes"

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
