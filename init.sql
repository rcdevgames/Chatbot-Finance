-- Enable UUID extension for generating UUID primary keys
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create default categories for the finance bot
INSERT INTO categories (name, type, icon, keywords) VALUES
('Makanan', 'expense', '🍔', '["makan", "makanan", "nasi", "roti", "jajanan", "restoran", "kafe", "minuman"]'),
('Transportasi', 'expense', '🚗', '["transportasi", "bensin", "parkir", "tol", "ojek", "grab", "gojek", "bus", "kereta"]'),
('Belanja', 'expense', '🛒', '["belanja", "shopping", "baju", "sepatu", "elektronik", "kebutuhan", "rumah tangga"]'),
('Hiburan', 'expense', '🎮', '["hiburan", "film", "game", "musik", "konser", "bioskop", "liburan", "rekreasi"]'),
('Kesehatan', 'expense', '🏥', '["kesehatan", "dokter", "rumah sakit", "obat", "vitamin", "gym", "olahraga"]'),
('Pendidikan', 'expense', '📚', '["pendidikan", "sekolah", "kuliah", "buku", "kursus", "les", "biaya kuliah"]'),
('Tagihan', 'expense', '📄', '["tagihan", "listrik", "air", "internet", "telepon", "pajak", "asuransi"]'),
('Gaji', 'income', '💰', '["gaji", "salary", "penghasilan", "upah", "honor"]'),
('Bonus', 'income', '🎁', '["bonus", "tunjangan", "THR", "insentif"]'),
('Usaha', 'income', '💼', '["usaha", "bisnis", "usaha sampingan", "jual", "dagang"]'),
('Investasi', 'income', '📈', '["investasi", "saham", "deposito", "bunga", "dividen"]'),
('Lainnya', 'expense', '📝', '["lainnya", "lain-lain", "dll"]')
ON CONFLICT DO NOTHING;