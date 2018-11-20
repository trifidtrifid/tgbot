package tgbot

import (
	"fmt"
	"github.com/trifidtrifid/tbotapi"
	"strconv"
	"sync"
	"time"
)

type Club struct {
	Users map[string]*User
	usersMtx sync.Mutex

	IoService IoService
	commonFund int
	fundMtx sync.Mutex

	credit int
	creditMtx sync.Mutex
	cfg *Config
}

func (club *Club) RunApDistribution() {

	go func (club *Club) {
		timer2 := time.NewTicker(1 * time.Second)
		for now := range timer2.C {
			if (now.Unix()) % 1800 == 0 {
				club.DistrubAp()
			}
		}
	}(club)
}



func (club *Club) DistrubAp() {
	club.usersMtx.Lock()
	defer club.usersMtx.Unlock()

	for _, user := range club.Users {
		user.DistrubAp()
	}
	club.SaveClub()
}

func (club *Club) SaveClub() {
	var v []UserCfg
	club.cfg.File.Users = v

	for _, user := range club.Users {
		var userCfg UserCfg
		userCfg.Salary = user.Salary
		userCfg.InCredit = user.InCredit
		userCfg.HoldAmount = user.HoldAmount
		userCfg.AP = user.AP
		userCfg.CreditLimit = user.CreditLimit
		userCfg.Chat = user.Chat
		club.cfg.File.Users = append(club.cfg.File.Users, userCfg)
	}

	club.cfg.Save()
}

func CreateClub(cfg *Config) *Club {
	club := new(Club)
	club.cfg = cfg
	club.Users = make(map[string]*User)

	for _, userCfg := range cfg.File.Users {
		user := club.GetUser(userCfg.Chat)
		user.Salary = userCfg.Salary

		user.InCredit = userCfg.InCredit
		club.CreditAdd(user.InCredit)

		user.HoldAmount = userCfg.HoldAmount
		club.FundAdd(user.HoldAmount)

		user.AP = userCfg.AP
		user.Chat = userCfg.Chat
		user.CreditLimit = userCfg.CreditLimit

	}

	return club
}

func (club *Club) GetUser(chat tbotapi.Chat) *User {
	club.usersMtx.Lock()
	defer club.usersMtx.Unlock()

	user, ok := club.Users[strconv.Itoa(chat.ID)]
	if !ok {
		user = new(User)
		user.Close = make(chan interface{})
		user.Msgs = make(chan UserMessage)
		user.Club = club
		user.Chat = chat
		club.Users[strconv.Itoa(chat.ID)] = user
		go user.Run()
	}

	return user
}

func (club *Club) FindUser(userId string) *User {
	club.usersMtx.Lock()
	defer club.usersMtx.Unlock()

	user, ok := club.Users[userId]
	if !ok {
		return nil
	}

	return user
}

func (club *Club) NotifyEveryone(text string, chat *tbotapi.Chat) {
	club.usersMtx.Lock()
	defer club.usersMtx.Unlock()

	for _, user := range club.Users {
		if (chat == nil) || (chat.ID != user.Chat.ID) {
			club.IoService.sendText(tbotapi.NewRecipientFromChat(user.Chat), text)
		}
	}
}

func (club *Club) ClubUsersInfo() string {
	club.usersMtx.Lock()
	defer club.usersMtx.Unlock()
	var retStr string
	for _, user := range club.Users {
		retStr += fmt.Sprintf("%s\n", user.getClubInfo())
	}
	return retStr
}

func (club *Club) FundAdd(i int) {
	club.fundMtx.Lock()
	club.commonFund += i
	club.fundMtx.Unlock()
}

func (club *Club) FundRemove(i int) {
	club.fundMtx.Lock()
	club.commonFund -= i
	club.fundMtx.Unlock()
}

func (club *Club) GetFund() int {
	club.fundMtx.Lock()
	fund := club.commonFund
	club.fundMtx.Unlock()
	return fund
}

func (club *Club) CreditAdd(i int) {
	club.fundMtx.Lock()
	club.credit += i
	club.fundMtx.Unlock()
}

func (club *Club) CreditRemove(i int) {
	club.fundMtx.Lock()
	club.credit -= i
	club.fundMtx.Unlock()
}

func (club *Club) GetCredit() int {
	club.fundMtx.Lock()
	credit := club.credit
	club.fundMtx.Unlock()
	return credit
}