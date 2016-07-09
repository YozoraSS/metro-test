package main

import (
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	os.Exit(Main())
}

// Main function for echo bot
func Main() int {
	var (
		channelID     1465528493
		channelSecret = os.Getenv("9d668b0329734daa0612c2d17b49f1b2")
		channelMID    = os.Getenv("ua878bfa47c8b079685b8915c4fa40ffa")
		err           error
	)

	// Setup bot client
	channelID, err = strconv.ParseInt(os.Getenv("CHANNEL_ID"), 10, 64)
	if err != nil {
		log.Print(err)
		return 1
	}
	bot, err := linebot.NewClient(channelID, channelSecret, channelMID)
	if err != nil {
		log.Print(err)
		return 1
	}

	// Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		received, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}
		for _, result := range received.Results {
			content := result.Content()
			if content != nil && content.IsMessage && content.ContentType == linebot.ContentTypeText {
				text, err := content.TextContent()
				_, err = bot.SendText([]string{content.From}, text.Text)
				if err != nil {
					log.Print(err)
				}
			}
		}
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Print(err)
		return 1
	}
	return 0
}
