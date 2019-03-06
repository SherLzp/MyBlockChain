package main

import "fmt"

func (cli *CLI) CreateBlockChain(address string) {
	bc := CreateBlockChain()

	if bc == nil {
		return
	}

	defer bc.db.Close()
}

func (cli *CLI) AddBlock(data string) {
	bc := NewBlockChain()

	if bc == nil {
		return
	}

	defer bc.db.Close()

	bc.AddBlock(data)
	fmt.Printf("添加区块成功！\n")
}

//正向打印
func (cli *CLI) PrinBlockChain() {
	bc := NewBlockChain()

	if bc == nil {
		return
	}

	defer bc.db.Close()
	bc.Printchain()
	fmt.Printf("打印区块链完成\n")
}

//反向打印
func (cli *CLI) PrinBlockChainReverse() {
	bc := NewBlockChain()

	if bc == nil {
		return
	}

	defer bc.db.Close()

	//创建迭代器
	it := bc.NewIterator()

	//调用迭代器，返回我们的每一个区块数据
	for {
		//返回区块，左移
		block := it.Next()

		fmt.Printf("===========================\n\n")
		fmt.Printf("版本号: %d\n", block.Version)
		fmt.Printf("前区块哈希值: %x\n", block.PrevHash)
		fmt.Printf("梅克尔根: %x\n", block.MerkelRoot)
		fmt.Printf("时间戳: %d\n", block.TimeStamp)
		fmt.Printf("难度值: %d\n", block.Difficulty)
		fmt.Printf("随机数 : %d\n", block.Nonce)
		fmt.Printf("当前区块哈希值: %x\n", block.Hash)
		fmt.Printf("区块数据 :%s\n", block.Data)

		if len(block.PrevHash) == 0 {
			fmt.Printf("区块链遍历结束！")
			break
		}
	}
}
