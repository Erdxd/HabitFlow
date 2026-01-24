package scheduler

import (
	"database/sql"
	"log"
	"time"
	"z/database"

	"github.com/go-co-op/gocron"
)

func ScheduleForReset(db *sql.DB) {
	anadyrLocation, err := time.LoadLocation("Asia/Anadyr")
	if err != nil {
		log.Fatal("Ошибка загрузки часового пояса Анадыря:", err)
	}

	s := gocron.NewScheduler(anadyrLocation)
	_, err = s.Every(1).Day().At("12:46").Do(func() {
		database.ResetAllStatus(db)

	})
	if err != nil {
		log.Fatal("Ошибка добавления задачи:", err)
	}

	s.StartAsync()
}
