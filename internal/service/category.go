package service

import (
	"strings"

	"chatbot/internal/model"
	"chatbot/internal/repository"
)

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
}

func NewCategoryService(categoryRepo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *CategoryService) GetAllCategories() ([]model.Category, error) {
	return s.categoryRepo.GetAllCategories()
}

func (s *CategoryService) GetCategoriesByType(categoryType string) ([]model.Category, error) {
	return s.categoryRepo.GetCategoriesByType(categoryType)
}

func (s *CategoryService) DetectCategory(description, transactionType string) (string, error) {
	if description == "" {
		if transactionType == "expense" {
			return "Lainnya", nil
		}
		return "Lainnya", nil
	}

	// Search for matching category by keywords
	categories, err := s.categoryRepo.SearchCategoryByKeyword(description)
	if err != nil {
		return "Lainnya", err
	}

	// Filter by transaction type
	for _, category := range categories {
		if category.Type == transactionType {
			return category.Name, nil
		}
	}

	// Try exact category name match
	allCategories, err := s.categoryRepo.GetAllCategories()
	if err == nil {
		descLower := strings.ToLower(description)
		for _, category := range allCategories {
			if category.Type == transactionType {
				if strings.Contains(descLower, strings.ToLower(category.Name)) {
					return category.Name, nil
				}
			}
		}
	}

	// If no match found, return default category
	if transactionType == "expense" {
		return "Lainnya", nil
	}
	return "Lainnya", nil
}

func (s *CategoryService) GetCategoryIcon(categoryName string) string {
	icons := map[string]string{
		"Makanan":     "🍔",
		"Transport":   "🚗",
		"Tagihan":     "📄",
		"Kos/Sewa":    "🏠",
		"Hiburan":     "🎮",
		"Belanja":     "🛒",
		"Kesehatan":   "💊",
		"Pendidikan":  "📚",
		"Lainnya":     "📦",
		"Gaji":        "💰",
		"Investasi":   "📈",
		"Freelance":   "💼",
		"Bonus":       "🎁",
	}

	if icon, exists := icons[categoryName]; exists {
		return icon
	}
	return "💰"
}