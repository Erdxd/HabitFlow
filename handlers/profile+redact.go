package handlers

import (
	"log"
	"net/http"
	"strconv"
	"z/account"
	encrypto "z/hashing"
	"z/model"
)

func (h *Handlers) ProfileHandler(w http.ResponseWriter, r *http.Request) {
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
		UserData, err := account.GetDataUser(h.DB, Id_user)
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
func (h *Handlers) RedactLoginHandler(w http.ResponseWriter, r *http.Request) {
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

		passwordfromdb, err := account.GetPasswordwithId(h.DB, Id_user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		Coincidence := encrypto.CheckPassword(passwordfromdb, password)
		if Coincidence {
			err := account.RedactLogin(h.DB, Id_user, newusername)
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
