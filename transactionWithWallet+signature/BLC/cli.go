package BLC

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct{}

//打印参数信息
func printUsage() {
	fmt.Println("Usage:")

	fmt.Println("\taddresslists -- 输出所有钱包地址.")
	fmt.Println("\tcreatewallet -- 创建钱包.")
	fmt.Println("\tcreateblockchain -address -- 交易数据.")
	fmt.Println("\tsend -from FROM -to TO -amount AMOUNT -- 交易明细.")
	fmt.Println("\tprintchain -- 输出区块信息.")
	fmt.Println("\tgetbalance -address -- 输出区块信息.")
}

//判断终端参数的个数
func validateArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(-1)
	}
}

func (cli *CLI) Run() {
	//判断终端参数的个数，如果没有参数，直接打印Usage信息并退出
	validateArgs()

	addresslistsCmd := flag.NewFlagSet("addresslists", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	sendBlockCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	getbalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)

	flagFrom := sendBlockCmd.String("from", "", "转账源地址......")
	flagTo := sendBlockCmd.String("to", "", "转账目的地地址......")
	flagAmount := sendBlockCmd.String("amount", "", "转账金额......")

	flagCreateBlockchainWithAddress := createBlockchainCmd.String("address", "", "创建创世区块的地址")
	getbalanceWithAdress := getbalanceCmd.String("address", "", "要查询某一个账号的余额.......")

	switch os.Args[1] {
	case "send":
		err := sendBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "addresslists":
		err := addresslistsCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getbalance":
		err := getbalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1)
	}

	if sendBlockCmd.Parsed() {

		if *flagFrom == "" || *flagTo == "" || *flagAmount == "" {
			printUsage()
			os.Exit(1)
		}

		from := JSONtoArray(*flagFrom)
		to := JSONtoArray(*flagTo)
		for index, fromAdress := range from {
			if !IsValidForAdress([]byte(fromAdress)) || !IsValidForAdress([]byte(to[index])) {
				fmt.Printf("地址无效......")
				printUsage()
				os.Exit(1)
			}
		}

		amount := JSONtoArray(*flagAmount)
		fmt.Println("from:")
		fmt.Println(from)
		fmt.Println("to:")
		fmt.Println(to)
		fmt.Println("amount:")
		fmt.Println(amount)
		cli.send(from, to, amount)
	}
	if printChainCmd.Parsed() {
		fmt.Println("printchain ....")
		cli.printChain()
	}
	if addresslistsCmd.Parsed() {
		cli.addresslists()
	}
	if createWalletCmd.Parsed() {
		fmt.Println("创建钱包")
		cli.createWallet()
	}
	if getbalanceCmd.Parsed() {
		if !IsValidForAdress([]byte(*getbalanceWithAdress)) {
			fmt.Println("地址无效....")
			printUsage()
			os.Exit(1)
		}
		cli.getBalance(*getbalanceWithAdress)
		// fmt.Printf("查询 %s 的余额:%d \n", *getbalanceWithAdress, cli.getBalance(*getbalanceWithAdress))
	}
	if createBlockchainCmd.Parsed() {

		if !IsValidForAdress([]byte(*flagCreateBlockchainWithAddress)) {
			fmt.Println("地址无效....")
			printUsage()
			os.Exit(1)
		}
		fmt.Println("createBlockchainCmd")
		cli.createGenesisBlockchain(*flagCreateBlockchainWithAddress)
	}
}
