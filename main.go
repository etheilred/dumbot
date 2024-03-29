package main

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"myProj5/pkg/utils"
	"os"
	"time"
)

var numOfUpdates = 0
var numOfMessages = 0
var numOfCallbackQueries = 0

func httpDaemon() {
	m := martini.Classic()
	m.Get("/", func() string {
		reqTime := time.Now().Format("2 Jan 2006 15:04:05")
		return fmt.Sprintf(
			"<h1>Stats on %s</h1>\nNumber of upds: %d<br/>\nNumber of msgs: %d<br/>\nNumber of callback queries: %d<br/>",
			reqTime,
			numOfUpdates,
			numOfMessages,
			numOfCallbackQueries)
	})
	m.Run()
}

var users = make(map[int64]utils.UserPref)

func handleCommands(bot *tgbotapi.BotAPI, upd tgbotapi.Update) {
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "h")
	if upd.Message.Text == "/start" {
		usr, ok := users[upd.Message.Chat.ID]
		if ok {
			msg.Text = "Here we go again..."
		} else {
			msg.Text = "Nice to meet you!"
			usr = utils.UserPref{
				ID: upd.Message.Chat.ID,
			}
			usr.SetInlineKeyboard(tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("press", "press"),
				tgbotapi.NewInlineKeyboardButtonData("don't press", "nopress"),
			),
			)
			msg.ReplyMarkup = usr.InlineKeyboard
			users[usr.ID] = usr
		}
	} else if upd.Message.Text == "/reset" {
		delete(users, upd.Message.Chat.ID)
		msg.Text = "Wow, seems that I have forgotten yous"
	}
	bot.Send(msg)
}

func handleText(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	if update.Message.Text == "hello" {
		msg.Text = "hi"
	} else {
		msg.Text = "((("
	}
	bot.Send(msg)
}

func main() {
	go httpDaemon()
	token := os.Getenv("TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal("Failed to get bot updates channel")
	}

	for update := range updates {
		numOfUpdates++
		if update.CallbackQuery != nil {
			numOfCallbackQueries++
			// log.Println("[CallbackData]", update.CallbackQuery.Data)
			// bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID,update.CallbackQuery.Data))
			bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data))
		}

		if update.Message != nil { // ignore any non-Message Updates
			numOfMessages++
			log.Printf("[%s] %s, is command %b",
				update.Message.From.UserName,
				update.Message.Text,
				update.Message.IsCommand(),
			)

			if update.Message.IsCommand() {
				go handleCommands(bot, update)
			} else {
				go handleText(bot, update)
			}
		}
	}
}
