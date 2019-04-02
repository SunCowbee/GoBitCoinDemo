package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

type ProofOfWork struct {
	// 区块
	block *Block
	// 目标值
	target *big.Int
}

func CreateProofOfWork(block *Block) *ProofOfWork {
	pow := ProofOfWork{
		block: block,
	}
	targetStr := "0000100000000000000000000000000000000000000000000000000000000000"
	bigIntTmp := big.Int{}
	bigIntTmp.SetString(targetStr, 16)
	pow.target = &bigIntTmp
	return &pow
}

func (pow *ProofOfWork) Run() (hash []byte, nonce uint64) {
	var hashArr [32]byte
	fmt.Printf("target : %x\n", pow.target)
	for {
		//[Size]byte
		hashArr =sha256.Sum256(pow.PrepareHash(nonce))

		tTmp := big.Int{}
		tTmp.SetBytes(hash[:])

		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		if tTmp.Cmp(pow.target) == -1 {
			//找到了，退出
			fmt.Printf("found hash : %x, %d\n", hashArr, nonce)
			break
		} else {
			//继续循环
			nonce++
		}
	}

	return hashArr[:], nonce
}

func (pow *ProofOfWork) PrepareHash(nonce uint64) []byte {
	block := pow.block
	tmp := [][]byte{
		Uint64ToBytes(block.Version),
		block.PreHash,
		//block.Hash,
		block.MerkelRoot,
		Uint64ToBytes(block.TimeStamp),
		Uint64ToBytes(block.Difficulity),
		Uint64ToBytes(nonce),
		block.Data,
	}

	data := bytes.Join(tmp, []byte{})
	return data
}