package main

import (
	"github.com/ShersBlockChain/bolt"
	"log"
)

type BlockChainIterator struct {
	db *bolt.DB
	//游标，用于不断索引
	currentHashPointer []byte
}

func (bc *BlockChain) NewIterator() *BlockChainIterator {
	return &BlockChainIterator{
		bc.db,
		//最初指向区块链的最后一个区块hash，随着next调用不断变换
		bc.tail,
	}
}

//迭代器是属于区块链的，但next方法是属于迭代器的
func (it *BlockChainIterator) Next() *Block {
	var block Block
	//1.返回当前区块
	it.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panic("迭代器遍历时，bucket不应该为空")
		}
		blockTmp := bucket.Get(it.currentHashPointer)
		//解码
		block = DeSerialize(blockTmp)
		//游标hash左移
		it.currentHashPointer = block.PrevHash
		return nil
	})
	return &block
}
