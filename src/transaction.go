package main

import (
	"bytes"
	"encoding/gob"
	"crypto/sha256"
	"fmt"
	"os"
)

const coinBaseReward = 12.5

type Transaction struct {
	TXID []byte
	TXInput []TXInput
	TXOutput []TXOutput
}

type TXInput struct {
	TXID []byte
	Vout int64
	Script string
}

func (input *TXInput) CanUnlockUTXOWith(script string) bool {
	return input.Script == script
}

type TXOutput struct {
	Value float64
	Script string
}

func (output TXOutput) CanUnlockUTXOWith(script string) bool {
	return output.Script == script
}

func (tx *Transaction) SetTXID() {
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	CheckError(err)
	data := buffer.Bytes()
	hash := sha256.Sum256(data)

	tx.TXID = hash[:]
}

// 挖矿奖励
func NewCoinBaseTransaction(address string, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("reward to %s %d btc", address, coinBaseReward)
	}

	tx := Transaction{
		TXID: []byte{},
		TXInput: []TXInput{{[]byte{}, -1, data}},
		TXOutput: []TXOutput{{Value: coinBaseReward, Script: address}},
	}
	tx.SetTXID()
	return &tx
}

func (tx *Transaction) IsCoinBaseTransaction() bool {
	return len(tx.TXInput) == 1 && len(tx.TXInput[0].TXID) == 0 && tx.TXInput[0].Vout == -1
}

func NewTransaction(from, to string, amount float64, blockChain *BlockChain) *Transaction {
	suitableUTXOArray, total := blockChain.FindSuitableUTXO(from, amount)

	var inputArray []TXInput
	var outputArray []TXOutput

	for _, utxo := range suitableUTXOArray {
		inputArray = append(inputArray, TXInput{

		})
	}

	outputArray = append(outputArray, TXOutput{
		amount,
		to,
	})

	if total > amount { // 找零
		outputArray = append(outputArray, TXOutput{
			total - amount,
			from,
		})
	}


	tx := Transaction{
		TXID: []byte{},
		TXInput: inputArray,
		TXOutput: outputArray,
	}
	tx.SetTXID()
	return &tx
}