package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"crypto/sha256"
	"bytes"

	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)
const walletFile = "wallet.dat"
const addressChecksumLen = 4

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey	[]byte
}


// 随机获取一对公钥，私钥
func newKeyPair()(ecdsa.PrivateKey,[]byte)  {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil{
		log.Panic(err)
	}
	pubKey := append(private.PublicKey.X.Bytes(),private.PublicKey.Y.Bytes()...)
	return *private,pubKey
}

// 生成一个钱包
func NewWallet() *Wallet  {
	private, public := newKeyPair()
	wallet := Wallet{private,public}
	return &wallet
}

// hash两遍pubKey, RIPEMD160(SHA256(PubKey))
func HashPubKey(pubKey []byte) []byte  {
	publicSHA256 := sha256.Sum256(pubKey)
	RIPEMD160HASHER := ripemd160.New()
	_,err := RIPEMD160HASHER.Write(publicSHA256[:])
	if err != nil{
		log.Panic(err)
	}
	publicRIPEMD160 := RIPEMD160HASHER.Sum(nil)
	return publicRIPEMD160
}

func checksum(payload []byte)[]byte  {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])
	return secondSHA[:addressChecksumLen]
}
/*
1.Take the public key and hash it twice with RIPEMD160(SHA256(PubKey)) hashing algorithms.
2.Prepend the version of the address generation algorithm to the hash.
3.Calculate the checksum by hashing the result of step 2 with SHA256(SHA256(payload)). The checksum is the first four bytes of the resulted hash.
4.Append the checksum to the version+PubKeyHash combination.
5.Encode the version+PubKeyHash+checksum combination with Base58.
 */
func (w *Wallet) GetAddress() []byte{
	pubKeyHash := HashPubKey(w.PublicKey)
	versionedPayload := append([]byte{version},pubKeyHash...)
	checksum := checksum(versionedPayload)
	fullPayload := append(versionedPayload,checksum...)
	address := Base58Encode(fullPayload)
	return address
}

func ValidateAddress(address string) bool{
	pubKeyHash := Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash) - addressChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1:len(pubKeyHash)-addressChecksumLen]
	targetChecksum := checksum(append([]byte{version},pubKeyHash...))
	return bytes.Compare(actualChecksum,targetChecksum) == 0
}