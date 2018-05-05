package main

import (
	"fmt"
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"encoding/hex"
)

const subsidy = 10   //发现一个新区块时的奖励


type TXOutput struct {
	value int		//  交易的数值
	ScriptPubkey string   //加密value的公玥
}

type TXInput struct {
	Txid []byte   //transaction id
	Vout int    // stores an index of an output in the transaction
	ScripSig string  // 解密ScriptPubkey的私玥
}

// 交易结构体
type Transaction struct {
	ID []byte	// 交易Id
	Vin []TXInput
	Vout []TXOutput
}

// no input and input.Vout == -1, we consider it a coinbase transaction
func (tx *Transaction)isCoinbase()  bool{
	return len(tx.Vin)==1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

//序列化transaction之后hash得到交易ID
func (tx *Transaction) SetID()  {
	var encoded bytes.Buffer
	var hash [32]byte

	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	if err != nil{
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func (in *TXInput)CanUnlockOutputWith(sig string) bool {
	return sig == in.ScripSig
}

func (out *TXOutput)CanUnlockInputWith(pubkey string) bool  {
	return pubkey == out.ScriptPubkey
}


/*先有input还是先有output是一个鸡生蛋蛋生鸡的问题，但在bitcoin里面，先有output，因为每个区块的
第一笔交易是矿工奖励给自己的，称为coinbase交易
*/
func NewCoinbaseTX(to, data string) *Transaction{
	if data == ""{
		data = fmt.Sprint("Reward to '%s'",to)
	}
	txin := TXInput{[]byte{},-1,data}
	txout := TXOutput{subsidy,to}
	transaction := Transaction{nil,[]TXInput{txin},[]TXOutput{txout}}
	transaction.SetID()
	return &transaction
}



func NewUTXOTransaction(from, to string,amount int,bc *BlockChain) *Transaction{
	var inputs []TXInput
	var outputs []TXOutput
	acc, validOutputs := bc.FindSpendableOutputs(from,amount)
	if acc < amount{
		log.Panic("ERROR:NOT ENOUGN FUNDS")
	}
	for txid,outs := range validOutputs{
		txidstring, err := hex.DecodeString(txid)
		if err != nil{
			log.Panic(err)
		}
		for _, out := range outs{
			input := TXInput{txidstring,out,from}
			inputs = append(inputs,input)
		}
	}
	outputs = append(outputs,TXOutput{amount,to})
	if acc > amount{
		outputs = append(outputs,TXOutput{acc-amount,from})
	}
	tx := Transaction{nil,inputs,outputs}
	tx.SetID()
	return &tx
}
