package main

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
	"bytes"
)

func test()  {
	db, err := bolt.Open(dbFile,0600,nil)
	if err!=nil{
		log.Panic(err)
	}
	//err = db.Update(func(tx *bolt.Tx) error {
	//	//cbtx := NewCoinbaseTX("Ivan", genesisCoinbaseData) //创世coinbase
	//	//genesisBlock := NewGenesisBlock(cbtx)              //创世块
	//
	//	bucket, err := tx.CreateBucket([]byte(blocksBucket))
	//	if err != nil {
	//		log.Panic(err)
	//	}
	//	//vv2 := genesisBlock.Transactions[0]
	//	//fmt.Println(vv2)
	//	err = bucket.Put([]byte("name"), []byte("value"))
	//	if err != nil {
	//		log.Panic(err)
	//	}
	//	return nil
	//})
	db.Close()
	db, err = bolt.Open(dbFile,0600,nil)
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip := b.Get([]byte("name"))
		//value := b.Get(tip)
		//valu1 := DeserializeBlock(tip)
		n := bytes.IndexByte(tip, 0)
		//vv2 := valu1.Transactions[0]
		fmt.Println(n)
		return nil
	})
}

func main() {
	cli := CLI{}
	cli.Run()
}
