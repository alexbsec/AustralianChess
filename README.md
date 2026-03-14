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
