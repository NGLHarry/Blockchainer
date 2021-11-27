package main

import (
	"Blockchainer/part14-seriAndDeserialBlock_cli/BLC"
)

func main() {
	blockchain := BLC.NewBlockChain()
	//创建CLI对象
	cli := BLC.CLI{blockchain}

	cli.Run()
}
