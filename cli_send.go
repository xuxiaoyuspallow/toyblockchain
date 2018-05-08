package main

import (
	"fmt"
	"log"
)

func (cli *CLI)send(from,to string,amount int)  {
	if !ValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}
	bc := NewBlockChain()
	utxoSet := UTXOSet{bc}
	defer bc.db.Close()

	tx := NewUTXOTransaction(from,to,amount,&utxoSet)
	cbTx := NewCoinbaseTX(from,"")
	txs := []*Transaction{cbTx,tx}

	newBlock := bc.MineBlock(txs)
	utxoSet.Update(newBlock)
	fmt.Println("Success")
}
