package handler

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jinnyohjinny/telexec/controller"
	"github.com/jinnyohjinny/telexec/utils"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	tele "gopkg.in/telebot.v4"
)

var bot *tele.Bot

const (
	maxOutputLength = 2000 // Batas maksimal output yang dikirim ke Telegram
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
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
		OnError: func(err error, ctx tele.Context) {
			log.Printf("Telebot error: %v", err)
		},
	})

	return err
}

func cmdHandler(log zerolog.Logger, cmdExec *controller.CmdOutputWriter) {
	log.Info().Msg("Registering command handlers")

	// Handler untuk /run
	bot.Handle("/run", func(ctx tele.Context) error {
		if !ctx.Message().Private() {
			return ctx.Reply("⚠️ Please use this command in private chat")
		}

		command := strings.TrimSpace(ctx.Message().Payload)
		if command == "" {
			return ctx.Reply("❌ Please specify a command\nExample: <code>/run ls -l</code>")
		}

		log.Info().Str("command", command).Msg("Executing command")

		// Eksekusi command dengan timeout
		out, errOut, err := cmdExec.ExecOutput(command)
		if err != nil {
			log.Error().
				Str("command", command).
				Str("error", string(errOut)).
				Msg("Command failed")

			// Format error output
			errorMsg := fmt.Sprintf(`
🚨 <b>Command Failed</b>
─────────────────
<b>Command:</b> <code>%s</code>
<b>Error:</b> %s
<b>Output:</b>
<pre>%s</pre>
`, command, err.Error(), truncateOutput(errOut))

			return ctx.Reply(errorMsg, tele.ModeHTML)
		}

		// Format success output
		successMsg := fmt.Sprintf(`
✅ <b>Command Executed</b>
─────────────────
<b>Command:</b> <code>%s</code>
<b>Output:</b>
<pre>%s</pre>
`, command, truncateOutput(out))

		return ctx.Reply(successMsg, tele.ModeHTML)
	})

	// Handler untuk /getfile
	bot.Handle("/getfile", func(ctx tele.Context) error {
		if !ctx.Message().Private() {
			return ctx.Reply("⚠️ Please use this command in private chat")
		}

		filename := strings.TrimSpace(ctx.Message().Payload)
		if filename == "" {
			return ctx.Reply("❌ Please specify filename\nExample: <code>/getfile main.go</code>")
		}

		// Validasi path untuk keamanan
		if strings.Contains(filename, "../") || strings.HasPrefix(filename, "/") {
			return ctx.Reply("⚠️ Invalid file path")
		}

		fullPath := filepath.Join(cmdExec.WorkDir, filename)
		log.Info().Str("filename", fullPath).Msg("File download requested")

		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			log.Error().Str("filename", fullPath).Msg("File not found")
			return ctx.Reply("❌ File not found")
		}

		file := &tele.Document{
			File:     tele.FromDisk(fullPath),
			FileName: filepath.Base(filename),
			MIME:     getMimeType(filename),
		}

		err := ctx.Reply(file)
		return err
	})

	// Handler untuk pesan teks biasa
	bot.Handle(tele.OnText, func(ctx tele.Context) error {
		log.Info().
			Str("text", ctx.Text()).
			Str("sender", ctx.Sender().Username).
			Msg("Received message")
		return nil
	})
}

// Helper function untuk memotong output yang terlalu panjang
func truncateOutput(output []byte) string {
	str := string(output)
	if len(str) > maxOutputLength {
		return str[:maxOutputLength] + "\n... (output truncated)"
	}
	return str
}

// Helper function untuk menentukan MIME type
func getMimeType(filename string) string {
	switch filepath.Ext(filename) {
	case ".txt", ".go", ".c", ".cpp", ".h", ".java", ".py", ".sh", ".md":
		return "text/plain"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".pdf":
		return "application/pdf"
	default:
		return "application/octet-stream"
	}
}

func Begin() {
	log.Println("Starting bot...")

	// Inisialisasi logger
	logger := utils.InitLog()

	// Inisialisasi bot
	if err := initBot(); err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize bot")
	}

	// Buat direktori output
	outDir := filepath.Join(".", "out")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		logger.Error().Err(err).Str("path", outDir).Msg("Failed to create output directory")
	}

	// Inisialisasi command executor
	cmdExec := controller.NewCmdOutputWriter(30, outDir) // Timeout 30 detik

	// Daftarkan handler
	cmdHandler(logger, cmdExec)

	// Info startup
	me := bot.Me
	logger.Info().
		Str("username", me.Username).
		Str("first_name", me.FirstName).
		Msg("Bot initialized successfully")

	logger.Info().Msg("Starting bot polling...")
	bot.Start()
}
