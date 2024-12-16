### Clone the Repository

```bash
git clone https://github.com/minhle003/Elsa-test.git
```

# Quiz Application Frontend

This is the frontend for the Quiz Application, built with React.js. The frontend allows users to create quizzes, join quiz sessions, submit answers, and view the leaderboard. It supports real-time updates using Firebase Firestore.

## Features

- Create quizzes with multiple questions and choices.
- Join quiz sessions and participate in real-time.
- Submit answers and view the leaderboard.
- Real-time updates using Websockets.

## Technologies Used

- React.js
- Firebase Firestore
- Axios for API requests
- Bootstrap for styling

## Prerequisites

- Node.js (v14 or higher)
- npm or yarn
- Firebase project with Firestore enabled

## Start the Development Server 

```bash
cd frontend
npm install
npm start
```

# Quiz Application Server

This is the backend server for the Quiz Application, built with Go. The server handles user authentication, quiz management, score calculation, and real-time communication using WebSockets and Firestore.

## Features

- User authentication and profile management.
- Quiz creation, updating, and retrieval.
- Real-time quiz sessions and leaderboard updates.

## Technologies Used

- Go
- Firestore
- WebSocket


## Prerequisites

- Go (v1.16 or higher)
- Firebase project with Firestore enabled

## Start the Development Server 

```bash
cd server
go mod tidy
go run cmd/server/server.go
```


