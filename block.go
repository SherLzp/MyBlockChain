package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"time"
)

//0. 定义结构
type Block struct {
	//1.版本号
	Version uint64
	//2. 前区块哈希
	PrevHash []byte
	//3. Merkel根（梅克尔根，这就是一个哈希值，我们先不管，我们后面v4再介绍）
	MerkelRoot []byte
	//4. 时间戳
	TimeStamp uint64
	//5. 难度值
	Difficulty uint64
	//6. 随机数，也就是挖矿要找的数据
	Nonce uint64

	//a. 当前区块哈希,正常比特币区块中没有当前区块的哈希，我们为了是方便做了简化！
	Hash []byte
	//b. 数据
	Data []byte
}

//1. 补充区块字段
//2. 更新计算哈希函数
//3. 优化代码

//实现一个辅助函数，功能是将uint64转成[]byte
func Uint64ToByte(num uint64) []byte {
	var buffer bytes.Buffer

	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buffer.Bytes()
}

//2. 创建区块
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

//3. 生成哈希
func (block *Block) SetHash() {
	//var blockInfo []byte
	//1. 拼装数据
	/*
		blockInfo = append(blockInfo, Uint64ToByte(block.Version)...)
		blockInfo = append(blockInfo, block.PrevHash...)
		blockInfo = append(blockInfo, block.MerkelRoot...)
		blockInfo = append(blockInfo, Uint64ToByte(block.TimeStamp)...)
		blockInfo = append(blockInfo, Uint64ToByte(block.Difficulty)...)
		blockInfo = append(blockInfo, Uint64ToByte(block.Nonce)...)
		blockInfo = append(blockInfo, block.Data...)
	*/
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

	//2. sha256
	//func Sum256(data []byte) [Size]byte {
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
}
