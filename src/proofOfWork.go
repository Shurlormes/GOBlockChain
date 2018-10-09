package main

import (
	"math/big"
	"bytes"
	"math"
	"crypto/sha256"
	"fmt"
	"time"
)

const targetBits = 24

type ProofOfWork struct {
	block *Block
	target *big.Int // 目标值
}

func NewProofOfWork(block *Block) *ProofOfWork {

	target := big.NewInt(1)
	target.Lsh(target, uint(256 - targetBits))

	proofOfWork := ProofOfWork{
		block: block,
		target: target,
	}

	return &proofOfWork
}

func (proofOfWork *ProofOfWork) PrepareData(nonce int64) []byte {
	block := proofOfWork.block
	block.MerkleRoot = block.HashTransactions()
	tmp := [][]byte{
		IntToByte(block.Version),
		block.PrevBlockHash,
		block.MerkleRoot,
		IntToByte(block.TimeStamp),
		IntToByte(targetBits),
		IntToByte(nonce),
	}

	return bytes.Join(tmp, []byte{})
}

func (proofOfWork *ProofOfWork) Run() (int64, []byte) {
	// 1.拼装数据
	var nonce int64 = 0
	var hash [32]byte
	var hashInt big.Int

	// 2.哈希转big.Int
	fmt.Println("Begin Mining...")
	fmt.Printf("Target hash:\t0000%x\n", proofOfWork.target.Bytes())
	timeBegin := time.Now().Unix()
	for nonce < math.MaxInt64 {
		data := proofOfWork.PrepareData(nonce)
		hash = sha256.Sum256(data)

		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(proofOfWork.target) == -1 {
			timeEnd := time.Now().Unix()
			fmt.Printf("found hash:\t%x, nonce: %d, usage time: %ds\n", hash[:], nonce, timeEnd - timeBegin)
			break
		} else {
			//fmt.Printf("not found, current hash: %x, current nonce: %d\n", hash[:], nonce)
			nonce++
		}
	}
	return nonce, hash[:]
}

func (proofOfWork *ProofOfWork) IsValid() bool {
	var hashInt big.Int

	data := proofOfWork.PrepareData(proofOfWork.block.Nonce)
	hash := sha256.Sum256(data)

	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(proofOfWork.target) == -1
}