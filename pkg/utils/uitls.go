package utils

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type UserPref struct {
	ID int64
	InlineKeyboard tgbotapi.InlineKeyboardMarkup
	Keyboard tgbotapi.ReplyKeyboardMarkup
}

func (usr *UserPref) SetInlineKeyboard(buttons []tgbotapi.InlineKeyboardButton) {
	usr.InlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(buttons)
}

func (usr *UserPref) SetKeyboard(buttons []tgbotapi.KeyboardButton) {
	usr.Keyboard = tgbotapi.NewReplyKeyboard(buttons)
}