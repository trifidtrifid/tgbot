package tgbot

import (
	"fmt"
	"github.com/mrd0ll4r/tbotapi"
	"strconv"
	"sync"
)

type Club struct {
	Users map[string]*User
	usersMtx sync.Mutex

	IoService IoService
	commonFund int
	fundMtx sync.Mutex

	credit int
	creditMtx sync.Mutex
}


func CreateClub() *Club {
	club := new(Club)
	club.Users = make(map[string]*User)
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

func (club *Club) NotifyEveryone(text string) {
	club.usersMtx.Lock()
	defer club.usersMtx.Unlock()

	for _, user := range club.Users {
		club.IoService.sendText(tbotapi.NewRecipientFromChat(user.Chat), text)
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