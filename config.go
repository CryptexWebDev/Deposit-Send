package main

import (
	"encoding/json"
	"errors"
	"github.com/CryptexWebDev/Deposit-Send/storage"
	"github.com/CryptexWebDev/Deposit-Send/tools/log"
)

var (
	ErrConfigStorageEmpty = errors.New("config storage not set")
)

type Config struct {
	storage                storage.BinStorage
	NodeUrl                string            `json:"nodeUrl"`
	NodePort               string            `json:"nodePort"`
	NodeUseSSl             bool              `json:"nodeUseSSL"`
	NodeUseIPC             bool              `json:"nodeUseIPC"`
	NodeIPCSocket          string            `json:"nodeIPCSocket"`
	RpcAddress             string            `json:"rpcAddress"`
	RpcPort                string            `json:"rpcPort"`
	DebugMode              bool              `json:"debug_mode"`
	ParamsFlags            map[string]bool   `json:"flags,omitempty"`
	ParamsString           map[string]string `json:"paramsString,omitempty"`
	ParamsInt              map[string]int    `json:"paramsInt,omitempty"`
	AdditionalHeaders      map[string]string `json:"additionalHeaders,omitempty"`
	DepositContractAddress string            `json:"depositContractAddress"`
	DepositDataPath        string            `json:"depositDataPath"`
	DepositAmount          json.Number       `json:"depositAmount"`
}

func _configDefaultStorage() storage.BinStorage {
	configStore, err := storage.NewBinFileStorage("Config", ".", ".", globalConfigPath)
	if err != nil {
		log.Error("Can not get default config storage:", err)
	}
	return configStore
}

func (c *Config) Load() (err error) {
	if !c.storage.IsExists() {
		err = c.coldStart()
		if err != nil {
			return err
		}
	}
	jsonBytes, err := c.storage.Load()
	if err != nil {
		return
	}
	err = json.Unmarshal(jsonBytes, c)
	return
}

func (c *Config) Save() (err error) {
	data, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return
	}
	err = c.storage.Save(data)
	return
}

func (c *Config) coldStart() (err error) {
	if c.storage == nil {
		return ErrConfigStorageEmpty
	}
	c.NodeUrl = "localhost"
	c.NodePort = "8545"
	c.NodeUseIPC = false
	c.NodeIPCSocket = ""
	c.RpcAddress = "localhost"
	c.RpcPort = "21080"
	c.DepositDataPath = "./validator_keys/deposit_data-*.json"
	c.DepositContractAddress = "0x00000000219ab540356cbb839cbe05303d7705fa"
	c.DepositAmount = "2"
	return c.Save()
}

func (c *Config) Flag(name string) bool {
	var changed bool
	defer func() {
		if changed {
			err := c.Save()
			if err != nil {
				log.Error("Can not save config:", err)
			}
		}
	}()
	if c.ParamsFlags == nil {
		changed = true
		c.ParamsFlags = make(map[string]bool)
		c.ParamsFlags[name] = false
		return false
	}
	if value, ok := c.ParamsFlags[name]; !ok {
		c.ParamsFlags[name] = false
		changed = true
		return false
	} else {
		return value
	}
}

func (c *Config) String(flagName string, defaultValue string) string {
	var changed bool
	defer func() {
		if changed {
			err := c.Save()
			if err != nil {
				log.Error("Can not save config:", err)
			}
		}
	}()
	if c.ParamsString == nil {
		changed = true
		c.ParamsString = make(map[string]string)
		c.ParamsString[flagName] = defaultValue
		return defaultValue
	}
	if value, ok := c.ParamsString[flagName]; !ok {
		c.ParamsString[flagName] = defaultValue
		changed = true
		return defaultValue
	} else {
		return value
	}
}

func (c *Config) Int(flagName string, defaultValue int) int {
	var changed bool
	defer func() {
		if changed {
			err := c.Save()
			if err != nil {
				log.Error("Can not save config:", err)
			}
		}
	}()
	if c.ParamsInt == nil {
		changed = true
		c.ParamsInt = make(map[string]int)
		c.ParamsInt[flagName] = defaultValue
		return defaultValue
	}
	if value, ok := c.ParamsInt[flagName]; !ok {
		c.ParamsInt[flagName] = defaultValue
		changed = true
		return defaultValue
	} else {
		return value
	}
}
