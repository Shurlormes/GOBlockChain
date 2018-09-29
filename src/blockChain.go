package main

import (
	"github.com/boltdb/bolt"
	"os"
	"fmt"
)

const dbFile = "blockChain.db"
const blockBucket = "bucket"
const lastHashKey = "lastHashKey"
const genesisInfo = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

type BlockChain struct {
	db *bolt.DB
	lastBlockHash []byte // 最后一个区块的哈希
}

func InitBlockChain(address string) *BlockChain {
	if isDBExists() {
		fmt.Println("BlockChain exists alreday!")
		os.Exit(1)
	}

	db, err := bolt.Open(dbFile, 0600, nil)
	CheckError(err)

	var lastBlockHash []byte

	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		coinBase := NewCoinBaseTransaction(address, genesisInfo)
		genesisBlock := NewGenesisBlock(coinBase)
		lastBlockHash = genesisBlock.Hash

		bucket, err = tx.CreateBucket([]byte(blockBucket))
		CheckError(err)
		err = bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
		CheckError(err)
		err = bucket.Put([]byte(lastHashKey), lastBlockHash)
		CheckError(err)

		return nil
	})

	return &BlockChain{db, lastBlockHash}
}

func GetBlockChainHandler() *BlockChain {
	if !isDBExists() {
		fmt.Println("Pls init blockChain first!")
		os.Exit(1)
	}

	db, err := bolt.Open(dbFile, 0600, nil)
	CheckError(err)

	var lastBlockHash []byte

	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))

		if bucket != nil {
			lastBlockHash = bucket.Get([]byte(lastHashKey))
		}

		return nil
	})

	return &BlockChain{db, lastBlockHash}
}

func isDBExists() bool {
	_, err := os.Stat(dbFile)
	return !os.IsNotExist(err)
}

func (blockChain *BlockChain) AddBlock(transactions []*Transaction) {
	block := NewBlock(transactions, blockChain.lastBlockHash)

	blockChain.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			os.Exit(1)
		}

		err := bucket.Put(block.Hash, block.Serialize())
		CheckError(err)
		err = bucket.Put([]byte(lastHashKey), block.Hash)
		CheckError(err)

		blockChain.lastBlockHash = block.Hash

		return nil
	})
}

func (blockChain *BlockChain) NewIterator() *BlockChainIterator {
	return &BlockChainIterator{blockChain.db,blockChain.lastBlockHash}
}

func (blockChain *BlockChain) FindUTXOTransaction(address string) []*Transaction {
	var UTXOTransaction []*Transaction
	var spentUTXO = make(map[string][]int64)

	it := blockChain.NewIterator()
	for {
		block := it.Next()

		for _, transaction := range block.Transactions {
			if !transaction.IsCoinBaseTransaction() {
				for _, input := range transaction.TXInput	{
					if input.CanUnlockUTXOWith(address) {
						spentUTXO[string(input.TXID)] = append(spentUTXO[string(input.TXID)], input.Vout)
					}
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	it = blockChain.NewIterator()
	OPTIONS:
	for {
		block := it.Next()

		for _, transaction := range block.Transactions {

			for index, output := range transaction.TXOutput	{
				if spentUTXO[string(transaction.TXID)] != nil {
					voutArray :=  spentUTXO[string(transaction.TXID)]
					for _, vout := range voutArray {
						if vout == int64(index) {
							continue OPTIONS
						}
					}
				}

				if output.CanUnlockUTXOWith(address) {
					UTXOTransaction = append(UTXOTransaction, transaction)
				}
			}

		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return UTXOTransaction
}

func (blockChain *BlockChain) FindUTXO(address string) []*TXOutput {
	var TXOutputArray []*TXOutput

	UTXOTransaction := blockChain.FindUTXOTransaction(address)

	for _, tx := range UTXOTransaction {
		for _, output := range tx.TXOutput {
			if output.CanUnlockUTXOWith(address) {
				TXOutputArray = append(TXOutputArray, &output)
			}
		}
	}

	return TXOutputArray
}

func (blockChain *BlockChain) FindSuitableUTXO(address string, amount float64) (map[string]int64, float64) {
	var suitableUTXO = make(map[string]int64)
	var total float64

	balance := blockChain.GetBalance(address)
	if amount > balance {
		fmt.Printf("No enough money at %s", address)
		os.Exit(1)
	}


	//transactions := blockChain.FindUTXOTransaction(address)
	//for _, transaction := range transactions {
	//	for index, output := range transaction.TXOutput {
	//		if output.CanUnlockUTXOWith()
	//	}
	//} TODO


	return suitableUTXO, total
}

func (blockChain *BlockChain) GetBalance(address string) float64 {
	TXOutputArray := blockChain.FindUTXO(address)

	var balance float64 = 0

	for _, utxo := range TXOutputArray {
		balance += utxo.Value
	}

	return balance
}


