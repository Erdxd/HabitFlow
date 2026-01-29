package handlers

import (
	"net/http"
	"z/admin"
	"z/model"
)

func (h *Handlers) CheckUsers(w http.ResponseWriter, r *http.Request) {
	Users, err := admin.GetDataAboutAllUsers(h.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	AllUsers := struct {
		Allusers []model.User
	}{
		Allusers: Users,
	}
	tmplAdminMain.Execute(w, AllUsers)

}
