package botfunc

import (
	"UkinoShop/internal/helper"
	"UkinoShop/pkg/db"
	"database/sql"
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	_ "github.com/lib/pq"
	"log"
	"strings"
)

func GetNumberOfUsets() (int64, error) {

	var count int64

	db, err := sql.Open("postgres", db.DbInfo)
	if err != nil {
		log.Println(err)
		return -443, err
	}

	defer db.Close()

	data := db.QueryRow(`SELECT COUNT(DISTINCT username) FROM users;`)
	err = data.Scan(&count)
	if err != nil {
		log.Println(err)
		return -444, err
	}

	return count, nil
}

func StartButtons(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	switch update.CallbackQuery.Data {
	case "shop":
		newKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Sneakers", "sneakers"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("T-shirts", "t-shirts"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("<-Back", "backToStartMenu"),
			),
		)
		helper.EditMessage(update, bot, "Choose Category:", newKeyboard)
	case "help":
		newKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("<-Back", "backToStartMenu"),
			),
		)
		helper.EditMessage(update, bot, "Admin: @ukinomata\nSend message and describe the problem.", newKeyboard)
	case "cart":
		cart, err := db.CartOfUser(update)
		if err != nil {
			log.Println(err)
		}
		list := strings.Join(cart, "\n")
		newKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("<-Back", "backToStartMenu"),
				tgbotapi.NewInlineKeyboardButtonData("Remove Cart", "removeCart"),
			),
		)
		helper.EditMessage(update, bot, fmt.Sprintf("Your cart, bruh))):\n\n%s", list), newKeyboard)
	}
}

func ChooseCategory(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	switch update.CallbackQuery.Data {
	case "sneakers":
		newKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Vans Authentic", "vansAuthentic"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Nike Air Force", "nikeAF"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Converse All Stars", "converseAllStars"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("<-Back", "backToChooseCategory"),
			),
		)

		helper.EditMessage(update, bot, "Choose Sneakers:", newKeyboard)
	case "t-shirts":
		newKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Oversize", "oversize"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Pink", "pink"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Crew", "crew"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("<-Back", "backToChooseCategory"),
			),
		)

		helper.EditMessage(update, bot, "Choose T-shirt:", newKeyboard)
	}
}

func ChooseThing(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	switch update.CallbackQuery.Data {
	case "vansAuthentic":
		if err := db.AppendThingToDB(update, 2); err != nil {
			return
		}
	case "nikeAF":
		if err := db.AppendThingToDB(update, 1); err != nil {
			return
		}
	case "converseAllStars":
		if err := db.AppendThingToDB(update, 3); err != nil {
			return
		}
	case "oversize":
		if err := db.AppendThingToDB(update, 4); err != nil {
			return
		}
	case "pink":
		if err := db.AppendThingToDB(update, 5); err != nil {
			return
		}
	case "crew":
		if err := db.AppendThingToDB(update, 6); err != nil {
			return
		}
	}
}

func Back(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	switch update.CallbackQuery.Data {
	case "backToStartMenu":
		newKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Shop", "shop"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Help", "help"),
				tgbotapi.NewInlineKeyboardButtonData("Cart", "cart"),
			),
		)

		helper.EditMessage(update, bot, "Choose Action:", newKeyboard)
	case "backToChooseCategory":
		newKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Sneakers", "sneakers"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("T-shirts", "t-shirts"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("<-Back", "backToStartMenu"),
			),
		)

		helper.EditMessage(update, bot, "Choose Category:", newKeyboard)
	}
}

func DefaultAnswer(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var callbackMessage string
	if update.CallbackQuery.Message.Chat.ID != 726556686 {
		callbackMessage = "I know that you have interesting,but wait)))"
	} else {
		callbackMessage = "Greetins my Lord!"
	}
	helper.CallBackAnswer(update, bot, callbackMessage)
}
