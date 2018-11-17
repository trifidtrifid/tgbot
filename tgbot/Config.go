package tgbot

import (
	"encoding/json"
	"fmt"
	"github.com/trifidtrifid/tbotapi"
	"io/ioutil"
	"os"
)

type Config struct {
	Path string
	File ConfigFile

}

type ConfigFile struct {
	Users []UserCfg
}

type UserCfg struct {
	Salary int
	HoldAmount int
	CreditLimit int
	AP float64
	InCredit int
	Chat tbotapi.Chat
}

func (cfg *Config) Load() bool {
	jsonFile, err := os.Open(cfg.Path)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Successfully Opened %s\n", cfg.Path)
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)


	json.Unmarshal([]byte(byteValue), &cfg.File)

	return true
}


func (cfg *Config) Save() bool {
	jsonFile, err := os.Open(cfg.Path)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Successfully Opened %s\n", cfg.Path)
	defer jsonFile.Close()


	byteValue, _ := json.Marshal(&cfg.File)

	ioutil.WriteFile(cfg.Path, byteValue, os.ModePerm)

	return true


}
