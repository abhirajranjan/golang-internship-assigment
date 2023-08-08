package main

import (
	"score/memorydb"
	"score/model"

	"github.com/gin-gonic/gin"
)

var (
	DB model.Repo
)

func main() {
	DB = memorydb.NewMemoryDB()

	router := gin.Default()
	router.POST("/player", createNewPlayer)
	router.PUT("/player/:id", updatePlayer)
	router.DELETE("/player/:id", deletePlayer)
	router.GET("/player", getAllPlayerRankwise)
	router.GET("/player/rank/:val", getPlayerByRank)
	router.GET("/player/random", getRandomPlayer)

	router.Run()
}
