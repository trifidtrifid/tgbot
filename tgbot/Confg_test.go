package tgbot

import "testing"

func TestSaveLoad(t *testing.T) {

	var cfg Config
	cfg.Path = "cfg.json"

	var userCfg UserCfg
	userCfg.InCredit = 1000
	userCfg.Salary = 2000
	userCfg.Chat.ID = 123
	userCfg.Chat.FirstName = new(string)
	*userCfg.Chat.FirstName = "test"


	cfg.File.Users = append(cfg.File.Users, userCfg)

	cfg.Save()
}