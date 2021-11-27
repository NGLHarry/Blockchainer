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

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.HashTransaction(),
			IntToHex(pow.block.TimeStamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

// ProofOfWork object function
func (pow *ProofOfWork) Run() (int, []byte) {

	fmt.Printf("RUN....")

	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	// fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)

	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Printf("\n\n")

	return nonce, hash[:]
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

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1
	return isValid
}
