package main

import (
	"github.com/boltdb/bolt"
	"log"
	"os"
	"fmt"
	"encoding/hex"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"


type BlockChain struct {
	tip []byte
	db *bolt.DB
}

type BlockChainIterator struct {
	currentHash []byte
	db *bolt.DB
}

func DbExist() bool {
	_,err := os.Stat(dbFile)
	if err != nil{
		return true
	}
	return false
}

//创建一条新的链
func CreateNewBlockChain(address string)  *BlockChain{
	if DbExist(){
		fmt.Println("Blockchain already exists")
		os.Exit(1)
	}
	var tip []byte
	db, err := bolt.Open(dbFile,0600,nil)
	if err!=nil{
		log.Panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		cbtx := NewCoinbaseTX(address,genesisCoinbaseData)  //创世coinbase
		genesisBlock := NewGenesisBlock(cbtx)   //创世块

		bucket, err := tx.CreateBucket([]byte(blocksBucket))
		if err!=nil{
			log.Panic(err)
		}
		err = bucket.Put(genesisBlock.Hash,genesisBlock.Serialize())
		if err != nil{
			log.Panic(err)
		}
		tip = genesisBlock.Hash
		return nil
	})
	if err != nil{
		log.Panic(err)
	}
	bc := BlockChain{tip:tip,db:db}
	return &bc
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


//不是太懂
func (bc *BlockChain) FindUnspentTransactions(address string) []*Transaction {
	var unspentTXs []*Transaction
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()
		for _,tx := range block.Transactions{
			txId := hex.EncodeToString(tx.ID)

			Outputs:
				for outIdx,out := range tx.Vout{
					if spentTXOs[txId] != nil{
						for _,spentOut := range spentTXOs[txId]{
							if spentOut == outIdx{
								continue Outputs
							}
						}
					}
					if out.CanUnlockInputWith(address){
						unspentTXs = append(unspentTXs,tx)
					}
				}
				if tx.isCoinbase() == false{
					for _,in := range tx.Vin{
						if in.CanUnlockOutputWith(address){
							inTxID := hex.EncodeToString(in.Txid)
							spentTXOs[inTxID] = append(spentTXOs[inTxID],in.Vout)
						}
					}
				}
		}
		if len(block.PreBlockHash) == 0{
			break
		}
	}
	return unspentTXs
}