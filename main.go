package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"time"
	"z/bot"
	"z/database"
	"z/handlers"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
)

var db *sql.DB

func main() {
	var err error
	godotenv.Load("Database.env", "bot.env")

	db, err = database.InitDb()
	if err != nil {
		log.Fatal(err)
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	bot1, err := bot.InitBot(token)
	if err != nil {
		log.Fatal(err)
	}
	go bot.HandleMessages(bot1, db)
	go bot.Schedule(bot1, db)
	go ScheduleForReset(db)

	http.HandleFunc("/", handlers.Mainpage)
	http.HandleFunc("/add", handlers.AddHabits)
	http.HandleFunc("/delete", handlers.DeleteHabits)
	http.HandleFunc("/change", handlers.ResetStatus)
	http.HandleFunc("/register", handlers.RegisterPageHandler)
	http.HandleFunc("/login", handlers.LoginPageHandler)
	http.HandleFunc("/profile", handlers.ProfileHandler)
	http.HandleFunc("/redact", handlers.RedactLoginHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Server is running on port:", port)
	err = http.ListenAndServe("0.0.0.0:"+port, nil)
	if err != nil {
		log.Fatal(err)
	}

}

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
