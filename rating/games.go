package rating

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"goj/domain"
	"net/http"
)

func getGameEndpoint(ctx *gin.Context) {
	id := ctx.Param("id")
	game, found := domain.GetGame(&id)
	if found {
		ctx.JSON(http.StatusOK, game)
	} else {
		ctx.Status(http.StatusNoContent)
	}
}

func listGamesEndpoint(ctx *gin.Context) {
	list, err := domain.GetGamesShort()
	if err == nil {
		ctx.JSON(http.StatusOK, bson.M{"items": list})
	} else {
		ctx.Status(http.StatusInternalServerError)
	}
}

func removeGameEndpoint(ctx *gin.Context) {

}

func updateGameEndpoint(ctx *gin.Context) {
	var game domain.Game
	err := ctx.BindJSON(&game)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
	}
	err = domain.StoreGame(&game)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
	} else {
		ctx.Status(http.StatusAccepted)
	}
}
