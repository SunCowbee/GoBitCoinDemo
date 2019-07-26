package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"log"
	"time"
)

// 区块结构
type Block struct {
	//1.版本号
	Version uint64
	//2. 前区块哈希
	PrevHash []byte
	//3. Merkel根
	MerkelRoot []byte
	//4. 时间戳
	TimeStamp uint64
	//5. 难度值
	Difficulty uint64
	//6. 随机数，也就是挖矿要找的数据
	Nonce uint64

	//a. 当前区块哈希,正常比特币区块中没有当前区块的哈希，我们为了是方便做了简化！
	Hash []byte
	//b. 交易数据
	Transactions []*Transaction
}

// 将uint64转成[]byte
func Uint64ToByte(num uint64) []byte {
	var buffer bytes.Buffer

	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buffer.Bytes()
}

// 创建区块
func NewBlock(txs []*Transaction, prevBlockHash []byte) *Block {
	block := Block{
		Version:      00,
		PrevHash:     prevBlockHash,
		MerkelRoot:   []byte{},
		TimeStamp:    uint64(time.Now().Unix()),
		Difficulty:   0, //随便填写的无效值
		Nonce:        0,
		Hash:         []byte{},
		Transactions: txs,
	}

	block.MerkelRoot = block.MakeMerkelRoot()

	// 创建挖矿对象
	pow := NewProofOfWork(&block)
	// 挖矿
	hash, nonce := pow.Run()

	block.Hash = hash
	block.Nonce = nonce

	return &block
}

// 序列化区块
func (block *Block) Serialize() []byte {
	var buffer bytes.Buffer
	//- 使用gob进行序列化（编码）得到字节流
	//1. 定义一个编码器
	//2. 使用编码器进行编码
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(&block)
	if err != nil {
		log.Panic("编码出错!")
	}
	return buffer.Bytes()
}

// 反序列化区块
func Deserialize(data []byte) Block {

	decoder := gob.NewDecoder(bytes.NewReader(data))

	var block Block
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic("解码出错!", err)
	}

	return block
}

// 模拟梅克尔根，只是对交易的数据做简单的拼接，而不做二叉树处理！
func (block *Block) MakeMerkelRoot() []byte {
	var info []byte
	for _, tx := range block.Transactions {
		info = append(info, tx.TXID...)
	}

	hash := sha256.Sum256(info)
	return hash[:]
}
