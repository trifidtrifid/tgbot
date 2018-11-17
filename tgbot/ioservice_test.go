package tgbot

import (
	"fmt"
	"github.com/mrd0ll4r/tbotapi"
	"testing"
	"time"
)


func TestTT(t *testing.T) {
	timer2 := time.NewTicker(1 * time.Second)


	prevTime := time.Now()
	for now := range timer2.C {

		if (now.Unix()) % 3 == 0 {
			fmt.Println("Timer 2 expired", time.Now())
		} else {
			fmt.Println("sdfsdf ", now.Unix() - prevTime.Unix())

		}
	}
}

func TestMenu(t *testing.T) {


	ioService := CreateBotIoService("707764774:AAGfSYmOolr0YBfiz10lCNkAupmWhvVttRA")

	rcpt := tbotapi.Recipient{}
	rcpt.ChatID = new(int)
	*rcpt.ChatID = 281259469

	ioService.sendTextWithMenu(rcpt, "test")

	time.Sleep(10 * time.Second)


	for {
		botUpdate := <-ioService.Api.Updates
		if botUpdate.Error() != nil {
			// TODO handle this properly
			fmt.Printf("Update error: %s\n", botUpdate.Error())
			continue
		}

		update := botUpdate.Update()
		switch update.Type() {
		case tbotapi.MessageUpdate:
			msg := update.Message
			typ := msg.Type()
			if typ.IsChatAction() {
				fmt.Println("Ignoring chat action")
				return
			}
			if msg.Chat.IsChannel() {
				fmt.Println("Ignoring channel message")
				return
			}

			fmt.Printf("<-%d, From:\t%s, Type: %s Text: %s \n", msg.ID, msg.Chat, typ, *msg.Text)
		case tbotapi.InlineQueryUpdate:
			fmt.Println("Ignoring received inline query: ", update.InlineQuery.Query)
		case tbotapi.ChosenInlineResultUpdate:
			fmt.Println("Ignoring chosen inline query result (ID): ", update.ChosenInlineResult.ID)
		case tbotapi.CallbackQueryUpdate:
			fmt.Printf("Ignoring chosen callback query result (ID): %s from %s, ID: %s\n", update.CallbackQuery.Data,
				update.CallbackQuery.From, update.CallbackQuery.ID)


			answ := ioService.Api.NewOutgoingCallbackQueryResponse(update.CallbackQuery.ID)
			answ.Text = "hello"

			err := answ.Send()
			if err != nil {
				fmt.Printf("Error sending text: %s\n", err)
				return
			}

		default:
			fmt.Println("Ignoring unknown Update type.")
		}
	}

}