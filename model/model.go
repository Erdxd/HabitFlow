package model

import "time"

type HabitFlow struct {
	Id           int       `db:"id"`
	Habit_Name   string    `db:"habit_name"`
	Status_Today bool      `db:"status_today"`
	Streak       int       `db:"streak"`
	User_Id      int       `db:"user_id"`
	Created_At   time.Time `db:"created_at"`
	Last_At      time.Time `db:"last_at"`
	Deleted_At   time.Time `db:"deleted_at"`
}
type HabitReset struct {
	Id    int
	Error error
}
type User struct {
	Id_user          int    `db:"id_user"`
	Username         string `db:"username"`
	Email            string `db:"email"`
	Password         string `db:"passworduser"`
	Telegram_chat_id int64  `db:"telegram_chat_id"`
}
type UserBaseView struct {
	Id_user  int    `db:"id_user"`
	Username string `db:"username"`
	Email    string `db:"email"`
	Password string `db:"passworduser"`
}
