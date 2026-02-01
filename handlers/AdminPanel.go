package handlers

import (
	"net/http"
	"strconv"

	"z/admin"
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
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}
