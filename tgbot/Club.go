package tgbot

import "sync"

type Club struct {
	Users map[string]*User
	IoService IoService
	commonFund int
	fundMtx sync.Mutex
}


func CreateClub() *Club {
	club := new(Club)
	club.Users = make(map[string]*User)
	return club
}

func (club *Club) GetUser(userName string) *User {
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