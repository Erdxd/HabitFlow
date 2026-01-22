package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"
	"z/account"
	"z/bot"
	"z/database"
	encrypto "z/hashing"
	"z/model"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
)

var db *sql.DB
var tmplmain = template.Must(template.ParseFiles("templates/main.html"))
var tmplreg = template.Must(template.ParseFiles("templates/register.html"))
var tmpllog = template.Must(template.ParseFiles("templates/login.html"))
var tmplrofile = template.Must(template.ParseFiles("templates/profile.html"))
var tmplredact = template.Must(template.ParseFiles("templates/redact.html"))

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

	http.HandleFunc("/", Mainpage)
	http.HandleFunc("/add", AddHabits)
	http.HandleFunc("/delete", DeleteHabits)
	http.HandleFunc("/change", ResetStatus)
	http.HandleFunc("/register", RegisterPageHandler)
	http.HandleFunc("/login", LoginPageHandler)
	http.HandleFunc("/profile", ProfileHandler)
	http.HandleFunc("/redact", RedactLoginHandler)

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
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	tmplreg.Execute(w, nil)
}
func LoginPageHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		Login := r.FormValue("login")
		Password := r.FormValue("password")
		passwordfromdb, err := account.GetPassword(db, Login)
		if err != nil {
			http.Error(w, "Неверный логин или пароль", http.StatusInternalServerError)
			return
		}
		Coincidence := encrypto.CHeckPassword(passwordfromdb, Password)
		if Coincidence {
			Id_user, err := account.GetUserId(db, Login)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name:  "id_user",
				Value: strconv.Itoa(Id_user),
				Path:  "/",
			})
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			http.Error(w, "Неверный логин или пароль", http.StatusInternalServerError)
			return
		}

	}

	tmpllog.Execute(w, nil)

}
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
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

	if r.Method == "GET" {
		log.Println("ProfileHandler called with Id_user:", Id_user)
		UserData, err := account.GetDataUser(db, Id_user)
		if err != nil {
			log.Println("Error getting user data:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("UserData length:", len(UserData))
		UserDataCheck := struct {
			DataAll []model.UserBaseView
		}{
			DataAll: UserData,
		}
		log.Println("Executing template with UserDataCheck")
		tmplrofile.Execute(w, UserDataCheck)
	}

}
func RedactLoginHandler(w http.ResponseWriter, r *http.Request) {
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

		newusername := r.FormValue("new_login")
		password := r.FormValue("password")

		passwordfromdb, err := account.GetPasswordwithId(db, Id_user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		Coincidence := encrypto.CHeckPassword(passwordfromdb, password)
		if Coincidence {
			err := account.RedactLogin(db, Id_user, newusername)
			if err != nil {
				http.Error(w, "Cant change your login", http.StatusBadRequest)
				return
			}

		} else {
			http.Error(w, "wrong password", http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/profile", http.StatusSeeOther)

	}
	tmplredact.Execute(w, nil)
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
