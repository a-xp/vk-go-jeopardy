package domain

import (
	"errors"
	"fmt"
	"goj/configuration"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"sync"
	"time"
)

var gamesLock sync.RWMutex
var activeGames map[GamePost]*Game
var admins Int64Slice
var adminsLock sync.RWMutex

var MockResponse bool
var VkKey string
var PublicAddr string
var RatingAppUrl string
var DefaultListenerName = "GOJ LISTENER"

var ratingUpdateCb []func(id string)
var gameUpdateCb []func(id string)

func AddGameUpdateCallback(f func(id string)) {
	gameUpdateCb = append(gameUpdateCb, f)
}

func AddRatingUpdateCallback(f func(id string)) {
	ratingUpdateCb = append(ratingUpdateCb, f)
}

func InitEngine(cfg *configuration.Configuration) {
	DAO = initDAO(cfg)
	rand.Seed(time.Now().UnixNano())
	MockResponse = cfg.MockResponse
	PublicAddr = cfg.Http.PublicAddr
	RatingAppUrl = cfg.VkApp.Url
	VkKey = cfg.VkApp.Key
	gamesLock = sync.RWMutex{}
	adminsLock = sync.RWMutex{}
	ReloadGames()
	ReloadAdmins()
}

func ReloadAdmins() {
	adminUsers, _ := DAO.getAdminUsers()
	adminsLock.Lock()
	defer adminsLock.Unlock()
	admins = make([]int64, len(adminUsers))
	for i, u := range adminUsers {
		admins[i] = u.Id
	}
}

func ReloadGames() {
	games, err := DAO.loadActiveGames()
	if err != nil {
		log.Fatal(err)
	}
	gamesLock.Lock()
	defer gamesLock.Unlock()
	activeGames = make(map[GamePost]*Game)
	for _, game := range games {
		if game.Active {
			activeGames[game.Post] = game
		}
	}
}

func StopEngine() {
	DAO.close()
}

func FindGroupById(id int64) (*Group, bool) {
	group, exists := DAO.findGroupById(id)
	return group, exists
}

func GetActiveGame(postOwnerId int64, postId int64) (*Game, bool) {
	gamesLock.RLock()
	defer gamesLock.RUnlock()
	game, ok := activeGames[GamePost{
		PostId:      postId,
		PostOwnerId: postOwnerId,
	}]
	return game, ok
}

func GetUserById(vkId int64) (*User, bool) {
	return DAO.findUserById(vkId)
}

func StoreUser(user *User) error {
	return DAO.storeUser(user)
}

func GetGameSession(userId int64, gameId *string) (*Answer, error) {
	return DAO.getGameSession(userId, gameId)
}

func StoreGameSession(session *Answer) error {
	for _, cb := range ratingUpdateCb {
		cb(*session.GameId)
	}
	return DAO.storeGameSession(session)
}

func GetTopRating(gameId *string) ([]*RatingEntry, error) {
	return DAO.getGameTop(gameId, 100)
}

func GetUserRating(gameId *string, userId int64) *RatingEntry {
	return DAO.getUserRating(gameId, userId)
}

func IsAdmin(userId int64) bool {
	adminsLock.RLock()
	defer adminsLock.RUnlock()
	return admins.Search(userId)
}

func GetGameName(gameId *string) (*string, bool) {
	game, ok := DAO.getGameById(gameId)
	if ok {
		return &game.Name, true
	} else {
		return nil, false
	}
}

func GetAdmins() ([]*AdminUser, error) {
	return DAO.getAdminUsers()
}

func RemoveAdmin(id int64) error {
	err := DAO.removeAdmin(id)
	if err == nil {
		ReloadAdmins()
	}
	return err
}

var idPattern = regexp.MustCompile("^(http://|https://)?(www.)?(vk\\.com|vkontakte\\.ru)/(id\\d+|[a-zA-Z0-9_.]+)$")

