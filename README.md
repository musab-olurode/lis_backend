# University Of Ilorin Faculty Of Law (UILSS) Backend

## Description

This is the backend for the University of Ilorin Faculty of Law Student Society (UILSS) website. It is written in Go and uses PostgreSQL as the database.

## Prerequisites

The following tools are required to set up this project:

- Go
- PostgreSQL
- sqlc
- Goose
- Make

## Installation

### Go

Download and install Go from the [official website](https://golang.org/dl/).

### PostgreSQL

Download and install PostgreSQL from the [official website](https://www.postgresql.org/download/).

### sqlc

Install sqlc by running the following command:

```sh
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

### Goose

Install Goose by running the following command:

```sh
go install github.com/pressly/goose/cmd/goose
```

### Make

Make is a build automation tool that can be installed from the package manager of your OS. For example,

on Ubuntu you can install it with:

```sh
sudo apt-get install make
```

on Windows you can install it with:

```sh
choco install make
```

## Setup

1. Clone the repository:

```sh
git clone https://github.com/musab-olurode/lis_backend.git
```

2. Navigate to the project directory:

```sh
cd lis_backend
```

3. Install the Go dependencies:

```sh
go mod tidy
```

4. Environment variables:

- Create a `.env` file in the root directory of the project.
- Copy the contents of `.env.example` into `.env`.
- Fill in the right values for the environment variables in `.env` (especially the database url and JWT secret).

5. Run the migrations:

```sh
make migrate
```

## Running the Project

To run the project, use the `run` command:

```sh
make run
```

## Other Make Commands

- `make tidy`: Cleans up the Go modules.
- `make migrate`: Runs the database migrations.
- `make migrate-down`: Rolls back the last database migration.
- `make sqlc`: Generates Go code from SQL.
