package bot

import (
	"database/sql"
	"fmt"
	"strconv"
	"z/account"
	"z/model"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func InitBot(token string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return bot, nil
}
func SendHabitNotification(bot *tgbotapi.BotAPI, chatId int64, habits []model.HabitFlow) error {
	if len(habits) == 0 {

		message := "🎉 У тебя пока нет активных привычек! Добавь первую привычку на сайте."
		msg := tgbotapi.NewMessage(chatId, message)
		_, err := bot.Send(msg)
		return err

	}
	message := "📋 Твои привычки на сегодня:\n\n"
	for _, habit := range habits {
		status := "❌"
		if habit.Status_Today {
			bot.Send(tgbotapi.NewMessage(chatId, message))
			status = "✅"
		}

		message += fmt.Sprintf("%s %s (серия: %d)\n", status, habit.Habit_Name, habit.Streak)
	}
	msg := tgbotapi.NewMessage(chatId, message)
	_, err := bot.Send(msg)
	return err
}
func HandleMessages(bot *tgbotapi.BotAPI, db *sql.DB) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			if update.Message.Command() == "start" {
				args := update.Message.CommandArguments()
				userID, err := strconv.Atoi(args)
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неверный формат. Используй: /start ТВОЙ_ID_ПОЛЬЗОВАТЕЛЯ")
					bot.Send(msg)
					continue
				}

				err = account.SaveTelegramChatID(db, userID, update.Message.Chat.ID)
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка привязки чата")
					bot.Send(msg)
					continue
				}

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "✅ Чат успешно привязан! Теперь ты будешь получать уведомления о привычках.")
				bot.Send(msg)
			}
		}
	}
}
