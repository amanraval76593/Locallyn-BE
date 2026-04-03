package user

import "context"

type Repository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	VerifyUser(ctx context.Context, userId string) error
}

type Service interface {
	SignUp(ctx context.Context, req SignUpRequest) error
	VerifyUser(ctx context.Context, req VerifyUserRequest) error
	Login(ctx context.Context, req LoginRequest) (string, error)
}
