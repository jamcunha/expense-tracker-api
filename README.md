# Expense Tracking API

## Overview

The Expense Tracking API is a simple API that allows users to track their expenses.
It allows users to perform CRUD operations on their expenses, categories and budgets.
This API can be used to build a web or mobile application to help users track their expenses.

The API was written in Golang, using the standard http/net package, PostgreSQL as the database, and JWT tokens for authentication.

## Features

- **User Management:** users can register and login
- **Category Management:** users can create, read, update and delete categories for their expenses
- **Expense Management:** users can create, read, update and delete expenses, and associate them with categories
- **Budget Management:** users can create, read, update and delete budgets, and associate them with categories.
The API will return the total amount spent in a category in a given interval, which can be compared with the given budget.
- **Security:** the API uses JWT tokens to authenticate users

## API Endpoints

- **Base Path:** `/api/v1`

### Health

- **Health Check:**
    - **Endpoint:** `/health`
    - **Method:** `GET`
    - **Description:** Check if the API is running
    - **Request Body:** `None`
    - **Successful Response:**
        ```json
        {
            "status": "ok"
        }
        ```

### User

- **Get User Info:**
    - **Endpoint:** `/user/{id}`
    - **Method:** `GET`
    - **Description:** Get the user's information
    - **Header:** `Authorization: Bearer <access_token>`
    - **Request Body:** `None`
    - **Response:**
        ```json
        {
            "id": "527fef18-e8f9-4899-b807-3c9c94415b31",
            "created_at": "2021-07-25T20:00:00Z",
            "updated_at": "2021-07-25T20:00:00Z",
            "name": "John Doe",
            "email": "john@doe.com"
        }
        ```
- **Register:**
    - **Endpoint:** `/user`
    - **Method:** `POST`
    - **Description:** Register a new user
    - **Request Body:**
        ```json
        {
            "name": "John Doe",
            "email": "john@doe.com",
            "password": "password"
        }
        ```
    - **Successful Response:**
        ```json
        {
            "id": "527fef18-e8f9-4899-b807-3c9c94415b31",
            "created_at": "2021-07-25T20:00:00Z",
            "updated_at": "2021-07-25T20:00:00Z",
            "name": "John Doe",
            "email": "john@doe.com"
        }
        ```

- **Delete User:**
    - **Endpoint:** `/user/{id}`
    - **Method:** `DELETE`
    - **Description:** Delete a user
    - **Header:** `Authorization: Bearer <access_token>`
    - **Request Body:** `None`
    - **Successful Response:** 
        ```json
        {
            "id": "527fef18-e8f9-4899-b807-3c9c94415b31",
            "created_at": "2021-07-25T20:00:00Z",
            "updated_at": "2021-07-25T20:00:00Z",
            "name": "John Doe",
            "email": "john@doe.com"
        }
        ```

### Token

- **Get tokens:**
    - **Endpoint:** `/token`
    - **Method:** `POST`
    - **Description:** Get the access and refresh JWT tokens
    - **Request Body:**
        ```json
        {
            "email": "john@doe.com",
            "password": "password"
        }
        ```
    - **Successful Response:**
        ```json
        {
            "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
            "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
        }
        ```

- **Refresh token:**
    - **Endpoint:** `/token/refresh`
    - **Method:** `POST`
    - **Description:** Refresh the access token
    - **Request Body:**
        ```json
        {
            "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
        }
        ```
    - **Successful Response:**
        ```json
        {
            "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
        }
        ```

### Category

