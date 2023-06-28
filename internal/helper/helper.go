package helper

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
)

// создание бота

func CreateBot() (tgbotapi.UpdatesChannel, *tgbotapi.BotAPI) {
	bot, err := tgbotapi.NewBotAPI("TOKEN")

	if err != nil {
		log.Println(err)
		return nil, nil
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Println(err)
		return nil, nil
	}

	return updates, bot
}

// Для ответа при нажатии кнопки

func CallBackAnswer(update tgbotapi.Update, bot *tgbotapi.BotAPI, answer string) {
	callbackMessage := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, answer)
	_, err := bot.AnswerCallbackQuery(tgbotapi.CallbackConfig{
		CallbackQueryID: update.CallbackQuery.ID,
		Text:            "",
		ShowAlert:       false,
	})
	if err != nil {
		log.Println(err)
	}
	_, err = bot.Send(callbackMessage)
	if err != nil {
		log.Panic(err)
	}
}

// Для замены клавиатуры и текста сообщения

func EditMessage(update tgbotapi.Update, bot *tgbotapi.BotAPI, answer string, newKeyboard tgbotapi.InlineKeyboardMarkup) {
	message := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, answer)

	message.ReplyMarkup = &newKeyboard
	bot.Send(message)
}

// Для замены клавиатуры

func EditKeyboard(update tgbotapi.Update, bot *tgbotapi.BotAPI, newKeyboard tgbotapi.InlineKeyboardMarkup) {
	editMessage := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, newKeyboard)
	_, err := bot.Send(editMessage)
	if err != nil {
		log.Println(err)
		return
	}
}
