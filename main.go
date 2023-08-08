package main

import (
	"score/memorydb"

	"github.com/gin-gonic/gin"
)

func main() {
	DB := memorydb.NewMemoryDB()

	router := gin.Default()
	router.POST("/players", func(ctx *gin.Context) {
		CreateNewPlayer(ctx, DB)
	})

	router.PUT("/players/:id", func(ctx *gin.Context) {
		UpdatePlayer(ctx, DB)
	})

	router.DELETE("/players/:id", func(ctx *gin.Context) {
		DeletePlayer(ctx, DB)
	})

	router.GET("/players", func(ctx *gin.Context) {
		GetAllPlayerRankwise(ctx, DB)
	})

	router.GET("/players/rank/:val", func(ctx *gin.Context) {
		GetPlayerByRank(ctx, DB)
	})

	router.GET("/players/random", func(ctx *gin.Context) {
		GetRandomPlayer(ctx, DB)
	})

	router.Run()
}
