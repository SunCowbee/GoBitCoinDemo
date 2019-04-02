package main

import (
	"fmt"
)

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
