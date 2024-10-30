package main

import (
	"encoding/json"
	"os"
)

type DepositData struct {
	Pubkey                string `json:"pubkey"`
	WithdrawalCredentials string `json:"withdrawal_credentials"`
	Amount                int64  `json:"amount"`
	Signature             string `json:"signature"`
	DepositMessageRoot    string `json:"deposit_message_root"`
	DepositDataRoot       string `json:"deposit_data_root"`
	ForkVersion           string `json:"fork_version"`
	NetworkName           string `json:"network_name"`
	DepositCliVersion     string `json:"deposit_cli_version"`
}

func preloadDepositData(depositDataPath string) error {
	dataBytes, err := os.ReadFile(depositDataPath)
	if err != nil {
		return err
	}
	return json.Unmarshal(dataBytes, &depositData)
}
