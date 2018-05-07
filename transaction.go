package main

import (
	"fmt"
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"encoding/hex"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/elliptic"
	"math/big"
)

const subsidy = 10   //发现一个新区块时的奖励

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

/*先有input还是先有output是一个鸡生蛋蛋生鸡的问题，但在bitcoin里面，先有output，因为每个区块的
第一笔交易是矿工奖励给自己的，称为coinbase交易
*/
func NewCoinbaseTX(to, data string) *Transaction{
	if data == ""{
		data = fmt.Sprint("Reward to '%s'",to)
	}
	txin := TXInput{[]byte{},-1,nil,[]byte(data)}
	txout := NewTXOutput(subsidy,to)
	transaction := Transaction{nil,[]TXInput{txin},[]TXOutput{*txout}}
	transaction.SetID()
	return &transaction
}



func NewUTXOTransaction(from, to string,amount int,bc *BlockChain) *Transaction{
	var inputs []TXInput
	var outputs []TXOutput

	wallets, err := NewWallets()
	if err != nil{
		log.Panic(err)
	}
	wallet := wallets.GetWallet(from)
	pubKeyHash := HashPubKey(wallet.PublicKey)
	acc, validOutputs := bc.FindSpendableOutputs(pubKeyHash,amount)
	if acc < amount{
		log.Panic("ERROR:NOT ENOUGN FUNDS")
	}
	for txid,outs := range validOutputs{
		txidstring, err := hex.DecodeString(txid)
		if err != nil{
			log.Panic(err)
		}
		for _, out := range outs{
			input := TXInput{txidstring,out,nil,wallet.PublicKey}
			inputs = append(inputs,input)
		}
	}
	outputs = append(outputs,*NewTXOutput(amount,to))
	if acc > amount{
		outputs = append(outputs,*NewTXOutput(acc-amount, from))
	}
	tx := Transaction{nil,inputs,outputs}
	tx.ID = tx.Hash()
	bc.SignTransaction(&tx,wallet.PrivateKey)
	return &tx
}

// Serialize Transaction with gob
func (tx Transaction)Serialize()[]byte  {
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	return encoded.Bytes()
}

// hash Transaction with sha256.Sum256
func (tx *Transaction) Hash()[]byte  {
	var hash [32]byte
	txCopy := *tx
	txCopy.ID = []byte{}
	hash = sha256.Sum256(txCopy.Serialize())
	return hash[:]
}

func (tx *Transaction)TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	for _, vin := range tx.Vin{
		inputs = append(inputs,TXInput{vin.Txid,vin.Vout,nil,nil})
	}
	for _, vout := range tx.Vout{
		outputs = append(outputs,TXOutput{vout.Value,vout.PubKeyHash})
	}
	txCopy := Transaction{tx.ID,inputs,outputs}
	return txCopy
}

func (tx *Transaction)Sign(privateKey ecdsa.PrivateKey,prevTxs map[string]Transaction)  {
	if tx.isCoinbase(){
		return
	}
	txCopy := tx.TrimmedCopy()

	for inID, vin := range txCopy.Vin{
		prevTx := prevTxs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Vin[inID].PubKey = nil

		r, s ,err := ecdsa.Sign(rand.Reader, &privateKey, txCopy.ID)
		if err != nil{
			log.Panic(err)
		}
		signature := append(r.Bytes(),s.Bytes()...)
		tx.Vin[inID].Signature = signature
	}
}

func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.isCoinbase(){
		return true
	}
	for _, vin := range tx.Vin{
		if prevTXs[hex.EncodeToString(vin.Txid)].ID == nil{
			log.Panic("Error: Previous transaction is not correct")
		}
	}
	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inID, vin := range tx.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Vin[inID].PubKey = nil

		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s) == false {
			return false
		}
	}

	return true
}
