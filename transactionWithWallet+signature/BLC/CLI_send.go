package BLC

import (
	"fmt"
	"os"
)

//转账
func (cli *CLI) send(from []string, to []string, amout []string) {
	if !dbExists() {
		fmt.Println("数据不存在...")
		os.Exit(1)
	}
	blockchain := BlockchainObject()
	defer blockchain.DB.Close()

	blockchain.MineNewBlock(from, to, amout)

}
