package BLC

import(
	"time"
	"strconv"
	"crypto/sha256"
	"bytes"
)
 
type Block struct{
	//时间戳，创建区块的时间
	TimeStamp int64
	//上个区块的hash
	PrevBlockHash []byte
	//Data 交易数据
	Data []byte
	// Hash 当前区块的hash
	Hash []byte
}

func (block *Block) SetHash(){
	//1、将时间转换为数组
	timeString := strconv.FormatInt(block.TimeStamp,2)
	timestamp := []byte(timeString)
	//2、拼接
	headers := bytes.Join([][]byte{block.PrevBlockHash, block.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)
	block.Hash = hash[:]
}

func NewBlock(data string, prevBlockHash []byte) *Block{
	//创建区块
	block := &Block{time.Now().Unix(),prevBlockHash,[]byte(data),[]byte{}}
	//设置当前区块hash
	block.SetHash()
	return block
}

//创建创世区块
func NewGenesisBlock() *Block{
	return NewBlock("Genenis Block",[]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
}