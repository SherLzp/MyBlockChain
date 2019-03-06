package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"math/big"
)

const reward = 12.5

//1.定义交易结构
type Transaction struct {
	TXID []byte //交易ID
	//交易输入数组
	TXInputs  []TXInput
	TXOutputs []TXOutput
}

//定义交易输入
type TXInput struct {
	//引用的交易ID
	TXid []byte
	//引用的output的索引
	Index int64
	//解锁脚本,我们用地址来模拟
	//Sig string
	//真正的数字签名，由r，s拼接成的字节数组
	Signature []byte
	//这里的PubKey不存储原始的公钥，而是存储X，Y拼接的字符串，在校验端重新拆分（参考r，s传递）
	PubKey []byte
}

//定义交易输出
type TXOutput struct {
	//转账金额
	Value float64
	//锁定脚本，我们用地址模拟
	//PubKeyHash string
	//收款方的公钥的hash
	PubKeyHash []byte
}

//由于现在存储的字段是地址的公钥hash，所以无法直接创建TXoutput
//为了能够得到公钥hash，我们需要处理一下，写一个Lock函数
func (output *TXOutput) Lock(address string) {
	pubKeyHash := GetPubKeyHashFromAddress(address)
	//真正的锁定动作
	output.PubKeyHash = pubKeyHash
}

//给TXOutput提供一个创建的方法，否则无法调用Lock
func NewTXOutput(value float64, address string) *TXOutput {
	output := TXOutput{
		Value: value,
	}
	output.Lock(address)
	return &output
}

//设置交易ID
func (tx *Transaction) SetHash() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	data := buffer.Bytes()
	hash := sha256.Sum256(data)
	tx.TXID = hash[:]
}

//实现一个函数，判断当前交易是否为挖矿交易
func (tx *Transaction) IsCoinbase() bool {
	//1.交易的input只有一个
	//2.交易id为nil
	//3.交易的index为-1
	if len(tx.TXInputs) == 1 && len(tx.TXInputs[0].TXid) == 0 && tx.TXInputs[0].Index == -1 {
		return true
	}
	return false
}

//2.创建交易
func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {
	//1.创建交易之后要进行数字签名，所以需要私钥->打开钱包(NewWallets())
	ws := NewWallets()
	//2.根据地址找到自己的wallet
	wallet := ws.WalletsMap[from]
	if wallet == nil {
		fmt.Println("没有找到该地址的钱包，交易创建失败!")
		return nil
	}
	//3.得到对应的公钥私钥
	pubKey := wallet.PubKey
	privateKey := wallet.Private
	//传递公钥的hash
	pubKeyHash := HashPubKey(pubKey)
	//1.找到最合理的UTXO集合 map[string]uint64
	utxos, resValue := bc.FindNeedUTXOS(pubKeyHash, amount)
	if resValue < amount {
		fmt.Println("余额不足，交易失败")
		return nil
	}

	var inputs []TXInput
	var outputs []TXOutput
	//2.将这些UTXO逐一转成inputs
	for id, indexArray := range utxos {
		for _, i := range indexArray {
			input := TXInput{[]byte(id), int64(i), nil, pubKey}
			inputs = append(inputs, input)
		}
	}
	//3.创建outputs
	//output := TXOutput{amount, to}
	output := NewTXOutput(amount, to)
	outputs = append(outputs, *output)
	//4.如果有零钱需要找零
	if resValue > amount {
		output = NewTXOutput(resValue-amount, from)
		outputs = append(outputs, *output)
	}

	tx := Transaction{[]byte{}, inputs, outputs}
	tx.SetHash()

	bc.SignTransaction(&tx, privateKey)
	return &tx
}

