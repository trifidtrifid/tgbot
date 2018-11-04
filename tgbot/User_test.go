package tgbot

import (
	"fmt"
	"github.com/mrd0ll4r/tbotapi"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestIoService struct {
	Out chan string
}

func CreateTestIoService() *TestIoService {
	var testIo *TestIoService
	testIo = new(TestIoService)
	testIo.Out = make(chan string)
	return testIo
}
func (bot *TestIoService) sendMainMenu(recipient tbotapi.Recipient) {
	fmt.Println("print main menu")
	bot.Out<-"main menu"
}
func (bot *TestIoService) sendText(recipient tbotapi.Recipient, text string) {
	fmt.Println("print text ", text)
	bot.Out<-text
}
func (bot *TestIoService) checkText(expect string) bool {
	rcvedText :=  <-bot.Out
	if expect != rcvedText {
		fmt.Printf("expected: %s but received %s", expect, rcvedText);
		return false
	}
	return true
}

func sendText(user *User, text string) {
	var msg UserMessage;
	msg.Message.Text = new(string)
	*msg.Message.Text = text
	user.Msgs<-msg
}

func TestUserHoldDeposit(t *testing.T) {

	club := CreateClub()
	testIo := CreateTestIoService()
	club.IoService = testIo

	user := club.GetUser("trifid")

	sendText(user, "/start")
	assert.True(t, testIo.checkText("main menu"))

	sendText(user, "Hold")
	assert.True(t, testIo.checkText(HowMuchHold))

	sendText(user, "1000")
	assert.True(t, testIo.checkText(fmt.Sprintf(HoldDone, 1000)))
	assert.Equal(t, 1000, user.HoldAmount)
	assert.Equal(t, 1000, club.GetFund())
}

func TestUserHoldDepositTwo(t *testing.T) {

	club := CreateClub()
	testIo := CreateTestIoService()
	club.IoService = testIo
	{
		user := club.GetUser("trifid")

		sendText(user, "/start")
		assert.True(t, testIo.checkText("main menu"))

		sendText(user, "Hold")
		assert.True(t, testIo.checkText(HowMuchHold))

		sendText(user, "1000")
		assert.True(t, testIo.checkText(fmt.Sprintf(HoldDone, 1000)))
		assert.Equal(t, 1000, user.HoldAmount)
		assert.Equal(t, 1000, club.GetFund())
	}
	{
		user := club.GetUser("usera")

		sendText(user, "/start")
		assert.True(t, testIo.checkText("main menu"))

		sendText(user, "Hold")
		assert.True(t, testIo.checkText(HowMuchHold))

		sendText(user, "1000")
		assert.True(t, testIo.checkText(fmt.Sprintf(HoldDone, 1000)))
		assert.Equal(t, 1000, user.HoldAmount)
		assert.Equal(t, 2000, club.GetFund())
	}
}



func TestUserSetSalary(t *testing.T) {

	club := CreateClub()
	testIo := CreateTestIoService()
	club.IoService = testIo

	user := club.GetUser("trifid")

	sendText(user, "/start")
	assert.True(t, testIo.checkText("main menu"))

	sendText(user, "Salary")
	assert.True(t, testIo.checkText(HowMuchSal))

	sendText(user, "2000")
	assert.True(t, testIo.checkText(fmt.Sprintf(SalAnswer, 2000, 3000)))
	assert.Equal(t, 2000, user.Salary)
	assert.Equal(t, 3000, user.CreditLimit)
}