package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/CryptexWebDev/Deposit-Send/clients/microrest"
	"github.com/CryptexWebDev/Deposit-Send/tools"
	"github.com/gosuri/uitable"
	"net/url"
	"os"
	"path"
	"strings"
)

const (
	VERSION                     = "0.0.1"
	DEFAULT_VALIDATORS_KEYS_DIR = "validator_keys"
	DEFAULT_REST_API_URL        = "http://localhost:3500"
	VALIDATOR_STATUS_REQUEST    = "/eth/v1/beacon/states/finalized/validators"
)

var (
	validatorKeysDir = DEFAULT_VALIDATORS_KEYS_DIR
	restApiUrl       = DEFAULT_REST_API_URL
)

func main() {
	fmt.Println("Validator balances view app", VERSION)
	depositDataFiles, err := tools.FilesSearchByMask(validatorKeysDir, "deposit_data-*.json")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if len(depositDataFiles) == 0 {
		fmt.Println("No deposit data files found")
		return
	}
	fmt.Println("Deposit data files:", depositDataFiles)
	rest := microrest.NewRestClient(restApiUrl)
	depositData, err := readDepositData(validatorKeysDir, depositDataFiles)
	table := uitable.New()
	table.Separator = " | "
	table.AddRow("Index", "Pubkey", "Status", "Withdrawal", "Balance")
	for i, row := range depositData {
		validatorInfo := &ValidatorResponse{}
		req, _ := url.JoinPath(VALIDATOR_STATUS_REQUEST, "0x"+row.Pubkey)
		err := rest.Get(req, validatorInfo)
		if err != nil {
			fmt.Println("Can not get validator info:", err)
			return
		}
		var status, balance, withdrawalCredentials string
		if validatorInfo.ErrorCode == 404 {
			status = "not shown in blockchain"
			balance = "n/a"
			withdrawalCredentials = "n/a"
		} else {
			status = validatorInfo.Data.Status
			balance = validatorInfo.Data.Balance
			if validatorInfo.Data.Validator.WithdrawalCredentials != "" {
				if strings.HasPrefix(validatorInfo.Data.Validator.WithdrawalCredentials, "0x01") {
					withdrawalCredentials = "set"
				} else {
					withdrawalCredentials = "not set"
				}
			}

		}
		table.AddRow(i, row.Pubkey, status, withdrawalCredentials, balance)
	}
	fmt.Println(table)
}

type depositRow struct {
	Pubkey                string `json:"pubkey"`
	WithdrawalCredentials string `json:"withdrawal_credentials"`
}

type ValidatorInfo struct {
	ErrorCode int    `json:"code"`
	Index     string `json:"index"`
	Balance   string `json:"balance"`
	Status    string `json:"status"`
	Validator struct {
		Pubkey                     string `json:"pubkey"`
		WithdrawalCredentials      string `json:"withdrawal_credentials"`
		EffectiveBalance           string `json:"effective_balance"`
		Slashed                    bool   `json:"slashed"`
		ActivationEligibilityEpoch string `json:"activation_eligibility_epoch"`
		ActivationEpoch            string `json:"activation_epoch"`
		ExitEpoch                  string `json:"exit_epoch"`
		WithdrawableEpoch          string `json:"withdrawable_epoch"`
	} `json:"validator"`
}

type ValidatorResponse struct {
	ErrorCode           int  `json:"code"`
	ExecutionOptimistic bool `json:"execution_optimistic"`
	Finalized           bool `json:"finalized"`
	Data                struct {
		Index     string `json:"index"`
		Balance   string `json:"balance"`
		Status    string `json:"status"`
		Validator struct {
			Pubkey                     string `json:"pubkey"`
			WithdrawalCredentials      string `json:"withdrawal_credentials"`
			EffectiveBalance           string `json:"effective_balance"`
			Slashed                    bool   `json:"slashed"`
			ActivationEligibilityEpoch string `json:"activation_eligibility_epoch"`
			ActivationEpoch            string `json:"activation_epoch"`
			ExitEpoch                  string `json:"exit_epoch"`
			WithdrawableEpoch          string `json:"withdrawable_epoch"`
		} `json:"validator"`
	} `json:"data"`
}

func init() {
	var help bool
	flag.StringVar(&validatorKeysDir, "keys-dir", DEFAULT_VALIDATORS_KEYS_DIR, "Directory with validators keys")
	flag.StringVar(&restApiUrl, "rest-api", DEFAULT_REST_API_URL, "REST API URL")
	flag.BoolVar(&help, "help", false, "Print this help message.")
	flag.Parse()
	if help {
		flag.Usage()
		os.Exit(0)
	}
}

func readDepositData(dir string, files []string) (depositData []*depositRow, err error) {
	for _, file := range files {
		dataBin, err := os.ReadFile(path.Join(dir, file))
		if err != nil {
			return nil, err
		}
		var depositDataParsed []*depositRow
		err = json.Unmarshal(dataBin, &depositDataParsed)
		if err != nil {
			return nil, err
		}
		depositData = append(depositData, depositDataParsed...)
	}
	return depositData, nil
}
