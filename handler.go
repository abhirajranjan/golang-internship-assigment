package main

import (
	"fmt"
	"net/http"
	"score/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// global var DB for database repo

// for consistent error response
type respErr struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type respErrContainer struct {
	Error respErr `json:"error"`
}

func getAllPlayerRankwise(ctx *gin.Context) {
	var (
		players []model.Player
		err     error
	)

	players, err = DB.GetAllPlayerRankwise()

	// no players
	if errors.Is(err, model.ErrNoPlayer) {
		// return empty array of players
		ctx.JSON(http.StatusOK, players)
	} else if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "getAllPlayerRankwise"))
		return
	}

	ctx.JSON(http.StatusOK, players)
}

// param: val
func getPlayerByRank(ctx *gin.Context) {
	// extract rank from param val
	rank_string := ctx.Param("val")
	rank, err := strconv.Atoi(rank_string)
	if err != nil {
		// cannot convert rank to integer
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// check for player
	player, err := DB.GetPlayerByRank(rank)
	if errors.Is(err, model.ErrRankDoensnotExist) {
		// no player in rank
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	} else if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusOK, player)
}

func getRandomPlayer(ctx *gin.Context) {
	player, err := DB.GetRandomPlayer()
	if errors.Is(err, model.ErrNoPlayer) {
		// not an error if no player; empty response
		ctx.Status(http.StatusNoContent)
	} else if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusOK, player)
}

// param: id
func deletePlayer(ctx *gin.Context) {
	// get the param id
	id_string := ctx.Param("id")
	id, err := strconv.Atoi(id_string)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}

	err = DB.DeletePlayer(id)
	if errors.Is(err, model.ErrInvalidPlayer) {
		// no user found to delete
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	} else if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}

func createNewPlayer(ctx *gin.Context) {
	// bind the post data to player
	var player model.Player
	if err := ctx.Bind(&player); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// generate a new player and set id for player
	if err := DB.CreateNewPlayer(&player); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	// specifies where new player can be located
	ctx.Header("Location", fmt.Sprintf("%s/players/%d", ctx.Request.Host, player.Id))
	ctx.JSON(http.StatusCreated, player)
}

func updatePlayer(ctx *gin.Context) {
	// bind the update request
	var player model.Player
	if err := ctx.Bind(&player); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// updating country is not allowed
	if player.Country != "" {
		ctx.JSON(http.StatusForbidden, apiError("domain error", "country cannot be modified"))
	}
	// player id cannot be 0 for update
	if player.Id == 0 {
		ctx.JSON(http.StatusForbidden, apiError("domain error", "invalid id"))
	}

	err := DB.UpdatePlayer(&player)
	// player id requested to update cannot be found
	if errors.Is(err, model.ErrInvalidPlayer) {
		ctx.AbortWithStatus(http.StatusNotFound)
	} else if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusOK, player)
}

func apiError(_type, message string) respErrContainer {
	return respErrContainer{
		Error: respErr{
			Type:    _type,
			Message: message,
		},
	}
}
