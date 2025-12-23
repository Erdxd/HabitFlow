package database_test

import (
	"testing"
	"z/database"
	"z/model"
)

func TestFuncResetTime(t *testing.T) {
	reset := make(chan model.HabitReset)

	go database.ResetStatus(1, reset)
	result := <-reset
	if result.Error != nil {

		return
	}
}
