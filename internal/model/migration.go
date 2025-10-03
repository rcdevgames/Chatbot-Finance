package model

import (
	"time"
)

// MigrationSchema tracks database migration versions
type MigrationSchema struct {
	Version   string    `gorm:"primaryKey;type:varchar(20)" json:"version"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Applied   bool      `gorm:"default:false" json:"applied"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (MigrationSchema) TableName() string {
	return "migration_schemas"
}

// Default license keys for seeding (empty - will be generated via CLI)
var DefaultLicenseKeys = []LicenseKey{}

// Default categories for seeding
var DefaultCategories = []Category{
	{
		Name:     "Gaji",
		Type:     "income",
		Icon:     "💰",
		Keywords: []string{"gaji", "salary", "income", "penghasilan", "upah"},
	},
	{
		Name:     "Investasi",
		Type:     "income",
		Icon:     "📈",
		Keywords: []string{"investasi", "dividend", "profit", "keuntungan", "bunga"},
	},
	{
		Name:     "Freelance",
		Type:     "income",
		Icon:     "💼",
		Keywords: []string{"freelance", "project", "client", "proyek", "sampingan"},
	},
	{
		Name:     "Bonus",
		Type:     "income",
		Icon:     "🎁",
		Keywords: []string{"bonus", "thr", "insentif", "tunjangan"},
	},
	{
		Name:     "Makanan",
		Type:     "expense",
		Icon:     "🍔",
		Keywords: []string{"makan", "food", "resto", "warteg", "lunch", "dinner", "nasi"},
	},
	{
		Name:     "Transport",
		Type:     "expense",
		Icon:     "🚗",
		Keywords: []string{"transport", "bensin", "grab", "gojek", "parkir", "ojek", "tol"},
	},
	{
		Name:     "Tagihan",
		Type:     "expense",
		Icon:     "📄",
		Keywords: []string{"tagihan", "listrik", "air", "internet", "wifi", "token", "telepon"},
	},
	{
		Name:     "Kos/Sewa",
		Type:     "expense",
		Icon:     "🏠",
		Keywords: []string{"kos", "sewa", "rent", "kontrakan", "rumah"},
	},
	{
		Name:     "Hiburan",
		Type:     "expense",
		Icon:     "🎮",
		Keywords: []string{"hiburan", "nonton", "game", "netflix", "spotify", "bioskop"},
	},
	{
		Name:     "Belanja",
		Type:     "expense",
		Icon:     "🛒",
		Keywords: []string{"belanja", "shopping", "beli", "tokped", "shopee", "baju"},
	},
	{
		Name:     "Kesehatan",
		Type:     "expense",
		Icon:     "💊",
		Keywords: []string{"kesehatan", "obat", "dokter", "hospital", "rumah sakit", "vitamin"},
	},
	{
		Name:     "Pendidikan",
		Type:     "expense",
		Icon:     "📚",
		Keywords: []string{"kursus", "buku", "course", "training", "sekolah", "kuliah"},
	},
	{
		Name:     "Lainnya",
		Type:     "expense",
		Icon:     "📦",
		Keywords: []string{"lainnya", "other", "misc", "dll"},
	},
}