func AddAdmin(idStr string) error {
	client, err := CreateClient(VkKey)
	if err != nil {
		return err
	}
	match := idPattern.FindStringSubmatch(idStr)
	if match == nil {
		return errors.New("invalid ID")
	}
	users, err := client.GetUser(match[4])

	if err != nil {
		return err
	}

	if len(users) == 0 {
		return errors.New("user not found")
	}

	user := AdminUser{
		Id:    int64(users[0].UID),
		Name:  users[0].FirstName + " " + users[0].LastName,
		Image: users[0].Photo,
	}

	err = DAO.addAdmin(&user)

	if err == nil {
		ReloadAdmins()
	}

	return err
}

func GetGroups() ([]*Group, error) {
	return DAO.getGroups()
}

func RemoveCallbackServer(client *VKExt, groupId int64) error {
	servers, err := client.GetCallbackServers(groupId)
	for _, s := range servers {
		if s.Title == DefaultListenerName {
			if client.DeleteCallbackServer(groupId, s.Id) != nil {
				return err
			}
		}
	}
	return nil
}

func AddGroup(apiKey string) error {
	client, err := CreateClient(apiKey)
	if err != nil {
		return err
	}
	groups, err := client.GetGroup()
	if err != nil {
		return err
	}
	if len(groups) != 1 {
		return errors.New("invalid token")
	}
	vkGroup := groups[0]

	err = RemoveCallbackServer(client, vkGroup.Id)
	if err != nil {
		return err
	}

	code, err := client.GetConfirmCode(vkGroup.Id)

	if err != nil {
		return err
	}

	secret := RandStringBytes(20)

	group := Group{
		Id:          vkGroup.Id,
		ApiKey:      apiKey,
		ConfirmCode: *code,
		Name:        vkGroup.Name,
		Secret:      secret,
		Active:      false,
		Image:       vkGroup.Photo,
	}

	err = DAO.addGroup(&group)
	if err != nil {
		return err
	}

	serverId, err := client.AddCallbackServer(vkGroup.Id, PublicAddr+"/api/callback", DefaultListenerName, secret)
	if err != nil {
		err2 := DAO.removeGroup(group.Id)
		log.Print(err2)
	}
	err = client.SetCallbackSettings(vkGroup.Id, serverId)
	if err != nil {
		err2 := DAO.removeGroup(group.Id)
		log.Print(err2)
	}
	return err
}

func GetGamesShort() ([]*GameHeader, error) {
	list, err := DAO.loadGameHeaders()
	if err != nil {
		return nil, err
	}
	for i, v := range list {
		adr := fmt.Sprintf("%s#%s", RatingAppUrl, *v.Id)
		list[i].RatingUrl = &adr
	}
	return list, nil
}

func GetGame(id *string) (*Game, bool) {
	return DAO.getGameById(id)
}

func StoreGame(game *Game) error {
	if game.Id != nil {
		for _, cb := range gameUpdateCb {
			cb(*game.Id)
		}
	}
	game.Name = strings.TrimSpace(game.Name)
	for i, t := range game.Topics {
		game.Topics[i].Name = FilterAnswer(t.Name)
		for j, q := range t.Q {
			game.Topics[i].Q[j].Text = strings.TrimSpace(q.Text)
			if game.Topics[i].Q[j].Ans != nil {
				for k, ans := range game.Topics[i].Q[j].Ans {
					game.Topics[i].Q[j].Ans[k] = FilterAnswer(ans)
				}
			}
		}
	}
	err := DAO.storeGame(game)
	if err == nil {
		ReloadGames()
	}
	return err
}

func RemoveGroup(groupId int64) error {
	group, found := DAO.findGroupById(groupId)
	if found {
		client, err := CreateClient(group.ApiKey)
		if err == nil {
			_ = RemoveCallbackServer(client, groupId)
		}
		err = DAO.removeGroup(groupId)
		return err
	}
	return nil
}
