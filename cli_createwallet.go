package main

import "fmt"

func (cli *CLI)CreateWallet()  {
	wallets, _ := NewWallets()
	address := wallets.CreateWallet()
	wallets.SaveToFile()

	fmt.Printf("Your new address:%s\n",address)
}