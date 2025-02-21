package main

import (
	"log"
	"os"
	"time"

	"gopkg.in/telebot.v3"
)

func main() {
	// Вставь сюда свой токен
	token := os.Getenv("TOKEN")
	
	// Настройки бота
	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	// Создаём бота
	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	commands := []telebot.Command{
		{Text: "help", Description: "Показать справку"},
		{Text: "kaban", Description: "Замутить пользователя на 1 час"},
	}

	// Устанавливаем команды
	err = bot.SetCommands(commands)
	if err != nil {
		log.Fatal("Не удалось зарегистрировать команды:", err)
	}

	var mutedUsers = make(map[int64]time.Time) // Хранение замученных пользователей

	bot.Handle("/help", func(c telebot.Context) error {
		user := c.Sender()

		// Проверяем, замучен ли пользователь
		if muteTime, ok := mutedUsers[user.ID]; ok {
			if time.Now().Before(muteTime) {
				// Удаляем команду, если пользователь замучен
				err := c.Delete()
				if err != nil {
					return c.Send("Не удалось удалить команду: " + err.Error())
				}
				return nil
			}
			delete(mutedUsers, user.ID) // Убираем из списка, если мут истек
		}

		return c.Send("Привет! Бот обладает всего одной командой - напиши в чат /kaban и узнай, как он работает:)")
	})

	bot.Handle("/kaban", func(c telebot.Context) error {
		user := c.Sender()

		// Проверяем, замучен ли пользователь
		if muteTime, ok := mutedUsers[user.ID]; ok {
			if time.Now().Before(muteTime) {
				// Удаляем команду, если пользователь замучен
				err := c.Delete()
				if err != nil {
					return c.Send("Не удалось удалить команду: " + err.Error())
				}
				return nil
			}
			delete(mutedUsers, user.ID) // Убираем из списка, если мут истек
		}

		// Проверяем, что команда вызвана в группе
		if c.Chat().Type != telebot.ChatGroup && c.Chat().Type != telebot.ChatSuperGroup {
			return c.Send("Эта команда работает только в группах!")
		}

		// Время окончания мута (через 1 час)
		mutedUsers[user.ID] = time.Now().Add(1 * time.Hour)

		// Отправляем сообщение о муте
		return c.Send(user.FirstName + " теперь в муте на 1 час! 🚫")
	})

	// Обработка всех сообщений
	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		user := c.Sender()

		// Проверяем, замучен ли пользователь
		if muteTime, ok := mutedUsers[user.ID]; ok {
			if time.Now().After(muteTime) {
				delete(mutedUsers, user.ID) // Убираем из списка, если мут истек
				return nil
			}

			// Если сообщение — это команда (начинается с "/"), удаляем его
			if c.Message().Text[0] == '/' {
				err := c.Delete()
				if err != nil {
					return c.Send("Не удалось удалить команду: " + err.Error())
				}
			}

			err := c.Delete()
			if err != nil {
				return c.Send("Не удалось удалить сообщение: " + err.Error())
			}
		}

		return nil
	})

	bot.Handle(telebot.OnSticker, func(c telebot.Context) error {
		user := c.Sender()

		// Проверяем, замучен ли пользователь
		if muteTime, ok := mutedUsers[user.ID]; ok {
			// Если время мута истекло, удаляем пользователя из списка
			if time.Now().After(muteTime) {
				delete(mutedUsers, user.ID)
				return nil
			}

			// Удаляем стикер
			err := c.Delete()
			if err != nil {
				return c.Send("Не удалось удалить стикер: " + err.Error())
			}
		}

		return nil
	})

	// Обработка голосовых сообщений
	bot.Handle(telebot.OnVoice, func(c telebot.Context) error {
		user := c.Sender()

		// Проверяем, замучен ли пользователь
		if muteTime, ok := mutedUsers[user.ID]; ok {
			// Если время мута истекло, удаляем пользователя из списка
			if time.Now().After(muteTime) {
				delete(mutedUsers, user.ID)
				return nil
			}

			// Удаляем голосовое сообщение
			err := c.Delete()
			if err != nil {
				return c.Send("Не удалось удалить голосовое сообщение: " + err.Error())
			}
		}

		return nil
	})

	bot.Handle(telebot.OnMedia, func(c telebot.Context) error {
		user := c.Sender()

		// Проверяем, замучен ли пользователь
		if muteTime, ok := mutedUsers[user.ID]; ok {
			// Если время мута истекло, удаляем пользователя из списка
			if time.Now().After(muteTime) {
				delete(mutedUsers, user.ID)
				return nil
			}

			// Удаляем медиа
			err := c.Delete()
			if err != nil {
				return c.Send("Не удалось удалить медиа: " + err.Error())
			}
		}

		return nil
	})

	bot.Handle(telebot.OnVideo, func(c telebot.Context) error {
		user := c.Sender()

		// Проверяем, замучен ли пользователь
		if muteTime, ok := mutedUsers[user.ID]; ok {
			// Если время мута истекло, удаляем пользователя из списка
			if time.Now().After(muteTime) {
				delete(mutedUsers, user.ID)
				return nil
			}

			// Удаляем медиа
			err := c.Delete()
			if err != nil {
				return c.Send("Не удалось удалить медиа: " + err.Error())
			}
		}

		return nil
	})

	// Запускаем бота
	bot.Start()
}