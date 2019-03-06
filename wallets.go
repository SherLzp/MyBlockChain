package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"github.com/btcsuite/btcutil/base58"
	"io/ioutil"
	"log"
	"os"
)

const walletFile = "wallet.dat"

//定义一个Wallets结构，保存所有的wallet以及它的地址
type Wallets struct {
	WalletsMap map[string]*Wallet
}

//创建方法
func NewWallets() *Wallets {
	var ws Wallets
	ws.WalletsMap = make(map[string]*Wallet)
	ws.LoadFile()
	return &ws
}

func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := wallet.NewAddress()
	//wallets.WalletsMap = make(map[string]*Wallet)
	ws.WalletsMap[address] = wallet
	ws.SaveToFile()
	return address
}

//保存方法，把新建的wallet添加进去
func (ws *Wallets) SaveToFile() {
	var buffer bytes.Buffer
	gob.Register(elliptic.P256())
	encode := gob.NewEncoder(&buffer)
	err := encode.Encode(ws)
	if err != nil {
		log.Panic(err)
	}
	ioutil.WriteFile(walletFile, buffer.Bytes(), 0600)
}

//读取文件方法，把所有的wallet读出来
func (ws *Wallets) LoadFile() {
	//读取之前确认文件是否存在
	_, err := os.Stat(walletFile)
	if os.IsNotExist(err) {
		//ws.WalletsMap = make(map[string]*Wallet)
		return
	}
	content, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}
	gob.Register(elliptic.P256())
	//解码
	decoder := gob.NewDecoder(bytes.NewReader(content))
	var wsLocal Wallets
	err = decoder.Decode(&wsLocal)
	if err != nil {
		log.Panic(err)
	}
	ws.WalletsMap = wsLocal.WalletsMap
}

func (ws *Wallets) ListAllAddresses() []string {
	var addresses []string
	for address := range ws.WalletsMap {
		addresses = append(addresses, address)
	}
	return addresses
}

//通过地址返回公钥哈希
func GetPubKeyHashFromAddress(address string) []byte {
	//1.解码
	addressByte := base58.Decode(address) //25字节
	//2.截取出公钥hash，去除version--1字节，去除校验码--4字节
	pubKeyHash := addressByte[1 : len(addressByte)-4]
	return pubKeyHash
}
