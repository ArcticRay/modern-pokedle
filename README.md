# Pokedle 🎮

A daily Pokémon guessing game – Wordle-style. Built with Go, Kafka, Redis, PostgreSQL, and deployed on AWS.

## Tech Stack

| Layer              | Technology                  |
| ------------------ | --------------------------- |
| **Language**       | Go 1.23                     |
| **HTTP**           | chi                         |
| **Auth**           | GitHub OAuth2 + JWT         |
| **Database**       | PostgreSQL 16 + pgx + sqlc  |
| **Cache/Session**  | Redis 7                     |
| **Messaging**      | Apache Kafka (KRaft)        |
| **Observability**  | Prometheus + Grafana + Loki |
| **API Docs**       | Swagger / OpenAPI           |
| **Testing**        | testify + Testcontainers    |
| **CI/CD**          | GitHub Actions              |
| **Infrastructure** | Terraform + AWS ECS Fargate |

## Getting Started

### Prerequisites

- [Go 1.23+](https://go.dev/dl/)
- [Docker + Docker Compose](https://docs.docker.com/get-docker/)
- [golangci-lint](https://golangci-lint.run/usage/install/)
- [golang-migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
- [sqlc](https://docs.sqlc.dev/en/latest/overview/install.html)

### Setup

```bash
# 1. Clone the repo
git clone https://github.com/yourusername/pokedle
cd pokedle

# 2. Start all infrastructure (Postgres, Redis, Kafka, Grafana...)
make infra-up

# 3. Create your .env file
cp .env.example .env
# Edit .env and fill in GITHUB_CLIENT_ID, GITHUB_CLIENT_SECRET, JWT_SECRET

# 4. Run the API
make run
```

### Available Services

| Service      | URL                                 |
| ------------ | ----------------------------------- |
| API          | http://localhost:8080               |
| Swagger Docs | http://localhost:8080/swagger       |
| Health Check | http://localhost:8080/health        |
| Metrics      | http://localhost:8080/metrics       |
| Kafka UI     | http://localhost:8090               |
| Prometheus   | http://localhost:9090               |
| Grafana      | http://localhost:3001 (admin/admin) |

## Development

```bash
make help        # show all available commands
make lint        # run golangci-lint
make test        # run unit tests
make test-integration  # run all tests (requires Docker)
make generate    # regenerate sqlc + swagger
```

## Project Structure

```
pokedle/
├── cmd/
│   ├── api/                    # API server entrypoint
│   ├── leaderboard-consumer/   # Kafka consumer: updates leaderboard
│   └── analytics-consumer/     # Kafka consumer: writes stats
├── internal/
│   ├── auth/                   # JWT + OAuth2
│   ├── config/                 # Environment-based configuration
│   ├── database/               # Connection pool + migrations
│   ├── game/                   # Core game logic
│   ├── kafka/                  # Producers and consumers
│   ├── middleware/              # HTTP middleware (auth, logging, metrics)
│   ├── observability/          # Logging + Prometheus metrics
│   ├── pokemon/                # PokéAPI client
│   └── server/                 # HTTP router + server lifecycle
├── migrations/                 # SQL migration files (golang-migrate)
├── deployments/
│   ├── docker/                 # Prometheus, Grafana, Promtail configs
│   └── terraform/              # AWS infrastructure (Phase 9)
├── .github/workflows/          # GitHub Actions CI/CD
├── docker-compose.yml          # Local infrastructure
├── Dockerfile                  # Multi-stage production build
└── Makefile                    # Dev commands
```
