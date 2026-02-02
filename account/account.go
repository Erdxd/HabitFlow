package account

import (
	"database/sql"
	"log"
	encrypto "z/hashing"
	"z/model"
)

func Register(db *sql.DB, Users model.User) error {
	Sqlstatement := (`INSERT INTO "users" (username,email,password) VALUES ($1, $2, $3)`)
	HashedPassword, err := encrypto.Encrypto(Users.Password)
	_, err = db.Exec(Sqlstatement, Users.Username, Users.Email, HashedPassword)
	if err != nil {
		return err
	}
	return nil

}
func GetPassword(db *sql.DB, Username string) (string, error) {
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
	err := db.QueryRow(Sqlstatement, Username).Scan(&Id_user)
	if err != nil {
		return 0, err
	}
	return Id_user, nil
}
func SaveTelegramChatID(db *sql.DB, userID int, chatID int64) error {
	Sqlstatement := (`UPDATE "users" SET telegram_chat_id = $1 WHERE id_user = $2`)
	_, err := db.Exec(Sqlstatement, chatID, userID)
	return err
}
func GetTelegramChatID(db *sql.DB, userID int, chatID int64) (int64, error) {
	var TgChat int
	Sqlstatement := (`SELECT telegram_chat_id FROM "users" WHERE id_user = $1`)
	err := db.QueryRow(Sqlstatement, userID).Scan(&TgChat)
	if err != nil {
		return 0, err
	}
	return int64(TgChat), nil

}
func GetUserIdByTgID(db *sql.DB, chatId int64) (int, error) {

	var User_Id int
	Sqlstatement := (`SELECT id_user FROM "users" WHERE telegram_chat_id = $1`)
	err := db.QueryRow(Sqlstatement, chatId).Scan(&User_Id)
	if err != nil {
		return 0, err
	}
	return User_Id, nil

}
func GetDataUser(db *sql.DB, id_user int) ([]model.UserBaseView, error) {
	log.Println("GetDataUser called with id_user:", id_user)
	rows, err := db.Query(`SELECT id_user,username,email,password FROM users WHERE id_user = $1`, id_user)

	if err != nil {
		log.Println("Can't SELECT data by your tables")
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	var BaseUser []model.UserBaseView
	for rows.Next() {
		var U model.UserBaseView
		err := rows.Scan(&U.Id_user, &U.Username, &U.Email, &U.Password)
		if err != nil {
			return nil, err
		}
		BaseUser = append(BaseUser, U)

	}

	return BaseUser, nil
}
func RedactLogin(db *sql.DB, user_id int, username string) error {
	SqlStatement := (`UPDATE "users" SET username = $1 WHERE id_user = $2`)
	_, err := db.Exec(SqlStatement, username, user_id)
	if err != nil {
		return err
	}
	return nil

}
func GetPasswordwithId(db *sql.DB, Id_user int) (string, error) {
	var Password string
	Sqlstatement := (`SELECT password FROM "users" WHERE id_user =$1`)
	err := db.QueryRow(Sqlstatement, Id_user).Scan(&Password)
	if err != nil {
		return "", err
	}

	return Password, nil

}
func RedactPassword(db *sql.DB, NewHashedPassword string, id_user int) error {
	SqlStatement := (`UPDATE users SET password = $1 WHERE id_user = $2`)
	_, err := db.Exec(SqlStatement, NewHashedPassword, id_user)
	if err != nil {

		return err
	}
	return nil
}
