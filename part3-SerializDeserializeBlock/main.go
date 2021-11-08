package main

import (
	"Blockchainer/part3-SerializDeserializeBlock/BLC"
	"fmt"
)

func main() {
	block := &BLC.Block{Data: []byte("Send 3 BTC to ZhangBozhi"), Nonce: 1000}
	fmt.Println(block)
	fmt.Printf("%s\n", block.Data)
	fmt.Printf("%d\n", block.Nonce)
	fmt.Printf("\n\n")

	bytes := block.Serialize()
	fmt.Println(bytes)
	fmt.Printf("\n\n")

	blc := BLC.DeserializeBlock(bytes)
	fmt.Println(blc)
	fmt.Printf("%s\n", blc.Data)
	fmt.Printf("%d\n", blc.Nonce)
	fmt.Printf("\n\n")
}
