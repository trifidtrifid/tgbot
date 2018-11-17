package tgbot

import (
	"fmt"
	"github.com/mrd0ll4r/tbotapi"
	"github.com/stretchr/testify/assert"
	"hash/fnv"
	"strconv"
	"testing"
)


var testCfg Config

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
	fmt.Printf("-> to %s: %s\n", strconv.Itoa(*recipient.ChatID), text)
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

func createTestClub() (*Club, *TestIoService) {
	testCfg.Path = "cfg.json"
	club := CreateClub(&testCfg)
	testIo := CreateTestIoService()
	club.IoService = testIo
	return club, testIo
}

func (bot *TestIoService) skipText() {
	<-bot.Out
}

func makeChat(name string) tbotapi.Chat {
	var chat tbotapi.Chat
	chat.Username = new(string)
	*chat.Username = name
	
	h := fnv.New32a()
	h.Write([]byte(name))
	chat.ID = int(h.Sum32())
	
	return chat
}

func sendText(user *User, text string) {
	var msg UserMessage;
	msg.Message.Text = new(string)
	msg.Message.Chat = user.Chat
	*msg.Message.Text = text
	user.Msgs<-msg
}

func TestUserHoldDeposit(t *testing.T) {

	testCfg.Path = "cfg.json"
	club := CreateClub(&testCfg)

	testIo := CreateTestIoService()
	club.IoService = testIo

	user := club.GetUser(makeChat("trifid"))

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

	testCfg.Path = "cfg.json"
	club := CreateClub(&testCfg)

	testIo := CreateTestIoService()
	club.IoService = testIo
	{
		user := club.GetUser(makeChat("trifid"))

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
		user := club.GetUser(makeChat("usera"))

		sendText(user, "/start")
		assert.True(t, testIo.checkText("main menu"))

		sendText(user, "Hold")
		assert.True(t, testIo.checkText(HowMuchHold))

		sendText(user, "1000")
		assert.True(t, testIo.checkText(fmt.Sprintf(HoldDone, 1000)))
		assert.Equal(t, 1000, user.HoldAmount)
		assert.Equal(t, 2000, club.GetFund())
		testIo.skipText()

	}

	club.SaveClub()

	club = CreateClub(&testCfg)
	{
		user := club.GetUser(makeChat("usera"))
		assert.Equal(t, 1000, user.HoldAmount)

		user = club.GetUser(makeChat("trifid"))
		assert.Equal(t, 1000, user.HoldAmount)

	}
	assert.Equal(t, 2000, club.GetFund())

}



func TestUserSetSalary(t *testing.T) {

	testCfg.Path = "cfg.json"
	club := CreateClub(&testCfg)

	testIo := CreateTestIoService()
	club.IoService = testIo

	user := club.GetUser(makeChat("trifid"))

	sendText(user, "/start")
	assert.True(t, testIo.checkText("main menu"))

	sendText(user, "Salary")
	assert.True(t, testIo.checkText(HowMuchSal))

	sendText(user, "2000")
	assert.True(t, testIo.checkText(fmt.Sprintf(SalAnswer, 2000, 3000)))
	assert.Equal(t, 2000, user.Salary)
	assert.Equal(t, 3000, user.CreditLimit)
}

func TestUserGetMoney(t *testing.T) {

	testCfg.Path = "cfg.json"
	club := CreateClub(&testCfg)

	testIo := CreateTestIoService()
	club.IoService = testIo
	{
		user := club.GetUser(makeChat("trifid"))

		sendText(user, "/start")
		assert.True(t, testIo.checkText("main menu"))

		sendText(user, "Hold")
		assert.True(t, testIo.checkText(HowMuchHold))

		sendText(user, "2000")
		assert.True(t, testIo.checkText(fmt.Sprintf(HoldDone, 2000)))
		assert.Equal(t, 2000, user.HoldAmount)
		assert.Equal(t, 2000, club.GetFund())
	}
	club.SaveClub()
	club = CreateClub(&testCfg)
	club.IoService = testIo
	{
		user := club.GetUser(makeChat("usera"))

		sendText(user, "/start")
		assert.True(t, testIo.checkText("main menu"))

		sendText(user, "Salary")
		assert.True(t, testIo.checkText(HowMuchSal))

		sendText(user, "1000")
		assert.True(t, testIo.checkText(fmt.Sprintf(SalAnswer, 1000, 1500)))
		assert.Equal(t, 1500, user.CreditLimit)
	}
	club.SaveClub()
	club = CreateClub(&testCfg)
	club.IoService = testIo

	{
		user := club.GetUser(makeChat("usera"))

		sendText(user, "/start")
		assert.True(t, testIo.checkText("main menu"))

		sendText(user, "Borrow")
		assert.True(t, testIo.checkText(fmt.Sprintf(HowMuchTake, 1500, 2000)))

		sendText(user, "700")
		assert.True(t, testIo.checkText(fmt.Sprintf(TakenSucc, 700)))
		assert.Equal(t, 700, user.InCredit)
		testIo.skipText() // skip get notification
	}
	club.SaveClub()
	club = CreateClub(&testCfg)
	club.IoService = testIo

	{
		user := club.GetUser(makeChat("usera"))

		sendText(user, "/start")
		assert.True(t, testIo.checkText("main menu"))

		sendText(user, "Return Money")
		assert.True(t, testIo.checkText(fmt.Sprintf(HowMuchReturn, 700)))

		sendText(user, "300")
		assert.True(t, testIo.checkText(fmt.Sprintf(ReturnSucc, 400)))
		assert.Equal(t, 400, user.InCredit)
		testIo.skipText() // skip return notification
	}
}

func TestUserApDistrub(t *testing.T) {
	club, testIo := createTestClub()

	user := club.GetUser(makeChat("trifid"))

	sendText(user, "Hold")
	assert.True(t, testIo.checkText(HowMuchHold))

	sendText(user, "2000")
	assert.True(t, testIo.checkText(fmt.Sprintf(HoldDone, 2000)))
	assert.Equal(t, 2000, user.HoldAmount)

	user.DistrubAp()
	assert.Equal(t, 2.74, user.AP)

	user.DistrubAp()
	assert.Equal(t, 5.48, user.AP)

	user.DistrubAp()
	assert.Equal(t, 8.22, user.AP)

}

func TestUserTimer(t *testing.T) {



}