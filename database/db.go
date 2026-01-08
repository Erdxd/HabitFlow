package database

import (
	"database/sql"
	"log"
	"time"
	"z/model"
)

func CheckHabits(Id_user int) ([]model.HabitFlow, error) {
	rows, err := db.Query(`SELECT id, habit_name, status_today, streak FROM "HabitFlow" WHERE user_id = $1`, Id_user)

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
	SqlStatement := (`INSERT INTO "HabitFlow" (id, habit_name, status_today, streak)  VALUES ($1,$2 ,$3,$4, $5)`)
	_, err := db.Exec(SqlStatement, Habits.Id, Habits.Habit_Name, Habits.Status_Today, Habits.Streak, Id_user)
	if err != nil {
		return err
	}
	return nil
}
func DeleteHabits(db *sql.DB, id int) error {
	SqlStatement := (`DELETE FROM "HabitFlow" WHERE id = $1 `)
	_, err := db.Exec(SqlStatement, id)
	if err != nil {
		return err
	}
	return nil

}
func ChangeStatusToday(db *sql.DB, id int) error {
	SqlStatement := (`UPDATE "HabitFlow" SET status_today = true WHERE id = $1`)
	_, err := db.Exec(SqlStatement, id)
	if err != nil {
		return err
	}
	return nil
}
func ResetStatus(db *sql.DB, id int, reset chan model.HabitReset) {

	go func() error {
		time.Sleep(24 * time.Hour)
		SqlStatement := (`UPDATE "HabitFlow" SET status_today = false WHERE id = $1`)
		_, err := db.Exec(SqlStatement, id)
		if err != nil {
			return err
		}
		reset <- model.HabitReset{Id: id, Error: nil}

		return nil
	}()
}
func AddStreak(db *sql.DB, id int, streak int) error {
	SqlStatement := (`UPDATE "HabitFlow" SET streak = $1 WHERE id = $2`)
	_, err := db.Exec(SqlStatement, streak+1, id)
	if err != nil {
		return err
	}
	return nil
}
