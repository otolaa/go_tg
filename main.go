package main

import (
	"fmt"
	"go_tg/config"
	"log"
	"net/http"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// color: 1 red, 2 green, 3 yello, 4 blue, 5 purple, 6 blue
func p(color int, str ...any) {
	suffixColor := "\033[3" + strconv.Itoa(color) + "m"
	fmt.Printf("%s%s%s", suffixColor, fmt.Sprint(str...), "\033[0m\n")
}

func connectWithTg(token string, url string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bot.Debug = false

	log.Printf("authorized on account ~ [%s]", bot.Self.UserName)

	p(3, url)
	p(3, token)

	whUrl := url + "/" + token
	wh, _ := tgbotapi.NewWebhook(whUrl)
	wh.AllowedUpdates = []string{"message", "edited_channel_post", "callback_query"}
	_, err = bot.Request(wh)
	if err != nil {
		return nil, err
	}

	commandStart := tgbotapi.BotCommand{
		Command:     "start",
		Description: "🚀 Start bot",
	}

	commandHi := tgbotapi.BotCommand{
		Command:     "hi",
		Description: "🌵 The version",
	}

	bc := tgbotapi.NewSetMyCommands(commandStart, commandHi)
	_, err = bot.Request(bc)
	if err != nil {
		return nil, err
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		return nil, err
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	return bot, nil
}

func setTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("it's ok, v" + config.VERSION))
}

func receiveUpdates(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		p(2, "[+] → ", update.Message.Chat.UserName, " → ", update.Message.Chat.ID, " → ", update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
		continue
	}
}

func main() {
	bot, err := connectWithTg(config.TOKEN, config.URL_BOT)
	if err != nil {
		log.Fatal(err)
	}

	updates := bot.ListenForWebhook("/" + config.TOKEN)
	http.HandleFunc("/", setTest)
	go http.ListenAndServe(":8080", nil)

	receiveUpdates(bot, updates)
}
