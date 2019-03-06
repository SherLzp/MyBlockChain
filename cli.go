package main

import (
	"fmt"
	"os"
	"strconv"
)

//接收命令行参数并且控制区块链操作的文件

type CLI struct {
	bc *BlockChain
}

const Usage = `
	printChain            "print all blockchain data"
	getBalance --address ADDRESS "获取指定地址的余额"
	send FROM TO AMOUNT MINER DATA "由from转amount给to 由miner挖矿同时写入data"
	newWallet "创建一个新的钱包（私钥公钥对）"
	listAddresses "列举所有的地址"
`

//接收参数的动作，放到一个函数中
func (cli *CLI) Run() {
	//得到所有的命令
	args := os.Args
	if len(args) < 2 {
		fmt.Println(Usage)
		return
	}
	//解析命令
	cmd := args[1]
	switch cmd {
	case "printChain":
		//打印区块
		cli.bc.PrintBlockChain()
	case "getBalance":
		fmt.Printf("获取余额\n")
		//确保命令有效
		if len(args) == 4 && args[2] == "--address" {
			//获取数据
			address := args[3]
			cli.GetBalance(address)
		} else {
			fmt.Printf("添加区块参数使用不当，请检查")
			fmt.Printf(Usage)
			return
		}
	case "send":
		fmt.Printf("转账开始\n")
		if len(args) != 7 {
			fmt.Printf("参数个数错误")
			fmt.Printf(Usage)
		}
		from := args[2]
		to := args[3]
		amount, _ := strconv.ParseFloat(args[4], 64)
		miner := args[5]
		data := args[6]
		cli.Send(from, to, amount, miner, data)
	case "newWallet":
		fmt.Printf("创建新的钱包....\n")
		cli.NewWallet()
	case "listAddresses":
		fmt.Printf("列举所有地址...\n")
		cli.ListAddresses()
	default:
		fmt.Printf(Usage)
	}

	//执行相应的action
}
