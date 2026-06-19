package rating

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"goj/domain"
	"net/http"
	"strconv"
)

func listGroupsEndpoint(ctx *gin.Context) {
	items, err := domain.GetGroups()
	if err == nil {
		for i, _ := range items {
			items[i].Secret = ""
			items[i].ApiKey = ""
			items[i].ConfirmCode = ""
		}
		ctx.JSON(http.StatusOK, bson.M{"items": items})
	} else {
		ctx.Status(http.StatusInternalServerError)
	}
}

func removeGroupEndpoint(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	err = domain.RemoveGroup(id)
	if err == nil {
		ctx.Status(http.StatusAccepted)
	} else {
		ctx.Status(http.StatusInternalServerError)
	}
}

type AddGroupRequest struct {
	ApiKey string `json:"apiKey" binding:"required"`
}

func addGroupEndpoint(ctx *gin.Context) {
	var request AddGroupRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	err = domain.AddGroup(request.ApiKey)
	if err == nil {
		ctx.Status(http.StatusAccepted)
	} else {
		ctx.Status(http.StatusBadRequest)
	}
}
