package main

import (
	"time"
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
)

type Block struct {
	Timestamp  int64
	Transactions  []*Transaction
	PreBlockHash []byte
	Hash []byte
	Nonce int
}

func NewGenesisBlock(coinbase *Transaction) *Block{
	return NewBlock([]*Transaction{coinbase}, []byte{})   //you can write anything on genesis block
}

func (b *Block)HashTransactions() []byte {
	var txHashes [][]byte
	var hash [32]byte
	for _,tx := range b.Transactions{
		txHashes = append(txHashes,tx.ID)
	}
	hash = sha256.Sum256(bytes.Join(txHashes,[]byte{}))
	return hash[:]
}

func NewBlock(tx []*Transaction,preblockhash []byte)  *Block{
	block := &Block{time.Now().Unix(),tx,preblockhash,[]byte{},0}
	pow := NewProofofWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}