package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// 工作量证明的结构
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// 创建挖矿对象
func NewProofOfWork(block *Block) *ProofOfWork {
	pow := ProofOfWork{
		block: block,
	}

	//我们指定的难度值，现在是一个string类型，需要进行转换
	targetStr := "0000100000000000000000000000000000000000000000000000000000000000"
	tmpInt := big.Int{}
	//将难度值赋值给big.int，指定16进制的格式
	tmpInt.SetString(targetStr, 16)

	pow.target = &tmpInt
	return &pow
}

// 挖矿
func (pow *ProofOfWork) Run() ([]byte, uint64) {

	var nonce uint64
	block := pow.block
	var hash [32]byte

	fmt.Println("开始挖矿...")
	for {
		tmp := [][]byte{
			Uint64ToByte(block.Version),
			block.PrevHash,
			block.MerkelRoot,
			Uint64ToByte(block.TimeStamp),
			Uint64ToByte(block.Difficulty),
			Uint64ToByte(nonce),
		}

		blockInfo := bytes.Join(tmp, []byte{})

		hash = sha256.Sum256(blockInfo)
		tmpInt := big.Int{}
		tmpInt.SetBytes(hash[:])

		//比较当前的哈希与目标哈希值，如果当前的哈希值小于目标的哈希值，就说明找到了，否则继续找
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		//
		if tmpInt.Cmp(pow.target) == -1 {
			fmt.Printf("挖矿成功！hash : %x, nonce : %d\n", hash, nonce)
			return hash[:], nonce
		} else {
			nonce++
		}
	}
}
