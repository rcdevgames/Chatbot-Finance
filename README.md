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
- Supabase account
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
# Telegram
TELEGRAM_BOT_TOKEN=your_bot_token_here

# Supabase
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your_supabase_anon_key

# Groq AI
GROQ_API_KEY=gsk_your_groq_api_key

# Server
PORT=8080
WEBHOOK_URL=https://yourdomain.com/webhook
```

### 3. Setup Database

1. Buat project baru di [Supabase](https://supabase.com)
2. Copy `database/init.sql` ke Supabase SQL editor
3. Run script SQL tersebut

### 4. Setup Telegram Bot

1. Chat dengan [@BotFather](https://t.me/botfather) di Telegram
2. Create new bot: `/newbot`
3. Copy token dan masukkan ke `TELEGRAM_BOT_TOKEN`

### 5. Setup Groq API

1. Daftar di [Groq](https://groq.com)
2. Dapatkan API key
3. Masukkan ke `GROQ_API_KEY`

## 🏃‍♂️ Run Locally

```bash
# Download dependencies
go mod tidy

# Run bot
go run cmd/bot/main.go
```

## 🚀 Deployment

### Railway

```bash
# Install Railway CLI
npm install -g @railway/cli

# Login
railway login

# Deploy
railway up
```

### Docker

```bash
# Build image
docker build -t telegram-finance-bot .

# Run container
docker run -p 8080:8080 --env-file .env telegram-finance-bot
```

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
| `TELEGRAM_BOT_TOKEN` | ✅ | - | Telegram bot token |
| `SUPABASE_URL` | ✅ | - | Supabase project URL |
| `SUPABASE_ANON_KEY` | ✅ | - | Supabase anon key |
| `GROQ_API_KEY` | ✅ | - | Groq API key |
| `PORT` | ❌ | 8080 | Server port |
| `WEBHOOK_URL` | ❌ | - | Telegram webhook URL |
| `DEFAULT_TIMEZONE` | ❌ | Asia/Jakarta | User timezone |
| `DEFAULT_LANGUAGE` | ❌ | id | Default language |

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
- Pastikan Supabase credentials benar
- Check SQL schema sudah di-run
- Verify RLS policies

### LLM tidak bekerja
- Check `GROQ_API_KEY` valid
- Verify quota tidak habis
- Check network connection