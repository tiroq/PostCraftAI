# PostCraft AI

PostCraft AI is an intelligent service that transforms long-form articles into concise, engaging LinkedIn posts using OpenAI's API. It features robust JWT-based authentication, an admin panel for user management and request tracking, and per-user rate limiting to control API usage.

## Features

- **Article Transformation:** Convert detailed articles into polished, LinkedIn-ready posts.
- **JWT Authentication:** Secure signup, login, and role-based access (user/admin).
- **Admin Panel:** Enable users with custom access expirations, update rate limits, and view detailed usage statistics.
- **Rate Limiting:** Prevent abuse by limiting OpenAI API calls to a configurable number of requests per minute.
- **Request Logging:** Track all user requests with timestamps for analysis and graphing.
- **Dockerized:** Easily build and deploy both backend (Go) and frontend (React) using Docker Compose.

## Project Structure

```bash
postcraft-ai/
├── backend/
│   ├── Dockerfile
│   ├── go.mod
│   ├── main.go
│   ├── handlers/
│   │   ├── auth.go
│   │   ├── admin.go
│   │   └── generate.go
│   └── models/
│       └── user.go
├── frontend/
│   ├── Dockerfile
│   ├── package.json
│   ├── public/
│   │   └── index.html
│   └── src/
│       ├── index.js
│       ├── App.js
│       ├── global.css
│       └── components/
│           ├── Login.js
│           ├── Signup.js
│           ├── GeneratePost.js
│           ├── AdminDashboard.js
│           └── UserDashboard.js
├── docker-compose.yml
├── Makefile
└── README.md
```

## Prerequisites

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Make](https://www.gnu.org/software/make/) (optional)

**Environment Variables:**

- `OPENAI_API_KEY`: Your OpenAI API key.
- `JWT_SECRET`: A secure secret for signing JWT tokens.
  > **Note:** Generate a secure secret (for example, run `openssl rand -base64 32`) and keep it confidential.

## Installation and Setup

1. **Clone the Repository:**

    ```bash
    git clone https://github.com/yourusername/postcraft-ai.git
    cd postcraft-ai
    ```

2. **Set Environment Variables:**

    Create a .env files or set the environment variables in Docker Compose.

3. **Build and Run:**

    Use the provided Makefile:

    ```bash
    make up
    ```

    Or run Docker Compose directly:

    ```bash
    docker compose up --build
    ```

4. **Access the Application:**

    Frontend: <http://localhost:3000>
    Backend API: <http://localhost:18080>

## Usage

### Endpoints Overview

- **User Endpoints:**

  - `POST /signup`: Register a new user.
  - `POST /login`: Authenticate and receive a JWT token.
  - `POST /generate-post`: Transform an article into a LinkedIn post (requires a valid JWT and enabled access).

- **Admin Endpoints (require admin JWT token):**

  - `POST /admin/enable-user`: Enable a user and set access expiration (in minutes; default 7 days).
  - `POST /admin/update-expiration`: Update an enabled user's expiration.
  - `GET /admin/list-users`: List all registered users.
  - `POST /admin/update-rate-limit`: Update the global rate limit (requests per minute).
  - `GET /admin/request-stats`: Retrieve logs of generate-post requests.

### Frontend

- **User Dashboard:** Displays the JWT token and an example cURL command to call /generate-post. Includes a link to the Generate Post form.
- **Admin Dashboard:** Allows management of users, updating expirations, rate limits, and viewing request logs. Includes a link to the Generate Post form.

### Example cURL Command

After logging in as a user (and obtaining a JWT token), run:

```bash
curl -X POST http://localhost:8080/generate-post \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"article": "Your article content here."}'
```

Replace `YOUR_JWT_TOKEN` with your actual token.

### Contributing

Contributions are welcome! Fork the repository and submit pull requests.

### License

This project is licensed under the MIT License.

Enjoy building and enhancing PostCraft AI!
