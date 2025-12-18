# Price is Right - Multiplayer Quiz Game

A real-time multiplayer game inspired by "The Price is Right", where players compete to guess the closest price of various items.

## Features

- **Real-time Multiplayer**: Play with friends instantly using WebSockets.
- **Live Leaderboard**: See round results and current standings after every guess.
- **Game Stages**: Lobby, Guessing Round, Round Results, and Final Game Over screen.
- **Host Controls**: Only the host (first player) can start or reset the game.
- **Dark Mode**: Toggle between light and dark themes.
- **Responsive Design**: Works on desktop and mobile.

## Tech Stack

- **Frontend**: React, Vite, Bootstrap (React-Bootstrap)
- **Backend**: Go (Golang), Gorilla WebSocket
- **Communication**: WebSockets

## Prerequisites

- [Go](https://go.dev/dl/) (1.16 or later)
- [Node.js](https://nodejs.org/) (16 or later)
- [npm](https://www.npmjs.com/)

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/yourusername/price-is-right.git
cd price-is-right
```

### 2. Start the Backend Server

```bash
cd server
go mod tidy
go run main.go
```
The server will start on `http://localhost:8080`.

### 3. Start the Frontend Client

Open a new terminal window:

```bash
cd client
npm install
npm run dev
```
The client will start on `http://localhost:5173`.

## Configuration

You can configure game settings (like timer duration) in `server/config.json`.

```json
{
    "guessingTime": 30,
    "resultTime": 10
}
```

## Deployment

For detailed deployment instructions on an Ubuntu server with Nginx, please refer to [DEPLOY.md](DEPLOY.md).

## License

This project is open source and available under the [MIT License](LICENSE).
