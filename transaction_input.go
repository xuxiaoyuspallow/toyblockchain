package main

import "bytes"

type TXInput struct {
	Txid []byte   //transaction id
	Vout int    // stores an index of an output in the transaction
	Signature []byte  // 解密ScriptPubkey的私玥
	PubKey    []byte
}


func (in *TXInput)UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PubKey)
	return bytes.Compare(lockingHash,pubKeyHash) == 0
}