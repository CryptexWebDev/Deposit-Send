package abi

func (m *SmartContractsManager) CallByMethod(contractAddress string, method string, params ...interface{}) (callDataHex string, err error) {
	contract, err := m.GetSmartContractByAddress(contractAddress)
	if err != nil {
		return "", err
	}
	methodAbi, err := contract.Abi.GetMethodByName(method)
	if err != nil {
		return "", err
	}
	callDataHex, err = methodAbi.encodeInputs(params...)
	if err != nil {
		return "", err
	}
	return callDataHex, nil

}
