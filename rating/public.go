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
)

func meEndpoint(ctx *gin.Context) {
	userId := ctx.GetInt64("userId")
	ctx.JSON(http.StatusOK, bson.M{"isAdmin": domain.IsAdmin(userId), "isValidClient": true})
}

func ratingEndpoint(ctx *gin.Context) {
	gameId := ctx.GetString("gameId")
	userId := ctx.GetInt64("userId")
	name, ok := domain.GetGameName(&gameId)
	if ok {
		list, err := domain.GetTopRating(&gameId)
		if err == nil {
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
		} else {
			ctx.Status(http.StatusInternalServerError)
		}
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
