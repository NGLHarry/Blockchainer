package BLC

import (
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

type CLI struct {
	BC *BlockChain
}

//打印参数信息
func (cli *CLI) printUsage() {
	fmt.Println("Usage")
	fmt.Println("Usage:")
	fmt.Println("\tgetbalance -address ADDRESS - Get balance of ADDRESS")
	fmt.Println("\tcreateblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("\tprintchain - Print all the blocks of the blockchain:")
	fmt.Println("\tsend -from FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO")
}

//判断终端参数的个数
func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(-1)
	}
}

func (cli *CLI) sendToken(from []string, to []string, amount []string) {
	//1. 10->mobuyi
	//2. 3->zhangjie
	// (1) 新建一个交易
	// var inputs []*Transaction

	// tx1 := NewUTXOTransaction("maobuyi", "zhangjie", 3, cli.BC, inputs)
	// fmt.Println("第一笔交易")
	// fmt.Println(tx1)
	// inputs = append(inputs, tx1)
	// fmt.Println("第一笔交易完成")

	// tx2 := NewUTXOTransaction("maobuyi", "feifei", 2, cli.BC, inputs)
	// fmt.Println("第二笔交易")
	// fmt.Println(tx2)
	// inputs = append(inputs, tx2)

	// tx3 := NewUTXOTransaction("maobuyi", "fengshuo", 2, cli.BC, inputs)
	// fmt.Println("第三笔交易")
	// fmt.Println(tx3)

	// tx4 := NewUTXOTransaction("zhangjie", "fengshuo", 2, cli.BC, inputs)
	// fmt.Println("第四笔交易")
	// fmt.Println(tx4)
	// cli.BC.MineBlock([]*Transaction{tx1, tx2, tx3, tx4})

	fmt.Println("from:")
	fmt.Println(from)
	fmt.Println("to:")
	fmt.Println(to)
	fmt.Println("amount:")
	fmt.Println(amount)

	blockchain := GetBlockchain()
	defer blockchain.DB.Close()

	var inputs []*Transaction
	for index, f := range from {
		num, err := strconv.Atoi(amount[index])
		if err != nil {
			log.Panic(err)
		}

		tx := NewUTXOTransaction(f, to[index], num, blockchain, inputs)
		inputs = append(inputs, tx)
	}

	blockchain.MineBlock(inputs)

}

// func (cli *CLI) addblock(data string) {
// 	// cli.BC.AddBlock(data)

// 	// fmt.Println("FindUnspentTransactions")

// 	// fmt.Println(cli.BC.FindUnspentTransations("maobuyi"))
// 	// count, outputMap := cli.BC.FindSpendableOutputs("maobuyi", 2)

// 	// fmt.Println(count)
// 	// fmt.Println(outputMap)
// 	cli.sendToken()
// }

func (cli *CLI) printChain() {

	//判断数据库是否存在
	if dbExists() == false {
		fmt.Println("the database is not exists!!!!")
		cli.printUsage()
		return
	}

	blockchain := GetBlockchain()
	defer blockchain.DB.Close()

	var blockChainIterator *BlockchainIterator
	blockChainIterator = blockchain.Iterator()
	var hashInt big.Int
	for {
		err := blockChainIterator.DB.View(func(tx *bolt.Tx) error {
			//获取表
			b := tx.Bucket([]byte(blocksBucket))
			//通过Hash获取区块字节数组
			blockBytes := b.Get(blockChainIterator.CurrentHash)
			block := DeSerialBlock(blockBytes)

			fmt.Println("=======================================")

			fmt.Printf("PrevBlockHash：%x \n", block.PrevBlockHash)
			fmt.Printf("Timestamp：%s \n", time.Unix(block.TimeStamp, 0).Format("2006-01-02 03:04:05 PM"))
			fmt.Printf("Hash：%x \n", block.Hash)
			fmt.Printf("Nonce：%d \n", block.Nonce)

			for _, tx := range block.Transaction {
				fmt.Println()
				fmt.Println(tx)
			}
			fmt.Println("=======================================")
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

	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	genenisAddress := createBlockchainCmd.String("address", "", "创建创世区块,并且将数据打包到数据库.")

	// addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	// addBlockData := addBlockCmd.String("data", "", "Block data")

	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	getBlanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	balanceAddress := getBlanceCmd.String("address", "", "查询余额...")

	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	sendFrom := sendCmd.String("from", "", "源地址...")
	sendTo := sendCmd.String("to", "", "目标地址...")
	sendAmount := sendCmd.String("amount", "", "转账金额...")

	switch os.Args[1] {
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getbalance":
		err := getBlanceCmd.Parse(os.Args[2:])
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
		os.Exit(1)
	}

	if createBlockchainCmd.Parsed() {
		if *genenisAddress == "" {
			cli.printUsage()
			os.Exit(1)
		}
		fmt.Println("创建创世区块并且存储到数据库....")
		cli.createBlockchain(*genenisAddress)
	}
	if getBlanceCmd.Parsed() {
		if *balanceAddress == "" {
			cli.printUsage()
			os.Exit(1)
		}
		fmt.Printf("查询 %s 的余额:%d \n", *balanceAddress, cli.getBalance(*balanceAddress))
	}
	if sendCmd.Parsed() {
		// if *sendFrom == "" || *sendTo == "" {
		// 	cli.printUsage()
		// 	os.Exit(1)
		// }
		// fmt.Printf("from:%s to:%s amount : %s\n", *sendFrom, *sendTo, *sendAmount)

		fromAddresses := JSONtoArray(*sendFrom)
		toAddresses := JSONtoArray(*sendTo)
		sendAmounts := JSONtoArray(*sendAmount)

		if len(fromAddresses) == len(toAddresses) && len(fromAddresses) == len(sendAmounts) {
			cli.sendToken(fromAddresses, toAddresses, sendAmounts)
		} else {
			fmt.Println("传入的参数有误....")
		}

		fmt.Println("from:")
		fmt.Println(fromAddresses)
		fmt.Println("to:")
		fmt.Println(toAddresses)
		fmt.Println("amount:")
		fmt.Println(sendAmounts)
	}
	if printChainCmd.Parsed() {
		fmt.Println("printchain ....")
		cli.printChain()
	}
}

func (cli *CLI) createBlockchain(genesis string) {
	if dbExists() {
		fmt.Println("创世区块已存在.....")
		os.Exit(1)
	}
	CreateGenenisBlockchain(genesis)
}

func (cli *CLI) getBalance(address string) int {
	blockchain := GetBlockchain()
	defer blockchain.DB.Close()

	// 查询某个地址所有的未花费交易输出
	txs := blockchain.FindUnspentTransations(address, nil)
	// 遍历数组所有未花费的值叠加
	balance := 0
	for _, tx := range txs {
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				balance += out.Value
			}
		}
	}
	//返回
	return balance
}
