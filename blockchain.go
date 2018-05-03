package main

import (
	"github.com/boltdb/bolt"
	"log"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"

type BlockChain struct {
	tip []byte
	db *bolt.DB
}

func (bc *BlockChain) AddBlock(data string)  {
	var lastHash []byte
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("1"))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	newBlock := NewBlock(data, lastHash)
	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash,newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = b.Put([]byte("1"),newBlock.Hash)
		bc.tip = newBlock.Hash
		return nil
	})
	//preBlock := bc.blocks[len(bc.blocks)-1]
	//newBlock := NewBlock(data,preBlock.Hash)
	//bc.blocks = append(bc.blocks,newBlock)
}

func NewGenesisBlock() *Block{
	return NewBlock("Genesis Block", []byte{})   //you can write anything on genesis block
}

func NewBlockChain() *BlockChain {
	var tip []byte
	db, err := bolt.Open(dbFile,0600,nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil {
			genesis := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}
			err = b.Put(genesis.Hash,genesis.Serialize())
			err = b.Put([]byte("1"),genesis.Hash)
			tip = genesis.Hash
		}else{
			tip = b.Get([]byte("1"))
		}
		return nil
	})
	bc := BlockChain{tip,db}
	return &bc
}

type BlockChainIterator struct {
	currentHash []byte
	db *bolt.DB
}

func (bc *BlockChain) Iterator() *BlockChainIterator {
	bci := &BlockChainIterator{bc.tip, bc.db}

	return bci
}

func (i *BlockChainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedblock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedblock)
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	i.currentHash = block.PreBlockHash
	return block
}