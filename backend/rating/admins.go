package rating

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"goj/domain"
	"log"
	"net/http"
	"strconv"
)

func listAdminsEndpoint(ctx *gin.Context) {
	list, err := domain.GetAdmins()
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
	} else {
		ctx.JSON(200, bson.M{"items": list})
	}
}

func removeAdminEndpoint(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	err = domain.RemoveAdmin(id)
	if err != nil {
		log.Print(err)
		ctx.Status(http.StatusInternalServerError)
	} else {
		ctx.Status(http.StatusAccepted)
	}
}

type AddAdminRequest struct {
	Link string `json:"link"`
}

func addAdminEndpoint(ctx *gin.Context) {
	var request AddAdminRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
	} else {
		if err = domain.AddAdmin(request.Link); err == nil {
			ctx.Status(http.StatusAccepted)
		} else {
			log.Print(err)
			ctx.Status(http.StatusInternalServerError)
		}
	}
}
