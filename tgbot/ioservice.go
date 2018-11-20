package tgbot

import (
	"fmt"
	"github.com/trifidtrifid/tbotapi"
	"log"
)

type IoService interface {
	sendMainMenu(recipient tbotapi.Recipient)
	sendText(recipient tbotapi.Recipient, text string)
	sendApRequest(recipient tbotapi.Recipient, text string, originator string)
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
		Keyboard:        [][]tbotapi.KeyboardButton{
			{{Text: "Hold"}},
			{{Text: "Unhold"}},
			{{Text: "Salary"}},
			{{Text: "Info"}},
			{{Text: "Borrow"}},
			{{Text: "Return Money"}},
			{{Text: "Users"}},
		    {{Text: "/ask_for_ap"}}},
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
	msg := bot.Api.NewOutgoingMessage(recipient, text)
	msg.ParseMode = tbotapi.ModeMarkdown
	outMsg, err := msg.Send()
	if err != nil {
		fmt.Printf("Error sending text: %s\n", err)
		return
	}
	fmt.Printf("send text ->%d, To:\t%s, Text: %s\n", outMsg.Message.ID, outMsg.Message.Chat, *outMsg.Message.Text)

}

func (bot *BotIoService) sendApRequest(recipient tbotapi.Recipient, text string, originator string) {

	// Now simply echo that back.
	msg := bot.Api.NewOutgoingMessage(recipient, text)
	msg.ParseMode = tbotapi.ModeMarkdown

	msg.SetInlineKeyboardMarkup(tbotapi.InlineKeyboardMarkup{
		InlineKeyboard:        [][]tbotapi.InlineKeyboardButton{
			{{Text: "send 20 AP", CallbackData: originator + " " + "20"}},
			{{Text: "send 100 AP", CallbackData: originator + " " + "100"}},
			{{Text: "send 500 AP", CallbackData: originator + " " + "500"}},
			{{Text: "send 5000 AP", CallbackData: originator + " " + "5000"}}},
	})

	outMsg, err := msg.Send()
	if err != nil {
		fmt.Printf("Error sending text: %s\n", err)
		return
	}
	fmt.Printf("send text with inline menu ->%d, To:\t%s, Text: %s Cmd: %s\n",
		outMsg.Message.ID, outMsg.Message.Chat, *outMsg.Message.Text, originator)
}
