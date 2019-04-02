package main

type BlockChain struct {
	blocks []*Block
}

func CreateBlockChain() *BlockChain {
	genesisBlock := CreateGenesisBlock()
	blockChain := BlockChain{
		blocks: []*Block{genesisBlock},
	}
	return &blockChain
}

func (blockChain *BlockChain) AddBlock(data string) {
	preBlock := blockChain.blocks[len(blockChain.blocks)-1]
	preHash := preBlock.Hash
	block := CreatBlock(data, preHash)
	blockChain.blocks = append(blockChain.blocks, block)
}
