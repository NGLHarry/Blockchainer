package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

const subsidy = 10

type Transaction struct {
	// 1、交易ID
	ID []byte
	// 2、交易输入
	Vin []TXInput
	// 3、交易输出
	Vout []TXOutput
}

// 交易输入
type TXInput struct {
	// 1、交易ID 上个区块的
	Txid []byte
	// 2、存储TXoutput在Vout里面的索引
	Vout int
	// 3、用户名 签名
	ScriptSig string
}

// 交易输出
type TXOutput struct {
	Value        int    //分
	ScriptPubKey string //
}

// 创建一个新的coinbase交易
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	// 特殊的输入
	txin := TXInput{[]byte{}, -1, data}
	// 创建输出
	txout := TXOutput{subsidy, to}
	// 创建交易
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()

	return &tx
}

//设置交易hash
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	// 将序列化后的字节数组生成256hash
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]

}
