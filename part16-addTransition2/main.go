package main

import (
	"Blockchainer/part16-addTransition2/BLC"
)

func main() {
	//blockchain := BLC.NewBlockChain()
	//创建CLI对象
	cli := BLC.CLI{}

	cli.Run()
}
