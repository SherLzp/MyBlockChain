package main

import "fmt"

func (cli *CLI) GetBalance(address string) {
	//1.校验地址
	if !IsValidAddress(address) {
		fmt.Printf("地址无效: %s\n", address)
		return
	}
	//2.生成公钥哈希
	pubKeyHash := GetPubKeyHashFromAddress(address)
	utxos := cli.bc.FindUTXOs(pubKeyHash)
	total := 0.0
	for _, utxo := range utxos {
		total += utxo.Value
	}
	fmt.Printf("\"%s\"的余额为: %f\n", address, total)
}

func (cli *CLI) Send(from, to string, amount float64, miner, data string) {
	//fmt.Printf("from: %s,to: %s,amount: %f,miner: %s,data: %s\n", from, to, amount, miner, data)
	//1.校验地址
	if !IsValidAddress(from) {
		fmt.Printf("地址无效 from: %s\n", from)
		return
	}
	//1.校验地址
	if !IsValidAddress(to) {
		fmt.Printf("地址无效 to: %s\n", to)
		return
	}
	//1.校验地址
	if !IsValidAddress(miner) {
		fmt.Printf("地址无效 miner: %s\n", miner)
		return
	}
	//1.创建挖矿交易
	coinbase := NewCoinbaseTX(miner, data)
	//2.创建普通交易
	tx := NewTransaction(from, to, amount, cli.bc)
	if tx == nil {
		fmt.Printf("无效的交易")
		return
	}
	//3.将交易添加到区块
	cli.bc.AddBlock([]*Transaction{coinbase, tx})
	fmt.Printf("转账成功!\n")
}

func (cli *CLI) NewWallet() {
	//wallet := NewWallet()
	//address := wallet.NewAddress()
	ws := NewWallets()
	address := ws.CreateWallet()
	fmt.Printf("地址: %s\n", address)
	//for address := range ws.WalletsMap {
	//	fmt.Printf("地址: %s\n", address)
	//}
	//fmt.Printf("私钥: %v\n", wallet.Private)
	//fmt.Printf("公钥: %v\n", wallet.PubKey)
}

func (cli *CLI) ListAddresses() {
	ws := NewWallets()
	addresses := ws.ListAllAddresses()
	for _, address := range addresses {
		fmt.Printf("地址: %s\n", address)
	}
}
