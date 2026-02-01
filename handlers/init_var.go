package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"z/database"
	"z/jwt"
)

var err error

func InitDBTohandlers() {
	db, err = database.InitDb()
	if err != nil {
		log.Fatal(err)
	}

}

var (
	tmplmain      = template.Must(template.ParseFiles("templates/main.html"))
	tmplreg       = template.Must(template.ParseFiles("templates/register.html"))
	tmpllog       = template.Must(template.ParseFiles("templates/login.html"))
	db            *sql.DB
	tmplrofile    = template.Must(template.ParseFiles("templates/profile.html"))
	tmplredact    = template.Must(template.ParseFiles("templates/redact.html"))
	tmplAdminMain = template.Must(template.ParseFiles("templates/admin.html"))
)

type Handlers struct {
	DB *sql.DB
}

func (h *Handlers) Cookie_userID(w http.ResponseWriter, r *http.Request) (int, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		log.Println("You are havent cookie")
		return 0, err
	}

	claims, err := jwt.ValidateToken(cookie.Value)
	if err != nil {
		log.Println("Your token is uncorrect")
		return 0, err
	}

	return claims.User_id, nil

}
