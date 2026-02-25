package repository

import (
	"context"
	"database/sql"

	"github.com/boostflashcards/backend/internal/models"
)

func (r *Repository) CreateUser(ctx context.Context, email, passwordHash, name, role string) (*models.User, error) {
	if role != "student" && role != "tutor" && role != "admin" {
		role = "student"
	}
	res, err := r.DB.ExecContext(ctx,
		`INSERT INTO users (email, password_hash, name, role, created_at) VALUES (?, ?, ?, ?, NOW())`,
		email, passwordHash, name, role)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.GetUserByID(ctx, id)
}

func (r *Repository) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	var u models.User
	err := r.DB.QueryRowContext(ctx,
		`SELECT id, email, password_hash, name, role, avatar_url, created_at FROM users WHERE id = ?`, id).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Role, &u.AvatarURL, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := r.DB.QueryRowContext(ctx,
		`SELECT id, email, password_hash, name, role, avatar_url, created_at FROM users WHERE email = ?`, email).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Role, &u.AvatarURL, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) CreateOAuthIdentity(ctx context.Context, userID int64, provider, providerUserID string) (*models.OAuthIdentity, error) {
	res, err := r.DB.ExecContext(ctx,
		`INSERT INTO oauth_identities (user_id, provider, provider_user_id, created_at) VALUES (?, ?, ?, NOW())`,
		userID, provider, providerUserID)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.GetOAuthIdentityByID(ctx, id)
}

func (r *Repository) GetOAuthIdentityByID(ctx context.Context, id int64) (*models.OAuthIdentity, error) {
	var o models.OAuthIdentity
	err := r.DB.QueryRowContext(ctx,
		`SELECT id, user_id, provider, provider_user_id, created_at FROM oauth_identities WHERE id = ?`, id).
		Scan(&o.ID, &o.UserID, &o.Provider, &o.ProviderUserID, &o.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *Repository) GetUserByOAuth(ctx context.Context, provider, providerUserID string) (*models.User, error) {
	var userID int64
	err := r.DB.QueryRowContext(ctx,
		`SELECT user_id FROM oauth_identities WHERE provider = ? AND provider_user_id = ?`,
		provider, providerUserID).Scan(&userID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return r.GetUserByID(ctx, userID)
}

func (r *Repository) UpdateUserAvatar(ctx context.Context, userID int64, avatarURL string) error {
	_, err := r.DB.ExecContext(ctx, `UPDATE users SET avatar_url = ? WHERE id = ?`, avatarURL, userID)
	return err
}
