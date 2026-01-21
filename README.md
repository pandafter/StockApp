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
- **Architecture**: Command Pattern for business logic.

### Infrastructure
- **Containerization**: [Docker](https://www.docker.com/) & Docker Compose

## 🚀 Getting Started

### Prerequisites
- Docker & Docker Compose
- Go 1.23+ (for local development)
- Node.js 18+ (for local development)

### Quick Start (Docker)
The easiest way to get the app running is using Docker Compose:

```bash
docker-compose up -d
```
This will start the CockroachDB instance. You can then run the backend and frontend locally or containerize them as well.

### Manual Setup

#### 1. Backend
```bash
cd backend
go mod download
go run cmd/server/main.go
```
The server will start at `http://localhost:8081`.

#### 2. Frontend
```bash
cd frontend
npm install
npm run dev
```
The dashboard will be available at `http://localhost:5173`.

## 🏗️ Project Structure

```text
.
├── backend/            # Go source code
│   ├── cmd/            # Entry points (main.go)
│   ├── internal/       # Private library code (business logic, DB, models)
│   └── ...
├── frontend/           # Vue.js source code
│   ├── src/            # Components, Views, Stores, Styles
│   └── ...
├── docker-compose.yml  # Infrastructure as code
└── .gitignore          # Repository hygiene
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
