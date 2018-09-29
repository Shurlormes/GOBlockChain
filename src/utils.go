package main

import (
	"encoding/binary"
	"bytes"
	"os"
)

func IntToByte(num int64) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	CheckError(err)
	return buffer.Bytes()
}

func CheckError(err error) {
	if err != nil {
		os.Exit(1)
	}
}
