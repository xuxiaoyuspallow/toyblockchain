package main

import "fmt"

func (cli * CLI) createBlockchain(address string)  {
	bc := CreateNewBlockChain(address)
	defer bc.db.Close()
	utxoSet := UTXOSet{bc}
	utxoSet.Reindex()
	fmt.Println("create block chain done")
}