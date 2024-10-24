package drlclient

import (
	"encoding/json"
	"github.com/CryptexWebDev/Deposit-Send/storage"
	"github.com/CryptexWebDev/Deposit-Send/tools/log"
	"github.com/CryptexWebDev/Deposit-Send/types"
)

type Config struct {
	storage       storage.BinStorage
	ChainName     string
	ChainId       string
	ChainSymbol   string
	Decimals      int
	Confirmations int  `json:"confirmations"`
	Debug         bool `json:"debug"`
	Tokens        []*types.TokenInfo
}

func _configDefaultStorage() storage.BinStorage {
	configStore, err := storage.NewBinFileStorage("Config", "data", "client", "config.json")
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
	c.ChainName = "Dorol"
	c.ChainId = "dorol"
	c.ChainSymbol = "DRL"
	c.Decimals = 18
	c.Confirmations = 10
	c.Tokens = []*types.TokenInfo{}
	return c.Save()
}
