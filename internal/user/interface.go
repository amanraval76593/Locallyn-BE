package user

import (
	"context"
	"locallyn-be/internal/common/auth"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	VerifyUser(ctx context.Context, userId string) error
	CreateUserProfile(ctx context.Context, userId string, userName string, displayName string) (*UserProfile, error)
}

type Service interface {
	SignUp(ctx context.Context, req SignUpRequest) error
	VerifyUser(ctx context.Context, req VerifyUserRequest) error
	Login(ctx context.Context, req LoginRequest) (string, error)
	CreateProfile(ctx context.Context, claims *auth.Claims, req CreateProfileRequest) (*UserProfile, error)
}
