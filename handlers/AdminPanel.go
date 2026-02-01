package handlers

import (
	"log"
	"net/http"
	"strconv"

	"z/admin"
	"z/database"
	encrypto "z/hashing"
	"z/model"
)

func (h *Handlers) CheckUsers(w http.ResponseWriter, r *http.Request) {
	Users, err := admin.GetDataAboutAllUsers(h.DB)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Allusers := struct {
		Allusers []model.User
	}{
		Allusers: Users,
	}
	tmplAdminMain.Execute(w, Allusers)

}
func (h *Handlers) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	NewPassword := r.FormValue("password")
	NewPasswordhashed, err := encrypto.Encrypto(NewPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id_user, err := strconv.Atoi(r.FormValue("Id_user"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = admin.ChangePasswordForUser(h.DB, id_user, NewPasswordhashed)
	log.Println(NewPasswordhashed)
	log.Println(NewPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}
func (h *Handlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id_user, err := strconv.Atoi(r.FormValue("Id_user"))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = database.DeleteAllHabit(h.DB, id_user)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not delete this account", http.StatusInternalServerError)
		return
	}
	err = admin.DeleteAccount(h.DB, id_user)
	if err != nil {
		http.Error(w, "Could not delete this account", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}
