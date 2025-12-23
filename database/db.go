package database

import (
	"database/sql"
	"log"
	"time"
	"z/model"
)

func ChechHabits() ([]model.HabitFlow, error) {
	rows, err := db.Query(`SELECT id, habit_name, status_today FROM "HabitFlow WHERE user_id = $1`)

	if err != nil {
		log.Println("Can't SELECT data by your tables")
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	var tasks []model.HabitFlow
	for rows.Next() {
		var t model.HabitFlow
		err := rows.Scan(&t.Id, &t.Habit_Name, &t.Status_Today, &t.Created_At, &t.Last_At, &t.Deleted_At)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil

}
func AddHabit(db *sql.DB, Habits model.HabitFlow) error {
	SqlStatement := (`INSERT INTO "HabitFlow" (id, habit_name, status_today, created_at)  VALUES ($1,$2 ,$3,$4) `)
	_, err := db.Exec(SqlStatement, Habits.Id, Habits.Habit_Name, Habits.Status_Today)
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
	SqlStatement := (`UPDATE "HabitFlow SET status_today = true WHERE id = $1"`)
	_, err := db.Exec(SqlStatement, id)
	if err != nil {
		return err
	}
	return nil
}
func ResetStatus(id int, reset chan model.HabitReset) {

	go func() error {
		time.Sleep(24 * time.Second)
		SqlStatement := (`UPDATE "HabitFlow" SET status_today = false WHERE id = $1"`)
		_, err := db.Exec(SqlStatement, id)
		if err != nil {
			return err
		}
		reset <- model.HabitReset{Id: id, Error: nil}

		return nil
	}()
}
