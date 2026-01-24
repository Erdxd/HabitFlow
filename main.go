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
	bot1, err := bottg.InitBot(token)
	if err != nil {
		log.Fatal(err)
	}
	go bottg.HandleMessages(bot1, db)
	go scheduler.Schedule(bot1, db)
	go scheduler.ScheduleForReset(db)

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
