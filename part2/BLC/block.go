package BLC

import(
	"time"
	"fmt"
	//"strconv"
	//"crypto/sha256"
	//"bytes"
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
	// Nonce随机数
	Nonce int
}

/*
func (block *Block) SetHash(){
	//1、将时间转换为数组
	timeString := strconv.FormatInt(block.TimeStamp,2)
	timestamp := []byte(timeString)
	//2、拼接
	headers := bytes.Join([][]byte{block.PrevBlockHash, block.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)
	block.Hash = hash[:]
}
*/

func NewBlock(data string, prevBlockHash []byte) *Block{
	//创建区块
	block := &Block{time.Now().Unix(),prevBlockHash,[]byte(data),[]byte{}, 0}
	//设置当前区块hash
	//block.SetHash()
	
	// 构建工作量证明
	// 将block作为参数，创建一个pow对象
	pow := NewProofOfWork(block)

	//Run() 执行一次工作量证明
	nonce,hash := pow.Run()

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
func NewGenesisBlock() *Block{
	return NewBlock("Genenis Block",[]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
}