package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"

	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)
const addressChecksumLen = 4

type Wallet struct {
	// 1、私钥
	privateKey ecdsa.PrivateKey
	// 2、公钥
	publicKey []byte
}

func IsValidForAddress(address []byte) bool {
	version_public_checksumBytes := Base58Decode(address)
	fmt.Println(version_public_checksumBytes)

	checkSumBytes := version_public_checksumBytes[len(version_public_checksumBytes)-addressChecksumLen:]

	version_ripemd160 := version_public_checksumBytes[:len(version_public_checksumBytes)-addressChecksumLen]
	fmt.Println(len(checkSumBytes))
	fmt.Println(len(version_ripemd160))

	checkBytes := CheckSum(version_ripemd160)
	if bytes.Compare(checkSumBytes, checkBytes) == 0 {
		return true
	}
	return false
}

// 获取地址
func (w *Wallet) GetAddress() []byte {
	// 1、hash160
	ripemd160Hash := w.Ripemd160Hash(w.publicKey)
	version_ripemd160Hash := append([]byte{version}, ripemd160Hash...)
	checkSumBytes := CheckSum(version_ripemd160Hash)
	bytes := append(version_ripemd160Hash, checkSumBytes...)
	return Base58Encode(bytes)
}

func CheckSum(payload []byte) []byte {
	hash1 := sha256.Sum256(payload)
	hash2 := sha256.Sum256(hash1[:])
	return hash2[:addressChecksumLen]
}
func (w *Wallet) Ripemd160Hash(publicKey []byte) []byte {
	// 256
	hash256 := sha256.New()
	hash256.Write(publicKey)
	hash := hash256.Sum(nil)

	//160
	ripemd160 := ripemd160.New()
	ripemd160.Write(hash)
	return ripemd160.Sum(nil)

}

// 创建钱包
func NewWallet() *Wallet {
	privateKey, publicKey := newKeyPair()
	return &Wallet{privateKey, publicKey}
}

// 通过私钥产生公钥
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	publicKey := append(private.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, publicKey
}