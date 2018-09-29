package main

import (
	"os"
	"fmt"
	"flag"
)

const usage = `
	i --address ADDRESS				"init blockChain"
	a --data DATA 					"add block to blockChain"
	s --from FROM --to TO --amount AMOUNT 		"send coin AMOUNT from FROM to TO"
	b --address ADDRESS				"get balance of the address"
	p 						"print blockChain"
`

const InitBlockChainCmdString = "i"
const AddBlockCmdString = "a"
const PrintChainCmdString = "p"
const GetBalanceCmdString = "b"
const SendCoinCmdString = "s"

type CommandLine struct {
}

func (commandLine *CommandLine) PrintUsage() {
	fmt.Println("invalid input!")
	fmt.Println(usage)
	os.Exit(1)
}

func (commandLine *CommandLine) CheckParams() {
	if len(os.Args) < 2 {
		commandLine.PrintUsage()
	}
}

func (commandLine *CommandLine) Run() {
	commandLine.CheckParams()

	addBlockCmd := flag.NewFlagSet(AddBlockCmdString, flag.ExitOnError)
	printChainCmd := flag.NewFlagSet(PrintChainCmdString, flag.ExitOnError)
	initChainCmd := flag.NewFlagSet(InitBlockChainCmdString, flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet(GetBalanceCmdString, flag.ExitOnError)
	//sendCoinCmd := flag.NewFlagSet(SendCoinCmdString, flag.ExitOnError)

	addBlockParam := addBlockCmd.String("data", "", "block transaction info")
	initChainParam := initChainCmd.String("address", "", "miner address")
	getBalanceParam := getBalanceCmd.String("address", "", "miner address")
	//sendCoinParam := sendCoinCmd.String("address", "", "miner address info")

	switch os.Args[1] {

	case InitBlockChainCmdString:
		err := initChainCmd.Parse(os.Args[2:])
		CheckError(err)
		if initChainCmd.Parsed() {
			if *initChainParam == "" {
				fmt.Println("address data should not be empty")
				commandLine.PrintUsage()
			}

			commandLine.InitBlockChain(*initChainParam)
		}

	case AddBlockCmdString:
		err := addBlockCmd.Parse(os.Args[2:])
		CheckError(err)
		if addBlockCmd.Parsed() {
			if *addBlockParam == "" {
				fmt.Println("addBlock data should not be empty")
				commandLine.PrintUsage()
			}
			commandLine.AddBlock(*addBlockParam)
		}

	case PrintChainCmdString:
		err := printChainCmd.Parse(os.Args[2:])
		CheckError(err)
		if printChainCmd.Parsed() {
			commandLine.PrintChain()
		}

	case GetBalanceCmdString:
		err := getBalanceCmd.Parse(os.Args[2:])
		CheckError(err)
		if getBalanceCmd.Parsed() {
			if *getBalanceParam == "" {
				fmt.Println("address data should not be empty")
				commandLine.PrintUsage()
			}

			commandLine.getBalance(*getBalanceParam)
		}

	//case SendCoinCmdString:
	//	err := printChainCmd.Parse(os.Args[2:])
	//	CheckError(err)
	//	if printChainCmd.Parsed() {
	//		commandLine.PrintChain()
	//	}

	default:
		commandLine.PrintUsage()
	}
}