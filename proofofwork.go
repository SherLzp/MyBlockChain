package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

//定义一个工作量证明的结构ProofOfWork
type ProofOfWork struct {
	//1. block
	block *Block
	//2. 目标值
	//采用big.Int数据结构
	target *big.Int
}

//2. 提供创建POW的函数
func NewProofOfWork(block *Block) *ProofOfWork {
	pow := ProofOfWork{
		block: block,
	}

	//hashString := "0001000000000000000000000000000000000000000000000000000000000000"
	//tmp := big.Int{}
	//tmp.SetString(hashString, 16)

	//pow.target = &tmp

	//目标值：
	// 0001000000000000000000000000000000000000000000000000000000000000
	//初始值
	// 0000000000000000000000000000000000000000000000000000000000000001
	//左移256位
	//10000000000000000000000000000000000000000000000000000000000000000
	//右移4位(16进制位数)
	// 0001000000000000000000000000000000000000000000000000000000000000

	targetLocal := big.NewInt(1)
	//targetLocal.Lsh(targetLocal, 256)
	//targetLocal.Rsh(targetLocal, difficulty)
	var difficulty uint
	difficulty = 4
	targetLocal.Lsh(targetLocal, 256-difficulty)
	pow.target = targetLocal

	return &pow
}

//3. 提供计算不断计算hash的哈数
func (pow *ProofOfWork) Run() ([]byte, uint64) {
	//1. 拼装数据（区块的数据，还有不断变化的随机数）
	//2. 做哈希运算
	//3. 与pow中的target进行比较
	//a. 找到了，退出返回
	//b. 没找到，继续找，随机数加1

	var nonce uint64
	//block := pow.block
	var hash [32]byte

	for {

		fmt.Printf("hash : %x\r", hash)

		//1. 拼装数据（区块的数据，还有不断变化的随机数）
		blockInfo := pow.PrepareData(nonce)
		//2. 做哈希运算
		//func Sum256(data []byte) [Size]byte {
		hash = sha256.Sum256(blockInfo)
		//3. 与pow中的target进行比较
		tmpInt := big.Int{}
		//将我们得到hash数组转换成一个big.int
		tmpInt.SetBytes(hash[:])

		//比较当前的哈希与目标哈希值，如果当前的哈希值小于目标的哈希值，就说明找到了，否则继续找

		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		//
		//func (x *Int) Cmp(y *Int) (r int) {
		if tmpInt.Cmp(pow.target) == -1 {
			//a. 找到了，退出返回
			fmt.Printf("挖矿成功！hash : %x, nonce : %d\n", hash, nonce)
			//break
			return hash[:], nonce
		} else {
			//b. 没找到，继续找，随机数加1
			nonce++
		}
	}

	//return []byte("HelloWorld"), 10
}

//
//4. 提供一个校验函数
//
//- IsValid()

func (pow *ProofOfWork) PrepareData(nonce uint64) []byte {
	block := pow.block
	tmp := [][]byte{
		Uint64ToByte(block.Version),
		block.PrevHash,
		block.MerkelRoot,
		Uint64ToByte(block.TimeStamp),
		Uint64ToByte(block.Difficulty),
		Uint64ToByte(nonce),
		block.Data,
	}

	//将二维的切片数组链接起来，返回一个一维的切片
	blockInfo := bytes.Join(tmp, []byte{})
	return blockInfo
}
