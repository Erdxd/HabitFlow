package handlers

import (
	"database/sql"
	"html/template"
)

var (
	tmplmain   = template.Must(template.ParseFiles("templates/main.html"))
	tmplreg    = template.Must(template.ParseFiles("templates/register.html"))
	tmpllog    = template.Must(template.ParseFiles("templates/login.html"))
	db         *sql.DB
	tmplrofile = template.Must(template.ParseFiles("templates/profile.html"))
	tmplredact = template.Must(template.ParseFiles("templates/redact.html"))
)
