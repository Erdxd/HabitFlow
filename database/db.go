package database

import (
	"database/sql"
	"log"
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
	SqlStatement := (`INSERT INTO tasks (id, task, taskstatus, comment, user_id)  VALUES ($1,$2 ,$3,$4,$5) `)
	_, err := db.Exec(SqlStatement, Habits.Id, Habits.Habit_Name, Habits.Status_Today)
	if err != nil {
		return err
	}
	return nil
}
