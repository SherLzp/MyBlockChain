package main

/**
1.定义结构
2.前区块hash
3.当前区块hash
4.数据
5.创建区块
6.生成hash
7.引入区块链
8.添加区块
9.重构代码
*/

func main() {
	bc := NewBlockChain("18fh8wzXAzP9kE433CwNCQ34e4rjeDZgZN")
	cli := CLI{bc}
	cli.Run()
	//bc.AddBlock("第二个区块")
	//bc.AddBlock("第三个区块")

	//调用迭代器，返回每一个区块数据
	//it := bc.NewIterator()
	//for {
	//	//返回区块，左移
	//	block := it.Next()
	//	fmt.Println("===================================")
	//	fmt.Printf("前区块的hash值： %x\n", block.PrevHash)
	//	fmt.Printf("当前区块的hash值： %x\n", block.Hash)
	//	fmt.Printf("区块数据:  %s\n", block.Data)
	//
	//	if len(block.PrevHash) == 0 {
	//		fmt.Printf("区块遍历结束")
	//		break
	//	}
	//}
	bc.db.Close()
}
