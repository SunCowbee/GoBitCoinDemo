package main

import "fmt"

type Block struct {
	PreHash []byte
	Hash    []byte
	Data    []byte
}

func NewBlock(data string, preHash []byte) *Block{
	block:=Block{
		PreHash: preHash,
		Hash:[]byte{},
		Data:[]byte(data),
	}
	return &block
}

func main() {
	block := NewBlock("helloworld", []byte{})
	fmt.Printf("%x\n",block.PreHash)
	fmt.Printf("%x\n",block.Hash)
	fmt.Printf("%x\n",block.Data)
}
