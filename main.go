package main

import "fmt"

func main() {
	bc := NewBlockChain()
	bc.AddBlock("A向B转账1BTC")
	bc.AddBlock("A又向B转账3BTC")

	for i, block := range bc.blocks {

		fmt.Printf("======== 当前区块高度： %d ========\n", i)
		fmt.Printf("前区块哈希值： %x\n", block.PrevHash)
		fmt.Printf("当前区块哈希值： %x\n", block.Hash)
		fmt.Printf("区块数据 :%s\n", block.Data)
	}
}
