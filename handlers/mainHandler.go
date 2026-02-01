package handlers

import (
	"log"
	"net/http"
	"strconv"
	"z/database"
	"z/model"
)

func (h *Handlers) Mainpage(w http.ResponseWriter, r *http.Request) {
	Id_user, err := h.Cookie_userID(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
	Habits, err := database.CheckHabits(h.DB, Id_user)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := struct {
		HabitAll []model.HabitFlow
	}{
		HabitAll: Habits,
	}
	tmplmain.Execute(w, data)
}
func (h *Handlers) AddHabits(w http.ResponseWriter, r *http.Request) {
	Id_user, err := h.Cookie_userID(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
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

		}
		Habit := model.HabitFlow{
			Id:           Id,
			Habit_Name:   Habit_Name,
			Status_Today: Status_Today,
			Streak:       Streak,
		}
		err = database.AddHabit(h.DB, Habit, Id_user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (h *Handlers) DeleteHabits(w http.ResponseWriter, r *http.Request) {
	Id_user, err := h.Cookie_userID(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
	if r.Method == "POST" {
		Id, err := strconv.Atoi(r.FormValue("Id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = database.DeleteHabits(h.DB, Id, Id_user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handlers) ResetStatus(w http.ResponseWriter, r *http.Request) {
	Id_user, err := h.Cookie_userID(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
	if r.Method == "POST" {
		Id, err := strconv.Atoi(r.FormValue("Id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err, streak := database.GetStreak(h.DB, Id_user, Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = database.ChangeStatusToday(h.DB, Id, Id_user, streak)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
