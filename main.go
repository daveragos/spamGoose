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
			log.Printf("Received message: %s", update.Message.Text)
			log.Printf("Message from user: %s", update.Message.From.UserName)
			log.Printf("Message from user ID: %d", update.Message.From.ID)
			log.Printf("Message from chat ID: %d", update.Message.Chat.ID)
			log.Printf("Message from chat type: %s", update.Message.Chat.Type)
			log.Printf("Message from chat title: %s", update.Message.Chat.Title)
			log.Printf("Message from chat username: %s", update.Message.Chat.UserName)
			log.Printf("Message from chat first name: %s", update.Message.Chat.FirstName)
			log.Printf("Message from chat last name: %s", update.Message.Chat.LastName)
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
