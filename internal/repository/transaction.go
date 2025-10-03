package repository

import (
	"fmt"

	"chatbot/internal/model"
)

type TransactionRepository struct {
	db *SupabaseClient
}

func NewTransactionRepository(db *SupabaseClient) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) CreateTransaction(transaction *model.Transaction) error {
	var result []model.Transaction
	err := r.db.Post("transactions", transaction, &result)
	if err != nil {
		return HandleSupabaseError(err)
	}

	if len(result) > 0 {
		transaction.ID = result[0].ID
	}

	return nil
}

func (r *TransactionRepository) GetTransactionsByUserID(userID string, limit ...int) ([]model.Transaction, error) {
	limitValue := 50
	if len(limit) > 0 {
		limitValue = limit[0]
	}

	var transactions []model.Transaction
	params := map[string]interface{}{
		"user_id": userID,
		"order":   "created_at.desc",
		"limit":   limitValue,
	}
	err := r.db.Get("transactions", params, &transactions)
	if err != nil {
		return nil, HandleSupabaseError(err)
	}

	return transactions, nil
}

func (r *TransactionRepository) GetTransactionsByUserAndDateRange(userID, startDate, endDate string) ([]model.Transaction, error) {
	var transactions []model.Transaction
	params := map[string]interface{}{
		"user_id": userID,
		"date":    fmt.Sprintf("gte.%s,lte.%s", startDate, endDate),
		"order":   "date.desc",
	}
	err := r.db.Get("transactions", params, &transactions)
	if err != nil {
		return nil, HandleSupabaseError(err)
	}

	return transactions, nil
}

func (r *TransactionRepository) GetLastTransaction(userID string) (*model.Transaction, error) {
	transactions, err := r.GetTransactionsByUserID(userID, 1)
	if err != nil {
		return nil, err
	}

	if len(transactions) == 0 {
		return nil, fmt.Errorf("no transactions found")
	}

	return &transactions[0], nil
}

func (r *TransactionRepository) UpdateTransaction(transactionID string, transaction *model.Transaction) error {
	var result []model.Transaction
	err := r.db.Patch("transactions", transactionID, transaction, &result)
	return HandleSupabaseError(err)
}

func (r *TransactionRepository) DeleteTransaction(transactionID string) error {
	var result []model.Transaction
	err := r.db.Delete("transactions", transactionID, &result)
	return HandleSupabaseError(err)
}

func (r *TransactionRepository) GetTransactionsByCategory(userID, category string, limit ...int) ([]model.Transaction, error) {
	limitValue := 50
	if len(limit) > 0 {
		limitValue = limit[0]
	}

	var transactions []model.Transaction
	params := map[string]interface{}{
		"user_id": userID,
		"category": category,
		"order":   "created_at.desc",
		"limit":   limitValue,
	}
	err := r.db.Get("transactions", params, &transactions)
	if err != nil {
		return nil, HandleSupabaseError(err)
	}

	return transactions, nil
}

func (r *TransactionRepository) GetTransactionsByType(userID, transactionType string, limit ...int) ([]model.Transaction, error) {
	limitValue := 50
	if len(limit) > 0 {
		limitValue = limit[0]
	}

	var transactions []model.Transaction
	params := map[string]interface{}{
		"user_id": userID,
		"type":    transactionType,
		"order":   "created_at.desc",
		"limit":   limitValue,
	}
	err := r.db.Get("transactions", params, &transactions)
	if err != nil {
		return nil, HandleSupabaseError(err)
	}

	return transactions, nil
}