package main

import (
	"score/memorydb"

	"github.com/gin-gonic/gin"
)

func main() {
	DB := memorydb.NewMemoryDB()

	router := gin.Default()
	router.POST("/player", func(ctx *gin.Context) {
		CreateNewPlayer(ctx, DB)
	})

	router.PUT("/player/:id", func(ctx *gin.Context) {
		UpdatePlayer(ctx, DB)
	})

	router.DELETE("/player/:id", func(ctx *gin.Context) {
		DeletePlayer(ctx, DB)
	})

	router.GET("/player", func(ctx *gin.Context) {
		GetAllPlayerRankwise(ctx, DB)
	})

	router.GET("/player/rank/:val", func(ctx *gin.Context) {
		GetPlayerByRank(ctx, DB)
	})

	router.GET("/player/random", func(ctx *gin.Context) {
		GetRandomPlayer(ctx, DB)
	})

	router.Run()
}
