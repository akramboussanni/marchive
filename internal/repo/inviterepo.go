package repo

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/akramboussanni/marchive/internal/model"
	"github.com/akramboussanni/marchive/internal/utils"
	"github.com/jmoiron/sqlx"
)

type InviteRepo struct {
	db *sqlx.DB
}

func NewInviteRepo(db *sqlx.DB) *InviteRepo {
	return &InviteRepo{db: db}
}

// CreateInvite creates a new invite for a user
func (r *InviteRepo) CreateInvite(ctx context.Context, inviterID int64) (*model.Invite, error) {
	// Check if user has invite tokens
	var tokens int
	err := r.db.GetContext(ctx, &tokens, "SELECT invite_tokens FROM users WHERE id = ?", inviterID)
	if err != nil {
		return nil, err
	}

	if tokens <= 0 {
		return nil, ErrNoInviteTokens
	}

	// Generate unique token
	token, err := r.generateUniqueToken(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()

	// Create invite
	invite := &model.Invite{
		Token:     token,
		InviterID: inviterID,
		CreatedAt: now,
	}

	// Insert invite
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO invites (id, token, inviter_id, created_at) 
		VALUES (?, ?, ?, ?)
	`, utils.GenerateSnowflakeID(), token, inviterID, now)
	if err != nil {
		return nil, err
	}

	// Deduct token from user
	_, err = r.db.ExecContext(ctx, "UPDATE users SET invite_tokens = invite_tokens - 1 WHERE id = ?", inviterID)
	if err != nil {
		return nil, err
	}

	return invite, nil
}

// GetInviteByToken gets an invite by its token
func (r *InviteRepo) GetInviteByToken(ctx context.Context, token string) (*model.Invite, error) {
	var invite model.Invite
	err := r.db.GetContext(ctx, &invite, `
		SELECT * FROM invites WHERE token = ? AND revoked_at IS NULL
	`, token)
	if err != nil {
		return nil, err
	}
	return &invite, nil
}

// UseInvite marks an invite as used and creates a new user
func (r *InviteRepo) UseInvite(ctx context.Context, token string, username string, passwordHash string) error {
	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get invite
	var invite model.Invite
	err = tx.GetContext(ctx, &invite, `
		SELECT * FROM invites WHERE token = ? AND revoked_at IS NULL AND used_at IS NULL
	`, token)
	if err != nil {
		return err
	}

	// Check if username is already taken
	var exists int
	err = tx.GetContext(ctx, &exists, "SELECT 1 FROM users WHERE username = ?", username)
	if err == nil {
		return ErrUsernameTaken
	}

	// Create user
	userID := utils.GenerateSnowflakeID()
	now := time.Now().Unix()
	_, err = tx.ExecContext(ctx, `
		INSERT INTO users (id, username, password_hash, created_at, user_role, invite_tokens) 
		VALUES (?, ?, ?, ?, 'user', 1)
	`, userID, username, passwordHash, now)
	if err != nil {
		return err
	}

	// Mark invite as used
	_, err = tx.ExecContext(ctx, `
		UPDATE invites SET invitee_username = ?, invitee_id = ?, used_at = ? WHERE token = ?
	`, username, userID, now, token)
	if err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit()
}

// RevokeInvite revokes an unused invite and returns the token to the user
func (r *InviteRepo) RevokeInvite(ctx context.Context, token string, inviterID int64) error {
	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get invite
	var invite model.Invite
	err = tx.GetContext(ctx, &invite, `
		SELECT * FROM invites WHERE token = ? AND inviter_id = ? AND revoked_at IS NULL AND used_at IS NULL
	`, token, inviterID)
	if err != nil {
		return err
	}

	// Mark invite as revoked
	now := time.Now().Unix()
	_, err = tx.ExecContext(ctx, `
		UPDATE invites SET revoked_at = ? WHERE token = ?
	`, now, token)
	if err != nil {
		return err
	}

	// Return token to user
	_, err = tx.ExecContext(ctx, `
		UPDATE users SET invite_tokens = invite_tokens + 1 WHERE id = ?
	`, inviterID)
	if err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit()
}

// GetUserInvites gets all invites for a user
func (r *InviteRepo) GetUserInvites(ctx context.Context, userID int64) ([]model.Invite, error) {
	var invites []model.Invite
	err := r.db.SelectContext(ctx, &invites, `
		SELECT * FROM invites WHERE inviter_id = ? ORDER BY created_at DESC
	`, userID)
	return invites, err
}

// GetUserInviteTokens gets the number of invite tokens a user has
func (r *InviteRepo) GetUserInviteTokens(ctx context.Context, userID int64) (int, error) {
	var tokens int
	err := r.db.GetContext(ctx, &tokens, "SELECT invite_tokens FROM users WHERE id = ?", userID)
	return tokens, err
}

// generateUniqueToken generates a unique invite token
func (r *InviteRepo) generateUniqueToken(ctx context.Context) (string, error) {
	for {
		// Generate 32-byte random token
		bytes := make([]byte, 32)
		if _, err := rand.Read(bytes); err != nil {
			return "", err
		}
		token := hex.EncodeToString(bytes)

		// Check if token already exists
		var exists int
		err := r.db.GetContext(ctx, &exists, "SELECT 1 FROM invites WHERE token = ?", token)
		if err != nil {
			// Token doesn't exist, we can use it
			return token, nil
		}
		// Token exists, try again
	}
}

// Custom errors
var (
	ErrNoInviteTokens = errors.New("no invite tokens available")
	ErrUsernameTaken  = errors.New("username already taken")
)
