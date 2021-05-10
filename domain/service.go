package domain

import (
	"errors"
	"goj/configuration"
	"log"
	"math/rand"
	"regexp"
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

func InitEngine(cfg *configuration.Configuration) {
	DAO = initDAO(cfg)
	rand.Seed(time.Now().UnixNano())
	MockResponse = cfg.MockResponse
	PublicAddr = cfg.Http.PublicAddr
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
	games, err := DAO.loadGames()
	if err != nil {
		log.Fatal(err)
	}
	gamesLock.Lock()
	defer gamesLock.Unlock()
	activeGames = make(map[GamePost]*Game)
	for _, game := range games {
		if game.Active {
			activeGames[game.Post] = &game
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

func GetGameSession(userId int64, gameId string) (*Answer, error) {
	return DAO.getGameSession(userId, gameId)
}

func StoreGameSession(session *Answer) error {
	return DAO.storeGameSession(session)
}

func GetTopRating(gameId string) (*[]RatingEntry, error) {
	return DAO.getGameTop(gameId, 100)
}

func GetUserRating(gameId string, userId int64) *RatingEntry {
	return DAO.getUserRating(gameId, userId)
}

func IsAdmin(userId int64) bool {
	adminsLock.RLock()
	defer adminsLock.RUnlock()
	return admins.Search(userId)
}

func GetGameName(gameId string) (*string, bool) {
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
	return DAO.removeAdmin(id)
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

	return DAO.addAdmin(&user)
}

func GetGroups() ([]*Group, error) {
	return DAO.getGroups()
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

	err = client.CreateListener(vkGroup.Id, PublicAddr+"/api/callback", "GOJ LISTENER", secret)
	if err != nil {
		DAO.removeGroup(group.Id)
	}
	return err
}
