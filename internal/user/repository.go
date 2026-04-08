package user

import (
	"context"
	"errors"
	"locallyn-be/pkg/database"

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

func (r *repository) CreateUserProfile(ctx context.Context, userId string, userName string, displayName string) (*UserProfile, error) {
	query := `
		INSERT INTO user_profiles(user_id, username, display_name)
		VALUES ($1, $2, $3)
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
		`

	var profile UserProfile

	err := database.DB.QueryRow(ctx, query, userId, userName, displayName).Scan(
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
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "user_profiles_username_key":
				return nil, ErrUsernameAlreadyExists
			case "user_profiles_pkey":
				return nil, ErrUserProfileExists
			}
		}

		return nil, err
	}

	return &profile, nil
}
