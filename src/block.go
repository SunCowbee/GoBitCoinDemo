package main

import "crypto/sha256"

type Block struct {
	PreHash []byte
	Hash    []byte
	Data    []byte
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
