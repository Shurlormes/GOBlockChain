package main

import "fmt"


func (commandLine *CommandLine) InitBlockChain(address string) {
	blockChain := InitBlockChain(address)
	defer blockChain.db.Close()
	fmt.Println("Initialize blockChain successfully")
}


func (commandLine *CommandLine) AddBlock(data string) {
	blockChain := GetBlockChainHandler()
	defer blockChain.db.Close()
	//GetBlockChainHandler().AddBlock(data)
}

func (commandLine *CommandLine) PrintChain() {
	blockChain := GetBlockChainHandler()
	defer blockChain.db.Close()

	it := blockChain.NewIterator()
	for {
		block := it.Next()
		block.PrintBlock()

		if len(block.PrevBlockHash) == 0 {
			fmt.Println("\nprint over!")
			break
		}
	}
}


func (commandLine *CommandLine) getBalance(address string) {
	blockChain := GetBlockChainHandler()
	defer blockChain.db.Close()
	fmt.Printf("The balance of %s is %f\n", address, blockChain.GetBalance(address))
}