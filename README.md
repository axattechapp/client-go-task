# Client GO task

## Table of Contents

- [Introduction](#introduction)
- [Prerequisites](#prerequisites)
- [Environment Setup](#environment-setup)
- [Commands](#commands)
  - [Start Database](#start-database)
  - [Stop Database](#stop-database)
  - [Database Migration](#database-migration)
    - [Migrate Up](#migrate-up)
    - [Migrate Down](#migrate-down)
  - [Run Server](#run-server)
  - [Run Tests](#run-tests)
- [API Testing with Postman](#api-testing-with-postman)

## Introduction

We worked on the development of a comprehensive web application using Go, PostgreSQL, and the Gin framework. The application will manage users, their profiles, and career-related data through both RESTful and GraphQL APIs. It will follow clean architecture principles to ensure modularity, testability, and maintainability. Key features include JWT-based authentication, concurrency management with goroutines, webhook notifications, and testing.

## Prerequisites

Before you begin, ensure you have met the following requirements:
- [Docker](https://docs.docker.com/get-docker/) installed.
- [Go](https://golang.org/dl/) installed.
- [Migrate CLI](https://github.com/golang-migrate/migrate) installed.
- PostgreSQL database set up and running.

## Database Setup

This project utilizes a Postgres database for data persistence. You can leverage Docker Compose to manage the database container.

**Run Docker Compose:**

  The project likely includes a `docker-compose.yml` file that defines the services (Postgres and potentially pgAdmin). To start the database container(s), run:

   ```sh
   docker-compose up -d or make startDatabase
  ```

## Environment Setup

1. Clone the repository:
    ```sh
    git clone https://github.com/axattechapp/client-go-task.git
    cd client-go-task
    ```

2. configure your `.env` file in `pkg/common/envs/`:
    ```sh
    pkg/common/envs/.env
    ```
    please config your email setting
    ```bash
    EMAIL_USERNAME=your_email_username
    EMAIL_PASSWORD=your_email_password
    EMAIL_SENDER=your_email_sender_address



## Commands

### Start Database

To start the database using Docker Compose, run:
```sh
make startDatabase
```

### Stop Database

To stop the database using Docker Compose, run:
```sh
make stopDatabase
```

### Database Migration

#### Migrate Up

To apply all up migrations, run:
```sh
make migrateUP
```

#### Migrate Down

To apply all down migrations, run:
```sh
make migrateDOWN
```

#### Database Code Generation (sqlc generate)
To generating Go code based on your SQL schema definitions
```sh
make sqlcGenerate
```

#### GraphQL Schema Generation (gqlGenerate)
To generate Go code for your GraphQL schema
```sh
make gqlGenerate
```

### Run Server

To start the server, run:

```sh
make server
```

### Run Tests

To run the tests, use:

```sh
make test
```
### API Testing with Postman
You can test the API endpoints using Postman. Import the Postman collection using the link below:
- The API server running locally
- [Postman](https://www.postman.com/orange-capsule-968775/workspace/golang-client-task)
- [GraphQL Playground](http://127.0.0.1:8080/graphql/)
- [PGAdmin](http://127.0.0.1:5050/login?next=/)
-pgadmin login details are given in .env file, and can be configured.
