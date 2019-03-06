package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

//1.定义一个结构
//block
//目标值
type ProofOfWork struct {
	//block
	block *Block
	//目标值
	target *big.Int
}

//2.创建POW
func NewProofOfWork(block *Block) *ProofOfWork {
	pow := ProofOfWork{
		block: block,
	}
	//指定难度值，String类型，需要进行转换
	targetStr := "0000100000000000000000000000000000000000000000000000000000000000"
	//引入辅助变量，将上面的难度值转成bigint
	tmpInt := big.Int{}
	//将难度值赋值给bigint，指定16进制格式
	tmpInt.SetString(targetStr, 16)
	pow.target = &tmpInt
	return &pow
}

//不断计算hash的函数
func (pow *ProofOfWork) Run() ([]byte, uint64) {
	//1.拼装数据(区块数据、随机数)
	//2.做hash运算
	//3.验证
	var nonce uint64
	block := pow.block
	var hash [32]byte
	fmt.Println("开始挖矿......")
	for {
		tmp := [][]byte{
			Uint64ToByte(block.Version),
			block.PrevHash,
			block.MerkelRoot,
			Uint64ToByte(block.TimeStamp),
			Uint64ToByte(block.Difficulty),
			Uint64ToByte(nonce),
			//只对区块头做hash值，通过MerkelRoot产生影响
			//block.Data,
		}
		//将二维的切片数组连接起来，返回一个一维的切片
		blockInfo := bytes.Join(tmp, []byte{})
		hash = sha256.Sum256(blockInfo)
		//将我们得到的hash数组转换成bigint
		tmInt := big.Int{}
		tmInt.SetBytes(hash[:])
		//比较当前的hash值与目标hash值，如果当前hash小于目标hash，则找到，否则继续运算
		if tmInt.Cmp(pow.target) == -1 {
			//找到hash值
			fmt.Printf("挖矿成功!hash: %x ,nonce: %d\n", hash, nonce)
			break
		} else {
			nonce++
		}

	}
	//return []byte("挖矿成功"), 0
	return hash[:], nonce
}

//校验函数
