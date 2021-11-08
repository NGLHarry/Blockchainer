package BLC

type BlockChain struct {
	Blocks []*Block //存储有序的区块
}

//新增区块功能
func (blockchain *BlockChain) AddBlock(data string) {
	//创建新的Block
	preBlock := blockchain.Blocks[len(blockchain.Blocks)-1]
	newBlock := NewBlock(data, preBlock.Hash)
	//将区块添加到Blocks中
	blockchain.Blocks = append(blockchain.Blocks, newBlock)
}

//创建一个带有创世区块节点的区块链
func NewBlockChain() *BlockChain {
	return &BlockChain{[]*Block{NewGenesisBlock()}}
}
