package bot

import (
	"database/sql"
	"log"
	"time"
	"z/database"
	"z/model"

	"github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Schedule(bot *tgbotapi.BotAPI, db *sql.DB) {
	anadyrLocation, err := time.LoadLocation("Asia/Anadyr")
	if err != nil {
		log.Fatal("Ошибка загрузки часового пояса Анадыря:", err)
	}

	s := gocron.NewScheduler(anadyrLocation)
	_, err = s.Every(1).Day().At("07:00").Do(func() {
		sendDailyNotifications(bot, db)
	})
	if err != nil {
		log.Fatal("Ошибка добавления задачи:", err)
	}

	s.StartAsync()
}

func sendDailyNotifications(bot *tgbotapi.BotAPI, db *sql.DB) {
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

		habits, err := database.CheckHabits(user.Id_user)
		if err != nil {
			log.Printf("Ошибка получения привычек для пользователя %d: %v", user.Id_user, err)
			continue
		}

		if err := SendHabitNotification(bot, user.Telegram_chat_id, habits); err != nil {
			log.Printf("Ошибка отправки уведомления пользователю %d: %v", user.Id_user, err)
		} else {
			log.Printf("Успешно отправлено уведомление пользователю %d", user.Id_user)
		}
	}
}
