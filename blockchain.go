package main

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ShersBlockChain/bolt"
	"log"
	"time"
)

//4.引入区块链
//使用boltDb代替数组 重构区块链
type BlockChain struct {
	//定义一个区块链数组
	//blocks []*Block
	db   *bolt.DB
	tail []byte //存储最后一个区块的hash
}

const blockChainDb = "blockChain.db"
const blockBucket = "blockBucket"
const blockLastHashKey = "lastHashKey"

//5.定义一个区块链
func NewBlockChain(address string) *BlockChain {
	//return &BlockChain{
	//	blocks: []*Block{genesisBlock},
	//}
	var lastHash []byte
	//1.打开数据库
	db, err := bolt.Open(blockChainDb, 0600, nil)
	if err != nil {
		log.Panic("打开数据库失败")
	}
	//defer db.Close()
	db.Update(func(tx *bolt.Tx) error {
		//2.找到bucket--如果没有就创建一个
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			//没有bucket，需要创建
			bucket, err = tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				log.Panic("Bucket创建失败")
			}
			//3.写数据
			//创建一个创世区块，并作为第一个区块加入区块链
			genesisBlock := GenesisBlock(address)
			//Hash作为key，block的字节流作为value
			bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
			//修改最后一个区块的hash
			bucket.Put([]byte(blockLastHashKey), genesisBlock.Hash)
			lastHash = genesisBlock.Hash
		} else {
			lastHash = bucket.Get([]byte(blockLastHashKey))
		}
		return nil
	})
	return &BlockChain{db, lastHash}
}

//创世区块
func GenesisBlock(address string) *Block {
	coinbase := NewCoinbaseTX(address, "创世区块")
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

//6.添加区块
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	for _, tx := range txs {
		if !bc.VerifyTransaction(tx) {
			fmt.Printf("矿工发现无效交易\n")
			return
		}
	}

	//区块链数据库
	db := bc.db
	//最后一个区块的hash
	lastHash := bc.tail
	db.Update(func(tx *bolt.Tx) error {
		//完成数据添加
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panic("bucket不应该为空，请检查")
		}

		block := NewBlock(txs, lastHash)
		//更新区块链数据库--写区块
		bucket.Put(block.Hash, block.Serialize())
		bucket.Put([]byte(blockLastHashKey), block.Hash)
		//更新内存中的区块链
		bc.tail = block.Hash
		return nil
	})
}

func (bc *BlockChain) PrintBlockChain() {
	it := bc.NewIterator()
	var blockHeight = 0
	for {
		//返回区块，左移
		block := it.Next()
		fmt.Printf("==================区块高度: %d ====================\n", blockHeight)
		blockHeight++
		fmt.Printf("版本号: %d\n", block.Version)
		fmt.Printf("前区块的hash值： %x\n", block.PrevHash)
		fmt.Printf("MerkelRoot: %x\n", block.MerkelRoot)
		timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
		fmt.Printf("时间戳: %s\n", timeFormat)
		fmt.Printf("难度值: %d\n", block.Difficulty)
		fmt.Printf("随机数: %d\n", block.Nonce)
		fmt.Printf("当前区块的hash值： %x\n", block.Hash)
		fmt.Printf("区块数据:  %s\n", block.Transactions[0].TXInputs[0].PubKey)

		if len(block.PrevHash) == 0 {
			fmt.Printf("区块遍历结束\n")
			break
		}
	}
}

//找到指定地址的所有UTXO
func (bc *BlockChain) FindUTXOs(senderPubKeyHash []byte) []TXOutput {
	var UTXO []TXOutput
	txs := bc.FindUTXOTransactions(senderPubKeyHash)
	for _, tx := range txs {
		for _, output := range tx.TXOutputs {
			//if address == output.PubKeyHash {
			if bytes.Equal(senderPubKeyHash, output.PubKeyHash) {
				UTXO = append(UTXO, output)
			}
		}
	}
	return UTXO
}
func (bc *BlockChain) FindNeedUTXOS(senderPubKeyHash []byte, amount float64) (map[string][]uint64, float64) {
	//找到合理的UTXO集合
	utxos := make(map[string][]uint64)
	//找到的UTXOS包含的钱的总数
	var calc float64
	txs := bc.FindUTXOTransactions(senderPubKeyHash)
	for _, tx := range txs {
		for i, output := range tx.TXOutputs {
			//if from == output.PubKeyHash {
			//直接比较是否相同，返回true或false
			if bytes.Equal(senderPubKeyHash, output.PubKeyHash) {
				//UTXO = append(UTXO, output)
				//我们要实现的逻辑：找到自己需要的最少UTXO
				//3.比较一下是否满足转账需求 -- a.满足直接返回 b.不满足继续创建
				if calc < amount {
					//1.把utxo加入集合
					utxos[string(tx.TXID)] = append(utxos[string(tx.TXID)], uint64(i))
					//2.统计utxo当前总额
					calc += output.Value
					//如果余额满足条件
					if calc >= amount {
						fmt.Printf("找到了满足的金额: %f\n", calc)
						return utxos, calc
					}
				} else {
					fmt.Printf("不满足转账金额，当前金额为: %f,目标金额为: %f\n", calc, amount)
				}
			}
		}
	}
	return utxos, calc
}

