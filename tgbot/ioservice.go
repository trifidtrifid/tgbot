package tgbot

import (
	"fmt"
	"github.com/mrd0ll4r/tbotapi"
	"log"
)

type IoService interface {
	sendMainMenu(recipient tbotapi.Recipient)
	sendText(recipient tbotapi.Recipient, text string)
}


type BotIoService struct {
	Api *tbotapi.TelegramBotAPI
}

func CreateBotIoService(token string) *BotIoService {

	var err error
	var bot *BotIoService = new(BotIoService)
	bot.Api, err = tbotapi.New(token)
	if err != nil {
		log.Print(err)
		return nil
	}

	fmt.Println("Starting...")

	// Just to show its working.
	fmt.Printf("User ID: %d\n", bot.Api.ID)
	fmt.Printf("Bot Name: %s\n", bot.Api.Name)
	fmt.Printf("Bot Username: %s\n", bot.Api.Username)

	return bot
}

func (bot *BotIoService) sendMainMenu(recipient tbotapi.Recipient) {

	toSend := bot.Api.NewOutgoingMessage(recipient, "Select action")
	toSend.SetReplyKeyboardMarkup(tbotapi.ReplyKeyboardMarkup{
		Keyboard:        [][]tbotapi.KeyboardButton{{{Text: "Hold"}},
			{{Text: "Salary"}},
			{{Text: "Info"}},
			{{Text: "Get money"}},
			{{Text: "Return Money"}}},
		OneTimeKeyboard: true,
	})

	// Send it.
	outMsg, err := toSend.Send()
	if err != nil {
		fmt.Printf("Error sending main menu: %s\n", err)
		return
	}
	fmt.Printf("->%d, To:\t%s, Text: %s\n", outMsg.Message.ID, outMsg.Message.Chat, *outMsg.Message.Text)
}

func (bot *BotIoService) sendText(recipient tbotapi.Recipient, text string) {

	// Now simply echo that back.
	outMsg, err := bot.Api.NewOutgoingMessage(recipient, text).Send()
	if err != nil {
		fmt.Printf("Error sending text: %s\n", err)
		return
	}
	fmt.Printf("send text ->%d, To:\t%s, Text: %s\n", outMsg.Message.ID, outMsg.Message.Chat, *outMsg.Message.Text)

}