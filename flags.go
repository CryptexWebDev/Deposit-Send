package main

import (
	"flag"
	"os"
)

var (
	depositAddressPrivateKey = ""
	dryRun                   bool
)

func init() {
	var help bool
	flag.StringVar(&depositAddressPrivateKey, "private-key", "", "Deposit address private key")
	flag.BoolVar(&dryRun, "dry-run", false, "Test before deposit processing")
	flag.BoolVar(&help, "help", false, "Show help")
	flag.Parse()
	if help {
		flag.Usage()
		os.Exit(0)
	}
}
