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

// Функция для получения информации о сервере
func serverInfo(s *discordgo.Session, channelID string, guildID string) {
	guild, err := s.Guild(guildID)
	if err != nil {
		s.ChannelMessageSend(channelID, "Ошибка получения информации о сервере: "+err.Error())
		return
	}

	members, _ := s.GuildMembers(guildID, "", 1000)
	channels, _ := s.GuildChannels(guildID)

	var textChannels, voiceChannels, categoryChannels int
	for _, channel := range channels {
		switch channel.Type {
		case discordgo.ChannelTypeGuildText:
			textChannels++
		case discordgo.ChannelTypeGuildVoice:
			voiceChannels++
		case discordgo.ChannelTypeGuildCategory:
			categoryChannels++
		}
	}

	createdAt, _ := discordgo.SnowflakeTimestamp(guild.ID)

	info := fmt.Sprintf("**Информация о сервере**\n"+
		"Название: %s\n"+
		"ID: %s\n"+
		"Владелец: <@%s>\n"+
		"Создан: %s\n"+
		"Участников: %d\n"+
		"Каналов: %d (Текстовых: %d, Голосовых: %d, Категорий: %d)\n"+
		"Ролей: %d\n"+
		"Уровень проверки: %d\n",
		guild.Name, guild.ID, guild.OwnerID, createdAt.Format("02.01.2006 15:04:05"),
		len(members), len(channels), textChannels, voiceChannels, categoryChannels,
		len(guild.Roles), guild.VerificationLevel)

	s.ChannelMessageSend(channelID, info)
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

		// Обработка команды !serverinfo
		if m.Content == prefix+"serverinfo" {
			serverInfo(s, m.ChannelID, m.GuildID)
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
