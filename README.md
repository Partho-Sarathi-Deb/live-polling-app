# Live Polling App

A real-time polling tool — create a poll, share the link, and watch votes update live as they come in, no refresh needed.

**Live app:** https://live-polling-app-sable.vercel.app
**Backend API:** https://live-polling-app-kaxm.onrender.com

## Tech Stack

- Frontend: React (Vite)
- Backend: Go (Gin)
- Database: MongoDB (Atlas)
- Real-time: Redis (Redis Cloud) — pub/sub + atomic counters

## How it works

- MongoDB stores users, polls, and individual vote records (source of truth)
- Redis holds live vote counts per poll (`HINCRBY` for atomic increments — safe under concurrent votes) and a pub/sub channel per poll
- When someone votes: the vote is written to Mongo, the Redis counter is incremented, and the new counts are published to that poll's channel
- Each connected browser has an open WebSocket subscribed to that channel, so updates push instantly with no polling or refresh

## Running locally

### Backend
```
cd backend
go mod tidy
```
Create a `.env` file in `backend/`:
```
MONGO_URI=your_mongodb_connection_string
JWT_SECRET=any_random_string
REDIS_URL=your_redis_connection_string
```
```
go run .
```
Runs on `http://localhost:8080`

### Frontend
```
cd frontend
npm install
```
Create a `.env` file in `frontend/`:
```
VITE_API_URL=http://localhost:8080
VITE_WS_URL=ws://localhost:8080
```
```
npm run dev
```
Runs on `http://localhost:5173`

## Key decisions

[Write 3-5 sentences here, in your own words: why Redis alongside Mongo instead of just Mongo counters, why JWT for auth, anything you'd do differently with more time. This is the part interviewers will actually ask about — make sure it's true to how you understand it.]

## Deployment

- Backend: Render (Go web service)
- Frontend: Vercel (static Vite build)
- Database: MongoDB Atlas (free tier)
- Redis: Redis Cloud (free tier)