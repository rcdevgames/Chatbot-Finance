# 🚀 Quick Start Guide

## 📋 Prerequisites

- Telegram Bot token dari [@BotFather](https://t.me/botfather)
- Docker & Docker Compose
- Groq API key (gratis di [groq.com](https://groq.com))

## ⚡ 5 Menit Setup

### 1. Create Telegram Bot
```
/start /newbot
Nama bot: Finance Bot Gua
Username: finance_bot_gua
```
Copy token yang diberikan.

### 2. Setup Database dengan Docker
```bash
# Clone repository
git clone <repository-url>
cd telegram-finance-bot

# Create network external
docker network create waw_bridge

# Setup environment
cp .env.example .env
# Edit .env file dengan credentials kamu
```

### 3. Get Groq API Key
1. Go to [groq.com](https://groq.com)
2. Sign up
3. Dashboard → API Keys
4. Create new key
5. Copy key

### 4. Configure Environment
Edit `.env` file dengan credentials kamu:
```bash
# Database
DATABASE_URL=postgres://chatbot:your_password_here@postgres:5432/chatbot?sslmode=disable
POSTGRES_PASSWORD=your_password_here

# Telegram
TELEGRAM_BOT_TOKEN=123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11

# Groq AI
GROQ_API_KEY=gsk_1234567890abcdef
```

### 5. Run Bot dengan Docker Compose
```bash
# Jalankan database dan aplikasi
docker-compose up -d

# Check logs
docker-compose logs -f app

# Database otomatis ter-setup dengan schema dan default categories
```

## 🎱 Test Bot

Chat dengan bot kamu di Telegram:
- `/start` - Welcome message
- `bayar makan 50rb` - Add transaction
- `pengeluaran hari ini?` - Query expenses
- `/help` - Lihat semua commands

## ✅ Check Status

Health check endpoint:
```bash
curl http://localhost:8080/health
```

Response seharusnya:
```json
{
  "status": "ok",
  "message": "Telegram Finance Bot is running",
  "database": "healthy"
}
```

## 🚀 Deploy ke Production

### Docker Compose (Recommended)
```bash
# Build dan run di background
docker-compose up -d --build

# Check status
docker-compose ps
docker-compose logs -f app
```

### Railway
```bash
# Install Railway CLI
npm install -g @railway/cli

# Login & deploy
railway login
railway init
railway up
```

Set `DATABASE_URL` dan environment variables lainnya di Railway dashboard!

## 🔧 Commands Lengkap

| Command | Deskripsi |
|---------|-----------|
| `/start` | Mulai bot |
| `/help` | Bantuan |
| `/summary` | Summary keuangan |

## 💬 Natural Language Examples

**Input Transaksi:**
- `bayar makan 50rb di warteg`
- `gaji masuk 8jt kemarin`
- `bayar kos 2juta`
- `belanja pulsa 25rb`

**Query Data:**
- `pengeluaran minggu ini?`
- `income bulan ini berapa?`
- `total gua keluar buat makan?`
- `brp duit gua sisa?`

**Update/Delete:**
- `salah harusnya 45rb` (update last)
- `hapus transaksi terakhir`

**Insights:**
- `gimana keuangan gua bulan ini?`
- `kasi summary keuangan dong`
- `ada tips ga?`

## 🎯 Tips & Tricks

1. **Format Angka**: 50rb, 50.000, 50000, 5jt
2. **Auto-detect**: Hari ini, kemarin, minggu ini, bulan lalu
3. **Kategori**: Auto-detect dari keywords (makan, transport, dll)
4. **Context**: Bot ingat chat sebelumnya
5. **Error Handling**: Bot kasih saran kalau input tidak jelas

## 🆘 Troubleshooting

### Bot tidak merespon
- Check `TELEGRAM_BOT_TOKEN` benar
- Pastikan tidak ada error di console
- Test ping ke Telegram API

### Database error
- Verify PostgreSQL connection
- Check container status: `docker-compose ps`
- Check database logs: `docker-compose logs postgres`
- Pastikan network `waw_bridge` sudah dibuat

### Groq tidak bekerja
- Verify `GROQ_API_KEY`
- Check quota tidak habis
- Test network connection

## 📞 Support

Jika ada masalah:
1. Check console logs
2. Verify environment variables
3. Test dengan examples di atas
4. Cek README.md untuk detail

Enjoy using your Finance Bot! 🎉