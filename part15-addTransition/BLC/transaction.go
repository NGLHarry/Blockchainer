package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
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

// 判断当前交易是否是CoinbaseTX
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && tx.Vin[0].Vout == -1 && len(tx.Vin[0].Txid) == 0
}

// 创建一个新的coinbase交易
func NewCoinbaseTX(to, data string) *Transaction {
	fmt.Printf("NewCoinbaseTX: %s\n", data)
	if data == " " {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	// 特殊的输入
	txin := TXInput{[]byte{}, -1, data}
	// 创建输出
	txout := TXOutput{subsidy, to}
	// 创建交易
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()
	fmt.Println("NewCoinbaseTx done!!!")

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

// 检查账号地址，解锁
func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	fmt.Printf("caCanUnlockOutputWith:%s\n", unlockingData)
	return in.ScriptSig == unlockingData
}

// 检查是否能够解锁账号
func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}

// 新建新的UTXO交易，转账
// 1. 先找到包含当前用户未花费输出的所有交易的集合
// 2. 找到用户足够的余额所对应的未花费输出
// 未花费输出：TXOutPut没有对应TXInput
// 3. 12,{“11111”:[1,2,3]}
// 4. 新建输入
// 5. 新建输出：
// (1)TXOutput{10,chenhaolin}
// (2)TXOutput{2,liyuechun}

// 建立转账交易
func NewUTXOTransaction(from, to string, amount int, bc *BlockChain) *Transaction {
	//输入
	var inputs []TXInput
	//输出
	var outputs []TXOutput

	//1.找到有效的可用的交易输出数据模型
	// 查询未花费的输出
	acc, validOutpus := bc.FindSpendableOutputs(from, amount)

	fmt.Println("NewUTXOTransaction")
	fmt.Println(acc)
	fmt.Println(validOutpus)

	if acc < amount {
		log.Panic("ERROR:Not enough fund")
	}

	//建立输入
	for txid, outs := range validOutpus {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}
		for _, out := range outs {
			//建立一个输入
			input := TXInput{txID, out, from}
			//将输入添加到inputs数组中
			inputs = append(inputs, input)
		}
	}
	//建立输出，转账
	output := TXOutput{amount, to}
	outputs = append(outputs, output)
	// 建立输出，找零
	output = TXOutput{acc - amount, from}
	//创建交易
	tx := Transaction{nil, inputs, outputs}
	tx.SetID()
	return &tx
}
