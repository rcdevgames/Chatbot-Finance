package repository

import (
	"fmt"
	"log"
)

type SupabaseClient struct {
	customClient *SupabaseCustomClient
}

func NewSupabaseClient(url, key string) *SupabaseClient {
	client := NewSupabaseCustomClient(url, key)
	return &SupabaseClient{customClient: client}
}

func (s *SupabaseClient) Get(table string, params map[string]interface{}, result interface{}) error {
	return s.customClient.Get(table, params, result)
}

func (s *SupabaseClient) Post(table string, data interface{}, result interface{}) error {
	return s.customClient.Post(table, data, result)
}

func (s *SupabaseClient) Patch(table string, id string, data interface{}, result interface{}) error {
	return s.customClient.Patch(table, id, data, result)
}

func (s *SupabaseClient) Delete(table string, id string, result interface{}) error {
	return s.customClient.Delete(table, id, result)
}

// Helper function to handle Supabase errors
func HandleSupabaseError(err error) error {
	if err != nil {
		log.Printf("Supabase error: %v", err)
		return fmt.Errorf("database operation failed: %w", err)
	}
	return nil
}