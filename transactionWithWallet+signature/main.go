package main

import (
	"Blockchainer/transactionWithWallet+signature/BLC"
)

func main() {
	//blockchain := BLC.NewBlockChain()
	//创建CLI对象
	cli := BLC.CLI{}

	cli.Run()
}
