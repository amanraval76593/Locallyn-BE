package user

import "time"

type User struct {
	Id           string
	Email        string
	PasswordHash string
	IsVerified   bool
	CreatedAt    time.Time
}
