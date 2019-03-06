package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
	"log"
)

//这里的钱包是一个结构，每一个钱包保存了公钥私钥对
type Wallet struct {
	Private *ecdsa.PrivateKey
	//PubKey *ecdsa.PublicKey
	//这里的PubKey不存储原始的公钥，而是存储X，Y拼接的字符串，在校验端重新拆分（参考r，s传递）
	PubKey []byte
}

//创建钱包
func NewWallet() *Wallet {
	curve := elliptic.P256()
	//生成私钥
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic()
	}
	//生成公钥
	pubKeyOrigin := privateKey.PublicKey
	pubKey := append(pubKeyOrigin.X.Bytes(), pubKeyOrigin.Y.Bytes()...)
	return &Wallet{
		Private: privateKey,
		PubKey:  pubKey,
	}
}

//生成地址
func (w *Wallet) NewAddress() string {
	pubKey := w.PubKey
	rip160HashValue := HashPubKey(pubKey)
	version := byte(00)
	//拼接version
	payload := append([]byte{version}, rip160HashValue...)
	//checksum
	checkCode := CheckSum(payload)
	//25字节数据
	payload = append(payload, checkCode...)
	//base58生成地址
	address := base58.Encode(payload)
	return address
}

func HashPubKey(data []byte) []byte {
	hash := sha256.Sum256(data)
	rip160hasher := ripemd160.New()
	_, err := rip160hasher.Write(hash[:])
	if err != nil {
		log.Panic()
	}
	//返回rip160的hash结果
	rip160HashValue := rip160hasher.Sum(nil)
	return rip160HashValue
}

func CheckSum(data []byte) []byte {
	//两次sha256
	hash1 := sha256.Sum256(data)
	hash2 := sha256.Sum256(hash1[:])
	//前4字节校验码
	checkCode := hash2[:4]

	return checkCode
}

func IsValidAddress(address string) bool {
	//1.解码
	addressByte := base58.Decode(address)
	if len(addressByte) < 4 {
		return false
	}
	//2.截取数据
	payload := addressByte[:len(addressByte)-4]
	checkSum1 := addressByte[len(addressByte)-4:]
	//3.做checkSum函数
	checkSum2 := CheckSum(payload)
	//4.比较
	return bytes.Equal(checkSum1, checkSum2)
}
