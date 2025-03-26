package handler

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jinnyohjinny/telexec/controller"
	"github.com/jinnyohjinny/telexec/utils"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	tele "gopkg.in/telebot.v4"
)

var bot *tele.Bot

const (
	formatOk = `
Running : %s
===========
Out :
===========
 %s`

	formatErr = `
Running : %s
===========
Err :
===========
 %s`
)

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

func cmdHandler(log zerolog.Logger, commandExec controller.CmdOutputWriter) {
	log.Info().Msg("Registering /cmd handler")

	bot.Handle("/cmd", func(ctx tele.Context) error {
		if !ctx.Message().Private() {
			_, err := ctx.Bot().Send(ctx.Chat(), "Please use this command in private chat")
			return err
		}

		commandName := ctx.Message().Payload
		log.Info().Str("state", "exec").Msg(commandName)

		log.Debug().Msg("Before execution")
		outOk, outErr, errExec := commandExec.ExecOutput(commandName)
		log.Debug().
			Str("output", string(outOk)).
			Str("error", string(outErr)).
			Err(errExec).
			Msg("After ExecOutput")

		if errExec != nil {
			log.Error().Str("state", "exec").Msg(errExec.Error())
			return ctx.Reply(fmt.Sprintf(formatErr, commandName, outErr))
		}
		fmt.Printf("Payload: %s\n", ctx.Message().Payload)
		return ctx.Reply(fmt.Sprintf(formatOk, commandName, outOk))
	})

	bot.Handle(tele.OnText, func(ctx tele.Context) error {
		log.Info().
			Str("text", ctx.Text()).
			Bool("private", ctx.Message().Private()).
			Msg("Received text message")
		return nil
	})
}

func Begin() {
	log.Println("Starting bot...")

	if err := initBot(); err != nil {
		log.Fatalf("Failed to initialize bot: %v", err)
	}

	log := utils.InitLog()
	cmdExec := controller.CmdOutputWriter{
		TimeoutSecond: 10,
	}

	cmdHandler(log, cmdExec)

	me := bot.Me
	log.Info().
		Str("username", me.Username).
		Str("first_name", me.FirstName).
		Msg("Bot initialized successfully")

	log.Info().Msg("Starting bot polling...")
	bot.Start()
}
