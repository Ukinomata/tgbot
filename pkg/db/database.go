package db

import (
	"UkinoShop/internal/helper"
	"database/sql"
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	_ "github.com/lib/pq"
	"log"
	"strings"
)

var (
	user     = "user"
	password = "password"
	host     = "host"
	dbname   = "dbname"
	sslmode  = "disable"
)

var DbInfo = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", user, password, host, dbname, sslmode) //DataSourceName

func CreateTable() error {
	db, err := sql.Open("postgres", DbInfo)

	if err != nil {
		log.Println(err)
		return err
	}

	defer db.Close()

	if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS users(
    ID SERIAL PRIMARY KEY ,
    TIMESTAMP TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    USERNAME TEXT,
    CHAT_ID INT
)`); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func RegisterUser(username string, chatid int64) error {

	db, err := sql.Open("postgres", DbInfo)

	if err != nil {
		log.Println(err)
		return err
	}

	defer db.Close()

	data := `INSERT INTO users(username,chat_id) SELECT $1,$2 WHERE NOT EXISTS(SELECT * FROM users WHERE chat_id = $2);`

	if _, err = db.Exec(data, `@`+username, chatid); err != nil {
		log.Println(err)
		return err
	}

	return err
}

func ClearData() {
	db, err := sql.Open("postgres", DbInfo)
	if err != nil {
		return
	}

	defer db.Close()

	if _, err = db.Exec(`TRUNCATE TABLE users`); err != nil {
		return
	}
}

func TryConnect() error {
	db, err := sql.Open("postgres", DbInfo)
	if err != nil {
		log.Println(err)
		return err
	}

	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Println(err)
		return err
	}

	log.Println("You connect to db")
	return nil
}

func AppendThingToDB(update tgbotapi.Update, id int) error {
	db, err := sql.Open("postgres", DbInfo)
	if err != nil {
		log.Println(err)
		return err
	}

	defer db.Close()

	data := `INSERT INTO cart(user_id, product_id)
VALUES (
        (SELECT id FROM users WHERE chat_id = $1),$2
       )`

	if _, err = db.Exec(data, update.CallbackQuery.Message.Chat.ID, id); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func CartOfUser(update tgbotapi.Update) ([]string, error) {

	db, err := sql.Open("postgres", DbInfo)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query(`SELECT products.product_name FROM cart
JOIN products ON cart.product_id = products.product_id
WHERE user_id = (select id FROM users
                           where chat_id = $1);`, update.CallbackQuery.Message.Chat.ID)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	var str []string
	for rows.Next() {
		var value string
		err = rows.Scan(&value)
		if err != nil {
			log.Fatal(err)
		}
		str = append(str, value)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return str, nil
}

func RemoveCart(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	db, err := sql.Open("postgres", DbInfo)
	if err != nil {
		log.Println(err)
		return
	}

	defer db.Close()

	data := `DELETE FROM cart
USING products,users
WHERE chat_id = $1`

	if _, err = db.Exec(data, update.CallbackQuery.Message.Chat.ID); err != nil {
		log.Println(err)
		return
	}

	cart, err := CartOfUser(update)
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
