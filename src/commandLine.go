package main

import (
	"os"
	"fmt"
	"flag"
)

const usage = `
	i --address ADDRESS				"init blockChain"
	s --from FROM --to TO --amount AMOUNT 		"send coin AMOUNT from FROM to TO"
	b --address ADDRESS				"get balance of the address"
	p 						"print blockChain"
`

const InitBlockChainCmdString = "i"
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

	printChainCmd := flag.NewFlagSet(PrintChainCmdString, flag.ExitOnError)
	initChainCmd := flag.NewFlagSet(InitBlockChainCmdString, flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet(GetBalanceCmdString, flag.ExitOnError)
	sendCoinCmd := flag.NewFlagSet(SendCoinCmdString, flag.ExitOnError)

	initChainParam := initChainCmd.String("address", "", "miner address")
	getBalanceParam := getBalanceCmd.String("address", "", "miner address")

	sendFromParam := sendCoinCmd.String("from", "", "sender address info")
	sendToParam := sendCoinCmd.String("to", "", "receiver address info")
	sendAmountParam := sendCoinCmd.Float64("amount", 0, "send amount info")

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

			commandLine.GetBalance(*getBalanceParam)
		}

	case SendCoinCmdString:
		err := sendCoinCmd.Parse(os.Args[2:])
		CheckError(err)
		if sendCoinCmd.Parsed() {
			if *sendFromParam == "" {
				fmt.Println("from address data should not be empty")
				commandLine.PrintUsage()
			}

			if *sendToParam == "" {
				fmt.Println("to address data should not be empty")
				commandLine.PrintUsage()
			}

			if *sendAmountParam == 0 {
				fmt.Println("amount should not be 0")
				commandLine.PrintUsage()
			}
		}

		commandLine.Send(*sendFromParam, *sendToParam, *sendAmountParam)

	default:
		commandLine.PrintUsage()
	}
}