> [!NOTE]
> All Endpoints require a valid JWT token in the Authorization header
> Example: `Authorization: Bearer <token>

- **Get Categories:**
    - **Endpoint:** `/category`
    - **Method:** `GET`
    - **Description:** Get all categories
    - **Request Body:** `None`
    - **Successful Response:**
        ```json
        [
            {
                "id": "527fef18-e8f9-4899-b807-3c9c94415b31",
                "created_at": "2021-07-25T20:00:00Z",
                "updated_at": "2021-07-25T20:00:00Z",
                "name": "Food",
                "user_id": "527fef18-e8f9-4899-b807-3c9c94415b31"
            },
            {
                "id": "527fef18-e8f9-4899-b807-3c9c94415b32",
                "created_at": "2021-07-25T20:00:00Z",
                "updated_at": "2021-07-25T20:00:00Z",
                "name": "Transport",
                "user_id": "527fef18-e8f9-4899-b807-3c9c94415b31"
            }
        ]
        ```

- **Get Category by ID:**
    - **Endpoint:** `/category/{id}`
    - **Method:** `GET`
    - **Description:** Get a category by ID
    - **Request Body:** `None`
    - **Successful Response:**
        ```json
        {
            "id": "527fef18-e8f9-4899-b807-3c9c94415b31",
            "created_at": "2021-07-25T20:00:00Z",
            "updated_at": "2021-07-25T20:00:00Z",
            "name": "Food",
            "user_id": "527fef18-e8f9-4899-b807-3c9c94415b31"
        }
        ```

- **Create Category:**
    - **Endpoint:** `/category`
    - **Method:** `POST`
    - **Description:** Create a new category
    - **Request Body:**
        ```json
        {
            "name": "Food"
        }
        ```
    - **Successful Response:**
        ```json
        {
            "id": "527fef18-e8f9-4899-b807-3c9c94415b31",
            "created_at": "2021-07-25T20:00:00Z",
            "updated_at": "2021-07-25T20:00:00Z",
            "name": "Food",
            "user_id": "527fef18-e8f9-4899-b807-3c9c94415b31"
        }
        ```

- **Update Category:**
    - **Endpoint:** `/category/{id}`
    - **Method:** `PUT`
    - **Description:** Update a category
    - **Request Body:**
        ```json
        {
            "name": "Books"
        }
        ```
    - **Successful Response:**
        ```json
        {
            "id": "527fef18-e8f9-4899-b807-3c9c94415b31",
            "created_at": "2021-07-25T20:00:00Z",
            "updated_at": "2021-07-25T20:00:00Z",
            "name": "Books",
            "user_id": "527fef18-e8f9-4899-b807-3c9c94415b31"
        }
        ```

- **Delete Category:**
    - **Endpoint:** `/category/{id}`
    - **Method:** `DELETE`
    - **Description:** Delete a category
    - **Request Body:** `None`
    - **Successful Response:**
        ```json
        {
            "id": "527fef18-e8f9-4899-b807-3c9c94415b31",
            "created_at": "2021-07-25T20:00:00Z",
            "updated_at": "2021-07-25T20:00:00Z",
            "name": "Books",
            "user_id": "527fef18-e8f9-4899-b807-3c9c94415b31"
        }
        ```

### Expense

> [!NOTE]
> All Endpoints require a valid JWT token in the Authorization header
> Example: `Authorization: Bearer <token>

- **Get Expenses:**
    - **Endpoint:** `/expense`
    - **Method:** `GET`
    - **Description:** Get all expenses
    - **Request Body:** `None`
    - **Successful Response:**
        ```json
        [
            {
                "id": "527fef18-e8f9-4899-b807-3c9c94415b31",
                "created_at": "2021-07-25T20:00:00Z",
                "updated_at": "2021-07-25T20:00:00Z",
                "amount": 10.0,
                "description": "Lunch",
                "category_id": "527fef18-e8f9-4899-b807-3c9c94415b31",
                "user_id": "527fef18-e8f9-4899-b807-3c9c94415b31"
            },
            {
                "id": "527fef18-e8f9-4899-b807-3c9c94415b32",
                "created_at": "2021-07-25T20:00:00Z",
                "updated_at": "2021-07-25T20:00:00Z",
                "amount": 5.0,
                "description": "Bus ticket",
                "category_id": "527fef18-e8f9-4899-b807-3c9c94415b32",
                "user_id": "527fef18-e8f9-4899-b807-3c9c94415b31"
            }
        ]
        ```

- **Get Expense by ID:**
    - **Endpoint:** `/expense/{id}`
    - **Method:** `GET`
    - **Description:** Get an expense by ID
    - **Request Body:** `None`
    - **Successful Response:**
        ```json
        {
            "id": "527fef18-e8f9-4899-b807-3c9c94415b31",
            "created_at": "2021-07-25T20:00:00Z",
            "updated_at": "2021-07-25T20:00:00Z",
            "amount": 10.0,
            "description": "Lunch",
            "category_id": "527fef18-e8f9-4899-b807-3c9c94415b31",
            "user_id": "527fef18-e8f9-4899-b807-3c9c94415b31"
        }
        ```

