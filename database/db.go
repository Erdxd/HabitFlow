package database

import (
	"database/sql"
	"log"
	"z/model"
)

func CheckHabits(db *sql.DB, Id_user int) ([]model.HabitFlow, error) {
	rows, err := db.Query(`SELECT id, habit_name, status_today, streak FROM "HabitFlow" WHERE user_id= $1`, Id_user)

	if err != nil {
		log.Println("Can't SELECT data by your tables")
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	var tasks []model.HabitFlow
	for rows.Next() {
		var t model.HabitFlow
		err := rows.Scan(&t.Id, &t.Habit_Name, &t.Status_Today, &t.Streak)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil

}
func AddHabit(db *sql.DB, Habits model.HabitFlow, Id_user int) error {
	SqlStatement := (`INSERT INTO "HabitFlow" (id, habit_name, status_today, streak, user_id) VALUES ($1,$2 ,$3,$4, $5)`)
	_, err := db.Exec(SqlStatement, Habits.Id, Habits.Habit_Name, Habits.Status_Today, Habits.Streak, Id_user)
	if err != nil {
		return err
	}
	return nil
}
func DeleteHabits(db *sql.DB, id int, user_id int) error {
	SqlStatement := (`DELETE FROM "HabitFlow" WHERE id = $1 AND user_id = $2 `)
	_, err := db.Exec(SqlStatement, id, user_id)
	if err != nil {
		return err
	}
	return nil

}
func ChangeStatusToday(db *sql.DB, id int, user_id, streak int) error {
	SqlStatement := (`UPDATE "HabitFlow" SET status_today = true, streak = $1 WHERE id = $2 AND user_id = $3`)
	_, err := db.Exec(SqlStatement, streak+1, id, user_id)
	if err != nil {
		return err
	}
	return nil
}

func GetStatus(db *sql.DB, user_id, id int) (error, bool) {
	var Status_Today bool
	SqlStatement := (`SELECT status_today FROM "HabitFlow" WHERE id = $1 AND user_id = $2`)
	err := db.QueryRow(SqlStatement, id, user_id).Scan(&Status_Today)
	if err != nil {
		return err, false
	}
	return nil, Status_Today

}
func GetStreak(db *sql.DB, user_id, id int) (error, int) {
	var streak int
	SqlStatement := (`SELECT streak FROM "HabitFlow" WHERE id = $1 AND user_id = $2`)
	err := db.QueryRow(SqlStatement, id, user_id).Scan(&streak)
	if err != nil {
		return err, 0
	}
	return nil, streak

}
func ResetAllStatus(db *sql.DB) error {
	SqlStatement := (`UPDATE "HabitFlow" SET status_today = false WHERE status_today = true`)
	_, err := db.Exec(SqlStatement)
	if err != nil {
		return err
	}
	return nil
}