func (bc *BlockChain) FindUTXOTransactions(senderPubKeyHash []byte) []*Transaction {
	//包含存储所有utxo交易的集合
	var txs []*Transaction
	//我们定义一个map来保存消费过的utxo，key是这个output的交易id，value是这个交易中索引的数组
	spentOutputs := make(map[string][]int64)
	//1.遍历区块
	//创建迭代器
	it := bc.NewIterator()
	for {
		block := it.Next()
		//2.遍历交易
		for _, tx := range block.Transactions {
			//fmt.Printf("current txid: %x\n", tx.TXID)
			//3.遍历output，找到和自己相关的utxo（在添加output之前检查是否已经消耗过）
		OUTPUT:
			for i, output := range tx.TXOutputs {
				//fmt.Printf("current index: %d\n", i)
				//做一个过滤，将所有消耗过的outputs和即将添加的output对比一下
				//如果相同，则跳过，否则添加
				//如果当前的交易id存在于我们已经标识的map，说明这个交易里面有消耗过得output
				//map[3333]=[]int64{0,1}
				if spentOutputs[string(tx.TXID)] != nil {
					for _, j := range spentOutputs[string(tx.TXID)] {
						if int64(i) == j {
							//当前准备添加的output已经消耗过了
							continue OUTPUT
						}
					}
				}
				//获取与目标地址相同的output，加到返回utxo数组中
				//if output.PubKeyHash == address {
				if bytes.Equal(output.PubKeyHash, senderPubKeyHash) {
					//UTXO = append(UTXO, output)
					//返回所有包含address的utxo的交易记录
					txs = append(txs, tx)
				}
			}
			//如果当前交易是挖矿交易的话，那么不做遍历，直接跳过
			if !tx.IsCoinbase() {
				//4.遍历input，找到自己花费过的utxo的集合（把自己消耗过得给标识出来）
				for _, input := range tx.TXInputs {
					//判断当前input和目标地址是否一致，如果相同，则是该地址消耗过的output，如果是就加入map
					//if input.Sig == address {
					if bytes.Equal(HashPubKey(input.PubKey), senderPubKeyHash) {
						indexArray := spentOutputs[string(input.TXid)]
						indexArray = append(indexArray, input.Index)
					}
				}
			} else {
				fmt.Printf("这是Coinbase 不做input遍历\n")
			}

		}
		if len(block.PrevHash) == 0 {
			fmt.Printf("区块遍历完成退出\n")
			break
		}
	}

	return txs
}

//根据id查找交易本身，需要遍历区块链
func (bc *BlockChain) FindTransactionByTXid(id []byte) (Transaction, error) {
	//1.遍历区块链
	it := bc.NewIterator()
	for {
		block := it.Next()
		//2.遍历交易
		for _, tx := range block.Transactions {
			//3.比较交易，找到了直接退出
			if bytes.Equal(tx.TXID, id) {
				return *tx, nil
			}
		}

		if len(block.PrevHash) == 0 {
			fmt.Printf("区块链遍历结束!\n")
			break
		}
	}
	//4.如果没找到，返回空Transaction同时返回错误状态
	return Transaction{}, errors.New("无效的交易id，请检查!")
}

func (bc *BlockChain) SignTransaction(tx *Transaction, privateKey *ecdsa.PrivateKey) {
	//签名，交易创建的最后进行签名
	prevTXs := make(map[string]Transaction)

	//找到所有引用的交易
	//1.根据inputs来找，有多少个input就遍历多少次
	//2.找到目标的交易，根据TXID来找
	//3.添加到prevTXs
	for _, input := range tx.TXInputs {
		//根据id查找交易本身，我们需要遍历整个区块链
		tx, err := bc.FindTransactionByTXid(input.TXid)
		if err != nil {
			log.Panic(err)
		}

		prevTXs[string(input.TXid)] = tx
	}

	tx.Sign(privateKey, prevTXs)
}

func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}
	//签名，交易创建的最后进行签名
	prevTXs := make(map[string]Transaction)

	//找到所有引用的交易
	//1.根据inputs来找，有多少个input就遍历多少次
	//2.找到目标的交易，根据TXID来找
	//3.添加到prevTXs
	for _, input := range tx.TXInputs {
		//根据id查找交易本身，我们需要遍历整个区块链
		tx, err := bc.FindTransactionByTXid(input.TXid)
		if err != nil {
			log.Panic(err)
		}

		prevTXs[string(input.TXid)] = tx
	}
	return tx.Verify(prevTXs)
}
