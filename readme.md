# player score management application
This project is a part of golang backend internship assignment.

## Usage
run in docker 
```
docker run -p 8080:8080 ghcr.io/abhirajranjan/player-score-management:latest
```

## Endpoints
1. POST /players – Creates a new entry for a player
2. PUT  /players/:id – Updates the player attributes. Only name and
score can be updated
3. DELETE   /players/:id – Deletes the player entry
4. GET  /players – Displays the list of all players in descending order
5. GET  /players/rank/:val – Fetches the player ranked “val”
6. GET  /players/random – Fetches a random player
