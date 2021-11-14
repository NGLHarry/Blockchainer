package main

import (
	"Blockchainer/part11-seriAndDeserialBlock_cli/BLC"
	"fmt"
)

func main() {
	blockchain := BLC.NewBlockChain()
	fmt.Println(blockchain)
	fmt.Printf("tip: %x\n", blockchain.Tip)

}
