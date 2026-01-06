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
func Login(db *sql.DB, Username string) (string, error) {
	var Password string
	Sqlstatement := (`SELECT password FROM "users" WHERE username =$1`)
	err := db.QueryRow(Sqlstatement, Username).Scan(&Password)
	if err != nil {
		return "", err
	}

	return Password, nil
}
func GetUserId(db *sql.DB, Username string) (int, error) {
	var Id_user int
	Sqlstatement := (`SELECT id_user FROM "users" WHERE username = $1`)
	err := db.QueryRow(Sqlstatement).Scan(&Id_user)
	if err != nil {
		return 0, err
	}
	return Id_user, nil
}
