package user

import (
	"context"
	"errors"
	"locallyn-be/pkg/database"
)

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(ctx context.Context, user *User) error {
	query := `
	INSERT INTO users(email,password_hash)
	VALUES ($1,$2)
	RETURNING id,created_at
	`

	err := database.DB.QueryRow(ctx, query,
		user.Email,
		user.PasswordHash,
	).Scan(&user.Id, &user.CreatedAt)

	return err
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT * 
		FROM users
		WHERE email=$1
	`
	var user User

	err := database.DB.QueryRow(ctx, query, email).Scan(
		&user.Id,
		&user.Email,
		&user.PasswordHash,
		&user.IsVerified,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, errors.New("User Not Found")
	}

	return &user, nil
}
