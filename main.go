package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	token string
)

func init() {
	token = os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		fmt.Println("Отсутствует токен Discord бота в переменных окружения")
		os.Exit(1)
	}
}

func main() {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Ошибка создания сессии Discord:", err)
		return
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("Ошибка открытия соединения:", err)
		return
	}

	fmt.Println("Бот запущен. Нажмите CTRL+C для выхода.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!status" {
		status := checkGorepostbotStatus()
		s.ChannelMessageSend(m.ChannelID, status)
	}
}

func checkGorepostbotStatus() string {
	cmd := exec.Command("systemctl", "is-active", "gorepostbot.service")
	output, err := cmd.Output()

	if err != nil {
		return fmt.Sprintf("Не удалось проверить статус gorepostbot: %v", err)
	}

	status := strings.TrimSpace(string(output))

	if status == "active" {
		uptime := getServiceUptime("gorepostbot.service")
		return fmt.Sprintf("✅ gorepostbot активен и работает\nВремя работы: %s", uptime)
	} else {
		return fmt.Sprintf("❌ gorepostbot не работает. Текущий статус: %s", status)
	}
}

func getServiceUptime(serviceName string) string {
	cmd := exec.Command("systemctl", "show", serviceName, "--property=ActiveEnterTimestamp")
	output, err := cmd.Output()

	if err != nil {
		return "неизвестно"
	}

	parts := strings.Split(string(output), "=")
	if len(parts) != 2 {
		return "неизвестно"
	}

	timestampStr := strings.TrimSpace(parts[1])
	timestamp, err := time.Parse("Mon 2006-01-02 15:04:05 MST", timestampStr)
	if err != nil {
		return timestampStr
	}

	duration := time.Since(timestamp)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60

	return fmt.Sprintf("%d часов, %d минут", hours, minutes)
}
