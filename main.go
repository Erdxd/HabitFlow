package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"z/database"
	"z/model"
)

var db *sql.DB
var tmplmain = template.Must(template.ParseFiles("templates/main.html"))

func main() {
	var err error
	db, err = database.InitDb()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", Mainpage)
	http.HandleFunc("/add", AddHabits)
	http.HandleFunc("/delete", DeleteHabits)
	http.HandleFunc("/change", ChangeStatusToday)
	http.HandleFunc("/resetstatus", ResetStatus)

	fmt.Println("localhost:8080")
	err = http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}
func Mainpage(w http.ResponseWriter, r *http.Request) {
	Habits, err := database.CheckHabits()
	if err != nil {
		http.Error(w, "Fail", http.StatusSeeOther)
		return
	}
	data := struct {
		HabitAll   []model.HabitFlow
		SearchTask []model.HabitFlow
	}{
		HabitAll: Habits,

		SearchTask: nil,
	}
	tmplmain.Execute(w, data)
}
func AddHabits(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		Id, err := strconv.Atoi(r.FormValue("Id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		Habit_Name := r.FormValue("Habit_Name")
		Status_Today := r.FormValue("Status_Today") == "on"
		Streak, err := strconv.Atoi(r.FormValue("Streak"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		Habit := model.HabitFlow{
			Id:           Id,
			Habit_Name:   Habit_Name,
			Status_Today: Status_Today,
			Streak:       Streak,
		}
		err = database.AddHabit(db, Habit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func DeleteHabits(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		Id, err := strconv.Atoi(r.FormValue("Id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = database.DeleteHabits(db, Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func ChangeStatusToday(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		Id, err := strconv.Atoi(r.FormValue("Id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = database.ChangeStatusToday(db, Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		streak, err := database.SelectStreak(db, Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		database.AddStreak(db, Id, streak)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)

}
func ResetStatus(w http.ResponseWriter, r *http.Request) {
	reset := make(chan model.HabitReset)
	if r.Method == "POST" {
		Id, err := strconv.Atoi(r.FormValue("Id"))
		Streak, err := strconv.Atoi((r.FormValue("Streak")))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		go database.ChangeStatusToday(db, Id)
		result := <-reset
		if result.Error != nil {
			http.Error(w, "Произошла ошибка", http.StatusInternalServerError)
			return
		} else {
			database.AddStreak(db, Id, Streak)
		}

	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
