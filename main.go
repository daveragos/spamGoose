package main

import (
	"log"
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

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)

	channelID := "@ragoosedumps"

	for update := range updates {
		if isValidMessage(strings.Trim(update.Message.Text, " ")) && !update.Message.IsCommand() {
			msg := tgbotapi.NewMessageToChannel(channelID, update.Message.Text+"\n\n#frombot")

			if _, err := bot.Send(msg); err != nil {
				log.Printf("Failed to send message to channel: %v", err)
			}
		}
	}
}

// rule for the message to be sent to the channel
func isValidMessage(msg string) bool {
	if msg == "" {
		return false
	}
	if len(msg) > 1 && msg[0] == '@' {
		return false
	}
	urlRegex := `(?i)\b((?:https?://|www\.|)\S+\.\S+)`
	if matched, _ := regexp.MatchString(urlRegex, msg); matched {
		return false
	}
	return true
}
