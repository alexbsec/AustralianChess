# Australian Chess

A multiplayer, real-time strategy chess variant played on an expanded **12x12 board**. This project features a high-performance Go backend and a modern TypeScript frontend, coordinated via WebSockets for seamless gameplay.



## 🇦🇺 Features
- **12x12 Expanded Board**: New strategic depths beyond the traditional 8x8.
- **Custom Pieces**: Includes unique movements like the **Kangaroo** and the **Oligarch**.
- **Real-time Multiplayer**: Powered by WebSockets for instant move synchronization.
- **Modern Tech Stack**: Go (Gin), TypeScript (Vite), PostgreSQL, and Nginx.

---

## 🏗 Project Structure

```text
.
├── backend/               # Go (Golang) REST API & WebSocket Server
│   ├── internal/chess/    # Core game engine logic
│   ├── migrations/        # Database schema (Goose)
│   └── ws/                # WebSocket orchestration
└── frontend/              # TypeScript / Vite / Nginx
    ├── src/engine/        # Shared game rules (ported to TS)
    ├── src/room/          # Room & Socket controllers
    └── public/pieces/     # Asset inventory
```

# Setup & Installation
## 1. Configure Environment Variables
### Copy the sample file to both the root and backend directories
```bash
cp backend/.env.sample .env
cp backend/.env.sample backend/.env
```

## 2. Clean Start with Docker
### This ensures any previous zombie containers or failed builds are cleared
```bash
docker-compose down --remove-orphans
```

## 3. Build and Launch
### Use --no-cache to ensure the 'npm not found' error is fully resolved
```bash
docker-compose build --no-cache
docker-compose up -d
```

## 4. Verify Services
### Check that all three containers (db, backend, frontend) are 'Up'
```bash
docker ps
```

# Access Point

### 🌐 Access Points

| Component       | URL                                   | Port (Host) |
| :-------------- | :------------------------------------ | :---------- |
| **Frontend UI** | [http://localhost](http://localhost)  | `80`        |
| **Backend API** | [http://localhost:8080/api/v1](http://localhost:8080/api/v1) | `8080`      |
| **Database** | `localhost:5432`                      | `5432`      |
