# Discord Status Bot

Бот для Discord, который показывает статус работы gorepostbot на сервере Linux.

## Возможности

- Проверка статуса сервиса gorepostbot через systemd
- Отображение времени работы сервиса
- Автоматический деплой через GitHub Actions

## Команды

- `!status` - показать текущий статус gorepostbot

## Настройка GitHub Actions

Для настройки автоматического деплоя необходимо добавить следующие секреты в репозиторий:

1. `DISCORD_BOT_TOKEN` - токен Discord бота
2. `USER` - имя пользователя на сервере, под которым будет запущен бот

## Локальная разработка

```bash
# Клонирование репозитория
git clone https://github.com/yourusername/godiscordbot.git
cd godiscordbot

# Установка зависимостей
go mod tidy

# Запуск с токеном
DISCORD_BOT_TOKEN=your_token_here go run main.go
```

## Установка на сервер через GitHub Actions

1. Настройте self-hosted runner для GitHub Actions на вашем сервере
2. Добавьте необходимые секреты в репозиторий GitHub
3. При пуше в ветку main, GitHub Actions автоматически выполнит сборку и деплой бота

## Ручная установка на сервер

1. Скопируйте исполняемый файл на сервер в `/opt/godiscordbot/`
2. Скопируйте файл `systemd/discord-status-bot.service` в `/etc/systemd/system/`
3. Отредактируйте файл сервиса, добавив токен бота
4. Выполните команды:

```bash
sudo systemctl daemon-reload
sudo systemctl enable discord-status-bot.service
sudo systemctl start discord-status-bot.service
```
