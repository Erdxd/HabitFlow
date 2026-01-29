package handlers

import (
	"log"
	"net/http"
	"z/account"
	"z/admin"
	encrypto "z/hashing"
	"z/jwt"
	"z/model"
)

func (h *Handlers) RegisterPageHandler(w http.ResponseWriter, r *http.Request) {

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
			err := account.Register(h.DB, User)
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
func (h *Handlers) LoginPageHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		Login := r.FormValue("login")
		Password := r.FormValue("password")
		passwordfromdb, err := account.GetPassword(h.DB, Login)
		if err != nil {
			log.Println(err)
			http.Error(w, "Неверный логин или пароль", http.StatusInternalServerError)
			return
		}
		Coincidence := encrypto.CheckPassword(passwordfromdb, Password)
		if Coincidence {
			Id_user, err := account.GetUserId(h.DB, Login)

			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			Admin, err := admin.GetAdmin(h.DB, Id_user)
			if err != nil {
				log.Println(err)
				http.Error(w, "Что-то пошло не так", http.StatusInternalServerError)
				return
			}
			tokenstring, err := jwt.GenerateToken(Id_user, Admin)
			if err != nil {
				log.Println(err)
				http.Error(w, "Что-то пошло не так", http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "token",
				Value:    tokenstring,
				Path:     "/",
				MaxAge:   86400,
				HttpOnly: true,
				Secure:   true,
			})
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			http.Error(w, "Неверный логин или пароль", http.StatusInternalServerError)
			return
		}

	}

	tmpllog.Execute(w, nil)

}
