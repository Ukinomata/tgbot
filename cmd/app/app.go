package main

import (
	"UkinoShop/internal/botfunc"
	"UkinoShop/internal/helper"
	"UkinoShop/pkg/db"
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

func main() {
	Start()
}

func Start() {
	// создаем бота
	updates, bot := helper.CreateBot()

	for update := range updates {
		if update.Message != nil {
			switch update.Message.Text {
			case "/start":
				db.RegisterUser(update.Message.Chat.UserName, update.Message.Chat.ID)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome to UkinoShop!")
				bot.Send(msg)

				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Choose Action:")

				keyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Shop", "shop"),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Help", "help"),
						tgbotapi.NewInlineKeyboardButtonData("Cart", "cart"),
					),
				)

				msg.ReplyMarkup = keyboard
				bot.Send(msg)
			case "/number_of_users":
				count, _ := botfunc.GetNumberOfUsets()
				answ := fmt.Sprintf("%d users used this bot", count)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, answ)
				bot.Send(msg)
			case "/test":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "THIS COMMAND FOR DEVELOPER!!!*Ukino's pet is angry*")

				btn1 := tgbotapi.NewInlineKeyboardButtonData("If you my master you can press this button", "MasterButton")

				row := tgbotapi.NewInlineKeyboardRow(btn1)

				keyboard := tgbotapi.NewInlineKeyboardMarkup(row)

				msg.ReplyMarkup = keyboard
				bot.Send(msg)
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "What do you want from me?")
				bot.Send(msg)
			}
		} else if update.CallbackQuery != nil {
			switch update.CallbackQuery.Data {
			case "shop", "help", "cart":
				botfunc.StartButtons(update, bot)
			case "sneakers", "t-shirts":
				botfunc.ChooseCategory(update, bot)
			case "vansAuthentic", "nikeAF", "converseAllStars", "oversize", "pink", "crew":
				botfunc.ChooseThing(update, bot)
			case "removeCart":
				db.RemoveCart(update, bot)
			case "backToStartMenu", "backToChooseCategory":
				botfunc.Back(update, bot)
			default:
				botfunc.DefaultAnswer(update, bot)
			}
		}
	}
}
