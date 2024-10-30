package main

import (
	"bufio"
	"fmt"
	"github.com/CryptexWebDev/Deposit-Send/abi"
	"github.com/CryptexWebDev/Deposit-Send/abi/abicoder"
	"github.com/CryptexWebDev/Deposit-Send/clients/drlclient"
	"github.com/CryptexWebDev/Deposit-Send/common/hexnum"
	"github.com/CryptexWebDev/Deposit-Send/storage"
	"github.com/CryptexWebDev/Deposit-Send/tools/log"
	"math/big"
	"os"
	"regexp"
	"strings"
)

var (
	client      *drlclient.Client
	depositData []*DepositData

	isValidPrivateKeyRegexp = regexp.MustCompile("^[A-Fa-f0-9]{64}$")
)

func main() {
	fmt.Println("Deposit contract CLI")
	log.SetLevel(5)
	log.Info("Node connection settings:")
	if !nodeUseIPC {
		log.Info("- Node Endpoint Connection : http-rpc")
		log.Info("- Node Url        :", nodeEndpoint)
	} else {
		log.Info("- Node Connection: ipc socket")
		log.Info("- Node ipc Socket Path:", nodeEndpoint)
	}
	storageManager, err := storage.NewStorageManager("data")
	if err != nil {
		log.Error("Can not init storage manager:", err)
		os.Exit(-1)
	}
	addressCodec := drlclient.GetAddressCodec()
	abiStorage := storageManager.GetModuleStorage("ABI", "abi")
	abiManager := abi.NewManager(
		abi.WithStorage(abiStorage.GetBinFileStorage("known_contracts.json")),
		abi.WithAddressCodec(addressCodec),
	)
	err = abiManager.Init()
	if err != nil {
		log.Error("Can not load abi manager:", err)
		os.Exit(-1)
	}
	clientStorage := storageManager.GetModuleStorage("Client", "client")
	var clientOptions = []drlclient.Option{
		drlclient.WithConfigStorage(clientStorage.GetBinFileStorage("config.json")),
		drlclient.WithAbiManager(abiManager),
	}
	if nodeUseIPC {
		clientOptions = append(clientOptions, drlclient.WithIPCClient(nodeEndpoint))
	} else {
		nodeEndpointParts := strings.Split(nodeEndpoint, ":")
		if len(nodeEndpointParts) != 2 {
			log.Error("Invalid node endpoint URL:", nodeEndpoint)
			os.Exit(-1)
		}
		nodeUrl := nodeEndpointParts[0]
		nodePort := nodeEndpointParts[1]
		clientOptions = append(clientOptions,
			drlclient.WithRpcClient(
				nodeUrl,
				nodePort,
				false,
				nil,
			),
		)
	}
	client = drlclient.NewClient(clientOptions...)
	err = client.Init()
	if err != nil {
		log.Error("Can not init chain client:", err)
		os.Exit(-1)
	}
	client.SetDebug(false)
	log.Info("Blockchain Info:")
	log.Info("- Chain Name:", client.GetChainName())
	log.Info("- Chain ID:", client.GetChainId())
	callData, err := abiManager.CallByMethod(depositContractAddress, "get_deposit_count")
	if err != nil {
		log.Error("Can not prepare deposit count call data:", err)
		os.Exit(-1)
	}
	_, err = client.Call(depositContractAddress, callData)
	if err != nil {
		log.Error("Can not call deposit count:", err)
		os.Exit(-1)
	}

	depositDataList, err := _searchFile(depositDataPath)
	if err != nil {
		log.Error("Can not load deposit data:", err)
		os.Exit(-1)
	}
	if len(depositDataList) == 0 {
		log.Warning("No deposit data found")
		os.Exit(0)
	}
	var depositDataFile string
	if len(depositDataList) == 1 {
		depositDataFile = depositDataList[0]
	} else {
		log.Info("TODO Select deposit data")
		os.Exit(0)
	}
	_, err = abiManager.CallByMethod(depositContractAddress, "get_deposit_root")
	if err != nil {
		log.Error("Can not prepare deposit count call data:", err)
		os.Exit(-1)
	}
	_, err = client.Call(depositContractAddress, callData)
	if err != nil {
		log.Error("Can not call deposit count:", err)
		os.Exit(-1)
	}
	log.Info("Preload deposit data")
	err = preloadDepositData(depositDataFile)
	if err != nil {
		log.Error("Can not preload deposit data:", err)
		os.Exit(-1)
	}
	log.Info("Deposit data loaded")
	for _, dd := range depositData {
		log.Info("- Pubkey:", dd.Pubkey)
	}
	var depositDataForSend *DepositData
	if len(depositData) > 1 {
		log.Warning("Multi deposit data found, todo: select deposit data")
		os.Exit(0)
	} else {
		depositDataForSend = depositData[0]
	}
	if depositAddressPrivateKey == "" {
		log.Warning("Please set deposit address private key:")
		depositAddressPrivateKey = askPrivateKey()
	}
	log.Info("Prepare deposit...")
	pkBytes, err := hexnum.ParseHexBytes(depositAddressPrivateKey)
	if err != nil {
		log.Error("Can not parse private key:", err)
		os.Exit(-1)
	}
	address, _, err := addressCodec.PrivateKeyToAddress(pkBytes)
	if err != nil {
		log.Error("Can not get address from private key:", err)
		os.Exit(-1)
	}
	log.Warning("Deposit address:", address)
	confirmed := confirm("Please check and confirm")
	if !confirmed {
		log.Warning("Deposit address is not confirmed, cancel deposit send procedure")
		os.Exit(0)
	}
	currentBalance, err := client.GetBalance(address)
	if err != nil {
		log.Error("Can not get balance for address:", err)
		os.Exit(-1)
	}
	log.Info("Current balance:", currentBalance)
	log.Debug("- Prepare deposit...")
	log.Debug("- Prepare deposit call data...")
	log.Debug("- Prepare first validator deposit...")
	gasPrice, err := client.GasPrice()
	log.Info("Expected Gas price:", gasPrice)
	nonce, err := client.PendingNonceAt(address)
	log.Info("Expected Nonce:", nonce)
	log.Warning("Validator pub key:", depositDataForSend.Pubkey)
	paramPubKey, err := hexnum.ParseHexBytes(depositDataForSend.Pubkey)
	if err != nil {
		log.Error("Can not parse validator pubkey:", err)
		os.Exit(-1)
	}
	paramWithdrawalCredentials, err := hexnum.ParseHexBytes(depositDataForSend.WithdrawalCredentials)
	if err != nil {
		log.Error("Can not parse withdrawal credentials:", err)
		os.Exit(-1)
	}
	paramSignature, err := hexnum.ParseHexBytes(depositDataForSend.Signature)
	if err != nil {
		log.Error("Can not parse signature:", err)
		os.Exit(-1)
	}
	paramDataRoot, err := hexnum.ParseHexBytes(depositDataForSend.DepositDataRoot)
	if err != nil {
		log.Error("Can not parse data root:", err)
		os.Exit(-1)
	}
	log.Info("Deposit Data:")
	log.Info("- Pubkey:", hexnum.BytesToHex(paramPubKey), ",", len(paramPubKey))
	log.Info("- Withdrawal Credentials:", hexnum.BytesToHex(paramWithdrawalCredentials), ",", len(paramWithdrawalCredentials))
	log.Info("- Data Root:", hexnum.BytesToHex(paramDataRoot), ",", len(paramDataRoot))
	callDataBytes, err := abicoder.EncodeWithSignature("deposit(bytes,bytes,bytes,bytes32)", paramPubKey, paramWithdrawalCredentials, paramSignature, paramDataRoot)
	if err != nil {
		log.Error("Can not prepare deposit call data:", err)
		os.Exit(-1)
	}
	callData = hexnum.BytesToHex(callDataBytes)
	var depositAmount int64 = depositDataForSend.Amount
	amountBig := big.NewInt(depositAmount).Mul(big.NewInt(depositAmount), big.NewInt(1000000000))
	log.Info("Deposit amount:", drlclient.WeiToEtherString(amountBig))
	//log.Debug("Deposit call data:", callData)
	client.SetDebug(false)
	gas, err := client.GetEstimatedGas(address, depositContractAddress, callData, amountBig)
	if err != nil {
		log.Error("Can not get estimated gas for deposit:", err)
		os.Exit(-1)
	} else {
		log.Info("Estimated gas:", gas)
	}
	if dryRun {
		os.Exit(0)
	}

	txId, err := client.SendTransactionByPrivateKey(depositAddressPrivateKey, address, depositContractAddress, callData, amountBig)
	if err != nil {
		log.Error("Can not send deposit transaction:", err)
		os.Exit(-1)
	}
	log.Info("Deposit transaction sent:", txId)
}

func askPrivateKey() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("private key: ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSuffix(text, "\n")
	if text == "" {
		return askPrivateKey()
	}
	if !isValidPrivateKey(text) {
		log.Error("Invalid private key, please enter valid private key")
		return askPrivateKey()
	}
	return text
}

func isValidPrivateKey(hex string) bool {
	return isValidPrivateKeyRegexp.MatchString(strings.TrimPrefix(hex, "0x"))
}

func confirm(message string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(message, "(y/n): ")
	text, _ := reader.ReadString('\n')
	if text[0] != 'y' && text[0] != 'Y' && text[0] != 'n' && text[0] != 'N' {
		return confirm(message)
	}
	return text[0] == 'y' || text[0] == 'Y'
}
