package admin

import (
	"database/sql"
	"log"
	"z/model"
)

func GetDataAboutAllUsers(db *sql.DB) ([]model.User, error) {
	SqlStatement := (`SELECT * FROM users`)
	rows, err := db.Query(SqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []model.User
	for rows.Next() {
		var u model.User
		err := rows.Scan(&u.Id_user, &u.Username, &u.Email, &u.Password, &u.Telegram_chat_id, &u.Admin)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
		if len(users) == 0 {
			log.Println("Zero users")
			return nil, err
		}
	}
	return users, nil

}
func GetAdmin(db *sql.DB, id_user int) (bool, error) {
	var Admin bool
	SqlSatement := (`SELECT admin from users WHERE id_user = $1`)
	err := db.QueryRow(SqlSatement, id_user).Scan(&Admin)
	if err != nil {
		return false, err
	}
	return Admin, nil
}
