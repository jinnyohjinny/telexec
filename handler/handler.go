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

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("token empty")
	}

	var errBot error

	bot, errBot = tele.NewBot(tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})

	return errBot
}

func streamChat() {
	bot.Handle("/halo", func(ctx tele.Context) error {
		return ctx.Send("Halo dunia")
	})
}

func Begin() {
	fmt.Println("Bot dimulai...")
	if err := initBot(); err != nil {
		log.Fatal(err)
	}

	streamChat()
	bot.Start()
}
