package main

import (
	"Blockchainer/wallet_address/BLC"
	"fmt"
)

func main() {
	wallet := BLC.NewWallet()
	address := wallet.GetAddress()
	fmt.Printf("这个地址为：%s\n", address)

	isValid := BLC.IsValidForAddress(address)
	fmt.Printf("%s 这个地址为 %v\n", address, isValid)

	wallets := BLC.NewWallets()

	fmt.Println(wallets.Wallets)
	wallets.CreateNewWallet()
	wallets.CreateNewWallet()

	fmt.Println(wallets.Wallets)
}
