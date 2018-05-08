package main

import (
	"fmt"
	"strconv"
)

func (cli *CLI) printChain(){
	bc := NewBlockChain()
	defer bc.db.Close()
	bci := bc.Iterator()
	for {
		block := bci.Next()

		fmt.Printf("============ Block %x ============\n", block.Hash)
		fmt.Printf("Prev. hash: %x\n", block.PreBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofofWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}
		fmt.Printf("\n\n")

		if len(block.PreBlockHash) == 0 {
			break
		}
	}
}
