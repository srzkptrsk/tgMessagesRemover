// It's a simple telegram bot for removing outdated messages from your groups.
package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	tgbotapi "gopkg.in/telegram-bot-api.v5"
	"log"
	"os"
	"strconv"
	"time"
)

// Database options in configuration.json
type Db struct {
	Dialect string
	Dsn     string
}

// Bot options in configuration.json
type Bot struct {
	Token   string
	Minutes string
}

// Entry point
func main() {
	conf := struct {
		Db  Db
		Bot Bot
	}{}

	file, _ := os.Open("configuration.json")
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&conf)
	if err != nil {
		log.Println("error:", err)
	}

	db, _ := gorm.Open(conf.Db.Dialect, conf.Db.Dsn)
	defer db.Close()
	db.AutoMigrate(&Message{})

	bot, err := tgbotapi.NewBotAPI(conf.Bot.Token)
	if err != nil {
		log.Panic(err)
	}

	// if command line argument is "cron"
	// get messages which older than (now - minutes in configuration)
	// and delete these messages one-by-one from telegram chat and database
	if len(os.Args) > 0 {
		for _, n := range os.Args[1:] {
			if n == "cron" {
				timeDuration, _ := strconv.ParseInt(conf.Bot.Minutes, 10, 32)
				timein := time.Now().Local().Add(-time.Minute * time.Duration(timeDuration))

				var messages []Message
				db.Where("created_at < ?", timein).Find(&messages)

				for _, message := range messages {
					fmt.Println(message.MessageId)
					_, err := bot.DeleteMessage(tgbotapi.NewDeleteMessage(message.ChatId, message.MessageId))
					if err != nil {
						fmt.Println("message was not deleted")
					}

					db.Unscoped().Delete(&message)
				}

				os.Exit(0)
			}
		}
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil && update.InlineQuery != nil {
			continue
		} else {
			// get all messages from your chats
			// and store chat id and message id (without any other details like content)
			// into your database
			if update.Message != nil {
				msg := Message{
					ChatId:    update.Message.Chat.ID,
					MessageId: update.Message.MessageID,
				}
				db.NewRecord(msg)
				db.Create(&msg)
			}
		}
	}
}
