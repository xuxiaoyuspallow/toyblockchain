package main

import (
	"fmt"
	"log"
)

func (cli *CLI) getBalance(address string)  {
	if !ValidateAddress(address){
		log.Panic("Error:Address is not valid")
	}
	bc := NewBlockChain(address)
	defer bc.db.Close()

	balance := 0
	pubKeyHash := Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1:len(pubKeyHash) - addressChecksumLen]
	UTOXs := bc.FindUTXO(pubKeyHash)
	for _, out := range UTOXs{
		balance += out.Value
	}
	fmt.Printf("Balance of '%s':%d\n",address,balance)
}
