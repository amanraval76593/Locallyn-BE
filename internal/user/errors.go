package user

import "errors"

var (
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrUserProfileExists     = errors.New("user profile already exists")
	ErrUserProfileNotFound   = errors.New("User profile does not exist")
	ErrNoProfileFieldsToEdit = errors.New("no profile fields provided for update")
	ErrInvalidProfileInput   = errors.New("invalid profile input")
)
