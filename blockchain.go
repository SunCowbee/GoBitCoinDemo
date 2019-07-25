package main

import (
	"./lib/bolt"
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"time"
)

// 区块链结构
type BlockChain struct {
	db   *bolt.DB
	tail []byte //存储最后一个区块的哈希
}

const blockChainDb = "blockChain.db"
const blockBucket = "blockBucket"

// 创建区块链
func NewBlockChain(address string) *BlockChain {

	// 最后一个区块的哈希
	var lastHash []byte

	// 打开数据库，没有会自己创建
	db, err := bolt.Open(blockChainDb, 0600, nil)
	//defer db.Close()

	if err != nil {
		log.Panic("打开数据库失败！")
	}

	// 操作数据库（改写）
	db.Update(func(tx *bolt.Tx) error {
		// 找到抽屉bucket
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			//没有抽屉，我们需要创建
			bucket, err = tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				log.Panic("创建bucket(b1)失败")
			}

			// 创建一个创世块
			genesisBlock := GenesisBlock(address)

			bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
			bucket.Put([]byte("LastHashKey"), genesisBlock.Hash)
			lastHash = genesisBlock.Hash

		} else {
			lastHash = bucket.Get([]byte("LastHashKey"))
		}
		return nil
	})

	return &BlockChain{db, lastHash}
}

//定义一个创世块
func GenesisBlock(address string) *Block {
	// 创建挖矿交易
	coinbase := NewCoinbaseTX(address, "Genesis Block")
	// 创建区块
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

// 区块链中添加区块
func (bc *BlockChain) AddBlock(txs []*Transaction) {

	// 遍历新增区块中的所有交易
	for _, tx := range txs {
		// 矿工校验交易中每个input的签名是否有效
		if !bc.VerifyTransaction(tx) {
			fmt.Printf("矿工发现无效交易!")
			return
		}
	}

	db := bc.db
	lastHash := bc.tail

	db.Update(func(tx *bolt.Tx) error {

		//完成数据添加
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panic("bucket 不应该为空，请检查!")
		}

		block := NewBlock(txs, lastHash)

		bucket.Put(block.Hash, block.Serialize())
		bucket.Put([]byte("LastHashKey"), block.Hash)

		bc.tail = block.Hash

		return nil
	})
}

// 遍历打印区块链
func (bc *BlockChain) Printchain() {

	blockHeight := 0
	bc.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("blockBucket"))

		//从第一个key-> value 进行遍历，到最后一个固定的key时直接返回
		b.ForEach(func(k, v []byte) error {
			if bytes.Equal(k, []byte("LastHashKey")) {
				return nil
			}

			block := Deserialize(v)
			fmt.Printf("=============== 区块高度: %d ==============\n", blockHeight)
			blockHeight++
			fmt.Printf("版本号: %d\n", block.Version)
			fmt.Printf("前区块哈希值: %x\n", block.PrevHash)
			fmt.Printf("梅克尔根: %x\n", block.MerkelRoot)
			timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
			fmt.Printf("时间戳: %s\n", timeFormat)
			fmt.Printf("难度值(随便写的）: %d\n", block.Difficulty)
			fmt.Printf("随机数 : %d\n", block.Nonce)
			fmt.Printf("当前区块哈希值: %x\n", block.Hash)
			fmt.Printf("区块数据 :%s\n", block.Transactions[0].TXInputs[0].PubKey)
			return nil
		})
		return nil
	})
}

// 找到指定地址的所有的utxo
func (bc *BlockChain) FindUTXOs(pubKeyHash []byte) []TXOutput {

	var UTXO []TXOutput
	// 遍历区块链，找到指定地址包含utxo的所有交易
	txs := bc.FindUTXOTransactions(pubKeyHash)

	for _, tx := range txs {
		for _, output := range tx.TXOutputs {
			if bytes.Equal(pubKeyHash, output.PubKeyHash) {
				UTXO = append(UTXO, output)
			}
		}
	}

	return UTXO
}

