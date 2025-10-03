package repository

import (
	"chatbot/internal/model"

	"gorm.io/gorm"
)

type ChatHistoryRepository struct {
	db *gorm.DB
}

func NewChatHistoryRepository(db *gorm.DB) *ChatHistoryRepository {
	return &ChatHistoryRepository{db: db}
}

func (r *ChatHistoryRepository) SaveChatHistory(history *model.ChatHistory) error {
	return r.db.Create(history).Error
}

func (r *ChatHistoryRepository) GetRecentChatHistory(userID string, limit int) ([]model.ChatHistory, error) {
	var history []model.ChatHistory
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&history).Error
	return history, err
}

func (r *ChatHistoryRepository) GetChatHistoryByUserID(userID string) ([]model.ChatHistory, error) {
	var history []model.ChatHistory
	err := r.db.Where("user_id = ?", userID).
		Order("created_at ASC").
		Find(&history).Error
	return history, err
}

func (r *ChatHistoryRepository) DeleteChatHistoryByUserID(userID string) error {
	return r.db.Delete(&model.ChatHistory{}, "user_id = ?", userID).Error
}