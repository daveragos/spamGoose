package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var userIDs = make(map[int64]string) // Map to store user-specific unique words
var userIDsMutex sync.Mutex          // Mutex to handle concurrent access to the map

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

	channelID := "@debugging_in_prod"

	words, err := loadWordsFromFile("dicts.json")
	if err != nil {
		log.Fatalf("Failed to load words from file: %v", err)
	}

	processUpdates(bot, updates, channelID, words)
}

func startHealthCheckServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func loadWordsFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data struct {
		Words []string `json:"words"`
	}
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, err
	}

	return data.Words, nil
}

func processUpdates(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel, channelID string, words []string) {
	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				handleCommand(bot, update.Message)
				continue
			}
			messageText := strings.TrimSpace(update.Message.Text)
			if isValidMessage(messageText) {
				sendMessageToChannel(bot, channelID, update.Message.From.ID, messageText, words)
			} else {
				warningMessage := "Your message contains a link or invalid content and cannot be posted."
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, warningMessage)
				if _, err := bot.Send(msg); err != nil {
					log.Printf("Failed to send warning message: %v", err)
				}
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

func generateUniqueWord(userID int64, words []string) string {
	usedWords := make(map[string]bool)
	for _, word := range userIDs {
		usedWords[word] = true
	}

	rand.Seed(time.Now().UnixNano() + userID)
	for {
		word := words[rand.Intn(len(words))]
		if !usedWords[word] {
			return word
		}
	}
}

func sendMessageToChannel(bot *tgbotapi.BotAPI, channelID string, userID int64, text string, words []string) {
	userIDsMutex.Lock()
	uniqueWord, exists := userIDs[userID]
	if !exists {
		uniqueWord = generateUniqueWord(userID, words)
		userIDs[userID] = uniqueWord
	}
	userIDsMutex.Unlock()

	msg := tgbotapi.NewMessageToChannel(channelID, text+"\n\n#"+uniqueWord+" #frombot")
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

	urlRegex := `((?:https?|ftp|ws|wss):\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)|(?:www\.)?[a-zA-Z0-9-]+\.[a-zA-Z]{2,6}(\/[a-zA-Z0-9()@:%_\+.~#?&//=]*)?)`
	allowedLinkRegex := `https://t\.me/debugging_in_prod/(\d+)`

	// Find all links in the message
	allLinks := regexp.MustCompile(urlRegex).FindAllString(msg, -1)
	if len(allLinks) > 0 {
		allowedLinkCount := 0
		for _, link := range allLinks {
			if regexp.MustCompile(allowedLinkRegex).MatchString(link) {
				allowedLinkCount++
			}
		}
		// If there are other links besides the allowed one, reject the message
		if allowedLinkCount != len(allLinks) {
			return false
		}
	}
	return true
}
