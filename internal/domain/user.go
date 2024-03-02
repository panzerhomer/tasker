package domain

type User struct {
	ID       int64
	Email    string `json:"email"`
	Password string `json:"password"`
}
