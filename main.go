package main

import (
	"flag"
	"score/model"

	"github.com/gin-gonic/gin"
)

var (
	// specify if we have redis connection
	// if unset, memory repo will be used
	redisConn string

	DB model.Repo
)

func main() {
	flag.StringVar(&redisConn, "redis dsn", "", "provide redis dsn. if unset memory database will be used")
	flag.Parse()

	DB = getDatabase()

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

func getDatabase() model.Repo {
}
