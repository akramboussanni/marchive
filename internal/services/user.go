package services

import (
	"context"
	"time"

	"github.com/akramboussanni/marchive/internal/model"
	"github.com/akramboussanni/marchive/internal/repo"
	"github.com/akramboussanni/marchive/internal/utils"
)

type UserService struct {
	userRepo *repo.UserRepo
}

func NewUserService(userRepo *repo.UserRepo) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

type CreateUserParams struct {
	Username       string
	PasswordHash   string
	Role           string
	RequestCredits int
}

// CreateUser creates a new user with the specified parameters
// Returns the created user or an error if username already exists
func (s *UserService) CreateUser(ctx context.Context, params CreateUserParams) (*model.User, error) {
	// Check if username already exists
	exists, err := s.userRepo.DuplicateName(ctx, params.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, repo.ErrUsernameTaken
	}

	// Set defaults
	if params.Role == "" {
		params.Role = "user"
	}

	// Create user model
	user := &model.User{
		ID:             utils.GenerateSnowflakeID(),
		Username:       params.Username,
		PasswordHash:   params.PasswordHash,
		Role:           params.Role,
		CreatedAt:      time.Now().Unix(),
		JwtSessionID:   utils.GenerateSnowflakeID(),
		InviteTokens:   0,
		RequestCredits: params.RequestCredits,
	}

	// Insert user into database
	err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
