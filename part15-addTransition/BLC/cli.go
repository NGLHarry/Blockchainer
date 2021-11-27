package BLC

import (
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

type CLI struct {
	BC *BlockChain
}

//打印参数信息
func (cli *CLI) printUsage() {
	fmt.Println("Usage")
	fmt.Println("\taddblock -data BLOCK_DATA - add a block to the blockchain")
	fmt.Println("\tprintchain - print all the blocks of the blockchain")
}

//判断终端参数的个数
func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(-1)
	}
}

func (cli *CLI) sendToken() {
	//1. 10->mobuyi
	//2. 3->zhangjie
	// (1) 新建一个交易
	tx1 := NewUTXOTransaction("maobuyi", "zhangjie", 3, cli.BC)
	tx2 := NewUTXOTransaction("maobuyi", "feifei", 2, cli.BC)

	cli.BC.MineBlock([]*Transaction{tx1, tx2})

}

func (cli *CLI) addblock(data string) {
	// cli.BC.AddBlock(data)

	// fmt.Println("FindUnspentTransactions")

	// fmt.Println(cli.BC.FindUnspentTransations("maobuyi"))
	// count, outputMap := cli.BC.FindSpendableOutputs("maobuyi", 2)

	// fmt.Println(count)
	// fmt.Println(outputMap)
	cli.sendToken()
}

func (cli *CLI) printChain() {
	var blockChainIterator *BlockchainIterator
	blockChainIterator = cli.BC.Iterator()
	var hashInt big.Int
	for {
		err := blockChainIterator.DB.View(func(tx *bolt.Tx) error {
			//获取表
			b := tx.Bucket([]byte(blocksBucket))
			//通过Hash获取区块字节数组
			blockBytes := b.Get(blockChainIterator.CurrentHash)
			block := DeSerialBlock(blockBytes)

			// fmt.Printf("Data：%s \n", string(block.Data))
			fmt.Printf("PrevBlockHash：%x \n", block.PrevBlockHash)
			fmt.Printf("Timestamp：%s \n", time.Unix(block.TimeStamp, 0).Format("2006-01-02 03:04:05 PM"))
			fmt.Printf("Hash：%x \n", block.Hash)
			fmt.Printf("Nonce：%d \n", block.Nonce)

			for _, tx := range block.Transaction {
				fmt.Println(tx)
			}
			fmt.Println()
			return nil
		})
		if err != nil {
			log.Panic(err)
		}

		blockChainIterator = blockChainIterator.Next()

		//将迭代器中的Hash存储到hashInt中
		hashInt.SetBytes(blockChainIterator.CurrentHash)

		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}
}

func (cli *CLI) Run() {
	//判断终端参数的个数，如果没有参数，直接打印Usage信息并退出
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	addBlockData := addBlockCmd.String("data", "", "Block data")

	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(0)
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			cli.printUsage()
			os.Exit(-1)
		}
		cli.addblock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		//通过迭代器遍历区块链中区块信息
		cli.printChain()
	}
}
