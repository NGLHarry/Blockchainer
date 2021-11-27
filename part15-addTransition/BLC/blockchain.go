package BLC

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/boltdb/bolt"
)

//database name
const dbFile = "blockchain.db"

//bucket
const blocksBucket = "blocks"

// 创世区块里面的数据信息
const genesisCoinBaseData = "first get bitcoin"

type BlockChain struct {
	Tip []byte   //区块链里最后一个区块的hash
	DB  *bolt.DB //数据库
}

//迭代器
type BlockchainIterator struct {
	CurrentHash []byte   // 当前正在遍历的区块的Hash
	DB          *bolt.DB // 数据库
}

//迭代器
func (blockchain *BlockChain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{blockchain.Tip, blockchain.DB}
}

//下一个迭代器
func (bi *BlockchainIterator) Next() *BlockchainIterator {
	var nextHash []byte
	//查询数据
	err := bi.DB.View(func(tx *bolt.Tx) error {
		//获取表
		b := tx.Bucket([]byte(blocksBucket))
		//通过当前的hash获取Block
		currentHashbytes := b.Get(bi.CurrentHash)
		//反序列化
		currentBlock := DeSerialBlock(currentHashbytes)
		nextHash = currentBlock.PrevBlockHash
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return &BlockchainIterator{nextHash, bi.DB}
}

// 先找到包含当前用户未花费输出的所有交易的集合
// 返回交易的数组
func (bc *BlockChain) FindUnspentTransations(address string) []Transaction {
	fmt.Println(address)
	// 输出未花费输出的交易
	var unspentTXs []Transaction
	//存储
	spentTXOs := make(map[string][]int)

	blockchainIterator := bc.Iterator()
	var hashInt big.Int
	for {
		err := blockchainIterator.DB.View(func(tx *bolt.Tx) error {
			// 获取表
			b := tx.Bucket([]byte(blocksBucket))
			// 通过Hash获取区块字节数组
			blockBytes := b.Get(blockchainIterator.CurrentHash)
			// 反序列化
			block := DeSerialBlock(blockBytes)

			for _, transction := range block.Transaction {
				fmt.Printf("TranscationHash:%x\n", transction.ID)
				// 将byte array转换成string
				txID := hex.EncodeToString(transction.ID)
			Outputs:
				for outIdx, out := range transction.Vout {
					// 是否已经花费？
					if spentTXOs[txID] != nil {
						for _, spentOut := range spentTXOs[txID] {
							if spentOut == outIdx {
								continue Outputs
							}
						}
					}
					fmt.Println(out.CanBeUnlockedWith(address))
					if out.CanBeUnlockedWith(address) {
						unspentTXs = append(unspentTXs, *transction)
						fmt.Println(unspentTXs)

					}
				}
				if !transction.IsCoinbase() {
					for _, in := range transction.Vin {
						if in.CanUnlockOutputWith(address) {
							intxID := hex.EncodeToString(in.Txid)
							spentTXOs[intxID] = append(spentTXOs[intxID], in.Vout)
						}
					}
				}
			}

			for _, transction := range block.Transaction {
				fmt.Printf("TransactionHash:%x\n", transction.ID)
			}

			fmt.Println()
			return nil
		})
		if err != nil {
			log.Panic(err)
		}
		// 获取下一个迭代器
		blockchainIterator = blockchainIterator.Next()

		// 将迭代器中的hash存储到hashInt
		hashInt.SetBytes(blockchainIterator.CurrentHash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}
	return unspentTXs
}

//查找可用的未消费的输出信息
func (bc *BlockChain) FindSpendableOutputs(address string, amout int) (int, map[string][]int) {
	//{"11111":[1,2,3],"00000":[2,3,5]}
	// 字典，存储交易id，Vout里面的未花费TXOutput的index
	unspentOutpus := make(map[string][]int)
	// 查看未花费
	unspentTXs := bc.FindUnspentTransations(address)

	accumulated := 0 //统计unspentoutpus里面对应的TXOutPut所对应的总量

Work:
	//遍历交易数组
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)
		// 遍历交易里面Vout
		for outIdx, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) && accumulated < amout {
				accumulated += out.Value
				unspentOutpus[txID] = append(unspentOutpus[txID], outIdx)

				if accumulated >= amout {
					break Work
				}
			}
		}
	}
	return accumulated, unspentOutpus
}

//创建一个带有创世区块节点的区块链
func NewBlockChain() *BlockChain {
	var tip []byte //获取最后一个区块的hash
	//1、尝试打开或创建数据库
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	//2、db.update更新数据库
	//1) 表是否存在？不存在，创建表
	//2) 创建创世区块
	//3) 将创世区块序列化
	//4) 把创世区块的Hash作为key，Block的序列化数据作为value存储到表中
	//5) 设置key,l,将hash作为value再次存储到数据库
	err = db.Update(func(tx *bolt.Tx) error {
		//判断这一张表是否存在数据库中
		b := tx.Bucket([]byte(blocksBucket))
		//说明表不存在
		if b == nil {
			fmt.Println("No existing blockchain found.creating a new one ...")
			//创建创世区块的交易对象
			cbtx := NewCoinbaseTX("maobuyi", genesisCoinBaseData)
			//创建创世区块
			genesisBlock := NewGenesisBlock(cbtx)
			//创建表
			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}
			//将创世区块序列化然后存入数据库中
			err = b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			//存储Hash
			err = b.Put([]byte("l"), genesisBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
			tip = genesisBlock.Hash

		} else {
			//key ：l
			//value:最后一个区块的Hash
			tip = b.Get([]byte("l"))
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
	return &BlockChain{tip, db}
}

// 根据交易的数组，打包新的区块
func (bc *BlockChain) MineBlock(txs []*Transaction) {
	err := bc.DB.Update(func(tx *bolt.Tx) error {
		// 新建区块
		newBlock := NewBlock(txs, bc.Tip)
		//将区块存储到数据库中
		b := tx.Bucket([]byte(blocksBucket))
		if b != nil {
			// key ： newBlock.hash
			// value: newBlock.Serialize()
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			//key :[]byte{"l"}
			//value: newBlock.Hash
			err = b.Put([]byte("l"), newBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
			//更新blockchain最新区块的hash
			bc.Tip = newBlock.Hash
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}
