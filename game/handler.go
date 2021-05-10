package game

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	vkapi "github.com/himidori/golang-vk-api"
	"goj/domain"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type Event struct {
	Type    string           `binding:"required"`
	GroupId int64            `json:"group_id" binding:"required"`
	Secret  string           `binding:"required"`
	Details WallReplyDetails `json:"object"`
}

type WallReplyDetails struct {
	Id           int64   `json:"id"`
	FromId       int64   `json:"from_id"`
	PostId       int64   `json:"post_id"`
	PostOwnerId  int64   `json:"post_owner_id"`
	ParentsStack []int64 `json:"parents_stack"`
	Text         string
}

func HandleVKEvent(c *gin.Context) {
	var event Event
	if err := c.ShouldBindBodyWith(&event, binding.JSON); err != nil {
		log.Printf("Failed to process %+v", err)
	}
	if event.Type == "confirmation" {
		group, exists := domain.FindGroupById(event.GroupId)
		if exists && group.Secret == event.Secret {
			c.String(http.StatusOK, group.ConfirmCode)
		} else {
			c.Status(http.StatusBadRequest)
		}
	} else {
		if event.Type == "wall_reply_new" && event.Details.FromId > 0 {
			go handleWallReply(&event)
		}
		c.String(http.StatusOK, "ok")
	}
}

func handleWallReply(event *Event) {
	text, ok := filterText(event.Details.Text)
	if !ok {
		return
	}
	group, exists := domain.FindGroupById(event.GroupId)
	if !exists || group.Secret != event.Secret {
		return
	}
	game, exists := domain.GetActiveGame(event.Details.PostOwnerId, event.Details.PostId)
	if !exists {
		return
	}
	client, err := vkapi.NewVKClientWithToken(group.ApiKey, nil, false)
	if err != nil {
		log.Print("Failed to create VK client during event processing", err)
		return
	}
	user, err := getUser(event.Details.FromId, client)
	if err != nil {
		log.Print("Can't request user data", err)
		return
	}
	session, err := getSession(user.Id, game.Id)
	if err != nil {
		log.Print("Can't get session", err)
		return
	}
	playSession(&processingContext{
		event:   event,
		text:    text,
		game:    game,
		user:    user,
		group:   group,
		session: session,
		client:  client,
	})
}

func getUser(vkId int64, client *vkapi.VKClient) (*domain.User, error) {
	user, exists := domain.GetUserById(vkId)
	if !exists {
		data, err := client.UsersGet([]int{int(vkId)})
		if err != nil {
			return nil, err
		}
		user = &domain.User{
			Id:       int64(data[0].UID),
			Img:      data[0].Photo,
			Name:     data[0].FirstName,
			Lastname: data[0].LastName,
		}
		if err = domain.StoreUser(user); err != nil {
			return nil, err
		}
	}
	return user, nil
}

func getSession(userId int64, gameId string) (*domain.Answer, error) {
	session, err := domain.GetGameSession(userId, gameId)
	if err != nil {
		return nil, err
	}
	if session == nil {
		session = &domain.Answer{
			Complete:     false,
			CurrentTopic: -1,
			Score:        0,
			GameId:       gameId,
			UserId:       userId,
			Topics:       nil,
		}
	}
	return session, nil
}

func filterText(original string) (string, bool) {
	replyPattern := regexp.MustCompile("^(.+,\\s*)")
	result := strings.ToLower(strings.TrimSpace(replyPattern.ReplaceAllString(original, "")))
	if len(result) > 0 {
		return result, true
	} else {
		return result, false
	}
}
