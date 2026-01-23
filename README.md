# 🐦 Chirpy

**Chirpy** is a short message service that allows users to post short messages that are visible to everyone.

Chirpy was originally started as part of the **Learn HTTP Servers in GO - boot.dev** course, where the primary focus was on backend handlers and fundamentals.

Building on that foundation, I continue to develop this project out of personal interest to deepen my skills beyond the course content.
Compared to **v1.0.0**, Chirpy already includes additional functionality that was not part of the boot.dev curriculum and is currently being extended with a full-featured frontend and further improvements.

At the moment, this includes a frontend-based authentication flow with login, registration, and logout using session ID cookies.

---

## Running the Project

### Prerequisites

Make sure the following tools are installed on your system:

* **Go**
* **templ**
* **Tailwind CSS**
* A running **PostgreSQL** instance

### Dependencies

The project uses the following Go dependencies:

* `github.com/a-h/templ`
* `github.com/alexedwards/argon2id`
* `github.com/golang-jwt/jwt/v5`
* `github.com/google/uuid`
* `github.com/joho/godotenv`
* `github.com/lib/pq`

Additional indirect dependencies are managed automatically via Go modules.

---

### Starting the Application

1. Generate templ files:

```bash
templ generate
```

Run this command from the project root directory.

2. Start the Go application:

```bash
go run .
```

After this, the application should be up and running.

---

## Environment Variables

The project uses environment variables for configuration.
An example file is provided as `.env.example`.

### Setup

1. Rename the example file:

```bash
cp .env.example .env
```

2. Adjust the values in `.env` according to your local setup:

```env
DB_URL="postgres://username:password@exampleUrl:5432/databasename?sslmode=disable"
```

### Variable Overview

* **DB_URL**
  PostgreSQL connection string used by the application.

---

## 📌 Notes

* This project is intentionally kept simple.
* The focus is on learning concepts rather than feature completeness.
