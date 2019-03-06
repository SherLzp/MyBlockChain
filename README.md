# 		从0到1构建一条区块链----v1.0

## 基础知识

1. 基本编程知识
2. Go语言基本语法
3. 区块链的基础知识（什么是区块、什么是区块链、什么是交易、什么是POW）[https://anders.com/blockchain/tokens.html](https://anders.com/blockchain/tokens.html)
4. 比特币总体流程与基本原理----这条链基于Bitcoin原理来实现

## 开发环境

1. OS：Win10
2. Golang
3. BoltDB（存储）
4. IDE：GoLand

## 初步构建思路

1. 构造结构体（简单版减少相关字段）
2. 创建区块
3. Hash生成
4. 定义区块链
5. 添加区块

## 具体实现

![bitcoin_transaction](.\imgs_v1\bitcoin_transaction.jpg)

### V1最终目录结构

![directory_v1](.\directory_v1.png)

### 构造Block结构体block.go（参考比特币区块结构先进行缩减构造简单版）

![bitcoin_genesis](.\imgs_v1\bitcoin_genesis.jpg)

```go
//1. 定义结构
type Block struct {
	//1.版本号
	Version uint64
	//2. 前区块哈希
	PrevHash []byte
	//3. Merkel根（是一个hash值，先不考虑后续再补充，现在置空即可）
	MerkelRoot []byte
	//4. 时间戳
	TimeStamp uint64
	//5. 难度值
	Difficulty uint64
	//6. 随机数，挖矿需要找到的随机数
	Nonce uint64

	//7. 当前区块hash，实际比特币的结构中是不存在的，后续用上了数据库再进行改进
	Hash []byte
	//8. 数据
	Data []byte
}
```

### 创建区块block.go

```go
//2.创建区块
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := Block{
		Version:    00,
		PrevHash:   prevBlockHash,
		MerkelRoot: []byte{},
		TimeStamp:  uint64(time.Now().Unix()),
		Difficulty: 0, //随便填写的无效值
		Nonce:      0, //同上
		Hash:       []byte{},
		Data:       []byte(data),
	}

	block.SetHash()

	return &block
}
```

### Hash生成block.go

我们现有的目的是为了生成当前区块的hash，而当前区块hash的计算是根据区块的数据计算而得，所以我们现在开始考虑实现一个Hash函数的生成方法。

在hash生成之前，我们需要一个辅助方法，能够将uint64转换为[]byte，这样方便数据拼装

```go
//辅助函数，功能是将uint64转成[]byte
func Uint64ToByte(num uint64) []byte {
	var buffer bytes.Buffer

	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buffer.Bytes()
}
```

有了辅助函数之后，再去考虑实现Hash的生成

```go
//3. 生成哈希
func (block *Block) SetHash() {
    //构建一个二维切片数组
	tmp := [][]byte{
		Uint64ToByte(block.Version),
		block.PrevHash,
		block.MerkelRoot,
		Uint64ToByte(block.TimeStamp),
		Uint64ToByte(block.Difficulty),
		Uint64ToByte(block.Nonce),
		block.Data,
	}

	//将二维的切片数组链接起来，返回一个一维的切片
	blockInfo := bytes.Join(tmp, []byte{})

	//2. sha256来计算hash
	//func Sum256(data []byte) [Size]byte
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
}
```

### 定义区块链结构体blockchain.go

```go
//4. 引入区块链
type BlockChain struct {
	//定一个区块链数组
	blocks []*Block
}
```

### 创建区块链blockchain.go

创建区块链之前，我们需要一个创世块

```go
//定义一个创世块
func GenesisBlock() *Block {
	return NewBlock("I'm genesis!", []byte{})
}
```

有了创世块之后，我们开始考虑实现新建区块链的方法

```go
//5. 创建区块链
func NewBlockChain() *BlockChain {
	//创建一个创世块，并作为第一个区块添加到区块链中
	genesisBlock := GenesisBlock()
	return &BlockChain{
		blocks: []*Block{genesisBlock},
	}
}
```

### 区块链中加入区块的方法blockchain.go

这里留有一个疑惑后续会进行解决，就是如何才能获取前一个区块的hash呢？

```go
//6. 添加区块
func (bc *BlockChain) AddBlock(data string) {
	//TODO 如何获取前区块的哈希呢？？

	//获取最后一个区块
	lastBlock := bc.blocks[len(bc.blocks)-1]
	prevHash := lastBlock.Hash

	//a. 创建新的区块
	block := NewBlock(data, prevHash)
	//b. 添加到区块链数组中
	bc.blocks = append(bc.blocks, block)
}
```



### V1版本实现结束，进行测试main.go

```go
func main() {
	bc := NewBlockChain()
	bc.AddBlock("A向B转账2BTC")
	bc.AddBlock("A向B转账3BTC")

	for i, block := range bc.blocks {

		fmt.Printf("======== 当前区块高度： %d ========\n", i)
		fmt.Printf("前区块哈希值： %x\n", block.PrevHash)
		fmt.Printf("当前区块哈希值： %x\n", block.Hash)
		fmt.Printf("区块数据 :%s\n", block.Data)
	}
}
```

