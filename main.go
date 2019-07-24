package main

func main() {
	// 创建区块链
	bc := NewBlockChain("14PxkwD8cTpzNAT1PYXRwK4qRNbkBVtgFP")
	// 创建命令行客户端
	cli := CLI{bc}
	// 运行命令行客户端
	cli.Run()
}