//3.创建挖矿交易
func NewCoinbaseTX(address string, data string) *Transaction {
	//挖矿交易特点
	//1.只有一个input
	//2.无需引用交易id
	//3.无需引用index
	//矿工由于挖矿时无需指定签名，所以PubKey字段可以由矿工自由填写
	//签名先填写为空
	input := TXInput{[]byte{}, -1, nil, []byte(data)}
	//output := TXOutput{reward, address}
	//新的创建方法
	output := NewTXOutput(reward, address)
	//对于挖矿交易来说只有一个input和一个output
	tx := Transaction{[]byte{}, []TXInput{input}, []TXOutput{*output}}
	tx.SetHash()

	return &tx
}

//4.根据交易调整程序

//签名的具体实现,参数为：私钥，inputs里面所有引用的交易的结构map[string]Transaction
//map[A] TransactionA
func (tx *Transaction) Sign(privateKey *ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	//对coinbase交易不签名
	if tx.IsCoinbase() {
		return
	}
	//1.创建一个当期交易的副本:txCopy，使用函数:TrimmedCopy：要把Signature和PubKey字段设置为nil
	txCopy := tx.TrimmedCopy()
	//2.循环遍历txCopy的inputs，得到input所引用的output的公钥哈希
	for i, input := range txCopy.TXInputs {
		prevTX := prevTXs[string(input.TXid)]
		if len(prevTX.TXID) == 0 {
			log.Panic("引用的交易无效")
		}
		//不要对input进行赋值，这是一个副本，要对txCopy.TXInputs[i]进行操作
		txCopy.TXInputs[i].PubKey = prevTX.TXOutputs[input.Index].PubKeyHash
		//所需要的三个数据都具备了，开始做哈希处理
		//3.生成要签名的数据，要签名的数据一定是哈希值
		txCopy.SetHash()
		//还原，以免影响后面input的签名
		//a.我们对每一个input都要签名一次，签名数据是由当前input引用的output的哈希+当前的outputs（都在当前tx的副本里）
		//b.要对拼好的txCopy进行哈希处理，SetHash得到TXID，这个TXID就是我们要签名的最终数据
		txCopy.TXInputs[i].PubKey = nil
		signDataHash := txCopy.TXID
		//4.进行签名动作
		r, s, err := ecdsa.Sign(rand.Reader, privateKey, signDataHash)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)
		//5.放到我们的input的Signature中
		txCopy.TXInputs[i].Signature = signature
	}

}

func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput
	for _, input := range tx.TXInputs {
		inputs = append(inputs, TXInput{input.TXid, input.Index, nil, nil})
	}

	for _, output := range tx.TXOutputs {
		outputs = append(outputs, output)
	}
	return Transaction{tx.TXID, inputs, outputs}
}

//分析校验过程
//所需要的数据：公钥、数据（txCopy、生成哈希）签名
//我们要对每一个签名过得input进行校验
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	//1.得到签名数据
	txCopy := tx.TrimmedCopy()
	for i, input := range tx.TXInputs {
		prevTX := prevTXs[string(input.TXid)]
		if len(prevTX.TXID) == 0 {
			log.Panic("引用的交易无效")
		}
		txCopy.TXInputs[i].PubKey = prevTX.TXOutputs[input.Index].PubKeyHash
		txCopy.SetHash()
		dataHash := txCopy.TXID
		//2.得到Signature,反推r，s
		signature := input.Signature //r,s
		pubKey := input.PubKey       //拆x,y
		r := big.Int{}
		s := big.Int{}

		r.SetBytes(signature[:len(signature)/2])
		s.SetBytes(signature[len(signature)/2:])

		//3.拆解PubKey，得到x，y
		X := big.Int{}
		Y := big.Int{}

		X.SetBytes(pubKey[:len(pubKey)/2])
		Y.SetBytes(pubKey[len(pubKey)/2:])
		pubKeyOrigin := ecdsa.PublicKey{elliptic.P256(), &X, &Y}
		//4.Verify
		if !ecdsa.Verify(&pubKeyOrigin, dataHash, &r, &s) {
			return false
		}
	}
	return true
}