- **Get Expense by Category:**
    - **Endpoint:** `/expense/category/{id}`
    - **Method:** `GET`
    - **Description:** Get all expenses in a category
    - **Request Body:** `None`
    - **Successful Response:**
        ```json
        [
            {
                "id": "527fef18-e8f9-4899-b807-3c9c94415b31",
                "created_at": "2021-07-25T20:00:00Z",
                "updated_at": "2021-07-25T20:00:00Z",
                "amount": 10.0,
                "description": "Lunch",
                "category_id": "527fef18-e8f9-4899-b807-3c9c94415b31",
                "user_id": "527fef18-e8f9-4899-b807-3c9c94415b31"
            },
            {
                "id": "527fef18-e8f9-4899-b807-3c9c94415b32",
                "created_at": "2021-07-25T20:00:00Z",
                "updated_at": "2021-07-25T20:00:00Z",
                "amount": 5.75,
                "description": "Dinner",
                "category_id": "527fef18-e8f9-4899-b807-3c9c94415b31",
                "user_id": "527fef18-e8f9-4899-b807-3c9c94415b31"
            }
        ]
        ```

- **Create Expense:**
    - **Endpoint:** `/expense`
    - **Method:** `POST`
    - **Description:** Create a new expense
    - **Request Body:**
        ```json
        {
            "amount": 10.0,
            "description": "Lunch",
            "category_id": "527fef18-e8f9-4899-b807-3c9c94415b31"
        }
        ```
    - **Successful Response:**
        ```json
        {
            "id": "527fef18-e8f9-4899-b807-3c9c94415b31",
            "created_at": "2021-07-25T20:00:00Z",
            "updated_at": "2021-07-25T20:00:00Z",
            "amount": 10.0,
            "description": "Lunch",
            "category_id": "527fef18-e8f9-4899-b807-3c9c94415b31",
            "user_id": "527fef18-e8f9-4899-b807-3c9c94415b31"
        }
        ```

- **Update Expense:**
    - **Endpoint:** `/expense/{id}`
    - **Method:** `PUT`
    - **Description:** Update an expense
    - **Request Body: (optional)**
        ```json
        {
            "amount": 15.0,
            "description": "Dinner",
            "category_id": "527fef18-e8f9-4899-b807-3c9c94415b31"
        }
        ```
    - **Successful Response:**
        ```json
        {
            "id": "527fef18-e8f9-4899-b807-3c9c94415b31",
            "created_at": "2021-07-25T20:00:00Z",
            "updated_at": "2021-07-25T20:00:00Z",
            "amount": 15.0,
            "description": "Dinner",
            "category_id": "527fef18-e8f9-4899-b807-3c9c94415b31",
            "user_id": "527fef18-e8f9-4899-b807-3c9c94415b31"
        }
        ```

- **Delete Expense:**
    - **Endpoint:** `/expense/{id}`
    - **Method:** `DELETE`
    - **Description:** Delete an expense
    - **Request Body:** `None`
    - **Successful Response:**
        ```json
        {
            "id": "527fef18-e8f9-4899-b807-3c9c94415b31",
            "created_at": "2021-07-25T20:00:00Z",
            "updated_at": "2021-07-25T20:00:00Z",
            "amount": 15.0,
            "description": "Dinner",
            "category_id": "527fef18-e8f9-4899-b807-3c9c94415b31",
            "user_id": "527fef18-e8f9-4899-b807-3c9c94415b31"
        }
        ```

### Budget

