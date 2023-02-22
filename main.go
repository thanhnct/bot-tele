package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/jasonlvhit/gocron"
	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func AppRecover() {
	if err := recover(); err != nil {
		cmd := exec.Command("myapp.exe")
		log.Println("restarting...")
		err := cmd.Run()

		if err != nil {
			log.Fatal(err)
		}
	}
}

func getIP() string {
	defer AppRecover()

	resp, err := http.Get("https://www.trackip.net/ip")
	if err != nil {
		log.Println(err)
		return ""
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return ""
	}

	return string(body)
}

func main() {
	fmt.Println("Bot starting....")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	chatID1, err := strconv.Atoi(os.Getenv("CHAT_ID_1"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	chatID2, err := strconv.Atoi(os.Getenv("CHAT_ID_2"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	botToken := ""
	bot, err := telego.NewBot(botToken)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	s := gocron.NewScheduler()
	s.Every(3).Hours().Do(func() {
		defer AppRecover()
		currentIp, err := os.ReadFile("current_ip.ip")
		if err != nil {
			log.Println(err)
		}

		if checkIp := getIP(); checkIp != string(currentIp) {
			f, err := os.Create("current_ip.ip")
			if err != nil {
				log.Println(err)
			}

			f.WriteString(checkIp)
			defer f.Close()

			_, _ = bot.SendMessage(tu.MessageWithEntities(tu.ID(int64((chatID1))),
				tu.Entity(fmt.Sprintf("Hi! Your IP has changed from %s to %s", currentIp, checkIp)).Bold(), tu.Entity("\n"),
			))

			_, _ = bot.SendMessage(tu.MessageWithEntities(tu.ID(int64((chatID2))),
				tu.Entity(fmt.Sprintf("Hi! Your IP has changed from %s to %s", currentIp, checkIp)).Bold(), tu.Entity("\n"),
			))
		}

	})

	<-s.Start()
}
