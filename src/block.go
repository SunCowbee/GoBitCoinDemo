package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"time"
)

type Block struct {
	// 版本号
	Version uint64
	// 前区块哈希值
	PreHash []byte
	// 默克尔根
	MerkelRoot []byte
	// 时间戳
	TimeStamp uint64
	// 难度值
	Difficulty uint64
	// 工作量证明随机数
	Nonce uint64
	// 本区块哈希值
	Hash []byte
	// 区块数据
	Data []byte
}

// 创建区块
func CreatBlock(data string, preHash []byte) *Block {
	block := Block{
		Version:    00,
		PreHash:    preHash,
		MerkelRoot: []byte{},
		TimeStamp:  uint64(time.Now().Unix()),
		Nonce:      0,
		Hash:       []byte{},
		Data:       []byte(data),
	}
	block.SetHash()
	return &block
}

// 将uint64的整形数据转换为[]byte
func Uint64ToBytes(num uint64)[]byte{
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buf.Bytes()
}

// 获取并设置本区块哈希值
func (block *Block) SetHash() {
	// 区块信息
	var blockInfo []byte
	blockInfo = append(blockInfo, Uint64ToBytes(block.Version)...)
	blockInfo = append(blockInfo, block.PreHash...)
	blockInfo = append(blockInfo, block.MerkelRoot...)
	blockInfo = append(blockInfo, Uint64ToBytes(block.TimeStamp)...)
	blockInfo = append(blockInfo, Uint64ToBytes(block.Nonce)...)
	blockInfo = append(blockInfo, Uint64ToBytes(block.Nonce)...)
	blockInfo = append(blockInfo, block.Hash...)
	blockInfo = append(blockInfo, block.Data...)
	// 获取本区块哈希值
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
}

// 创建创世区块
func CreateGenesisBlock() *Block {
	return CreatBlock("GenesisBlock", []byte{})
}
