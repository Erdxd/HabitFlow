package scheduler

import (
	"log"
	"time"
	"z/database"

	"github.com/go-co-op/gocron"
)

func (sh *Scheduler) ScheduleForReset() {
	anadyrLocation, err := time.LoadLocation("Asia/Anadyr")
	if err != nil {
		log.Fatal("Ошибка загрузки часового пояса Анадыря:", err)
	}

	s := gocron.NewScheduler(anadyrLocation)
	_, err = s.Every(1).Day().At("17:05").Do(func() {
		database.ResetAllStatus(sh.DB)

	})
	if err != nil {
		log.Fatal("Ошибка добавления задачи:", err)
	}

	s.StartAsync()

}
