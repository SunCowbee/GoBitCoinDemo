package main

// 区块链
type BlockChain struct {
	// 区块切片
	blocks []*Block
}

// 创建区块链
func CreateBlockChain() *BlockChain {
	// 创建创世区块
	genesisBlock := CreateGenesisBlock()
	blockChain := BlockChain{
		blocks: []*Block{genesisBlock},
	}
	return &blockChain
}

// 添加区块
func (blockChain *BlockChain) AddBlock(data string) {
	preBlock := blockChain.blocks[len(blockChain.blocks)-1]
	preHash := preBlock.Hash
	block := CreatBlock(data, preHash)
	blockChain.blocks = append(blockChain.blocks, block)
}
