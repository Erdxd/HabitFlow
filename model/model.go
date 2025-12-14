package model

import "time"

type HabitFlow struct {
	Id           int       `db:"id"`
	Habit_Name   string    `db:"habit_name"`
	Status_Today bool      `db:"status_today"`
	User_Id      int       `db:"user_id"`
	Created_At   time.Time `db:"created_at"`
	Last_At      time.Time `db:"last_at"`
	Deleted_At   time.Time `db:"deleted_at"`
}
