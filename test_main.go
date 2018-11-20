package main

import (
	"fmt"
	"github.com/trifidtrifid/tbotapi"
	"github.com/trifidtrifid/tgbot/tgbot"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

func main() {

	if len(os.Args) != 2 {
		log.Fatal("Specify config file path")
	}

	var cfg tgbot.Config
	cfg.Path = os.Args[1]

	if !cfg.Load() {
		log.Fatalf("Cannot load config file %s", cfg.Path)
	}

	club := tgbot.CreateClub(&cfg)

	ioService := tgbot.CreateBotIoService("707764774:AAGfSYmOolr0YBfiz10lCNkAupmWhvVttRA") //test
	//ioService := tgbot.CreateBotIoService("768558434:AAHJnCN-A4k-kzc3DdlywUP8tuH8rs8ni4Q")
	club.IoService = ioService

	club.RunApDistribution()

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
					user := club.GetUser(msg.Chat)
					if user == nil {
						fmt.Printf("Cannot find or create user for chat: %s\n", msg.Chat)
						return

					}
					user.Msgs <- tgbot.UserMessage{Message: *msg}
				case tbotapi.InlineQueryUpdate:
					fmt.Println("Ignoring received inline query: ", update.InlineQuery.Query)
				case tbotapi.ChosenInlineResultUpdate:
					fmt.Println("Ignoring chosen inline query result (ID): ", update.ChosenInlineResult.ID)
				case tbotapi.CallbackQueryUpdate:
					fmt.Printf("CallbackQueryUpdate. Data: %s, From %s, ID: %s\n", update.CallbackQuery.Data,
						update.CallbackQuery.From, update.CallbackQuery.ID)

					userFrom := club.FindUser(strconv.Itoa(update.CallbackQuery.From.ID))

					var userIdStr string
					var apAmount int
					fmt.Sscanf(update.CallbackQuery.Data, "%s %d", &userIdStr, &apAmount)

					userTo := club.FindUser(userIdStr)

					answ := ioService.Api.NewOutgoingCallbackQueryResponse(update.CallbackQuery.ID)

					if userFrom.AP < float64(apAmount) {
						answ.Text = fmt.Sprintf("not enough ap. you have only %.2f", userFrom.AP)
					} else {
						userTo.AP += float64(apAmount)
						userFrom.AP -= float64(apAmount)
						answ.Text = fmt.Sprintf("You send %d ap. You have %.2f", apAmount, userFrom.AP)
						userTo.SendText(fmt.Sprintf("You receive %d AP from %s", apAmount, update.CallbackQuery.From))
					}

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
