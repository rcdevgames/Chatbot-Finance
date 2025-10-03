package repository

import (
	"fmt"

	"chatbot/internal/model"
)

type UserRepository struct {
	db *SupabaseClient
}

func NewUserRepository(db *SupabaseClient) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *model.User) error {
	var result []model.User
	err := r.db.Post("users", user, &result)
	if err != nil {
		return HandleSupabaseError(err)
	}

	if len(result) > 0 {
		user.ID = result[0].ID
	}

	return nil
}

func (r *UserRepository) GetUserByTelegramID(telegramUserID int64) (*model.User, error) {
	var users []model.User
	params := map[string]interface{}{
		"telegram_user_id": telegramUserID,
		"limit":            1,
	}
	err := r.db.Get("users", params, &users)
	if err != nil {
		return nil, HandleSupabaseError(err)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return &users[0], nil
}

func (r *UserRepository) UpdateUser(user *model.User) error {
	var result []model.User
	err := r.db.Patch("users", user.ID, user, &result)
	return HandleSupabaseError(err)
}

func (r *UserRepository) UserExists(telegramUserID int64) (bool, error) {
	var users []model.User
	params := map[string]interface{}{
		"telegram_user_id": telegramUserID,
		"select":           "id",
		"limit":            1,
	}
	err := r.db.Get("users", params, &users)
	if err != nil {
		return false, HandleSupabaseError(err)
	}

	return len(users) > 0, nil
}