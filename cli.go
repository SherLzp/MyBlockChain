package main

import (
	"fmt"
	"os"
)

//这是一个用来接收命令行参数并且控制区块链操作的文件

type CLI struct {
	//bc *BlockChain
}

const Usage = `
	createBlockChain 地址     "创建区块链"
	addBlock --data DATA     "添加区块"
	printChain               "正向打印区块链"
	printChainR              "反向打印区块链"
`

//接受参数的动作，我们放到一个函数中
func (cli *CLI) Run() {

	//./block printChain
	//./block addBlock --data "HelloWorld"
	//1. 得到所有的命令
	args := os.Args
	if len(args) < 2 {
		fmt.Printf(Usage)
		return
	}

	//2. 分析命令
	cmd := args[1]
	switch cmd {
	case "createBlockChain":
		fmt.Printf("创建区块链!\n")
		address := os.Args[2]
		cli.CreateBlockChain(address)
	case "addBlock":
		//3. 执行相应动作
		fmt.Printf("添加区块\n")

		//确保命令有效
		if len(args) == 4 && args[2] == "--data" {
			//获取命令的数据
			//a. 获取数据
			data := args[3]
			//b. 使用bc添加区块AddBlock
			cli.AddBlock(data)
		} else {
			fmt.Printf("添加区块参数使用不当，请检查")
			fmt.Printf(Usage)
		}
	case "printChain":
		fmt.Printf("正向打印区块\n")
		cli.PrinBlockChain()
	case "printChainR":
		fmt.Printf("反向打印区块\n")
		cli.PrinBlockChainReverse()
	default:
		fmt.Printf("无效的命令，请检查!\n")
		fmt.Printf(Usage)
	}
}
