package database

import (
	"log"
	"z/model"
)

func ChechHabits() ([]model.HabitFlow, error) {
	rows, err := db.Query(`SELECT id, user_id, task, taskstatus, comment FROM tasks WHERE user_id = $1`)

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
func AddHabit() {

}
