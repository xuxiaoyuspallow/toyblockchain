package main

import "fmt"

func (cli * CLI) createBlockchain(address string)  {
	bc := CreateNewBlockChain(address)
	bc.db.Close()
	fmt.Println("create block chain done")
}