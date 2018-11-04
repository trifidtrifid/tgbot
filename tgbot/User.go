package tgbot

import (
	"fmt"
	"github.com/mrd0ll4r/tbotapi"
)

const (
	HowMuchSal = "How much salary?"
	SalAnswer = "According your salary %d, your available credit limit is %d"

	HowMuchHold = "How much you want to hold?"
	HoldDone = "Your $%d held. You are a good and kind person!"
	ErrGotoMain = "Error. Try /start again"
	InfoMsg = "Club Info" +
		"Common Club Fund for zero percent credits: %d" +
		"Your credit limit %d" +
		"APs: %d" +
		"Your holded funds: %d"
)

type User struct {
	Msgs chan UserMessage
	Close chan interface{}
	Salary int
	HoldAmount int
	Club *Club
	CreditLimit int
	AP int
}

type UserMessage struct {
	Message tbotapi.Message
}

func (user *User) Run() {
	for {
		select {
		case <-user.Close:
			return
		case userMsg := <-user.Msgs:
			user.processMsg(&userMsg)
		}
	}
}

func (user *User) processMsg(msg *UserMessage) {

	fmt.Printf("receive message %s\n", *msg.Message.Text)
	if *msg.Message.Text == "/start" {
		recipient := tbotapi.NewRecipientFromChat(msg.Message.Chat)
		user.Club.IoService.sendMainMenu(recipient)
		userMsg := <-user.Msgs
		switch *userMsg.Message.Text {
		case "Hold":
			user.Club.IoService.sendText(recipient, HowMuchHold)
			userMsg = <-user.Msgs
			n, _ := fmt.Sscanf(*userMsg.Message.Text, "%d", &user.HoldAmount)
			if n != 1 {
				user.Club.IoService.sendText(recipient, ErrGotoMain)
				return
			}
			user.Club.FundAdd(user.HoldAmount)
			user.Club.IoService.sendText(recipient, fmt.Sprintf(HoldDone, user.HoldAmount))

		case "Salary":
			user.Club.IoService.sendText(recipient, HowMuchSal)
			userMsg = <-user.Msgs
			n, _ := fmt.Sscanf(*userMsg.Message.Text, "%d", &user.Salary)
			if n != 1 {
				user.Club.IoService.sendText(recipient, ErrGotoMain)
				return
			}
			user.CreditLimit = int(float64(user.Salary) * 1.5)
			user.Club.IoService.sendText(recipient, fmt.Sprintf(SalAnswer, user.Salary, user.CreditLimit))
		case "Info":
			user.Club.IoService.sendText(recipient, fmt.Sprintf(InfoMsg, user.Club.GetFund(),
				user.CreditLimit, user.AP, user.HoldAmount))
		case "Get money":

		case "Return Money":

		default:
			fmt.Println("error. start from very begining")
		}
	}
}


