package main

import (
	"fmt"
	//"time"
)

// 正向打印区块链
func (cli *CLI) PrinBlockChain() {
	cli.bc.Printchain()
	fmt.Printf("打印区块链完成\n")
}

// 反向打印区块链
func (cli *CLI) PrinBlockChainReverse() {

	bc := cli.bc
	it := bc.NewIterator()

	for {

		block := it.Next()

		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}

		if len(block.PrevHash) == 0 {
			fmt.Printf("区块链遍历结束！")
			break
		}
	}
}

// 获取指定地址余额
func (cli *CLI) GetBalance(address string) {

	// 校验地址
	if !IsValidAddress(address) {
		fmt.Printf("地址无效 : %s\n", address)
		return
	}

	// 生成公钥哈希
	pubKeyHash := GetPubKeyFromAddress(address)
	// 找到指定公钥哈希对应的所有utxo
	utxos := cli.bc.FindUTXOs(pubKeyHash)

	total := 0.0
	for _, utxo := range utxos {
		total += utxo.Value
	}

	fmt.Printf("\"%s\"的余额为：%f\n", address, total)
}

// 转账
func (cli *CLI) Send(from, to string, amount float64, miner, data string) {

	if !IsValidAddress(from) {
		fmt.Printf("地址无效 from: %s\n", from)
		return
	}
	if !IsValidAddress(to) {
		fmt.Printf("地址无效 to: %s\n", to)
		return
	}
	if !IsValidAddress(miner) {
		fmt.Printf("地址无效 miner: %s\n", miner)
		return
	}

	//1. 创建挖矿交易
	coinbase := NewCoinbaseTX(miner, data)
	//2. 创建一个普通交易
	tx := NewTransaction(from, to, amount, cli.bc)
	if tx == nil {
		return
	}
	//3. 添加到区块
	cli.bc.AddBlock([]*Transaction{coinbase, tx})
	fmt.Printf("转账成功！")
}

// 创建新钱包
func (cli *CLI) NewWallet() {
	ws := NewWallets()
	address := ws.CreateWallet()
	fmt.Printf("地址：%s\n", address)

}

// 遍历钱包中所有地址
func (cli *CLI) ListAddresses() {
	ws := NewWallets()
	addresses := ws.ListAllAddresses()
	for _, address := range addresses {
		fmt.Printf("地址：%s\n", address)
	}
}
