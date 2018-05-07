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
		err = bucket.Put([]byte("l"),genesisBlock.Hash)
		if err != nil{
			log.Panic(err)
		}
		tip = genesisBlock.Hash
		return nil
	})
	if err != nil{
		log.Panic(err)
	}
	bc := BlockChain{tip,db}
	return &bc
}


func NewBlockChain(address string) *BlockChain {
	if DbExist() == false{
		fmt.Println("No existing blockchain found.Create one first")
		os.Exit(1)
	}
	var tip []byte
	db, err := bolt.Open(dbFile,0600,nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
			tip = b.Get([]byte("l"))
			value := b.Get(tip)
			valu1 := DeserializeBlock(value)
			fmt.Println(valu1)
			return nil
	})
	if err != nil{
		log.Panic(err)
	}
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


//get all transactions that have unspent output
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


//get all outputs which can unlocked by private key from transactions that have unspent output
func (bc *BlockChain) FindUTXO(address string) []TXOutput{
	var UTXOs []TXOutput
	unspenttarnsactions := bc.FindUnspentTransactions(address)
	for _,tx := range unspenttarnsactions{
		for _, out := range tx.Vout{
			if out.CanUnlockInputWith(address){
				UTXOs = append(UTXOs,out)
			}
		}
	}
	return UTXOs
}

func (bc *BlockChain) FindSpendableOutputs(address string,amount int) (int, map[string][]int)  {
	unspentOutputs := make(map[string][]int)
	unspentTxs := bc.FindUnspentTransactions(address)
	accunmulated := 0

	WORK:
		for _, tx := range unspentTxs{
			txID := hex.EncodeToString(tx.ID)
			for outIdx, out := range tx.Vout{
				if out.CanUnlockInputWith(address) && accunmulated < amount{
					accunmulated += out.Value
					unspentOutputs[txID] = append(unspentOutputs[txID],outIdx)
					if accunmulated >= amount{
						break WORK
					}
				}
			}
		}
		return accunmulated, unspentOutputs
}

// MineBlock mines a new block with the provided transactions
func (bc *BlockChain) MineBlock(transactions []*Transaction) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(transactions, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}

		bc.tip = newBlock.Hash

		return nil
	})
}