package BLC

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

var (
	//define Nonce max
	maxNonce = math.MaxInt64
)

const targetBits = 20

type ProofOfWork struct {
	block  *Block   //current block to validate
	target *big.Int //Big number storage, block difficulty
}

func (pow *ProofOfWork) prepareData(nonce int64) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.HashTransactions(),
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

func (proofOfWork *ProofOfWork) Run() ([]byte, int64) {

	//1. 将Block的属性拼接成字节数组

	//2. 生成hash

	//3. 判断hash有效性，如果满足条件，跳出循环

	var nonce int64
	nonce = 0

	var hashInt big.Int // 存储我们新生成的hash
	var hash [32]byte

	for {
		//准备数据
		dataBytes := proofOfWork.prepareData(nonce)

		// 生成hash
		hash = sha256.Sum256(dataBytes)
		fmt.Printf("\r%x", hash)

		// 将hash存储到hashInt
		hashInt.SetBytes(hash[:])

		//判断hashInt是否小于Block里面的target
		// Cmp compares x and y and returns:
		//
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		if proofOfWork.target.Cmp(&hashInt) == 1 {
			break
		}

		nonce = nonce + 1
	}

	return hash[:], int64(nonce)
}

func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	// fmt.Printf("--------------")
	// fmt.Printf("%b\n", target)
	//fmt.Printf("--------------")

	target.Lsh(target, uint(256-targetBits))
	//fmt.Printf("------target.Lsh------")
	//fmt.Printf("%b\n", target)

	pow := &ProofOfWork{block, target}

	return pow
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(int64(pow.block.Nonce))
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1
	return isValid
}
