package main

import (
	"crypto/sha256"
	"fmt"
)

type Block struct {
	PreHash []byte
	Hash    []byte
	Data    []byte
}

type BlockChain struct {
	blocks []*Block
}

func main() {
	blockChain := CreateBlockChain()
	blockChain.AddBlock("second block")
	blockChain.AddBlock("third block")
	for index, block := range blockChain.blocks {
		fmt.Printf("BlockIndex:\t%d\n", index)
		fmt.Printf("PreHash:\t%x\n", block.PreHash)
		fmt.Printf("Hash:\t%x\n", block.Hash)
		fmt.Printf("Data:\t%x\n", block.Data)
	}
}

func CreatBlock(data string, preHash []byte) *Block {
	block := Block{
		PreHash: preHash,
		Hash:    []byte{},
		Data:    []byte(data),
	}
	block.SetHash()
	return &block
}

func (block *Block) SetHash() {
	blockInfo := append(block.PreHash, block.Data...)
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
}

func CreateGenesisBlock() *Block {
	return CreatBlock("GenesisBlock", []byte{})
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
