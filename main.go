package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"
	"z/database"
	"z/model"
)

var db *sql.DB
var tmplmain = template.Must(template.ParseFiles("templates.main.html"))

func main() {
	var err error
	db, err = database.InitDb()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/main", Mainpage)
	http.HandleFunc("/add", AddHabits)
	err = http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
func Mainpage(w http.ResponseWriter, r *http.Request) {
	Habits, err := database.ChechHabits()
	if err != nil {
		http.Error(w, "Fail", 0)
		return
	}
	data := struct {
		TasksAll   []model.HabitFlow
		SearchTask []model.HabitFlow
	}{
		TasksAll:   Habits,
		SearchTask: nil,
	}
	tmplmain.Execute(w, data)
}
func AddHabits(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		Id, _ := strconv.Atoi(r.FormValue("Id"))

		Habit_Name := r.FormValue("Habiit_Name")
		Status_Today := r.FormValue("Status_Today") == "on"
		Habit := model.HabitFlow{
			Id:           Id,
			Habit_Name:   Habit_Name,
			Status_Today: Status_Today,
			Created_At:   time.Now(),
		}
		err := database.AddHabit(db, Habit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
	http.Redirect(w, r, "/main", http.StatusSeeOther)
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
	http.Redirect(w, r, "/main", http.StatusSeeOther)
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
	}
	http.Redirect(w, r, "/main", http.StatusSeeOther)

}
func ReesetStatus(w http.ResponseWriter, r *http.Request) {
	reset := make(chan model.HabitReset)
	if r.Method == "POST" {
		Id, err := strconv.Atoi(r.FormValue("Id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		go database.ChangeStatusToday(db, Id)
		result := <-reset
		if result.Error != nil {
			http.Error(w, "Произошла ошибка", http.StatusInternalServerError)
			return
		}

	}
	http.Redirect(w, r, "/main", http.StatusSeeOther)
}
