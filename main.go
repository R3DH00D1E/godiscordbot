package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"godiscordbot/config"

	"github.com/bwmarrin/discordgo"
)

var (
	token  string
	prefix = "!"
)

func init() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Ошибка загрузки конфигурации: %v", err)
	}
	var token = cfg.DSToken
	if token == "" {
		fmt.Println("Отсутствует токен Discord бота в переменных окружения")
		os.Exit(1)
	}
}

func main() {
	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Ошибка создания сессии Discord:", err)
		return
	}

	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		if m.Content == "привет" {
			s.ChannelMessageSend(m.ChannelID, "мир!")
		}
	})

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = sess.Open()
	if err != nil {
		fmt.Println("Ошибка открытия соединения:", err)
		return
	}
	defer sess.Close()

	fmt.Println("Админ-бот запущен. Нажмите CTRL+C для выхода.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
