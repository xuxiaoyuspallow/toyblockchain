package main

import (
	"os"
	"bytes"
	"encoding/binary"
	"log"
)

func DbExist() bool {
	_,err := os.Stat(dbFile)
	if err != nil{
		return false
	}
	return true
}


func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
