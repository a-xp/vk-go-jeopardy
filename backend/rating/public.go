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
	"strconv"
	"strings"
	"sync"
)

func meEndpoint(ctx *gin.Context) {
	userId := ctx.GetInt64("userId")
	ctx.JSON(http.StatusOK, bson.M{"isAdmin": domain.IsAdmin(userId), "isValidClient": true})
}

var ratings sync.Map
var names sync.Map

func DropNamesCache(groupId string) {
	names.Delete(groupId)
}

func DropRatingCache(groupId string) {
	ratings.Delete(groupId)
}

func ratingEndpoint(ctx *gin.Context) {
	gameId := ctx.GetString("gameId")
	userId := ctx.GetInt64("userId")
	v, ok := names.Load(gameId)
	var name string
	if ok {
		name = v.(string)
	} else {
		str, ok := domain.GetGameName(&gameId)
		if !ok {
			ctx.Status(http.StatusNoContent)
			return
		}
		name = *str
		names.Store(gameId, name)
	}
	v, ok = ratings.Load(gameId)
	var list []*domain.RatingEntry
	if !ok {
		var err error
		list, err = domain.GetTopRating(&gameId)
		if err != nil {
			ctx.Status(http.StatusInternalServerError)
			return
		}
		ratings.Store(gameId, list)
	} else {
		list = v.([]*domain.RatingEntry)
	}
	var userRating *domain.RatingEntry
	for _, e := range list {
		if e.UserId == userId {
			userRating = e
			break
		}
	}
	if userRating == nil {
		userRating = domain.GetUserRating(&gameId, userId)
	}
	ctx.JSON(http.StatusOK, bson.M{"rating": list, "name": name, "userRating": userRating})

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
			filteredQuery := values.Encode()
			mac := hmac.New(sha256.New, []byte(vkKey))
			mac.Write([]byte(filteredQuery))
			expected := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(mac.Sum(nil))
			if !validateSignature || expected == signature {
				userIdStr := values.Get("vk_user_id")
				userId, err := strconv.ParseInt(userIdStr, 10, 64)
				if err == nil {
					ctx.Set("userId", userId)
					ctx.Set("gameId", gameId)
					return
				}
			}
		}
	} else {
		auth := ctx.GetHeader("X-RootKey")
		if auth == rootPass {
			ctx.Set("userId", int64(1))
			return
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
