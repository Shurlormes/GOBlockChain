package main

import (
	"github.com/boltdb/bolt"
	"os"
)

type BlockChainIterator struct {
	db *bolt.DB
	currentHash []byte
}

func (blockChainIterator *BlockChainIterator) Next() *Block {
	var currentBlock *Block
	blockChainIterator.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			os.Exit(1)
		}

		data := bucket.Get(blockChainIterator.currentHash)
		currentBlock = DeSerialize(data)

		blockChainIterator.currentHash = currentBlock.PrevBlockHash

		return nil
	})

	return currentBlock
}