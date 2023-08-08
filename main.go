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
	player := router.Group("/player")
	{
		player.POST("/", createNewPlayer)
		player.PUT("/:id", updatePlayer)
		player.DELETE("/:id", deletePlayer)

		player.GET("/", getAllPlayerRankwise)
		player.GET("/rank/:val", getPlayerByRank)
		player.GET("/random", getRandomPlayer)
	}

	router.Run()
}
