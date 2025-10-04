package repository

import (
	"chatbot/internal/model"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) CreateTransaction(transaction *model.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *TransactionRepository) GetTransactionsByUserID(userID string, limit ...int) ([]model.Transaction, error) {
	limitValue := 50
	if len(limit) > 0 {
		limitValue = limit[0]
	}

	var transactions []model.Transaction
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limitValue).
		Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) GetTransactionsByUserAndDateRange(userID, startDate, endDate string) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.db.Where("user_id = ? AND date BETWEEN ? AND ?", userID, startDate, endDate).
		Order("date DESC").
		Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) GetLastTransaction(userID string) (*model.Transaction, error) {
	var transaction model.Transaction
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		First(&transaction).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &transaction, nil
}

func (r *TransactionRepository) UpdateTransaction(transaction *model.Transaction) error {
	return r.db.Save(transaction).Error
}

func (r *TransactionRepository) DeleteTransaction(transactionID string) error {
	return r.db.Delete(&model.Transaction{}, "id = ?", transactionID).Error
}

func (r *TransactionRepository) GetTransactionsByCategory(userID, category string, limit ...int) ([]model.Transaction, error) {
	limitValue := 50
	if len(limit) > 0 {
		limitValue = limit[0]
	}

	var transactions []model.Transaction
	err := r.db.Where("user_id = ? AND category = ?", userID, category).
		Order("created_at DESC").
		Limit(limitValue).
		Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) GetTransactionsByType(userID, transactionType string, limit ...int) ([]model.Transaction, error) {
	limitValue := 50
	if len(limit) > 0 {
		limitValue = limit[0]
	}

	var transactions []model.Transaction
	err := r.db.Where("user_id = ? AND type = ?", userID, transactionType).
		Order("created_at DESC").
		Limit(limitValue).
		Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) GetTransactionByID(transactionID string) (*model.Transaction, error) {
	var transaction model.Transaction
	err := r.db.Where("id = ?", transactionID).First(&transaction).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &transaction, nil
}

func (r *TransactionRepository) DeleteTransactionsByUserIDAndMonth(userID string, year int, month int) error {
	return r.db.Delete(&model.Transaction{}, "user_id = ? AND EXTRACT(YEAR FROM date) = ? AND EXTRACT(MONTH FROM date) = ?", userID, year, month).Error
}

func (r *TransactionRepository) DeleteTransactionsByUserIDAndYear(userID string, year int) error {
	return r.db.Delete(&model.Transaction{}, "user_id = ? AND EXTRACT(YEAR FROM date) = ?", userID, year).Error
}

func (r *TransactionRepository) DeleteAllTransactionsByUserID(userID string) error {
	return r.db.Delete(&model.Transaction{}, "user_id = ?", userID).Error
}