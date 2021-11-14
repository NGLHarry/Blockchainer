package main

import (
	"Blockchainer/part13-seriAndDeserialBlock_cli_queryDatabase/BLC"
	"fmt"
	"math/big"
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

	var blockchainIterator *BLC.BlockchainIterator
	blockchainIterator = blockchain.Iterator()

	var hashInt big.Int

	for {

		fmt.Printf("%x\n", blockchainIterator.CurrentHash)

		// 获取下一个迭代器
		blockchainIterator = blockchainIterator.Next()

		// 将迭代器中的hash存储到hashInt
		hashInt.SetBytes(blockchainIterator.CurrentHash)

		/*
			// Cmp compares x and y and returns:
			//
			//   -1 if x <  y
			//    0 if x == y
			//   +1 if x >  y
		*/
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}

	}
}
