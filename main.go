package main

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is not set")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	go startHealthCheckServer()

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)

	channelID := "@ragoose_dumps"

	processUpdates(bot, updates, channelID)
}

func startHealthCheckServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func processUpdates(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel, channelID string) {
	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				handleCommand(bot, update.Message)
				continue
			}
			if isValidMessage(strings.TrimSpace(update.Message.Text)) {
				sendMessageToChannel(bot, channelID, update.Message.Text)
			}
		}
	}
}

func handleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	switch message.Command() {
	case "start":
		startMessage := `Welcome! This bot is inspired by Jack Rhysider's experiment, allowing you to post directly to @ragoose_dumps Telegram channel. 
Simply send your message, and it will appear @ragoose_dumps with the hashtag #frombot. 
Please note that links and files are not supported. for obvious reasons.`
		msg := tgbotapi.NewMessage(message.Chat.ID, startMessage)
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send start message: %v", err)
		}
	}
}

func sendMessageToChannel(bot *tgbotapi.BotAPI, channelID, text string) {
	msg := tgbotapi.NewMessageToChannel(channelID, text+"\n\n#frombot")
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Failed to send message to channel: %v", err)
	}
}

func isValidMessage(msg string) bool {
	if msg == "" {
		return false
	}
	if strings.Contains(msg, "@") {
		return false
	}

	urlRegex := `((?:https?|ftp|ws|wss):\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)|(?:[a-zA-Z0-9_]+\.)?t\.me\/[a-zA-Z0-9_]+(\/[a-zA-Z0-9_]+)?)`
	if matched, _ := regexp.MatchString(urlRegex, msg); matched {
		return false
	}
	return true
}
