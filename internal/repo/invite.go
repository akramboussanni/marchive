package repo

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/akramboussanni/marchive/internal/model"
	"github.com/akramboussanni/marchive/internal/utils"
	"github.com/jmoiron/sqlx"
)

type InviteRepo struct {
	Columns
	db *sqlx.DB
}

func NewInviteRepo(db *sqlx.DB) *InviteRepo {
	repo := &InviteRepo{db: db}
	repo.Columns = ExtractColumns[model.Invite]()
	return repo
}

// CreateInvite creates a new invite for a user
func (r *InviteRepo) CreateInvite(ctx context.Context, inviterID int64) (*model.Invite, error) {
	// Check if user has invite tokens
	var tokens int
	err := r.db.GetContext(ctx, &tokens, "SELECT invite_tokens FROM users WHERE id = $1", inviterID)
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
	invite.ID = utils.GenerateSnowflakeID()
	query := fmt.Sprintf(`
		INSERT INTO invites (%s) 
		VALUES (%s)
	`, r.AllRaw, r.AllPrefixed)
	_, err = r.db.NamedExecContext(ctx, query, invite)
	if err != nil {
		return nil, err
	}

	// Deduct token from user
	_, err = r.db.ExecContext(ctx, "UPDATE users SET invite_tokens = invite_tokens - 1 WHERE id = $1", inviterID)
	if err != nil {
		return nil, err
	}

	return invite, nil
}

// GetInviteByToken gets an invite by its token
func (r *InviteRepo) GetInviteByToken(ctx context.Context, token string) (*model.Invite, error) {
	var invite model.Invite
	query := fmt.Sprintf(`SELECT %s FROM invites WHERE token = $1 AND revoked_at IS NULL`, r.AllRaw)
	err := r.db.GetContext(ctx, &invite, query, token)
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
	query := fmt.Sprintf(`SELECT %s FROM invites WHERE token = $1 AND revoked_at IS NULL AND used_at IS NULL`, r.AllRaw)
	err = tx.GetContext(ctx, &invite, query, token)
	if err != nil {
		return err
	}

	// Check if username is already taken
	var count int
	err = tx.GetContext(ctx, &count, "SELECT COUNT(*) FROM users WHERE username = $1", username)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrUsernameTaken
	}

	// Create user
	userID := utils.GenerateSnowflakeID()
	now := time.Now().Unix()
	_, err = tx.ExecContext(ctx, `
		INSERT INTO users (id, username, password_hash, created_at, user_role, invite_tokens) 
		VALUES ($1, $2, $3, $4, 'user', 1)
	`, userID, username, passwordHash, now)
	if err != nil {
		return err
	}

	// Mark invite as used
	_, err = tx.ExecContext(ctx, `
		UPDATE invites SET invitee_username = $1, invitee_id = $2, used_at = $3 WHERE token = $4
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
	query := fmt.Sprintf(`SELECT %s FROM invites WHERE token = $1 AND inviter_id = $2 AND revoked_at IS NULL AND used_at IS NULL`, r.AllRaw)
	err = tx.GetContext(ctx, &invite, query, token, inviterID)
	if err != nil {
		return err
	}

	// Mark invite as revoked
	now := time.Now().Unix()
	_, err = tx.ExecContext(ctx, `
		UPDATE invites SET revoked_at = $1 WHERE token = $2
	`, now, token)
	if err != nil {
		return err
	}

	// Return token to user
	_, err = tx.ExecContext(ctx, `
		UPDATE users SET invite_tokens = invite_tokens + 1 WHERE id = $1
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
	query := fmt.Sprintf(`SELECT %s FROM invites WHERE inviter_id = $1 ORDER BY created_at DESC`, r.AllRaw)
	err := r.db.SelectContext(ctx, &invites, query, userID)
	return invites, err
}

// GetUserInviteTokens gets the number of invite tokens a user has
func (r *InviteRepo) GetUserInviteTokens(ctx context.Context, userID int64) (int, error) {
	var tokens int
	err := r.db.GetContext(ctx, &tokens, "SELECT invite_tokens FROM users WHERE id = $1", userID)
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
		err := r.db.GetContext(ctx, &exists, "SELECT 1 FROM invites WHERE token = $1", token)
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
