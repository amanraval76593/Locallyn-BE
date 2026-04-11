package user

import (
	"context"
	"errors"
	"fmt"
	"locallyn-be/internal/common/auth"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repo       Repository
	cacheRepo  UserCache
	jwtSecret  string
	expiryTime int
}

func NewService(repo Repository, cacheRepo UserCache, jwtSecret string, expiryTime int) Service {
	return &service{repo: repo, cacheRepo: cacheRepo, jwtSecret: jwtSecret, expiryTime: expiryTime}
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

	err = s.repo.Create(ctx, user)

	if err != nil {
		return err
	}
	fmt.Printf("userId:%s\n", user.Id)
	return s.cacheRepo.setVerifyUserCode(user.Id, email)
}

func (s *service) VerifyUser(ctx context.Context, req VerifyUserRequest) error {

	code := req.Token

	data, err := s.cacheRepo.getVerifyUserData(code)

	if err != nil {
		return err
	}

	err = s.repo.VerifyUser(ctx, data.UserId)

	if err != nil {
		return err
	}
	log.Printf("User Data : %s", data)

	return nil

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

func (s *service) CreateProfile(ctx context.Context, claims *auth.Claims, req CreateProfileRequest) error {
	if claims == nil || claims.UserId == "" {
		return errors.New("invalid authenticated user")
	}

	err := s.repo.CreateUserProfile(ctx, claims.UserId, req.Username, req.DisplayName)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetProfile(ctx context.Context, claims *auth.Claims, req GetProfileRequest) (*UserProfile, error) {
	if claims == nil || claims.UserId == "" {
		return nil, errors.New("invalid authenticated user")
	}

	profile, err := s.repo.GetUserProfile(ctx, req.UserName)

	if err != nil {
		return nil, err
	}

	return profile, nil

}

func (s *service) UpdateProfile(ctx context.Context, claims *auth.Claims, req UpdateProfileRequest) (*UserProfile, error) {
	if claims == nil || claims.UserId == "" {
		return nil, errors.New("invalid authenticated user")
	}

	if req.Username == nil && req.DisplayName == nil && req.AvatarURL == nil {
		return nil, ErrNoProfileFieldsToEdit
	}

	if req.Username != nil {
		username := strings.TrimSpace(*req.Username)
		if len(username) < 4 {
			return nil, ErrInvalidProfileInput
		}
		req.Username = &username
	}

	if req.DisplayName != nil {
		displayName := strings.TrimSpace(*req.DisplayName)
		if len(displayName) < 4 {
			return nil, ErrInvalidProfileInput
		}
		req.DisplayName = &displayName
	}

	if req.AvatarURL != nil {
		avatarURL := strings.TrimSpace(*req.AvatarURL)
		req.AvatarURL = &avatarURL
	}

	profile, err := s.repo.UpdateUserProfile(ctx, claims.UserId, req)
	if err != nil {
		return nil, err
	}

	return profile, nil
}
