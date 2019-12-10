package main

import (
	"github.com/go-martini/martini"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func httpDaemon() {
	m := martini.Classic()
	m.Get("/", func() string {
		return "Hello world!"
	})
	m.Run()
}

type userPref struct {
	ID int64
}

var users = make(map[int64]userPref)

func handleCommands(bot *tgbotapi.BotAPI, upd tgbotapi.Update) {
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "")
	if upd.Message.Text == "start" {
		usr, ok := users[upd.Message.Chat.ID]
		if ok {
			msg.Text = "Here we go again..."
		} else {
			msg.Text = "Nice to meet you!"
			usr = userPref{
				ID: upd.Message.Chat.ID,
			}
			users[usr.ID] = usr
		}
	} else if upd.Message.Text == "reset" {
		delete(users, upd.Message.Chat.ID)
		msg.Text = "Wow, seems that I have forgotten yous"
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

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			handleCommands(bot, update)
		}
	}
}
