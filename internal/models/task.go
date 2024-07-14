package models

import "time"

type Task struct {
	ID      int
	UserID  int
	Since   time.Time
	Until   time.Time
	Minutes int
}

func NewTask(userID int) *Task {
	return &Task{
		UserID: userID,
		Since:  time.Now(),
	}
}
