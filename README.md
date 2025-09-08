# Chirpy - A Simple Social Media Backend

## Overview

Chirpy is a lightweight, Go-based backend service for a Twitter-like social media platform. It provides a RESTful API for managing users, posts (chirps), and authentication.

## Features

- **User Management**: Create, update, and authenticate users
- **Chirp Handling**: Create, retrieve, and delete short posts
- **Authentication**: JWT-based authentication with refresh tokens
- **Premium Features**: Support for premium user accounts (Chirpy Red)
- **Content Moderation**: Automatic filtering of profane content

## API Endpoints

### Authentication
- `POST /api/login`: Authenticate user and get access tokens
- `POST /api/refresh`: Refresh an access token
- `POST /api/revoke`: Revoke a refresh token (logout)

### Users
- `POST /api/users`: Create a new user
- `PUT /api/users`: Update an existing user
- `POST /api/polka/webhooks`: Upgrade a user to Chirpy Red (webhook)

### Chirps
- `GET /api/chirps`: Get all chirps (with optional author_id and sort parameters)
- `GET /api/chirps/{chirpID}`: Get a specific chirp
- `POST /api/chirps`: Create a new chirp
- `DELETE /api/chirps/{chirpID}`: Delete a chirp (user must be the author)

### Metrics
- `GET /api/metrics`: Get server metrics

## Technology Stack

- **Language**: Go
- **Database**: PostgreSQL
- **Authentication**: JWT with refresh tokens
- **API**: RESTful JSON API

## Getting Started

### Prerequisites
- Go 1.21+
- PostgreSQL

### Setup
1. Clone the repository
2. Create a `.env` file with the following variables:
JWT_SECRET=your-jwt-secret POLKA_KEY=your-polka-key PLATFORM=dev # Use 'prod' for production DATABASE_URL='your-postgres-url'
3. Initialize the database with migrations
4. Build and run the server


## Development

### Database Migrations
- Run migrations up: `./migrate up`
- Roll back migrations: `./migrate down`

### Reset (Development Only)
- `GET /api/reset`: Reset the database and metrics (only works in dev mode)

## Testing

Run the test suite with:
go test ./...

