# Project Documentation

## Overview
This documentation provides information about the project and its components.

## Table of Contents
- [Installation](#installation)
- [Usage](#usage)
- [Running with Docker Compose](#running-with-docker-compose)
- [Configuration](#configuration)

## Installation
To install and run the project, follow these steps:

1. Clone the repository: `git clone github.com/tz3/go-simple`
2. Change to the project directory: `cd go-simple`
3. Install dependencies: `go mod download`
4. Change your sql connection url
5. Start the project: ` expoexport DB_CONNECTION_STRING={{YOUR_CONNECTION_STRING}}
   go run cmd/main.go`

## Usage
Once the project is running, you can access it at `http://localhost:8080`. Use the provided APIs to interact with the application.

## Running with Docker Compose
To run the project using Docker Compose, make sure you have Docker and Docker Compose installed on your machine. Then, follow these steps:

1. Build the Docker image: `docker-compose build`
2. Start the project using `Docker Compose: docker-compose up`
3. To Use end points with default port example:- `GET http://localhost:8080/users?userID=1`

## Configuration
The project uses environment variables for configuration. The following variables can be set:

- `DB_CONNECTION_STRING`: Connection string for the PostgreSQL database.

## API Reference
The project exposes the following API endpoints:

- `GET /users?userID=1`: Get a user by ID.
- `POST /users`: Create a new user. 

To create a new user, make a `POST` request to `/users` with the following JSON body:

```json
{
 "firstName": "Moutaz",
 "lastName": "Chaara"   
}