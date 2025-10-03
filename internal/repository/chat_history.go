package repository

import (
	"chatbot/internal/model"
)

type ChatHistoryRepository struct {
	db *SupabaseClient
}

func NewChatHistoryRepository(db *SupabaseClient) *ChatHistoryRepository {
	return &ChatHistoryRepository{db: db}
}

func (r *ChatHistoryRepository) SaveChatHistory(history *model.ChatHistory) error {
	var result []model.ChatHistory
	err := r.db.Post("chat_history", history, &result)
	if err != nil {
		return HandleSupabaseError(err)
	}

	if len(result) > 0 {
		history.ID = result[0].ID
	}

	return nil
}

func (r *ChatHistoryRepository) GetRecentChatHistory(userID string, limit int) ([]model.ChatHistory, error) {
	var history []model.ChatHistory
	params := map[string]interface{}{
		"user_id": userID,
		"order":   "created_at.desc",
		"limit":   limit,
	}
	err := r.db.Get("chat_history", params, &history)
	if err != nil {
		return nil, HandleSupabaseError(err)
	}

	return history, nil
}