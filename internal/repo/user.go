package repo

import (
	"context"
	"fmt"

	"github.com/akramboussanni/marchive/internal/model"
	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	Columns
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	repo := &UserRepo{db: db}
	repo.Columns = ExtractColumns[model.User]()
	return repo
}

func (r *UserRepo) CreateUser(ctx context.Context, user *model.User) error {
	query := fmt.Sprintf(
		"INSERT INTO users (%s) VALUES (%s)",
		r.AllRaw,
		r.AllPrefixed,
	)
	_, err := r.db.NamedExecContext(ctx, query, user)
	return err
}

func (r *UserRepo) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	query := fmt.Sprintf("SELECT %s FROM users WHERE id = $1", r.AllRaw)
	err := r.db.GetContext(ctx, &user, query, id)
	return &user, err
}

func (r *UserRepo) GetUserByIDSafe(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	query := fmt.Sprintf("SELECT %s FROM users WHERE id = $1", r.SafeRaw)
	err := r.db.GetContext(ctx, &user, query, id)
	return &user, err
}

func (r *UserRepo) DuplicateName(ctx context.Context, username string) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists, "SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)", username)
	return exists, err
}

func (r *UserRepo) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	query := fmt.Sprintf("SELECT %s FROM users WHERE username=$1", r.AllRaw)
	err := r.db.GetContext(ctx, &user, query, username)
	return &user, err
}

func (r *UserRepo) DeleteUser(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *UserRepo) ChangeUserPassword(ctx context.Context, newPasswordHash string, userID int64) error {
	query := `
		UPDATE users
		SET password_hash = $1
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, newPasswordHash, userID)
	return err
}

func (r *UserRepo) ChangeJwtSessionID(ctx context.Context, userID int64, newID int64) error {
	query := `
		UPDATE users
		SET jwt_session_id = $1
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, newID, userID)
	return err
}

func (r *UserRepo) CountUsers(ctx context.Context) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM users")
	return count, err
}

// GiveEveryoneInviteToken gives 1 invite token to all users
func (r *UserRepo) GiveEveryoneInviteToken(ctx context.Context) error {
	query := `UPDATE users SET invite_tokens = invite_tokens + 1`
	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *UserRepo) GetUsersWithStats(ctx context.Context, limit, offset int) ([]model.UserWithStats, error) {
	query := `
		SELECT u.id, u.username, u.user_role, u.created_at,
		       COALESCE(dl.download_count, 0) as download_count,
		       COALESCE(u.request_credits, 0) as request_credits
		FROM users u
		LEFT JOIN (
			SELECT user_id, COUNT(*) as download_count
			FROM downloadrequests
			GROUP BY user_id
		) dl ON u.id = dl.user_id
		ORDER BY u.created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.UserWithStats
	for rows.Next() {
		var user model.UserWithStats
		err := rows.Scan(
			&user.ID, &user.Username, &user.Role,
			&user.CreatedAt, &user.DownloadCount, &user.RequestCredits,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepo) SearchUsers(ctx context.Context, query, role string, limit, offset int) ([]model.UserWithStats, error) {
	var sqlQuery string
	var args []interface{}

	if query != "" && role != "" {
		sqlQuery = `
			SELECT u.id, u.username, u.user_role, u.created_at,
			       COALESCE(dl.download_count, 0) as download_count,
			       COALESCE(u.request_credits, 0) as request_credits
			FROM users u
			LEFT JOIN (
				SELECT user_id, COUNT(*) as download_count
				FROM downloadrequests
				GROUP BY user_id
			) dl ON u.id = dl.user_id
			WHERE u.username ILIKE $1 AND u.user_role = $2
			ORDER BY u.created_at DESC
			LIMIT $3 OFFSET $4
		`
		searchPattern := "%" + query + "%"
		args = []interface{}{searchPattern, role, limit, offset}
	} else if query != "" {
		sqlQuery = `
			SELECT u.id, u.username, u.user_role, u.created_at,
			       COALESCE(dl.download_count, 0) as download_count,
			       COALESCE(u.request_credits, 0) as request_credits
			FROM users u
			LEFT JOIN (
				SELECT user_id, COUNT(*) as download_count
				FROM downloadrequests
				GROUP BY user_id
			) dl ON u.id = dl.user_id
			WHERE u.username ILIKE $1
			ORDER BY u.created_at DESC
			LIMIT $2 OFFSET $3
		`
		searchPattern := "%" + query + "%"
		args = []interface{}{searchPattern, limit, offset}
	} else if role != "" {
		sqlQuery = `
			SELECT u.id, u.username, u.user_role, u.created_at,
			       COALESCE(dl.download_count, 0) as download_count,
			       COALESCE(u.request_credits, 0) as request_credits
			FROM users u
			LEFT JOIN (
				SELECT user_id, COUNT(*) as download_count
				FROM downloadrequests
				GROUP BY user_id
			) dl ON u.id = dl.user_id
			WHERE u.user_role = $1
			ORDER BY u.created_at DESC
			LIMIT $2 OFFSET $3
		`
		args = []interface{}{role, limit, offset}
	} else {
		return r.GetUsersWithStats(ctx, limit, offset)
	}

	rows, err := r.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.UserWithStats
	for rows.Next() {
		var user model.UserWithStats
		err := rows.Scan(
			&user.ID, &user.Username, &user.Role,
			&user.CreatedAt, &user.DownloadCount, &user.RequestCredits,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepo) CountSearchUsers(ctx context.Context, query, role string) (int, error) {
	var sqlQuery string
	var args []interface{}

	if query != "" && role != "" {
		sqlQuery = "SELECT COUNT(*) FROM users WHERE username ILIKE $1 AND user_role = $2"
		searchPattern := "%" + query + "%"
		args = []interface{}{searchPattern, role}
	} else if query != "" {
		sqlQuery = "SELECT COUNT(*) FROM users WHERE username ILIKE $1"
		searchPattern := "%" + query + "%"
		args = []interface{}{searchPattern}
	} else if role != "" {
		sqlQuery = "SELECT COUNT(*) FROM users WHERE user_role = $1"
		args = []interface{}{role}
	} else {
		return r.CountUsers(ctx)
	}

	var count int
	err := r.db.GetContext(ctx, &count, sqlQuery, args...)
	return count, err
}

func (r *UserRepo) UpdateUser(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users
		SET username = $1, user_role = $2
		WHERE id = $3
	`
	_, err := r.db.ExecContext(ctx, query, user.Username, user.Role, user.ID)
	return err
}
