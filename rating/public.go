package rating

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"goj/domain"
	"net/http"
	"net/url"
	"strings"
)

type GameHeader struct {
	Id   string
	Name string
}

type ProfileDTO struct {
	IsAdmin bool `json:"isAdmin"`
	Games   []*GameHeader
}

func meEndpoint(ctx *gin.Context) {
	userId := ctx.GetInt64("userId")
	ctx.JSON(http.StatusOK, bson.M{"isAdmin": domain.IsAdmin(userId)})
}

func ratingEndpoint(ctx *gin.Context) {
	gameId := ctx.Param("gameId")
	name, ok := domain.GetGameName(gameId)
	if ok {
		ctx.JSON(http.StatusOK, bson.M{"rating": domain.GetTopRating(gameId), "name": name})
	} else {
		ctx.Status(http.StatusNoContent)
	}
}

func FilterPublic(ctx *gin.Context) {
	params := ctx.GetHeader("X-VK-PARAMS")
	if len(params) > 0 {
		values, err := url.ParseQuery(params)
		if err == nil {
			signature := values.Get("sign")
			gameId := values.Get("game")
			for k, _ := range values {
				if !strings.HasPrefix(k, "vk_") {
					values.Del(k)
				}
			}
			mac := hmac.New(sha256.New, []byte(vkKey))
			mac.Write([]byte(values.Encode()))
			expected := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(mac.Sum(nil))
			if expected == signature {
				ctx.Set("userId", values.Get("vk_user_id"))
				ctx.Set("gameId", gameId)
				return
			}
		}
	}
	ctx.Status(http.StatusBadRequest)
	ctx.Abort()
}

func FilterAdmin(ctx *gin.Context) {
	userId := ctx.GetInt64("userId")
	if !domain.IsAdmin(userId) {
		ctx.Status(http.StatusForbidden)
		ctx.Abort()
	}
}
