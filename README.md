# ACT_GO

Automated Crypto Trading platform built in Go. Connects to multiple exchanges (BingX, Binance, Phemex) for algorithmic trading with a web dashboard and Vue.js frontend.

## Architecture

```
ACT_GO/
├── main/          # Application entry point, logging setup, autorestart
├── web/           # Gin HTTP server (port 8080), REST API, authentication
│   └── act-vue-src/   # Vue 3 + Vite + Tailwind frontend
├── bngx/          # BingX exchange – bot logic, order API, WebSocket streams
├── bnce/          # Binance exchange – kline listener
├── phmx/          # Phemex exchange – account WebSocket
├── db/            # PostgreSQL via GORM – logging, user management
│   └── entities/  # Database models (Log, User)
├── common/        # Shared types (IOrder, Position)
├── utils/         # Generic Publisher[T], helper functions, build versioning
└── publish/       # Deployment scripts
```

## Tech Stack

| Layer        | Technology                                     |
| ------------ | ---------------------------------------------- |
| Language     | Go 1.22                                        |
| Web Server   | [Gin](https://github.com/gin-gonic/gin)        |
| Database     | PostgreSQL ([GORM](https://gorm.io))            |
| Frontend     | Vue 3 + Vite + Tailwind CSS                    |
| Exchanges    | BingX, Binance, Phemex (REST + WebSocket)      |
| Logging      | Lumberjack (rotating file logs on Linux)       |
| Hot Reload   | [autorestart](https://github.com/tillberg/autorestart) |

## Prerequisites

- **Go** 1.22+
- **PostgreSQL** with a database named `act`
- **Node.js** 18+ (for the Vue frontend)

## Getting Started

### 1. Clone

```bash
git clone https://github.com/paulsanswrk/act-go.git
cd act-go
```

### 2. Database

Create a PostgreSQL database and ensure connection settings in `db/db.go` match your environment:

```
host=localhost user=postgres password=qwerty dbname=act sslmode=disable
```

### 3. Run the backend

```bash
go run ./main
```

The API server starts on **port 8080**.

### 4. Run the frontend

```bash
cd web/act-vue-src
npm install
npm run dev
```

The dev server starts on `http://localhost:5173`.

## API Endpoints

| Method | Path            | Description                |
| ------ | --------------- | -------------------------- |
| GET    | `/`             | Health check               |
| GET    | `/ping`         | Pong                       |
| GET    | `/start-long`   | Start Bot1 long strategy   |
| GET    | `/stop-long`    | Stop Bot1 long strategy    |
| GET    | `/start-short`  | Start Bot1 short strategy  |
| GET    | `/stop-short`   | Stop Bot1 short strategy   |
| POST   | `/login`        | Authenticate (JSON body)   |

## Build

```bash
go build -o act_go ./main
```

On Linux the app writes rotating logs to `/home/ubuntu/act/log/act_go.log`.

## License

Private – All rights reserved.
