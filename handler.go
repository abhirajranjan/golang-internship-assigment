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

func GetAllPlayerRankwise(ctx *gin.Context, DB model.Repo) {
	var (
		players []model.Player
		err     error
	)

	players, err = DB.GetAllPlayerRankwise()

	// no players
	if errors.Is(err, model.ErrNoPlayer) {
		// return empty array of players
		ctx.JSON(http.StatusOK, players)
		return
	} else if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "getAllPlayerRankwise"))
		return
	}

	ctx.JSON(http.StatusOK, players)
}

// param: val
func GetPlayerByRank(ctx *gin.Context, DB model.Repo) {
	// extract rank from param val
	rank_string := ctx.Param("val")
	rank, err := strconv.Atoi(rank_string)
	if err != nil {
		// cannot convert rank to integer
		ctx.AbortWithStatusJSON(http.StatusBadRequest, apiError("domain error", "rank should be integer"))
		return
	}

	if rank <= 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, apiError("domain error", "rank cannot be 0 or negetive"))
		return
	}

	// check for player
	player, err := DB.GetPlayerByRank(rank)
	if errors.Is(err, model.ErrRankDoensnotExist) {
		// no player in rank
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	} else if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "getPlayerByRank"))
		return
	}

	ctx.JSON(http.StatusOK, player)
}

func GetRandomPlayer(ctx *gin.Context, DB model.Repo) {
	player, err := DB.GetRandomPlayer()
	if errors.Is(err, model.ErrNoPlayer) {
		// not an error if no player; empty response
		ctx.Status(http.StatusNoContent)
		return
	} else if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "getRandomPlayer"))
		return
	}

	ctx.JSON(http.StatusOK, player)
}

// param: id
func DeletePlayer(ctx *gin.Context, DB model.Repo) {
	// get the param id
	id_string := ctx.Param("id")
	id, err := strconv.Atoi(id_string)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// id cannot be 0 as it is default for unset
	if id <= 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, apiError("domain error", "invalid id"))
		return
	}

	err = DB.DeletePlayer(id)
	if errors.Is(err, model.ErrInvalidPlayer) {
		// no user found to delete
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	} else if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "deletePlayer"))
		return
	}

	ctx.Status(http.StatusOK)
}

func CreateNewPlayer(ctx *gin.Context, DB model.Repo) {
	// bind the post data to player
	var player model.Player
	if err := ctx.Bind(&player); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, errors.Wrap(err, "createNewPlayer"))
		return
	}

	// validates model except id parameter
	if err := model.Validate(player); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, apiError("domain error", err.Error()))
		return
	}

	// generate a new player and set id for player
	if err := DB.CreateNewPlayer(&player); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// specifies where new player can be located
	ctx.Header("Location", fmt.Sprintf("%s/players/%d", ctx.Request.Host, player.Id))
	ctx.JSON(http.StatusCreated, player)
}

// param: id
func UpdatePlayer(ctx *gin.Context, DB model.Repo) {
	var (
		id    int
		name  string
		score int
	)

	// player id should be present
	_id := ctx.Param("id")
	if _id == "" {
		ctx.JSON(http.StatusBadRequest, apiError("domain error", "missing id"))
		return
	}

	// player id should be integer
	id, err := strconv.Atoi(_id)
	if err != nil {
		ctx.JSON(http.StatusForbidden, apiError("domain error", "invalid id"))
		return
	}

	// player id should be not be <=0
	if id <= 0 {
		ctx.JSON(http.StatusForbidden, apiError("domain error", "invalid id"))
		return
	}

	// bind the update request
	var player map[string]interface{}
	if err := ctx.Bind(&player); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, errors.Wrap(err, "updatePlayer: ctx.Bind"))
		return
	}

	// updating country is not allowed
	if _, ok := player["country"]; ok {
		ctx.JSON(http.StatusForbidden, apiError("domain error", "country cannot be modified"))
		return
	}

	// check if score exists
	s, ok := player["score"]
	if ok {
		ss, ok := s.(float64)
		// score is not integer type
		if !ok {
			ctx.JSON(http.StatusForbidden, apiError("domain error", "invalid score"))
			return
		}
		score = int(ss)
	}

	// check if name exists
	n, ok := player["name"]
	if ok {
		nn, ok := n.(string)
		// name is not string type
		if !ok {
			ctx.JSON(http.StatusForbidden, apiError("domain error", "invalid name"))
			return
		}
		name = nn
	}

	resplayer, err := DB.UpdatePlayer(id, name, score, player)
	// player id requested to update cannot be found
	if errors.Is(err, model.ErrInvalidPlayer) {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	} else if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "updatePlayer"))
		return
	}

	ctx.JSON(http.StatusOK, resplayer)
}

func apiError(_type, message string) respErrContainer {
	return respErrContainer{
		Error: respErr{
			Type:    _type,
			Message: message,
		},
	}
}
