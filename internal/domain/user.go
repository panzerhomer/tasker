package domain

import (
	"fmt"
	"regexp"
)

const (
	UserRole int8 = iota + 1
	EditorRole
	AdminRole
)

const emailRgxString = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

type User struct {
	ID       int64
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     int8
}

func (u *User) Validate() error {
	if len(u.Password) < 4 || len(u.Password) > 16 {
		return fmt.Errorf("password must be 4-16 symbols")
	}

	emailRegex := regexp.MustCompile(emailRgxString)
	if !emailRegex.MatchString(u.Email) {
		return fmt.Errorf("email is wrong")
	}

	return nil
}

func (u *User) Sanitize() {
	u.Password = ""
}
