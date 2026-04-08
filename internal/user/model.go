package user

import "time"

type User struct {
	Id           string
	Email        string
	PasswordHash string
	IsVerified   bool
	CreatedAt    time.Time
}

type UserProfile struct {
	UserId             string
	Username           string
	DisplayName        string
	AvatarURL          *string
	TrustScore         float64
	TotalPosts         int
	TotalConfirmations int
	TotalReports       int
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
