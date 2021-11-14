package main

import (
	"Blockchainer/part12-seriAndDeserialBlock_cli_addBlock/BLC"
	"fmt"
)

func main() {
	blockchain := BLC.NewBlockChain()
	fmt.Println(blockchain)
	fmt.Printf("tip: %x\n", blockchain.Tip)
	fmt.Println("Send 100 BTC To shaolin!!")
	blockchain.AddBlock("Send 100 BTC To shaolin!!")
	fmt.Printf("tip:%x\n", blockchain.Tip)
	fmt.Println("Send 155 BTC To shaolin!!")
	blockchain.AddBlock("Send 155 BTC To shaolin!!")
	fmt.Printf("tip:%x\n", blockchain.Tip)
}
