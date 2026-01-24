package scheduler

import (
	"log"
	"time"
	bottg "z/BotTg"
	"z/database"
	"z/model"

	"github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (sh *Scheduler) Schedule(bot *tgbotapi.BotAPI) {
	anadyrLocation, err := time.LoadLocation("Asia/Anadyr")
	if err != nil {
		log.Fatal("Ошибка загрузки часового пояса Анадыря:", err)
	}

	s := gocron.NewScheduler(anadyrLocation)
	_, err = s.Every(1).Day().At("07:00").Do(func() {
		sh.sendDailyNotifications(bot)
	})
	if err != nil {
		log.Fatal("Ошибка добавления задачи:", err)
	}

	s.StartAsync()
}

func (s *Scheduler) sendDailyNotifications(bot *tgbotapi.BotAPI) {
	rows, err := db.Query(`SELECT id_user, telegram_chat_id FROM "users" WHERE telegram_chat_id IS NOT NULL`)
	if err != nil {
		log.Printf("Ошибка получения пользователей: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.Id_user, &user.Telegram_chat_id)
		if err != nil {
			log.Printf("Ошибка сканирования пользователя: %v", err)
			continue
		}

		habits, err := database.CheckHabits(s.DB, user.Id_user)
		if err != nil {
			log.Printf("Ошибка получения привычек для пользователя %d: %v", user.Id_user, err)
			continue
		}

		if err := bottg.SendHabitNotification(bot, user.Telegram_chat_id, habits); err != nil {
			log.Printf("Ошибка отправки уведомления пользователю %d: %v", user.Id_user, err)
		} else {
			log.Printf("Успешно отправлено уведомление пользователю %d", user.Id_user)
			bottg.SendKeyboard(bot, user.Telegram_chat_id)
		}
	}
}
