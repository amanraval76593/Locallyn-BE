package user

import (
	"context"
	"errors"
	"fmt"
	"locallyn-be/pkg/database"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func (r *repository) VerifyUser(ctx context.Context, userId string) error {
	query := `
			UPDATE users
			SET is_verified=$1
			WHERE id=$2
	`

	_, err := database.DB.Exec(ctx, query,
		true,
		userId,
	)

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

func (r *repository) CreateUserProfile(ctx context.Context, userId string, userName string, displayName string) error {
	query := `
		INSERT INTO user_profiles(user_id, username, display_name)
		VALUES ($1, $2, $3)
		`

	_, err := database.DB.Exec(ctx, query, userId, userName, displayName)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "user_profiles_username_key":
				return ErrUsernameAlreadyExists
			case "user_profiles_pkey":
				return ErrUserProfileExists
			}
		}

		return err
	}

	return nil
}

func (r *repository) GetUserProfile(ctx context.Context, userName string) (*UserProfile, error) {
	query := `
		SELECT
			user_id,
			username,
			display_name,
			avatar_url,
			trust_score,
			total_posts,
			total_confirmations,
			total_reports,
			created_at,
			updated_at
		FROM user_profiles
		WHERE username=$1
	`

	var profile UserProfile

	err := database.DB.QueryRow(ctx, query, userName).Scan(
		&profile.UserId,
		&profile.Username,
		&profile.DisplayName,
		&profile.AvatarURL,
		&profile.TrustScore,
		&profile.TotalPosts,
		&profile.TotalConfirmations,
		&profile.TotalReports,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err != nil {
		return nil, ErrUserProfileNotFound
	}

	return &profile, nil

}

func (r *repository) UpdateUserProfile(ctx context.Context, userId string, req UpdateProfileRequest) (*UserProfile, error) {
	setClauses := make([]string, 0, 4)
	args := make([]any, 0, 4)
	argPos := 1

	if req.Username != nil {
		setClauses = append(setClauses, fmt.Sprintf("username=$%d", argPos))
		args = append(args, *req.Username)
		argPos++
	}

	if req.DisplayName != nil {
		setClauses = append(setClauses, fmt.Sprintf("display_name=$%d", argPos))
		args = append(args, *req.DisplayName)
		argPos++
	}

	if req.AvatarURL != nil {
		setClauses = append(setClauses, fmt.Sprintf("avatar_url=$%d", argPos))
		args = append(args, *req.AvatarURL)
		argPos++
	}

	if len(setClauses) == 0 {
		return nil, ErrNoProfileFieldsToEdit
	}

	setClauses = append(setClauses, "updated_at=NOW()")
	args = append(args, userId)

	query := fmt.Sprintf(`
		UPDATE user_profiles
		SET %s
		WHERE user_id=$%d
		RETURNING
			user_id,
			username,
			display_name,
			avatar_url,
			trust_score,
			total_posts,
			total_confirmations,
			total_reports,
			created_at,
			updated_at
	`, strings.Join(setClauses, ", "), argPos)

	var profile UserProfile

	err := database.DB.QueryRow(ctx, query, args...).Scan(
		&profile.UserId,
		&profile.Username,
		&profile.DisplayName,
		&profile.AvatarURL,
		&profile.TrustScore,
		&profile.TotalPosts,
		&profile.TotalConfirmations,
		&profile.TotalReports,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "user_profiles_username_key" {
			return nil, ErrUsernameAlreadyExists
		}

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserProfileNotFound
		}

		return nil, err
	}

	return &profile, nil
}
