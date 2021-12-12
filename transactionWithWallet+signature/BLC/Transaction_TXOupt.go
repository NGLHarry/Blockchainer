package BLC

import "bytes"

// 交易输出
type TXOutput struct {
	Value         int64  //分
	Ripemd160Hash []byte //用户名
}

func (tXOutput *TXOutput) Lock(address string) {

	publicKeyHash := Base58Decode([]byte(address))

	tXOutput.Ripemd160Hash = publicKeyHash[1 : len(publicKeyHash)-4]

}

func NewTXOutput(value int64, address string) *TXOutput {
	txoutput := &TXOutput{value, nil}

	//设置Ripemd160Hash
	txoutput.Lock(address)
	return txoutput
}

// 解锁
func (txOutput *TXOutput) UnLockScriptPubKeyWithAddress(address string) bool {
	publicKeyHash := Base58Decode([]byte(address))
	hash160 := publicKeyHash[1 : len(publicKeyHash)-4]
	return bytes.Compare(txOutput.Ripemd160Hash, hash160) == 0
}
