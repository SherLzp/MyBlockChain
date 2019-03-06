package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"log"
	"time"
)

//1.定义结构
type Block struct {
	//版本号
	Version uint64
	//前区块hash
	PrevHash []byte
	//Merkel根--hash值
	MerkelRoot []byte
	//时间戳
	TimeStamp uint64
	//难度值
	Difficulty uint64
	//随机数
	Nonce uint64
	//当前区块hash
	Hash []byte
	//data
	//Data []byte
	Transactions []*Transaction
}

//1.补充区块字段
//2.更新hash函数
//3.优化代码

//实现辅助函数--将uint64转成[]byte
func Uint64ToByte(num uint64) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}

//2.创建区块
func NewBlock(txs []*Transaction, prevBlockHash []byte) *Block {
	block := Block{
		Version:      00,
		PrevHash:     prevBlockHash,
		MerkelRoot:   []byte{},
		TimeStamp:    uint64(time.Now().Unix()),
		Difficulty:   0,
		Nonce:        0,
		Hash:         []byte{},
		Transactions: txs,
	}

	block.MerkelRoot = block.MakeMerkelRoot()
	//block.SetHash()
	//创建一个pow对象
	pow := NewProofOfWork(&block)
	//查找随机数，进行hash运算
	hash, nonce := pow.Run()
	//根据挖矿结果对区块数据进行补充
	block.Hash = hash
	block.Nonce = nonce
	return &block
}

//实现将block转换为字节流--序列化
func (block *Block) Serialize() []byte {
	//编码的数据丢入buffer
	var buffer bytes.Buffer
	//序列化
	//定义编码器
	encoder := gob.NewEncoder(&buffer)
	//使用编码器进行编码
	err := encoder.Encode(&block)
	if err != nil {
		log.Panic("编码出错")
	}
	//fmt.Printf("编码后的数据: %x\n", buffer.Bytes())
	return buffer.Bytes()
}

//反序列化
func DeSerialize(data []byte) Block {
	//反序列化
	//定义解码器
	decoder := gob.NewDecoder(bytes.NewReader(data))
	//使用解码器进行解码
	var block Block
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic("解码出错")
	}
	return block
}

//3.生成hash
/*
func (block *Block) SetHash() {
	//var blockInfo []byte
	//1.拼装数据
	/*
		blockInfo = append(blockInfo, Uint64ToByte(block.Version)...)
		blockInfo = append(blockInfo, block.PrevHash...)
		blockInfo = append(blockInfo, block.MerkelRoot...)
		blockInfo = append(blockInfo, Uint64ToByte(block.TimeStamp)...)
		blockInfo = append(blockInfo, Uint64ToByte(block.Difficulty)...)
		blockInfo = append(blockInfo, Uint64ToByte(block.Nonce)...)
		blockInfo = append(blockInfo, block.Data...)
*/
/*
tmp := [][]byte{
	Uint64ToByte(block.Version),
	block.PrevHash,
	block.MerkelRoot,
	Uint64ToByte(block.TimeStamp),
	Uint64ToByte(block.Difficulty),
	Uint64ToByte(block.Nonce),
	block.Data,
}
//将二维的切片数组连接起来，返回一个一维的切片
blockInfo := bytes.Join(tmp, []byte{})
//2.sha256
hash := sha256.Sum256(blockInfo)
block.Hash = hash[:]
}
*/

//模拟MerkelRoot，只是对交易数据做简单拼接，不做二叉树处理
func (block *Block) MakeMerkelRoot() []byte {
	var info []byte
	//var finalInfo [][]byte
	for _, tx := range block.Transactions {
		//将交易的hash值拼接起来，再整体求一个hash
		info = append(info, tx.TXID...)
	}
	hash := sha256.Sum256(info)
	return hash[:]
}
