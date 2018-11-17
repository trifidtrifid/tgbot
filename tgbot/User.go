package tgbot

import (
	"fmt"
	"github.com/trifidtrifid/tbotapi"
	"math"
	"time"
)

const (
	HowMuchSal = "Specify your salary in RUB"
	SalAnswer = "According your salary %d RUB, your available credit limit is %d RUB"

	HowMuchHold = "How much you want to hold in RUB for your teammates?"
	HoldDone = "Your %d RUB held. You are a good and kind person!"
	ErrGotoMain = "Error. Try /start again"
	InfoMsg = "*Club Info for user %s*\n" +
		"Club Fund for zero percent credits: %d\n" +
		"Credit limit: %d RUB\n" +
		"APs: %d\n" +
		"Holded funds: %d RUB\n" +
		"Credit: %d RUB\n" +
		"Every 30min = 1 day"
	HowMuchTake = "How much you want to take from common fund? max available credit for you %d RUB. Common fund is %d RUB"
	TakenSucc = "You successfully borrow %d RUB"
	HowMuchReturn = "How much you want to return to common fund? your credit is %d RUB"
	ReturnSucc = "Now your credit is %d"

)

type User struct {
	Msgs chan UserMessage
	Close chan interface{}
	Salary int
	HoldAmount int
	Club *Club
	CreditLimit int
	AP float64
	InCredit int
	Chat tbotapi.Chat
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
			user.processMsg(userMsg)
		}
	}
}

func (user *User) sendMainMenu() {
	user.Club.IoService.sendMainMenu(tbotapi.NewRecipientFromChat(user.Chat))
}
func (user *User) sendText(text string) {
	user.Club.IoService.sendText(tbotapi.NewRecipientFromChat(user.Chat), text)
}

func (user *User) getClubInfo() string {
	return fmt.Sprintf(InfoMsg,
		user.Chat,
		user.Club.GetFund() - user.Club.GetCredit(),
		user.CreditLimit - user.InCredit,
		int(user.AP),
		user.HoldAmount,
		user.InCredit)
}

func (user *User) DistrubAp() {
	if user.InCredit != 0 {
		fmt.Printf("Active credit %d. No angel points today\n", user.InCredit)
		return
	}

	if user.HoldAmount == 0 {
		fmt.Printf("Hold amount %d. No angel points today\n", user.HoldAmount)
		return
	}

	//50% ap annual

	ap := float64(user.HoldAmount) * 0.5
	user.AP += ap / 365

	user.AP = math.Round(user.AP*100) / 100

	if user.HoldAmount == 0 {
		fmt.Printf("Hold amount %d. No angel points today\n", user.HoldAmount)
		return
	}

	fmt.Printf("%v Charge %f AP to %s\n", time.Now(), math.Round(ap/365) / 100, user.Chat)

}

func (user *User) processMsg(userMsg UserMessage) {

	fmt.Printf("receive message %s\n", *userMsg.Message.Text)
	switch *userMsg.Message.Text {
	case "/start":
		user.sendMainMenu()

	case "Hold":
		user.sendText(HowMuchHold)
		userMsg = <-user.Msgs
		var i int
		n, _ := fmt.Sscanf(*userMsg.Message.Text, "%d", &i)
		if n != 1 {
			user.sendText(ErrGotoMain)
			return
		}
		user.HoldAmount += i
		user.Club.FundAdd(i)
		user.sendText(fmt.Sprintf(HoldDone, user.HoldAmount))
		user.Club.NotifyEveryone(fmt.Sprintf("user %s hold %d RUB", userMsg.Message.Chat, i), &userMsg.Message.Chat)

	case "Salary":
		user.sendText( HowMuchSal)
		userMsg = <-user.Msgs
		n, _ := fmt.Sscanf(*userMsg.Message.Text, "%d", &user.Salary)
		if n != 1 {
			user.sendText( ErrGotoMain)
			return
		}
		user.CreditLimit = int(float64(user.Salary) * 1.5)
		user.sendText( fmt.Sprintf(SalAnswer, user.Salary, user.CreditLimit))
	case "Info":
		user.sendText( user.getClubInfo())
	case "Borrow":
		user.sendText( fmt.Sprintf(HowMuchTake, user.CreditLimit - user.InCredit, user.Club.GetFund() - user.Club.GetCredit()))

		for {
			userMsg = <-user.Msgs
			var i int
			n, _ := fmt.Sscanf(*userMsg.Message.Text, "%d", &i)
			if n != 1 {
				user.sendText(
					fmt.Sprintf("input just whole number lower then %d", user.CreditLimit - user.InCredit))
				continue
			}

			if i > (user.CreditLimit - user.InCredit) {
				user.sendText(
					fmt.Sprintf("You cannot borrow more then %d", user.CreditLimit - user.InCredit))
				continue
			}

			if i > (user.CreditLimit - user.InCredit) {
				user.sendText(
					fmt.Sprintf("You cannot borrow more then %d", user.CreditLimit - user.InCredit))
				continue
			}

			if i > (user.Club.GetFund() - user.Club.GetCredit()) {
				user.sendText(
					fmt.Sprintf("You can borrow only %d", user.Club.GetFund() - user.Club.GetCredit()))
				continue
			}

			if i > int(user.AP) {
				user.sendText(
					fmt.Sprintf("You have %f AP. You can borrow only %d RUB", user.AP, int(user.AP)))
				continue
			}

			user.AP -= float64(i)
			user.InCredit += i
			user.Club.CreditAdd(i)
			user.sendText(fmt.Sprintf(TakenSucc, i))

			user.Club.NotifyEveryone(fmt.Sprintf("user %s take from fund %d RUB", userMsg.Message.Chat, i), &userMsg.Message.Chat)
			break
		}

	case "Return Money":
		user.sendText( fmt.Sprintf(HowMuchReturn, user.InCredit))

		for {
			userMsg = <-user.Msgs
			var i int
			n, _ := fmt.Sscanf(*userMsg.Message.Text, "%d", &i)
			if n != 1 {
				user.sendText(
					fmt.Sprintf("input just whole number lower then %d", user.InCredit))
				continue
			}

			if i > (user.InCredit) {
				user.sendText(
					fmt.Sprintf("You shouldn't return more then %d", user.InCredit))
				continue
			}

			user.InCredit -= i
			user.Club.CreditRemove(i)
			user.sendText( fmt.Sprintf(ReturnSucc, user.InCredit))

			user.Club.NotifyEveryone(fmt.Sprintf("user %s return to fund %d RUB", userMsg.Message.Chat, i), &userMsg.Message.Chat)
			break
		}
	case "Users" :
		user.sendText(user.Club.ClubUsersInfo())
	default:
		user.sendText("try /start")
	}

	user.Club.SaveClub()
}


