package domain

import (
	"fmt"
	"time"
	"unicode"
)

const (
	StatusInProcess = iota + 1
	StatusDone
)

type Task struct {
	ID          int64
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
}

func (t *Task) Validate() error {
	if len(t.Name) < 1 {
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

	// layout := "02/01/2006"

	// parsedDate, err := time.Parse(layout, t.Deadline.String())
	// if err != nil {
	// 	fmt.Println("Ошибка парсинга даты:", err)
	// 	return nil
	// }

	// fmt.Println("Распарсенная дата:", time.Date())

	return nil
}
