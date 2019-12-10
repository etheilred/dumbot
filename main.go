package main

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"myProj5/pkg/utils"
)

func httpDaemon() {
	m := martini.Classic()
	m.Get("/", func() string {
		return "Hello world!"
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

func handleText(bot *tgbotapi.BotAPI, update tgbotapi.Update)  {
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
	bot, err := tgbotapi.NewBotAPI("1069764716:AAFkM-JdVVuA5nsh_gwhFGBO30Oc_kwjQVE")
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
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s, is command %b",
			update.Message.From.UserName,
			update.Message.Text,
			update.Message.IsCommand(),
		)
		if update.CallbackQuery != nil {
			log.Println("[CallbackData]", update.CallbackQuery.Data)
			bot.Send(
				tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("done: %s", update.CallbackQuery.Data)),
				)
		} else if update.Message.IsCommand() {
			go handleCommands(bot, update)
		} else {
			go handleText(bot, update)
		}
	}
}
