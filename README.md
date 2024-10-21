# Overview

Application is an API library management system.

# Installation

Open .env file in code editor and copy the code below:

```
# Database
DB_USER=user
DB_PASSWORD=password
DB_HOST=127.0.0.1
DB_PORT=5432
DB_NAME=postgres

# redis
REDIS_ADDRESS=localhost:6379
```

Open terminal and install migrate CLI tool:

```
brew install golang-migrate
```

# Database and Redis setup

Open this project in code editor and copy the code below to setup database and redis:

```
docker-compose up --build 
```

# Migrate to databse

Check if postgres and redis are running and then open code in code editor and copy the code below to migrate database:

```
make migrate-up
```

# Run App

To run app copy code below:
```
make run
```

# Endpoints

### User

 - Request: `GET /users`
 
 - Request: `POST /register`

### Book

 - Request: `GET /books`

 - Request: `POST /book/borrow`

 - Request: `POST /book/return`



