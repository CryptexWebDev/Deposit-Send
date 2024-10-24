package main

import (
	"fmt"
	"github.com/CryptexWebDev/Deposit-Send/abi"
	"github.com/CryptexWebDev/Deposit-Send/abi/abicoder"
	"github.com/CryptexWebDev/Deposit-Send/clients/drlclient"
	"github.com/CryptexWebDev/Deposit-Send/common/hexnum"
	"github.com/CryptexWebDev/Deposit-Send/storage"
	"github.com/CryptexWebDev/Deposit-Send/tools/log"
	"math/big"
	"os"
)

var (
	globalConfigPath = "config.json"
	config           = &Config{
		storage: _configDefaultStorage(),
	}
	client *drlclient.Client

	depositData []*DepositData
)

func main() {
	fmt.Println("Deposit contract CLI")
	configStorage, err := storage.NewBinFileStorage("Config", "", "", globalConfigPath)
	if err != nil {
		log.Error("Can not get init config storage:", err)
		os.Exit(-1)
	}
	config.storage = configStorage
	err = config.Load()
	if err != nil {
		log.Error("Can not load config:", err)
		os.Exit(-1)
	}
	log.SetLevel(5)
	log.Info("Node connection settings:")
	if !config.NodeUseIPC {
		log.Info("- Node Connection : http-rpc")
		log.Info("- Node Url        :", config.NodeUrl)
		log.Info("- Node Port       :", config.NodePort)
	} else {
		log.Info("- Node Connection: ipc socket")
		log.Info("- Node ipc Socket Path:", config.NodeIPCSocket)
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
	if config.NodeUseIPC {
		clientOptions = append(clientOptions, drlclient.WithIPCClient(config.NodeIPCSocket))
	} else {
		clientOptions = append(clientOptions,
			drlclient.WithRpcClient(
				config.NodeUrl,
				config.NodePort,
				config.NodeUseSSl,
				config.AdditionalHeaders,
			),
		)
	}
	client = drlclient.NewClient(clientOptions...)
	err = client.Init()
	if err != nil {
		log.Error("Can not init chain client:", err)
		os.Exit(-1)
	}
	log.Info("Blockchain Info:")
	log.Info("- Chain Name:", client.GetChainName())
	log.Info("- Chain ID:", client.GetChainId())
	callData, err := abiManager.CallByMethod(config.DepositContractAddress, "get_deposit_count")
	if err != nil {
		log.Error("Can not prepare deposit count call data:", err)
		os.Exit(-1)
	}
	_, err = client.Call(config.DepositContractAddress, callData)
	if err != nil {
		log.Error("Can not call deposit count:", err)
		os.Exit(-1)
	}

	_searchFile(config.DepositDataPath)

	_, err = abiManager.CallByMethod(config.DepositContractAddress, "get_deposit_root")
	if err != nil {
		log.Error("Can not prepare deposit count call data:", err)
		os.Exit(-1)
	}
	_, err = client.Call(config.DepositContractAddress, callData)
	if err != nil {
		log.Error("Can not call deposit count:", err)
		os.Exit(-1)
	}
	log.Info("Preload deposit data")
	err = preloadDepositData()
	if err != nil {
		log.Error("Can not preload deposit data:", err)
		os.Exit(-1)
	}
	log.Info("Deposit data loaded")
	for _, dd := range depositData {
		log.Info("- Pubkey:", dd.Pubkey)
	}
	if depositAddressPrivateKey == "" {
		os.Exit(0)
	}
	log.Info("Prepare deposit")
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
	log.Info("Deposit address:", address)
	currentBalance, err := client.GetBalance(address)
	if err != nil {
		log.Error("Can not get balance for address:", err)
		os.Exit(-1)
	}
	log.Info("Current balance:", currentBalance)
	log.Debug("Prepare deposit...")
	log.Debug("Prepare deposit call data...")
	log.Debug("Prepare first validator deposit...")
	gasPrice, err := client.GasPrice()
	log.Info("Expected Gas price:", gasPrice)
	nonce, err := client.PendingNonceAt(address)
	log.Info("Expected Nonce:", nonce)
	validatorDepositData := depositData[0]
	log.Warning("Validator pub key:", validatorDepositData.Pubkey)
	paramPubKey, err := hexnum.ParseHexBytes(validatorDepositData.Pubkey)
	if err != nil {
		log.Error("Can not parse validator pubkey:", err)
		os.Exit(-1)
	}
	paramWithdrawalCredentials, err := hexnum.ParseHexBytes(validatorDepositData.WithdrawalCredentials)
	if err != nil {
		log.Error("Can not parse withdrawal credentials:", err)
		os.Exit(-1)
	}
	paramSignature, err := hexnum.ParseHexBytes(validatorDepositData.Signature)
	if err != nil {
		log.Error("Can not parse signature:", err)
		os.Exit(-1)
	}
	paramDataRoot, err := hexnum.ParseHexBytes(validatorDepositData.DepositDataRoot)
	if err != nil {
		log.Error("Can not parse data root:", err)
		os.Exit(-1)
	}
	log.Info("Deposit Data:")
	log.Info("- Pubkey:", hexnum.BytesToHex(paramPubKey), ",", len(paramPubKey))
	log.Info("- Withdrawal Credentials:", hexnum.BytesToHex(paramWithdrawalCredentials), ",", len(paramWithdrawalCredentials))
	log.Info("- Signature:", hexnum.BytesToHex(paramSignature), ",", len(paramSignature))
	log.Info("- Data Root:", hexnum.BytesToHex(paramDataRoot), ",", len(paramDataRoot))
	callDataBytes, err := abicoder.EncodeWithSignature("deposit(bytes,bytes,bytes,bytes32)", paramPubKey, paramWithdrawalCredentials, paramSignature, paramDataRoot)
	//callData, err = abiManager.CallByMethod(config.DepositContractAddress, "deposit",
	//	paramPubKey,
	//	paramWithdrawalCredentials,
	//	paramSignature,
	//	paramDataRoot,
	//)
	if err != nil {
		log.Error("Can not prepare deposit call data:", err)
		os.Exit(-1)
	}
	callData = hexnum.BytesToHex(callDataBytes)
	var depositAmount int64 = validatorDepositData.Amount
	amountBig := big.NewInt(depositAmount).Mul(big.NewInt(depositAmount), big.NewInt(1000000000))
	log.Info("Deposit amount:", drlclient.WeiToEtherString(amountBig))
	log.Debug("Deposit call data:", callData)
	client.SetDebug(true)
	gas, err := client.GetEstimatedGas(address, config.DepositContractAddress, callData, amountBig)
	if err != nil {
		log.Error("Can not get estimated gas for deposit:", err)
		os.Exit(-1)
	} else {
		log.Info("Estimated gas:", gas)
	}
	if dryRun {
		os.Exit(0)
	}
	//privateKey, from, to, data string, amount *big.Int)
	txId, err := client.SendTransactionByPrivateKey(depositAddressPrivateKey, address, config.DepositContractAddress, callData, amountBig)
	if err != nil {
		log.Error("Can not send deposit transaction:", err)
		os.Exit(-1)
	}
	log.Info("Deposit transaction sent:", txId)
}
