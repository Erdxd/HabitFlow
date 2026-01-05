package account

import (
	"database/sql"
	"z/model"
)

func Register(db *sql.DB, Users model.User) error {
	Sqlstatement := (`INSERT INTO "users" (username,email,password) VALUES ($1, $2, $3)`)
	_, err := db.Exec(Sqlstatement, Users.Username, Users.Email, Users.Password)
	if err != nil {
		return err
	}
	return nil

}
