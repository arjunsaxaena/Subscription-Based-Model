# Subscription-Based Model

A microservices-based subscription management system built with Go, implementing a clean architecture pattern. This system manages users, subscription plans, and user subscriptions.

## Architecture

The project is structured as a microservices architecture with three main services:

1. **User Service**: Handles user management and authentication
2. **Plan Service**: Manages subscription plans and their details
3. **Subscription Service**: Handles user subscriptions and their lifecycle
4. **Pkg**: Contains all the common stuff which may be used by many services

Each service follows a clean architecture pattern with:
- Controllers: Handle HTTP requests and responses
- Repository: Manages data persistence
- Models: Define data structures
- CMD: Contains service entry points

## Tech Stack

- **Language**: Go 1.24.2
- **Web Framework**: Gin
- **Database**: PostgreSQL
- **Authentication**: JWT
- **Environment Variables**: godotenv
- **Database Migration**: SQL files
- **UUID**: google/uuid
- **SQL Builder**: huandu/go-sqlbuilder

## Project Structure
```
Subscription-Based-Model
├─ README.md
├─ docs
│  └─ QL.postman_collection.json
├─ go.mod
├─ go.sum
├─ migrations
│  ├─ 1_database.sql
│  ├─ 2_user.sql
│  ├─ 3_plan.sql
│  └─ 4_subscription.sql
├─ pkg
│  ├─ auth
│  │  ├─ auth.go
│  │  └─ middleware.go
│  ├─ database
│  │  └─ db.go
│  ├─ models
│  │  ├─ basentity.go
│  │  ├─ login.go
│  │  ├─ plan.go
│  │  ├─ subscription.go
│  │  └─ user.go
│  └─ utils
│     └─ validate.go
├─ plan_service
│  ├─ cmd
│  │  └─ main.go
│  ├─ controllers
│  │  └─ plan.go
│  └─ repository
│     └─ plan.go
├─ subscription_service
│  ├─ cmd
│  │  └─ main.go
│  ├─ controllers
│  │  └─ subscription.go
│  └─ repository
│     └─ subscription.go
└─ user_service
   ├─ cmd
   │  └─ main.go
   ├─ controllers
   │  └─ user.go
   └─ repository
      └─ user.go

```

## Getting Started

### Prerequisites

- Go 1.24.2 or later
- PostgreSQL

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/arjunsaxaena/Subscription-Based-Model.git
   cd Subscription-Based-Model
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up environment variables:
   Create a `.env` file in the root directory with the following variables:
   ```
   DB_URL=
   JWT_SECRET=
   USER_SERVICE_PORT=
   PLAN_SERVICE_PORT=
   SUBSCRIPTION_SERVICE_PORT=
   INTERNAL_SERVER_TOKEN=
   USER_SERVICE_URL=
   PLAN_SERVICE_URL=
   ```

4. Run database migrations:
   ```bash
   # Execute SQL files in the migrations directory in order
   ```

5. Start the services:
   ```bash
   # Start each service in a separate terminal
   cd user_service/cmd
   go run main.go

   cd plan_service/cmd
   go run main.go

   cd subscription_service/cmd
   go run main.go
   ```

## API Documentation

Detailed API documentation can be found in the `docs` directory.