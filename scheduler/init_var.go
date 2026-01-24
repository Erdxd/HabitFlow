package scheduler

import "database/sql"

var db *sql.DB

type Scheduler struct {
	DB *sql.DB
}
