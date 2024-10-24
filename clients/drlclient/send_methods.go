package drlclient

import (
	"github.com/CryptexWebDev/Deposit-Send/address"
	"github.com/CryptexWebDev/Deposit-Send/clients/urpc"
	"github.com/CryptexWebDev/Deposit-Send/common/hexnum"
	"github.com/CryptexWebDev/Deposit-Send/crypto"
	"github.com/CryptexWebDev/Deposit-Send/tools/log"
	"math/big"
)

func (c *Client) sendRawByPrivateKeyUnsafe(fromPrivateKey []byte, from, to string, amount, gasPrice *big.Int, gas int64) (txHash string, err error) {
	toBytes, err := c.addressCodec.DecodeAddressToBytes(to)
	if err != nil {
		return "", err
	}
	pk, _ := crypto.ECDSAKeysFromPrivateKeyBytes(fromPrivateKey)
	nonce, err := c.PendingNonceAt(from)
	if err != nil {
		return "", err
	}
	netId, err := c.GetNetId()
	if err != nil {
		return "", err
	}
	chainID := big.NewInt(netId)
	if c.debug {
		log.Warning("ChainID:", chainID)
	}
	txSigner := &crypto.EthTxSigner{
		Nonce:    uint64(nonce),
		GasPrice: gasPrice,
		Gas:      uint64(gas),
		To:       &toBytes,
		Value:    amount,
	}
	txSigner.SetChainId(chainID)
	sign := txSigner.Sign(pk)
	if len(sign) == 0 {
		log.Error("Can not sign transaction: sign is empty")
		return "", ErrTransactionSignError
	}
	txSignedBytes, err := txSigner.EncodeRPL()
	if err != nil {
		log.Error("Can not encode signed transaction:", err)
		return "", err

	}
	if c.debug {
		log.Warning("txSignedBytes", hexnum.BytesToHex(txSignedBytes))
	}
	txHash, err = c.SendRawTransaction(hexnum.BytesToHex(txSignedBytes))
	if err != nil {
		log.Error("Can not broadcast transaction:", err)
		return "", err
	}
	return txHash, nil
}

func (c *Client) sendRawByPrivateKeyWithDataUnsafe(fromPrivateKey []byte, from, to, data string, amount, gasPrice *big.Int, gas int64) (txHash string, err error) {
	dataBytes, err := hexnum.ParseHexBytes(data)
	if err != nil {
		return "", err
	}
	toBytes, err := c.addressCodec.DecodeAddressToBytes(to)
	if err != nil {
		return "", err
	}
	pk, _ := crypto.ECDSAKeysFromPrivateKeyBytes(fromPrivateKey)
	nonce, err := c.PendingNonceAt(from)
	if err != nil {
		return "", err
	}
	netId, err := c.GetNetId()
	if err != nil {
		return "", err
	}
	chainID := big.NewInt(netId)
	if c.debug {
		log.Warning("ChainID:", chainID)
	}
	txSigner := &crypto.EthTxSigner{
		Nonce:    uint64(nonce),
		GasPrice: gasPrice,
		Gas:      uint64(gas),
		To:       &toBytes,
		Data:     dataBytes,
		Value:    amount,
	}
	txSigner.SetChainId(chainID)
	sign := txSigner.Sign(pk)
	if len(sign) == 0 {
		log.Error("Can not sign transaction: sign is empty")
		return "", ErrTransactionSignError
	}
	txSignedBytes, err := txSigner.EncodeRPL()
	if err != nil {
		log.Error("Can not encode signed transaction:", err)
		return "", err

	}
	if c.debug {
		log.Warning("txSignedBytes", hexnum.BytesToHex(txSignedBytes))
	}
	txHash, err = c.SendRawTransaction(hexnum.BytesToHex(txSignedBytes))
	if err != nil {
		log.Error("Can not broadcast transaction:", err)
		return "", err
	}
	return txHash, nil
}

// SendRawTransaction sends the signed and RPL encoded transaction to the network.
// In fact, any transaction - transfer of funds, call of a smart contract function
// or deployment of a smart contract is carried out by calling this function
func (c *Client) SendRawTransaction(data string) (txHash string, err error) {
	req := urpc.NewRequest(ethSendRawTransaction)
	req.AddParams(data)
	result, err := c.rpcClient.Call(req)
	if err != nil {
		return "", err
	}
	err = result.ParseResult(&txHash)
	if err != nil {
		return "", err
	}
	return txHash, nil
}

// SendTransactionByPrivateKey sends the signed transaction to the network.
func (c *Client) SendTransactionByPrivateKey(privateKey, from, to, data string, amount *big.Int) (txId string, err error) {
	pkBytes, err := hexnum.ParseHexBytes(privateKey)
	if err != nil {
		return "", err
	}
	fromAddress, _, err := c.addressCodec.PrivateKeyToAddress(pkBytes)
	if err != nil {
		return "", err
	}
	if from != fromAddress {
		return "", address.ErrAddressPrivateKeyMismatch
	}
	gasPrice, err := c.GasPrice()
	gas, err := c.GetEstimatedGas(from, to, data, amount)
	if err != nil {
		return "", err
	}
	txId, err = c.sendRawByPrivateKeyWithDataUnsafe(pkBytes, from, to, data, amount, gasPrice, gas)
	if err != nil {
		return "", err
	}
	return txId, nil
}
