package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/akramboussanni/marchive/internal/model"
	"github.com/akramboussanni/marchive/internal/utils"
	"github.com/jmoiron/sqlx"
)

type RequestCreditsRepo struct {
	Columns
	db *sqlx.DB
}

func NewRequestCreditsRepo(db *sqlx.DB) *RequestCreditsRepo {
	repo := &RequestCreditsRepo{db: db}
	repo.Columns = ExtractColumns[model.RequestCreditsLog]()
	return repo
}

// GetUserRequestCredits gets the current request credits for a user
func (r *RequestCreditsRepo) GetUserRequestCredits(ctx context.Context, userID int64) (int, error) {
	var credits int
	query := `SELECT request_credits FROM users WHERE id = $1`
	err := r.db.GetContext(ctx, &credits, query, userID)
	return credits, err
}

// UpdateUserRequestCredits updates the request credits for a user
func (r *RequestCreditsRepo) UpdateUserRequestCredits(ctx context.Context, userID int64, amount int) error {
	query := `UPDATE users SET request_credits = request_credits + $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, amount, userID)
	return err
}

// LogCreditChange logs a change to user credits
func (r *RequestCreditsRepo) LogCreditChange(ctx context.Context, log *model.RequestCreditsLog) error {
	query := fmt.Sprintf(
		"INSERT INTO request_credits_log (%s) VALUES (%s)",
		r.AllRaw,
		r.AllPrefixed,
	)
	_, err := r.db.NamedExecContext(ctx, query, log)
	return err
}

// GrantCredits grants credits to a user and logs the action
func (r *RequestCreditsRepo) GrantCredits(ctx context.Context, userID int64, amount int, reason string, adminUserID int64) error {
	// Start a transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update user credits
	updateQuery := `UPDATE users SET request_credits = request_credits + $1 WHERE id = $2`
	_, err = tx.ExecContext(ctx, updateQuery, amount, userID)
	if err != nil {
		return fmt.Errorf("failed to update user credits: %w", err)
	}

	// Log the credit grant
	log := &model.RequestCreditsLog{
		ID:          utils.GenerateSnowflakeID(),
		UserID:      userID,
		Action:      model.RequestCreditActionGranted,
		Amount:      amount,
		Reason:      reason,
		AdminUserID: &adminUserID,
		CreatedAt:   time.Now().Unix(),
	}

	logQuery := fmt.Sprintf(
		"INSERT INTO request_credits_log (%s) VALUES (%s)",
		r.AllRaw,
		r.AllPrefixed,
	)
	_, err = tx.NamedExecContext(ctx, logQuery, log)
	if err != nil {
		return fmt.Errorf("failed to log credit grant: %w", err)
	}

	// Commit transaction
	return tx.Commit()
}

// UseCredits uses credits for a user and logs the action
func (r *RequestCreditsRepo) UseCredits(ctx context.Context, userID int64, amount int) error {
	// Start a transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if user has enough credits
	var currentCredits int
	checkQuery := `SELECT request_credits FROM users WHERE id = $1`
	err = tx.GetContext(ctx, &currentCredits, checkQuery, userID)
	if err != nil {
		return fmt.Errorf("failed to get user credits: %w", err)
	}

	if currentCredits < amount {
		return fmt.Errorf("insufficient request credits: %d available, %d needed", currentCredits, amount)
	}

	// Update user credits
	updateQuery := `UPDATE users SET request_credits = request_credits - $1 WHERE id = $2`
	_, err = tx.ExecContext(ctx, updateQuery, amount, userID)
	if err != nil {
		return fmt.Errorf("failed to update user credits: %w", err)
	}

	// Log the credit usage
	log := &model.RequestCreditsLog{
		ID:        utils.GenerateSnowflakeID(),
		UserID:    userID,
		Action:    model.RequestCreditActionUsed,
		Amount:    amount,
		Reason:    "Download request beyond daily limit",
		CreatedAt: time.Now().Unix(),
	}

	logQuery := fmt.Sprintf(
		"INSERT INTO request_credits_log (%s) VALUES (%s)",
		r.AllRaw,
		r.AllPrefixed,
	)
	_, err = tx.NamedExecContext(ctx, logQuery, log)
	if err != nil {
		return fmt.Errorf("failed to log credit usage: %w", err)
	}

	// Commit transaction
	return tx.Commit()
}

// GetUserCreditHistory gets the credit history for a user
func (r *RequestCreditsRepo) GetUserCreditHistory(ctx context.Context, userID int64, limit, offset int) ([]model.RequestCreditsLog, error) {
	var logs []model.RequestCreditsLog
	query := fmt.Sprintf(`
		SELECT %s FROM request_credits_log 
		WHERE user_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3
	`, r.AllRaw)
	err := r.db.SelectContext(ctx, &logs, query, userID, limit, offset)
	return logs, err
}

// CountUserCreditHistory gets the total count of credit history entries for a user
func (r *RequestCreditsRepo) CountUserCreditHistory(ctx context.Context, userID int64) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM request_credits_log WHERE user_id = $1`
	err := r.db.GetContext(ctx, &count, query, userID)
	return count, err
}

// GetAllUsersCredits gets request credits for all users (admin function)
func (r *RequestCreditsRepo) GetAllUsersCredits(ctx context.Context) (map[int64]int, error) {
	credits := make(map[int64]int)
	query := `SELECT id, request_credits FROM users`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userID int64
		var userCredits int
		if err := rows.Scan(&userID, &userCredits); err != nil {
			return nil, err
		}
		credits[userID] = userCredits
	}

	return credits, nil
}
