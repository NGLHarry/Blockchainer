package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	//1. 区块高度
	Height int64
	//2. 上一个区块HASH
	PrevBlockHash []byte
	//3. 交易数据
	Txs []*Transaction
	//4. 时间戳
	Timestamp int64
	//5. Hash
	Hash []byte
	// 6. Nonce
	Nonce int64
}

// 需要将Txs转换为[]byte
// 提供给挖矿使用
// 将区块里面所有的交易ID拼接，并且生成hash
func (block *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range block.Txs {
		txHashes = append(txHashes, tx.TxHash)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

// 将区块序列化成字节数组
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

// 反序列化
func DeSerialBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}

func NewBlock(txs []*Transaction, height int64, prevBlockHash []byte) *Block {
	//创建区块
	block := &Block{height, prevBlockHash, txs, time.Now().Unix(), nil, 0}

	// 构建工作量证明
	// 将block作为参数，创建一个pow对象
	pow := NewProofOfWork(block)

	//Run() 执行一次工作量证明
	hash, nonce := pow.Run()

	//设置区块Hash
	block.Hash = hash[:]
	// 设置Nonce值
	block.Nonce = nonce

	//验证合法性
	// isValid := pow.Validate()

	// fmt.Println(isValid)

	return block
}

//创建创世区块
func CreateGenesisBlock(txs []*Transaction) *Block {
	return NewBlock(txs, 1, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}
