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
	InfoMsg = "*Club Info*\n" +
		"Common Club Fund for zero percent credits: %d\n" +
		"Your credit limit: %d\n" +
		"APs: %d\n" +
		"Your holded funds: %d\n" +
		"You borrow: %d\n"
	HowMuchTake = "How much you want to take from common fund? max available credit for you %d. Common fund is %d"
	TakenSucc = "You successfully borrow $%d"
	HowMuchReturn = "How much you want to return to common fund? your credit %d"
	ReturnSucc = "now you credit %d"

)

type User struct {
	Msgs chan UserMessage
	Close chan interface{}
	Salary int
	HoldAmount int
	Club *Club
	CreditLimit int
	AP int
	InCredit int
	 Recipient *tbotapi.Recipient
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

func (user *User) getClubInfo() string {
	return fmt.Sprintf(InfoMsg,
		user.Club.GetFund() - user.Club.GetCredit(),
		user.CreditLimit - user.InCredit,
		user.AP,
		user.HoldAmount,
		user.InCredit)
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
			var i int
			n, _ := fmt.Sscanf(*userMsg.Message.Text, "%d", &i)
			if n != 1 {
				user.Club.IoService.sendText(recipient, ErrGotoMain)
				return
			}
			user.HoldAmount += i
			user.Club.FundAdd(i)
			user.Club.IoService.sendText(recipient, fmt.Sprintf(HoldDone, user.HoldAmount))
			user.Club.NotifyEveryone(fmt.Sprintf("user %s hold %d", msg.Message.Chat, i))

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
			user.Club.IoService.sendText(recipient, fmt.Sprintf(InfoMsg, user.Club.GetFund() - user.Club.GetCredit(),
				user.CreditLimit - user.InCredit, user.AP, user.HoldAmount, user.InCredit))
		case "Get money":
			user.Club.IoService.sendText(recipient, fmt.Sprintf(HowMuchTake, user.CreditLimit - user.InCredit, user.Club.GetFund() - user.Club.GetCredit()))

			for {
				userMsg = <-user.Msgs
				var i int
				n, _ := fmt.Sscanf(*userMsg.Message.Text, "%d", &i)
				if n != 1 {
					user.Club.IoService.sendText(recipient,
						fmt.Sprintf("input just whole number lower then %d", user.CreditLimit - user.InCredit))
					continue
				}

				if i > (user.CreditLimit - user.InCredit) {
					user.Club.IoService.sendText(recipient,
						fmt.Sprintf("You cannot borrow more then %d", user.CreditLimit - user.InCredit))
					continue
				}

				if i > (user.CreditLimit - user.InCredit) {
					user.Club.IoService.sendText(recipient,
						fmt.Sprintf("You cannot borrow more then %d", user.CreditLimit - user.InCredit))
					continue
				}

				if i > (user.Club.GetFund() - user.Club.GetCredit()) {
					user.Club.IoService.sendText(recipient,
						fmt.Sprintf("You can borrow only %d", user.Club.GetFund() - user.Club.GetCredit()))
					continue
				}

				user.InCredit += i
				user.Club.CreditAdd(i)
				user.Club.IoService.sendText(recipient, fmt.Sprintf(TakenSucc, i))
				break
			}

		case "Return Money":
			user.Club.IoService.sendText(recipient, fmt.Sprintf(HowMuchReturn, user.InCredit))

			for {
				userMsg = <-user.Msgs
				var i int
				n, _ := fmt.Sscanf(*userMsg.Message.Text, "%d", &i)
				if n != 1 {
					user.Club.IoService.sendText(recipient,
						fmt.Sprintf("input just whole number lower then %d", user.InCredit))
					continue
				}

				if i > (user.InCredit) {
					user.Club.IoService.sendText(recipient,
						fmt.Sprintf("You shouldn't return more then %d", user.InCredit))
					continue
				}

				user.InCredit -= i
				user.Club.CreditRemove(i)
				user.Club.IoService.sendText(recipient, fmt.Sprintf(ReturnSucc, user.InCredit))
				break
			}

		default:
			fmt.Println("error. start from very begining")
		}

//		user.Club.IoService.sendMainMenu(recipient)
	}


}


