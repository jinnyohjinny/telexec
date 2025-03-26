package handler

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
)

var bot *tele.Bot

func initBot() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		return fmt.Errorf("bot token is empty")
	}

	var err error
	bot, err = tele.NewBot(tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 30 * time.Second},
	})

	return err
}

func cmdHandler() {
	bot.Handle("/cmd", func(ctx tele.Context) error {
		log.Printf("Received command from %d", ctx.Sender().ID)

		if !ctx.Message().Private() {
			_, err := ctx.Bot().Send(ctx.Chat(), "Please use this command in private chat")
			return err
		}

		fmt.Printf("Payload: %s\n", ctx.Message().Payload)
		return ctx.Send("Hello world!")
	})

	bot.Handle(tele.OnText, func(ctx tele.Context) error {
		log.Printf("Received text: %s", ctx.Text())
		return nil
	})
}

func Begin() {
	log.Println("Starting bot...")

	if err := initBot(); err != nil {
		log.Fatalf("Failed to initialize bot: %v", err)
	}

	cmdHandler()

	log.Println("Handlers registered, starting bot...")

	bot.Start()
}
