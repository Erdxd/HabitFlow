package handlers

import (
	"log"
	"net/http"
	"strconv"
	"z/account"
	encrypto "z/hashing"
	"z/model"
)

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
				log.Println(err)
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
			log.Println(err)
			http.Error(w, "Неверный логин или пароль", http.StatusInternalServerError)
			return
		}
		Coincidence := encrypto.CHeckPassword(passwordfromdb, Password)
		if Coincidence {
			Id_user, err := account.GetUserId(db, Login)

			if err != nil {
				log.Println(err)
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
