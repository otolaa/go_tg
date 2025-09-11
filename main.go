package main

import (
	"fmt"
	"go_tg/config"
	"go_tg/stivenking"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	PL     string = "[+]"
	RS     string = " → "
	DS     string = " / "
	NS     string = "\n"
	PS     string = " "
	Suffix string = "~"
)

var (
	Emoji      []string = []string{"🌻", "🌶️", "🌵", "🚀", "👾", "🍎", "⚙️", "🎲", "🎯", "🏀", "⚽", "🎳", "♥️", "♠️", "♦️", "♣️"}
	SuffixLine string   = strings.Repeat(Suffix, 39)
)

// color: 1 red, 2 green, 3 yello, 4 blue, 5 purple, 6 blue
func p(color int, sep string, str ...any) {
	newStr := []any{}
	for index, v := range str {
		if index == 0 {
			newStr = append(newStr, v)
		} else {
			newStr = append(newStr, sep, v)
		}
	}

	suffixColor := "\033[3" + strconv.Itoa(color) + "m"
	fmt.Printf("%s%s%s", suffixColor, fmt.Sprint(newStr...), "\033[0m\n")
}

func connectWithTg(token string, url string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bot.Debug = false

	p(3, " ~ ", PL, bot.Self.UserName, url)

	whUrl := url + "/" + token
	wh, _ := tgbotapi.NewWebhook(whUrl)
	wh.AllowedUpdates = []string{"message", "edited_channel_post", "callback_query"}
	_, err = bot.Request(wh)
	if err != nil {
		return nil, err
	}

	commandStart := tgbotapi.BotCommand{
		Command:     "start",
		Description: Emoji[3] + " Start bot",
	}

	commandHi := tgbotapi.BotCommand{
		Command:     "settings",
		Description: Emoji[6] + " The settings",
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

func handleButton(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {

	mid := callback.Message.Chat.ID

	// Извлечь данные обратного вызова
	data := callback.Data
	commandParams := strings.Split(data, "_")

	uid64, err := strconv.Atoi(commandParams[1])
	if err != nil {
		return
	}

	sending := true
	switch commandParams[0] {
	case "active":
		sending = true
	case "disable":
		sending = false
	}

	p(4, " ~ ", PL, mid, data)

	us := config.SetUserSending(uint(uid64), sending)
	nameButton, valueButton, callbackButton := getButtonSending(&us)

	// Ответить на запрос обратного вызова
	callbackMess := tgbotapi.NewCallback(callback.ID, callbackButton)
	bot.Request(callbackMess)

	// Опционально — отредактировать сообщение, чтобы отразить выбор
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(nameButton, valueButton),
		),
	)

	edit := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, callbackButton)
	edit.ReplyMarkup = &markup
	bot.Send(edit)
}

func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {

	p(2, " → ", PL, message.Chat.UserName, message.Chat.ID, message.Text)

	// ~~~ add user DB
	userName := message.From.UserName
	if message.Chat.Type == "group" {
		userName = message.Chat.Title
	}

	user := config.SetUser(message.Chat.ID, userName)
	// ~~~ end

	if setStartCommand(bot, message) {
		return
	}

	if setSettingsCommand(bot, message, &user) {
		return
	}

	if setDefaultMessage(bot, message) {
		return
	}
}

func getButtonSending(user *config.User) (string, string, string) {
	nameButton := "❌ Выключить рассылку"
	valueButton := fmt.Sprintf("disable_%d", user.ID)
	callbackButton := "👍 Ваша рассылка включена."

	if !user.Sending {
		nameButton = "✅ Включить рассылку"
		valueButton = fmt.Sprintf("active_%d", user.ID)
		callbackButton = "✋ Ваша рассылка выключена."
	}

	return nameButton, valueButton, callbackButton
}

func setSettingsCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, user *config.User) bool {
	if !strings.HasPrefix(message.Text, "/settings") {
		return false
	}

	msgArr := []string{
		"🎲 → Случайные цитаты из книг.",
		SuffixLine,
		fmt.Sprintf("📌 Вы → @%s", message.From.UserName),
		fmt.Sprintf("🏀 Ваш id → %d", message.From.ID),
		SuffixLine,
		fmt.Sprintf("🕜 → %s", time.Now().Format("15:04 ~ 02.01.2006")),
		fmt.Sprintf("✉️ → рассылка ↓ по часовому поясу %s", time.Now().Format("MST")),
		fmt.Sprintf("⏰ → %s часы", "10,11,12,13,14,15,16,17,18,19"),
		SuffixLine,
		fmt.Sprintf("%s → %s ~ версия", Emoji[15], config.VERSION),
	}

	nameButton, valueButton, _ := getButtonSending(user)

	msg := tgbotapi.NewMessage(message.Chat.ID, strings.Join(msgArr, NS))
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(nameButton, valueButton),
		),
	)
	bot.Send(msg)

	return true
}

// command start
func setStartCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) bool {
	if !strings.HasPrefix(message.Text, "/start") {
		return false
	}

	msgArr := []string{
		"Задайте свой вопрос, не забудьте `?`",
		"Ответом на вопрос будит цитата из книг.",
		SuffixLine,
		Emoji[4] + RS + "Стивена Кинга",
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, strings.Join(msgArr, NS))
	bot.Send(msg)

	return true
}

// default message
func setDefaultMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) bool {
	if !strings.Contains(message.Text, "?") {
		return false
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, stivenking.GetQuote())
	msg.ReplyToMessageID = message.MessageID
	bot.Send(msg)
	return true
}

func main() {
	bot, err := connectWithTg(config.TOKEN, config.URL_BOT)
	if err != nil {
		log.Fatal(err)
	}

	updates := bot.ListenForWebhook("/" + config.TOKEN)
	http.HandleFunc("/", setTest)
	go http.ListenAndServe(":8080", nil)

	for update := range updates {
		switch {
		// Handle messages
		case update.Message != nil:
			handleMessage(bot, update.Message)

			// Handle button clicks
		case update.CallbackQuery != nil:
			handleButton(bot, update.CallbackQuery)
		}
	}
}
