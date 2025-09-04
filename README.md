# Kinance

Kinance is a modern family finance management API built with Go, Gin, GORM, and PostgreSQL. It provides endpoints for user authentication, profile management, budgeting, transactions, and receipt processing, with comprehensive OpenAPI documentation and containerized deployment.

## Features
- **User Authentication:** Register, login, and manage JWT tokens securely.
- **Profile Management:** View and update user profiles, manage family members.
- **Budgeting:** Create, update, delete, and list budgets for personal or family use.
- **Transactions:** Record, update, and delete financial transactions.
- **Receipts:** Upload and process receipt images using AI, retrieve receipt details.
- **API Documentation:** Interactive docs at `/docs` powered by Redoc and OpenAPI.
- **Dockerized:** Easy local development and deployment with Docker and Docker Compose.

## Project Structure
```
go.mod
Makefile
api/
  docs/
    openapi.yaml
cmd/
  api/
    main.go
internal/
  auth/
  budget/
  common/
  investment/
  notification/
  receipt/
  transaction/
  user/
pkg/
  cache/
  config/
  database/
  logger/
  middleware/
  queue/
  validator/
scripts/
tests/
```

## Getting Started

### Prerequisites
- Go 1.24+
- Docker & Docker Compose
- PostgreSQL

### Local Development
1. **Clone the repository:**
   ```sh
   git clone https://github.com/pastorenue/kinance.git
   cd kinance
   ```
2. **Copy and edit environment variables:**
   ```sh
   cp .env.example .env
   # Edit .env as needed
   ```
3. **Start services with Docker Compose:**
   ```sh
   make docker-run
   ```
4. **Run database migrations:**
   ```sh
   make into_db
   ```
5. **Run the application:**
   ```sh
   make run-app
   ```
6. **Access API docs:**
   - Open [http://localhost:8080/docs](http://localhost:8080/docs) in your browser.

### Manual Go Run
If you prefer not to use Docker:
```sh
go mod tidy
go run cmd/api/main.go
```

## API Documentation
- **OpenAPI Spec:** Located at `api/docs/openapi.yaml`
- **Redoc UI:** Served at `/docs`
- **Endpoints:**
  - `/api/v1/auth/register` - Register user
  - `/api/v1/auth/login` - Login
  - `/api/v1/auth/refresh` - Refresh token
  - `/api/v1/users/profile` - Get/Update profile
  - `/api/v1/users/family` - Get family members
  - `/api/v1/budgets` - List/Create budgets
  - `/api/v1/budgets/{id}` - Update/Delete budget
  - `/api/v1/transactions` - List/Create transactions
  - `/api/v1/transactions/{id}` - Update/Delete transaction
  - `/api/v1/receipts/upload` - Upload receipt
  - `/api/v1/receipts` - List receipts
  - `/api/v1/receipts/{id}` - Get receipt details

## Environment Variables
See `.env.example` for all required variables, including database, JWT, and AI config.

## Testing
Run unit and integration tests:
```sh
make test
```

## Contributing
1. Fork the repo
2. Create your feature branch (`git checkout -b feature/fooBar`)
3. Commit your changes (`git commit -am 'Add some fooBar'`)
4. Push to the branch (`git push origin feature/fooBar`)
5. Create a new Pull Request

## License
MIT

## Maintainers
- Emmanuel Pastor ([GitHub](https://github.com/pastorenue))

---
For questions or support, open an issue or contact the maintainer.
