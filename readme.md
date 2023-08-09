# player score management application
This project is a part of golang backend internship assignment.

## Usage
run in docker 
```
docker run -p 8080:8080 ghcr.io/abhirajranjan/player-score-management:latest
```

## Endpoints
1. POST http://localhost:8080/players - {"name": string, "country": string, "score": int} – Creates a new entry for a player
2. PUT  http://localhost:8080/players/:id - {"name": string, "score": int} – Updates the player attributes. Only name and
score can be updated
3. DELETE   http://localhost:8080/players/:id – Deletes the player entry
4. GET  http://localhost:8080/players – Displays the list of all players in descending order
5. GET  http://localhost:8080/players/rank/:val – Fetches the player ranked “val”
6. GET  http://localhost:8080/players/random – Fetches a random player
