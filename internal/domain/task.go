package domain

import (
	"fmt"
	"time"
	"unicode"
)

const (
	ToDo = iota + 1
	InProgress
	InReview
	Done
	Blocked
)

type Task struct {
	ID             int64
	Name           string `json:"name"`
	Description    string `json:"description"`
	Deadline       string `json:"deadline"`
	AssignedUserID int64  `json:"assigned_user_id"`
	AuthorID       int64  `json:"author_id"`
	Status         int8   `json:"status"`
}

func (t *Task) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("task name is too short")
	}

	checkName := func(s string) bool {
		for _, r := range s {
			if !unicode.IsLetter(r) {
				return false
			}
		}
		return true
	}

	if checkName(t.Name) {
		return fmt.Errorf("name must contain letters only")
	}

	layout := "1/2/2006" // Specify the layout of the date string

	parsedTime, err := time.Parse(layout, t.Deadline)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}

	t.Status = ToDo

	return nil
}

func parseDateFrom(date string) time.Time
