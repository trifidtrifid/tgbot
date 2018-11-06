package main

import (
	"fmt"
	"github.com/mrd0ll4r/tbotapi"
	"github.com/trifidtrifid/tgbot/tgbot"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {

	club := tgbot.CreateClub()
	ioService := tgbot.CreateBotIoService("768558434:AAHJnCN-A4k-kzc3DdlywUP8tuH8rs8ni4Q")
	club.IoService = ioService


	closing := make(chan struct{})
	closed := make(chan struct{})
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-closed:
				return
			case botUpdate := <-ioService.Api.Updates:
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
					if msg.Chat.Username != nil {
						user := club.GetUser(*msg.Chat.Username)
						user.Msgs <- tgbot.UserMessage{*msg}
					} else {
						user := club.GetUser(string(msg.Chat.ID))
						user.Msgs <- tgbot.UserMessage{*msg}
					}
				case tbotapi.InlineQueryUpdate:
					fmt.Println("Ignoring received inline query: ", update.InlineQuery.Query)
				case tbotapi.ChosenInlineResultUpdate:
					fmt.Println("Ignoring chosen inline query result (ID): ", update.ChosenInlineResult.ID)
				default:
					fmt.Println("Ignoring unknown Update type.")
				}
			}
		}
	}()

	// Ensure a clean shutdown.
	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-shutdown
		close(closing)
	}()

	fmt.Println("Bot started. Press CTRL-C to close...")

	// Wait for the signal.
	<-closing
	fmt.Println("Closing...")

	// Always close the API first, let it clean up the update loop.
	// This might take a while.
	ioService.Api.Close()
	close(closed)
	wg.Wait()

}
