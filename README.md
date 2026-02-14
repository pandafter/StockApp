# StockApp: Market Dashboard 🚀

## 🌟 Features

- **Real-Time Data Ingestion**: Background worker that syncs with external market APIs.
- **Interactive Visualizations**: Dynamic price history charts using Chart.js.
- **Personalized Watchlist**: Add/remove stocks to your dashboard with persistent storage.
- **Smart Recommendations**: Algorithm-based stock suggestions based on market trends.
- **Premium UI**: Glassmorphism design system with smooth animations and responsive layout.
- **Resilient Backend**: Built with Go, featuring a robust command pattern architecture and background processing.

## 🛠️ Tech Stack

### Frontend
- **Framework**: [Vue.js 3](https://vuejs.org/) (Composition API)
- **Store**: [Pinia](https://pinia.vuejs.org/)
- **Styling**: [Tailwind CSS](https://tailwindcss.com/)
- **Charts**: [Chart.js](https://www.chartjs.org/)
- **Build Tool**: [Vite](https://vitejs.dev/)

### Backend
- **Language**: [Go](https://go.dev/) (1.23+)
- **API Framework**: [Gin](https://gin-gonic.com/)
- **ORM**: [GORM](https://gorm.io/)
- **Database**: [CockroachDB](https://www.cockroachlabs.com/) (PostgreSQL compatible)
- **Architecture**: Service/Repository pattern (reginald-style).

### Infrastructure
- **Containerization**: [Docker](https://www.docker.com/) & Docker Compose

## 🚀 Getting Started

### Prerequisites
- Docker & Docker Compose
- Go 1.23+ (for local development)
- Node.js 18+ (for local development)

### Quick Start (Docker)
The easiest way to run the app is using Docker Compose. Everything runs in containers:

```bash
docker-compose up -d
```

- **Backend API**: http://localhost:8081
- **Frontend**: http://localhost:3000 (proxies API to backend)
- **CockroachDB**: localhost:26257

The frontend acts as a "mask" of the backend—all API calls go through the frontend nginx proxy to the backend.

### Manual Setup

#### 1. Start CockroachDB (if not using full Docker)
```bash
docker-compose up -d cockroachdb
```

#### 2. Backend
```bash
cd backend
go mod download
go run ./cmd/server
```
The server will start at `http://localhost:8081`.

#### 3. Frontend
```bash
cd frontend
npm install
npm run dev
```
The dashboard will be at `http://localhost:5173` (Vite proxies `/api` to backend).

## 🏗️ Project Structure

```text
.
├── backend/                 # Go backend (reginald-style layout)
│   ├── cmd/server/          # Entry point
│   ├── internal/            # Services: shared, stocks, stocksync
│   ├── pkg/                 # database, middleware, utils
│   └── Dockerfile
├── frontend/                # Vue.js (mask of backend)
│   ├── src/
│   ├── nginx.conf           # API proxy to backend
│   └── Dockerfile
└── docker-compose.yml       # Backend + Frontend + CockroachDB
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
