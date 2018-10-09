package main

import (
	"time"
	"fmt"
	"encoding/gob"
	"bytes"
	"crypto/sha256"
)

type Block struct {
	Version int64 // 版本
	PrevBlockHash []byte // 前区块哈希
	MerkleRoot []byte // 梅克尔跟
	TimeStamp int64 // 时间戳
	Bits int64 // 难度值
	Nonce int64 // 随机数

	Hash         []byte         // 当前区块哈希
	Transactions []*Transaction // 交易信息
}

func (block *Block) Serialize() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(block)
	CheckError(err)
	return buffer.Bytes()
}

func DeSerialize(data [] byte) *Block {
	if len(data) == 0 {
		return nil
	}

	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	CheckError(err)
	return &block
}

func NewGenesisBlock(coinBase *Transaction) *Block {
	return NewBlock([]*Transaction{coinBase}, []byte{})
}

func NewBlock(data []*Transaction, prevBlockHash []byte) *Block {
	var block Block

	block = Block{
		Version:       1,
		PrevBlockHash: prevBlockHash,
		MerkleRoot:    []byte{},
		TimeStamp:     time.Now().Unix(),
		Bits:          targetBits,
		Nonce:         0,
		Transactions:  data,
		Hash:          []byte{},
	}

	proofOfWork := NewProofOfWork(&block)
	nonce, hash := proofOfWork.Run()
	block.Nonce = nonce
	block.Hash = hash

	return &block
}

func (block *Block) PrintBlock() {
	fmt.Println()
	fmt.Printf("Version: %d\n", block.Version)
	fmt.Printf("PrevBlockHash: %x\n", block.PrevBlockHash)
	fmt.Printf("MerkleRoot: %x\n", block.MerkleRoot)
	fmt.Printf("TimeStamp: %d\n", block.TimeStamp)
	fmt.Printf("Bits: %d\n", block.Bits)
	fmt.Printf("Nonce: %d\n", block.Nonce)
	//fmt.Printf("Data: %s\n", block.Data)
	fmt.Printf("Hash: %x\n", block.Hash)
	fmt.Printf("IsValid: %v\n", NewProofOfWork(block).IsValid())
}

func (block *Block) HashTransactions() []byte {
	var txHash [][]byte
	for _, tx := range block.Transactions {
		txHash = append(txHash, tx.TXID)
	}

	data := bytes.Join(txHash, []byte{})
	hash := sha256.Sum256(data)
	return hash[:]

}