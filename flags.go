package main

import (
	"flag"
	"os"
)

var (
	depositAddressPrivateKey = ""
	dryRun                   bool
	depositDataPath          = ""
	nodeEndpoint             = ""
	nodeUseIPC               = false
	depositContractAddress   = ""
)

func init() {
	var help bool
	flag.StringVar(&depositAddressPrivateKey, "private-key", "", "Deposit address private key")
	flag.BoolVar(&dryRun, "dry-run", false, "Test before deposit processing")
	flag.BoolVar(&help, "help", false, "Show help")
	flag.StringVar(&depositDataPath, "deposit-data-path", "./validator_keys/deposit_data-*.json", "Path to deposit data files")
	flag.StringVar(&nodeEndpoint, "node-endpoint", "localhost:8545", "Node Endpoint")
	flag.BoolVar(&nodeUseIPC, "node-use-socket", false, "Use Node IPC socket")
	flag.StringVar(&depositContractAddress, "deposit-contract-address", "0x00000000219ab540356cbb839cbe05303d7705fa", "Deposit contract address")
	flag.Parse()
	if help {
		flag.Usage()
		os.Exit(0)
	}
}