> [!NOTE]
> All Endpoints require a valid JWT token in the Authorization header
> Example: `Authorization: Bearer <token>

- **Get Budgets:**
    - **Endpoint:** `/budget`
    - **Method:** `GET`
    - **Description:** Get all budgets
    - **Request Body:** `None`
    - **Successful Response:**
        ```json
        [
            {
                "id": "527fef18-e8f9-4899-b807-3c9c94415b31",
                "created_at": "2021-07-25T20:00:00Z",
                "updated_at": "2021-07-25T20:00:00Z",
                "amount": 100.0,
                "goal": 450.0,
                "start_date": "2021-07-01T00:00:00Z",
                "end_date": "2021-07-31T23:59:59Z",
                "user_id": "527fef18-e8f9-4899-b807-3c9c94415b31"
                "category_id": "527fef18-e8f9-4899-b807-3c9c94415b31",
            },
            {
                "id": "527fef18-e8f9-4899-b807-3c9c94415b32",
                "created_at": "2021-07-25T20:00:00Z",
                "updated_at": "2021-07-25T20:00:00Z",
                "amount": 50.0,
                "goal": 200.0,
                "start_date": "2021-07-01T00:00:00Z",
                "end_date": "2021-07-31T23:59:59Z",
                "user_id": "527fef18-e8f9-4899-b807-3c9c94415b31"
                "category_id": "527fef18-e8f9-4899-b807-3c9c94415b32",
            }
        ]
        ```

- **Get Budget by ID:**
    - **Endpoint:** `/budget/{id}`
    - **Method:** `GET`
    - **Description:** Get a budget by ID
    - **Request Body:** `None`
    - **Successful Response:**
        ```json
        {
            "id": "527fef18-e8f9-4899-b807-3c9c94415b31",
            "created_at": "2021-07-25T20:00:00Z",
            "updated_at": "2021-07-25T20:00:00Z",
            "amount": 100.0,
            "goal": 450.0,
            "start_date": "2021-07-01T00:00:00Z",
            "end_date": "2021-07-31T23:59:59Z",
            "user_id": "527fef18-e8f9-4899-b807-3c9c94415b31",
            "category_id": "527fef18-e8f9-4899-b807-3c9c94415b31",
        }
        ```

- **Create Budget:**
    - **Endpoint:** `/budget`
    - **Method:** `POST`
    - **Description:** Create a new budget
    - **Request Body:**
        ```json
        {
            "amount": 100.0,
            "goal": 450.0,
            "start_date": "2021-07-01T00:00:00Z",
            "end_date": "2021-07-31T23:59:59Z",
            "category_id": "527fef18-e8f9-4899-b807-3c9c94415b31"
        }
        ```
    - **Successful Response:**
        ```json
        {
            "id": "527fef18-e8f9-4899-b807-3c9c94415b31",
            "created_at": "2021-07-25T20:00:00Z",
            "updated_at": "2021-07-25T20:00:00Z",
            "amount": 100.0,
            "goal": 450.0,
            "start_date": "2021-07-01T00:00:00Z",
            "end_date": "2021-07-31T23:59:59Z",
            "user_id": "527fef18-e8f9-4899-b807-3c9c94415b31",
            "category_id": "527fef18-e8f9-4899-b807-3c9c94415b31",
        }
        ```

- **Delete Budget:**
    - **Endpoint:** `/budget/{id}`
    - **Method:** `DELETE`
    - **Description:** Delete a budget
    - **Request Body:** `None`
    - **Successful Response:**
        ```json
        {
            "id": "527fef18-e8f9-4899-b807-3c9c94415b31",
            "created_at": "2021-07-25T20:00:00Z",
            "updated_at": "2021-07-25T20:00:00Z",
            "amount": 100.0,
            "goal": 450.0,
            "start_date": "2021-07-01T00:00:00Z",
            "end_date": "2021-07-31T23:59:59Z",
            "user_id": "527fef18-e8f9-4899-b807-3c9c94415b31",
            "category_id": "527fef18-e8f9-4899-b807-3c9c94415b31",
        }
        ```

## Installation

>[!NOTE]
> If you are using Nix, just do `nix develop`.
> If you have `direnv` installed, you can just do `direnv allow`.
> Docker is still needed.

- **Prerequisites:**
    - [Go](https://golang.org/doc/install)
    - [Docker](https://docs.docker.com/get-docker/)
    - [Docker Compose](https://docs.docker.com/compose/install/)
    - [Goose](https://github.com/pressly/goose)
    - [sqlc](https://sqlc.dev)

- **Clone the repository:**
    ```bash
    git clone git@github.com:jamcunha/expense-tracking-api.git
    cd expense-tracking-api
    ```

- **Set up environment variables:**
    ```bash
    cp .env.example .env
    # Edit the .env file and set the environment variables
    ```

- **Set up the database:**
    ```bash
    docker compose up -d
    ```

- **Add the database schema:**
    ```bash
    make migration/up
    ```

- **Generate Sqlc Queries:**
    ```bash
    sqlc generate
    ```

- **Run the API:**
    ```bash
    make run
    ```

## Environment Variables

The following environment variables are required to run the API:

- **PORT:** the port the API will run on
- **DB_URL:** the URL to the PostgreSQL database
- **JWT_SECRET:** the secret used to sign the JWT tokens
- **JWT_EXPIRATION:** the expiration time for the JWT tokens in seconds

A [`.env.example`](./.env.example) file is provided.

## TODOs and Improvements

- write tests for API
- uniformize the errors log (and what to send to the client)
- add some kind of attempt limit to the login (use redis to store attempts and block the user)
- filter expenses by time interval
