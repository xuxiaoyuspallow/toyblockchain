package main

import (
	"time"
)

type Block struct {
	Timestamp  int64
	Data  []byte
	PreBlockHash []byte
	Hash []byte
	Nonce int
}


func NewBlock(data string,preblockhash []byte)  *Block{
	block := &Block{time.Now().Unix(),[]byte(data),preblockhash,[]byte{},0}
	pow := NewProofofWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}
