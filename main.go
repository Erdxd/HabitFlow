package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"z/account"
	"z/database"
	"z/model"
)

var db *sql.DB
var tmplmain = template.Must(template.ParseFiles("templates/main.html"))
var tmplreg = template.Must(template.ParseFiles("templates/register.html"))
var tmpllog = template.Must(template.ParseFiles("templates/login.html"))

func main() {
	var err error
	db, err = database.InitDb()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", Mainpage)
	http.HandleFunc("/add", AddHabits)
	http.HandleFunc("/delete", DeleteHabits)
	http.HandleFunc("/change", ResetStatus)
	http.HandleFunc("/register", RegisterPageHandler)
	http.HandleFunc("/login", LoginPageHandler)

	fmt.Println("localhost:8080")
	err = http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}
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

	Habits, err := database.CheckHabits(Id_user)
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
		Streak, err := strconv.Atoi(r.FormValue("Streak"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = database.ChangeStatusToday(db, Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		go func(Id int, Streak int) {
			reset := make(chan model.HabitReset)
			database.ResetStatus(db, Id, reset)
			result := <-reset
			if result.Error == nil {
				database.AddStreak(db, Id, Streak)
			}

		}(Id, Streak)

	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		NameUser := r.FormValue("username")
		Email := r.FormValue("email")
		Password := r.FormValue("password")
		PasswordRep := r.FormValue("password_repeat")
		if Password != PasswordRep {
			http.Error(w, "Пароли не совпадают", http.StatusInternalServerError)
			return
		} else if Password == PasswordRep {

			User := model.User{
				Username: NameUser,
				Email:    Email,
				Password: Password,
			}
			err := account.Register(db, User)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

	}

	tmplreg.Execute(w, nil)
	http.Redirect(w, r, "/login", http.StatusSeeOther)

}
func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		Login := r.FormValue("login")
		Password := r.FormValue("password")
		passwordfromdb, err := account.Login(db, Login)
		if err != nil {
			http.Error(w, "Пользователя не существует", http.StatusInternalServerError)
			return
		}
		if passwordfromdb != Password {
			http.Error(w, "Неверный пароль", http.StatusInternalServerError)
			return
		} else if passwordfromdb == Password {
			Id_user, err := account.GetUserId(db, Login)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name:  "Id_user",
				Value: strconv.Itoa(Id_user),
				Path:  "/",
			})

		}

	}
	tmpllog.Execute(w, nil)
	http.Redirect(w, r, "/", http.StatusSeeOther)

}