// 根据转账金额找到合理的utxo
func (bc *BlockChain) FindNeedUTXOs(senderPubKeyHash []byte, amount float64) (map[string][]uint64, float64) {
	// 需要的utxos集合
	// map[交易id]output索引
	utxos := make(map[string][]uint64)
	var calc float64

	// 遍历区块链，找到指定地址包含utxo的所有交易
	txs := bc.FindUTXOTransactions(senderPubKeyHash)

	for _, tx := range txs {
		for i, output := range tx.TXOutputs {
			if bytes.Equal(senderPubKeyHash, output.PubKeyHash) {
				if calc < amount {
					utxos[string(tx.TXID)] = append(utxos[string(tx.TXID)], uint64(i))
					calc += output.Value
					if calc >= amount {
						fmt.Printf("找到了满足的金额：%f\n", calc)
						return utxos, calc
					}
				} else {
					fmt.Printf("不满足转账金额,当前总额：%f， 目标金额: %f\n", calc, amount)
				}
			}
		}
	}
	return utxos, calc
}

// 遍历区块链，找到指定地址包含utxo的所有交易
func (bc *BlockChain) FindUTXOTransactions(senderPubKeyHash []byte) []*Transaction {
	var txs []*Transaction //存储所有包含utxo交易集合
	//map[交易id][]int64
	spentOutputs := make(map[string][]int64)

	it := bc.NewIterator()

	for {
		//1.遍历区块
		block := it.Next()

		//2. 遍历交易
		for _, tx := range block.Transactions {

		OUTPUT:
			for i, output := range tx.TXOutputs {
				//在这里做一个过滤，将所有消耗过的outputs和当前的所即将添加output对比一下
				//如果相同，则跳过，否则添加
				//如果当前的交易id存在于我们已经表示的map，那么说明这个交易里面有消耗过的output
				if spentOutputs[string(tx.TXID)] != nil {
					for _, j := range spentOutputs[string(tx.TXID)] {
						if int64(i) == j {
							continue OUTPUT
						}
					}
				}

				//这个output和我们目标的地址相同，满足条件，加到返回UTXO数组中
				if bytes.Equal(output.PubKeyHash, senderPubKeyHash) {
					txs = append(txs, tx)
				}
			}

			//如果当前交易是挖矿交易的话，那么不做遍历，直接跳过
			if !tx.IsCoinbase() {
				//4. 遍历input，找到自己花费过的utxo的集合(把自己消耗过的标示出来)
				for _, input := range tx.TXInputs {
					pubKeyHash := HashPubKey(input.PubKey)
					if bytes.Equal(pubKeyHash, senderPubKeyHash) {
						spentOutputs[string(input.TXid)] = append(spentOutputs[string(input.TXid)], input.Index)
					}
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
			fmt.Printf("区块遍历完成退出!")
		}
	}

	return txs
}

// 遍历整个区块链，根据交易id找到对应交易
func (bc *BlockChain) FindTransactionByTXid(id []byte) (Transaction, error) {

	it := bc.NewIterator()

	for {
		block := it.Next()
		for _, tx := range block.Transactions {
			if bytes.Equal(tx.TXID, id) {
				return *tx, nil
			}
		}

		if len(block.PrevHash) == 0 {
			fmt.Printf("区块链遍历结束!\n")
			break
		}
	}

	return Transaction{}, errors.New("无效的交易id，请检查!")
}

func (bc *BlockChain) SignTransaction(tx *Transaction, privateKey *ecdsa.PrivateKey) {

	prevTXs := make(map[string]Transaction)

	//找到所有引用的交易
	//1. 根据inputs来找，有多少input, 就遍历多少次
	//2. 找到目标交易，（根据TXid来找）
	//3. 添加到prevTXs里面
	for _, input := range tx.TXInputs {
		//根据id查找交易本身，需要遍历整个区块链
		tx, err := bc.FindTransactionByTXid(input.TXid)

		if err != nil {
			log.Panic(err)
		}

		prevTXs[string(input.TXid)] = tx

	}

	tx.Sign(privateKey, prevTXs)
}

// 矿工校验交易中每个input的签名是否有效
func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {

	if tx.IsCoinbase() {
		// 挖矿交易没有签名，无需校验
		return true
	}

	// 存储当前交易中所有input引用的utxo所在的交易和交易id
	prevTXs := make(map[string]Transaction)

	//找到所有引用的交易
	//1. 根据inputs来找，有多少input, 就遍历多少次
	//2. 找到目标交易，（根据TXid来找）
	//3. 添加到prevTXs里面
	for _, input := range tx.TXInputs {
		// 遍历整个区块链，根据交易id找到对应交易
		tx, err := bc.FindTransactionByTXid(input.TXid)

		if err != nil {
			log.Panic(err)
		}
		// input引用的utxo所在的交易和交易id
		prevTXs[string(input.TXid)] = tx

	}
	// 对每一个签名过得input进行校验
	return tx.Verify(prevTXs)
}
