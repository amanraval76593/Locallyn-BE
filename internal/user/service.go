package user

import (
	"context"
	"errors"
	"locallyn-be/internal/common/auth"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repo       Repository
	jwtSecret  string
	expiryTime int
}

func NewService(repo Repository, jwtSecret string, expiryTime int) Service {
	return &service{repo: repo, jwtSecret: jwtSecret, expiryTime: expiryTime}
}

func (s *service) SignUp(ctx context.Context, req SignUpRequest) error {
	email := strings.ToLower(req.Email)

	_, err := s.repo.FindByEmail(ctx, email)

	if err == nil {
		return errors.New("user already exist")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user := &User{
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	return s.repo.Create(ctx, user)
}

func (s *service) Login(ctx context.Context, req LoginRequest) (string, error) {

	email := strings.ToLower(req.Email)

	user, err := s.repo.FindByEmail(ctx, email)

	if err != nil {
		return "", errors.New("user does not exist")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))

	if err != nil {
		return "", errors.New("Incorrect Password")
	}

	if !user.IsVerified {
		return "", errors.New("User is not Verified")
	}

	token, err := auth.GenerateAccessToken(user.Id, user.Email, s.jwtSecret, s.expiryTime)

	if err != nil {
		return "", err
	}

	return token, nil

}
