package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	TGBOT_TOKEN := os.Getenv("TGBOT_TOKEN")
	fmt.Println(TGBOT_TOKEN)
	bot, e := tgbotapi.NewBotAPI(TGBOT_TOKEN)
	if e != nil {
		log.Printf("Auth error: %s", e)
		panic(e)
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			//mt.Println("recived")
			if update.Message.Photo != nil {
				//fmt.Println("it`s photo")
				photos := update.Message.Photo
				photo_url, _ := bot.GetFileDirectURL(photos[len(photos)-1].FileID)
				pu, _ := url.Parse(photo_url)
				fn := path.Base(pu.Path)
				response, e := http.Get(photo_url)
				if e != nil {
					panic(e)
				}
				fp := path.Join("/tmp", fn)
				file, e := os.Create(fp)

				if e != nil {
					panic(e)
				}
				io.Copy(file, response.Body)

				response.Body.Close()
				//cs := fmt.Sprintf("/bin/tesseract %s stdout -l rus+rus", fn)
				cmd := exec.Command("tesseract", fp, "stdout", "-l", "rus+eng")
				b, e := cmd.CombinedOutput()
				if e != nil {
					fmt.Println(e)
				}
				file.Close()

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, string(b))
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)

			}
		}
	}
}
