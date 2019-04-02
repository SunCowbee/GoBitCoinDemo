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

func main() {
	block := CreatBlock("helloworld", []byte{})
	fmt.Printf("PreHash:\t%x\n",block.PreHash)
	fmt.Printf("Hash:\t%x\n",block.Hash)
	fmt.Printf("Data:\t%x\n",block.Data)
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
