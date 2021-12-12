package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"time"
)

type Block struct {
	//时间戳，创建区块的时间
	TimeStamp int64
	//上个区块的hash
	PrevBlockHash []byte
	//Data 交易数据
	Transaction []*Transaction
	// Hash 当前区块的hash
	Hash []byte
	// Nonce随机数
	Nonce int
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

func DeSerialBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}

// 提供给挖矿使用
// 将区块里面所有的交易ID拼接，并且生成hash
func (b *Block) HashTransaction() []byte {
	var txHashes [][]byte
	var txhash [32]byte

	for _, tx := range b.Transaction {
		txHashes = append(txHashes, tx.ID)
	}
	txhash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txhash[:]
}

func NewBlock(transaction []*Transaction, prevBlockHash []byte) *Block {
	//创建区块
	block := &Block{time.Now().Unix(), prevBlockHash, transaction, []byte{}, 0}
	//设置当前区块hash
	//block.SetHash()

	// 构建工作量证明
	// 将block作为参数，创建一个pow对象
	pow := NewProofOfWork(block)

	//Run() 执行一次工作量证明
	nonce, hash := pow.Run()

	//设置区块Hash
	block.Hash = hash[:]
	// 设置Nonce值
	block.Nonce = nonce

	//验证合法性
	isValid := pow.Validate()

	fmt.Println(isValid)

	return block
}

//创建创世区块
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}
