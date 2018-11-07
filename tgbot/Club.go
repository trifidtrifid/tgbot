package tgbot

import (
	"github.com/mrd0ll4r/tbotapi"
	"log"
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

func (club *Club) GetUserByChat(chat tbotapi.Chat) *User {

	var user *User
	if chat.Username != nil {
		user = club.GetUser(*chat.Username)
	} else {
		user = club.GetUser(string(chat.ID))
	}
	if user.Recipient == nil {
		user.Recipient = new(tbotapi.Recipient)
		*user.Recipient = tbotapi.NewRecipientFromChat(chat)
	}
	return user
}

func (club *Club) GetUser(userName string) *User {
	club.usersMtx.Lock()
	defer club.usersMtx.Unlock()

	user, ok := club.Users[userName]
	if !ok {
		user = new(User)
		user.Close = make(chan interface{})
		user.Msgs = make(chan UserMessage)
		user.Club = club
		club.Users[userName] = user
		go user.Run()
	}

	return user
}

func (club *Club) NotifyEveryone(text string) {
	club.usersMtx.Lock()
	defer club.usersMtx.Unlock()

	for _, user := range club.Users {
		if user.Recipient != nil {
			club.IoService.sendText(*user.Recipient, text)
		} else {
			log.Printf("Empti recipient. Cannot broadcast text (%s)", text)
		}
	}
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