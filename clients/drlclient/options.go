package drlclient

import (
	"fmt"
	"github.com/CryptexWebDev/Deposit-Send/abi"
	"github.com/CryptexWebDev/Deposit-Send/clients/urpc"
	"github.com/CryptexWebDev/Deposit-Send/storage"
)

func WithRpcClient(nodeAddress, nodePort string, useSSL bool, headers map[string]string) Option {
	urlMaks := "http://%s:%s"
	if useSSL {
		urlMaks = "https://%s:%s"
	}
	endpointUrl := fmt.Sprintf(urlMaks, nodeAddress, nodePort)
	return func(client *Client) {
		rpcClient := urpc.NewClient(urpc.WithHTTPRpc(endpointUrl, headers))
		client.rpcClient = rpcClient
	}
}

func WithIPCClient(ipcPath string) Option {
	return func(client *Client) {
		rpcClient := urpc.NewClient(urpc.WithRpcIPCSocket(ipcPath))
		client.rpcClient = rpcClient
	}
}

func WithAbiManager(abiManager *abi.SmartContractsManager) Option {
	return func(client *Client) {
		client.abi = abiManager
	}
}

func WithConfigStorage(storage storage.BinStorage) Option {
	return func(client *Client) {
		client.config.storage = storage
	}
}
