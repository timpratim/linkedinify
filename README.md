# LinkedInify API

LinkedInify is an API starter kit designed to transform regular text into professional LinkedIn-style posts using AI.

## API Versioning

The API now uses versioning with all endpoints available under the `/api/v1/` prefix.

### Base URL

```bash
http://localhost:8080/api/v1
```

## Endpoints

### Authentication

- **Register a new user**
  - `POST /api/v1/auth/register`
  - Request body: `{ "email": "user@example.com", "password": "password" }`
  - Response: `{ "token": "jwt-token" }`

- **Login**
  - `POST /api/v1/auth/login`
  - Request body: `{ "email": "user@example.com", "password": "password" }`
  - Response: `{ "token": "jwt-token" }`

### LinkedInify

All LinkedInify endpoints require authentication via JWT token in the Authorization header.

- **Transform text to LinkedIn style**
  - `POST /api/v1/linkedinify`
  - Headers: `Authorization: Bearer your-jwt-token`
  - Request body: `{ "text": "Your text to transform" }`
  - Response: `{ "post": "Transformed LinkedIn-style text" }`

- **View transformation history**
  - `GET /api/v1/linkedinify`
  - Headers: `Authorization: Bearer your-jwt-token`
  - Response: Array of transformation items

## Environment Variables

The application requires the following environment variables:

- `DATABASE_DSN` - PostgreSQL connection string (default: "postgres://pratimbhosale@localhost:5432/linkedinify?sslmode=disable")
- `JWT_SECRET` - Secret for JWT token generation
- `OPENAI_TOKEN` - OpenAI API token for text transformation
- `TREBLLE_SDK_TOKEN` and `TREBLLE_API_KEY` - For API monitoring (optional)

## Running the Application

The backend runs on port 8080 by default.

```bash
go run cmd/api/main.go
```

## Testing with curl

### Register a new user
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

### Transform text (replace YOUR_TOKEN with the token from login/register)
```bash
curl -X POST http://localhost:8080/api/v1/linkedinify \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"text":"Here is my regular text that I want to transform into a LinkedIn post."}'
```

### Get history (replace YOUR_TOKEN with the token from login/register)
```bash
curl -X GET http://localhost:8080/api/v1/linkedinify \
  -H "Authorization: Bearer YOUR_TOKEN"
```
