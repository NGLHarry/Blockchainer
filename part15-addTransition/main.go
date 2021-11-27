package main

import (
	"Blockchainer/part15-addTransition/BLC"
)

func main() {
	blockchain := BLC.NewBlockChain()
	//创建CLI对象
	cli := BLC.CLI{blockchain}

	cli.Run()
}
