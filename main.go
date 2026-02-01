package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	bottg "z/BotTg"
	"z/database"
	"z/handlers"
	"z/scheduler"

	"z/middleware"

	"github.com/joho/godotenv"
)

var db *sql.DB

func main() {

	var err error
	godotenv.Load("Database.env", "bot.env", "jwt.env")

	db, err = database.InitDb()
	if err != nil {
		log.Fatal(err)
	}
	s := &scheduler.Scheduler{
		DB: db,
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	bot1, err := bottg.InitBot(token)
	if err != nil {
		log.Fatal(err)
	}
	go bottg.HandleMessages(bot1, db)
	go s.Schedule(bot1)
	go s.ScheduleForReset()

	h := &handlers.Handlers{
		DB: db,
	}

	http.HandleFunc("/", h.Mainpage)
	http.HandleFunc("/add", h.AddHabits)
	http.HandleFunc("/delete", h.DeleteHabits)
	http.HandleFunc("/change", h.ResetStatus)
	http.HandleFunc("/register", h.RegisterPageHandler)
	http.HandleFunc("/login", h.LoginPageHandler)
	http.HandleFunc("/profile", h.ProfileHandler)
	http.HandleFunc("/redact", h.RedactLoginHandler)
	http.HandleFunc("/admin/users", middleware.AdminOnly(h.CheckUsers))
	http.HandleFunc("/admin/update-password", middleware.AdminOnly(h.UpdatePassword))

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
