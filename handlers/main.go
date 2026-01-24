package handlers

import (
	"net/http"
	"strconv"
	"z/database"
	"z/model"
)

func Mainpage(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("id_user")
	if err != nil {
		http.Error(w, "Не смогли извлечь куки", http.StatusBadRequest)
		return
	}
	Id_user, err := strconv.Atoi(cookie.Value)
	if err != nil {
		http.Error(w, "Не смогли извлечь куки", http.StatusBadRequest)
		return
	}

	if Id_user == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	Habits, err := database.CheckHabits(db, Id_user)
	if err != nil {
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
func AddHabits(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("id_user")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	Id_user, err := strconv.Atoi(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if Id_user == 0 {
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
		err = database.AddHabit(db, Habit, Id_user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func DeleteHabits(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("id_user")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	Id_user, err := strconv.Atoi(cookie.Value)
	if err != nil {
		http.Error(w, "Не смогли извлечь куки", http.StatusBadRequest)
		return
	}

	if Id_user == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	if r.Method == "POST" {
		Id, err := strconv.Atoi(r.FormValue("Id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = database.DeleteHabits(db, Id, Id_user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ResetStatus(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("id_user")
	if err != nil {
		http.Error(w, "Не смогли извлечь куки", http.StatusBadRequest)
		return
	}
	Id_user, err := strconv.Atoi(cookie.Value)
	if err != nil {
		http.Error(w, "Не смогли извлечь куки", http.StatusBadRequest)
		return
	}

	if Id_user == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	if r.Method == "POST" {
		Id, err := strconv.Atoi(r.FormValue("Id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err, streak := database.GetStreak(db, Id_user, Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = database.ChangeStatusToday(db, Id, Id_user, streak)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err, status := database.GetStatus(db, Id_user, Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if status == true {

		}

	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
