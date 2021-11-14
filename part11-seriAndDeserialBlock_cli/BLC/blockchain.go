package BLC

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

//database name
const dbFile = "blockchain.db"

//bucket
const blocksBucket = "blocks"

type BlockChain struct {
	Tip []byte   //区块链里最后一个区块的hash
	DB  *bolt.DB //数据库
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
			//创建创世区块
			genesisBlock := NewGenesisBlock()
			//创建表
			b, err := tx.CreateBucket([]byte("blocks"))
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
