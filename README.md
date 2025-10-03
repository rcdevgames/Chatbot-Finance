# Telegram Finance Bot

Telegram bot untuk tracking keuangan personal dengan natural language processing.

## 🚀 Fitur

- **Natural Language Input**: Input transaksi pake bahasa sehari-hari
- **Smart Query**: Tanya data keuangan dengan bahasa natural
- **Auto Insights**: Dapat financial advice otomatis
- **CRUD via Chat**: Update/delete transaksi lewat chat
- **Multi-format Support**: Support berbagai format angka (50rb, 50.000, 5jt)

## 📋 Persyaratan

- Go 1.21+
- PostgreSQL 15+
- Docker & Docker Compose (untuk deployment)
- Groq API key
- Telegram Bot token

## 🛠️ Setup

### 1. Clone Repository

```bash
git clone <repository-url>
cd telegram-finance-bot
```

### 2. Setup Environment

```bash
cp .env.example .env
```

Edit `.env` file:

```bash
# Database
DATABASE_URL=postgres://chatbot:your_password_here@postgres:5432/chatbot?sslmode=disable
POSTGRES_DB=chatbot
POSTGRES_USER=chatbot
POSTGRES_PASSWORD=your_password_here
POSTGRES_PORT=5432

# Telegram
TELEGRAM_BOT_TOKEN=your_bot_token_here

# Groq AI
GROQ_API_KEY=gsk_your_groq_api_key

# Server
PORT=8080
WEBHOOK_URL=https://yourdomain.com/webhook
```

### 3. Setup Database

#### Opsi 1: Docker Compose (Recommended)

```bash
# Buat network external waw_bridge
docker network create waw_bridge

# Jalankan database dan aplikasi
docker-compose up -d

# Database akan otomatis ter-setup dengan schema yang dibutuhkan
```

#### Opsi 2: PostgreSQL Manual

1. Install PostgreSQL 15+
2. Buat database:
   ```sql
   CREATE DATABASE chatbot;
   CREATE USER chatbot WITH PASSWORD 'your_password_here';
   GRANT ALL PRIVILEGES ON DATABASE chatbot TO chatbot;
   ```
3. Jalankan schema:
   ```bash
   psql -U chatbot -d chatbot -f init.sql
   ```

### 4. Setup Telegram Bot

1. Chat dengan [@BotFather](https://t.me/botfather) di Telegram
2. Create new bot: `/newbot`
3. Copy token dan masukkan ke `TELEGRAM_BOT_TOKEN`

### 5. Setup Groq API

1. Daftar di [Groq](https://groq.com)
2. Dapatkan API key
3. Masukkan ke `GROQ_API_KEY`

## 🏃‍♂️ Run Locally

### Dengan Docker Compose (Recommended)

```bash
# Buat network external waw_bridge
docker network create waw_bridge

# Jalankan aplikasi dengan database
docker-compose up -d

# Check logs
docker-compose logs -f app
```

### Manual (Development)

```bash
# Install PostgreSQL dan setup database manual (lihat langkah 3)

# Download dependencies
go mod tidy

# Run bot
go run cmd/bot/main.go
```

## 🚀 Deployment

### Docker Compose (Production)

```bash
# Setup network external
docker network create waw_bridge

# Build dan run di background
docker-compose -f docker-compose.yml up -d --build

# Check health endpoint
curl http://localhost:8080/health
```

### Manual Docker

```bash
# Build image
docker build -t telegram-finance-bot .

# Run container
docker run -p 8080:8080 \
  --env-file .env \
  --network waw_bridge \
  telegram-finance-bot
```

### Environment Variable Production

Pastikan environment variables berikut di-set di production:
- `DATABASE_URL` (PostgreSQL connection string)
- `TELEGRAM_BOT_TOKEN`
- `GROQ_API_KEY`
- `WEBHOOK_URL` (jika menggunakan webhook)

## 📱 Cara Pakai

### Input Transaksi
```
- "bayar makan 50rb di warteg"
- "gaji masuk 8 juta kemarin"
- "bayar kos 2juta"
```

### Query Data
```
- "pengeluaran minggu ini?"
- "income bulan ini berapa?"
- "total gua keluar buat makan?"
```

### Update/Delete
```
- "salah harusnya 45rb" → Update transaksi terakhir
- "hapus transaksi terakhir" → Delete
```

### Insights
```
- "gimana keuangan gua bulan ini?"
- "kasi summary keuangan dong"
```

## 📁 Struktur Proyek

```
cmd/bot/                # Main entry point
internal/
  config/              # Environment configuration
  handler/             # HTTP handlers (Telegram webhook)
  model/               # Data structures & types
  repository/          # Database operations
  service/             # Business logic
  util/                # Utility functions
pkg/telegram/          # Telegram API client
database/              # SQL schemas
```

## 🔧 Konfigurasi

Environment variables:

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DATABASE_URL` | ✅ | - | PostgreSQL connection string |
| `TELEGRAM_BOT_TOKEN` | ✅ | - | Telegram bot token |
| `GROQ_API_KEY` | ✅ | - | Groq API key |
| `PORT` | ❌ | 8080 | Server port |
| `WEBHOOK_URL` | ❌ | - | Telegram webhook URL |
| `DEFAULT_TIMEZONE` | ❌ | Asia/Jakarta | User timezone |
| `DEFAULT_LANGUAGE` | ❌ | id | Default language |
| `GROQ_MODEL` | ❌ | llama-3.1-70b-versatile | Groq AI model |
| `GROQ_TEMPERATURE` | ❌ | 0.3 | AI response temperature |
| `GROQ_MAX_TOKENS` | ❌ | 1000 | Maximum response tokens |

## 🤝 Contributing

1. Fork repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 💡 Tips & Tricks

- Bot support bahasa Indonesia yang santai
- Format angka: 50rb, 50.000, 50000, 5jt
- Tanggal auto-detect: "hari ini", "kemarin", "minggu ini"
- Kategori auto-detect dari keywords
- Chat history disimpan untuk context

## 🆘 Troubleshooting

### Bot tidak merespon
- Check `TELEGRAM_BOT_TOKEN` benar
- Pastikan webhook sudah diset dengan benar
- Check server logs

### Error database
- Pastikan PostgreSQL credentials benar
- Check database sudah running dan accessible
- Verify UUID extension sudah di-install
- Check schema sudah di-migrate dengan benar

### LLM tidak bekerja
- Check `GROQ_API_KEY` valid
- Verify quota tidak habis
- Check network connection