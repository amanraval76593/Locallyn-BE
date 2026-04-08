package user

import "errors"

var (
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrUserProfileExists     = errors.New("user profile already exists")
)
