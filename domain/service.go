package domain

import (
	"goj/configuration"
	"log"
	"sync"
)

var gamesLock sync.RWMutex
var activeGames map[GamePost]*Game
var admins Int64Slice
var adminsLock sync.RWMutex

var MockResponse bool

func InitEngine(cfg *configuration.Configuration) {
	DAO = initDAO(cfg)
	MockResponse = cfg.MockResponse
	gamesLock = sync.RWMutex{}
	adminsLock = sync.RWMutex{}
	ReloadGames()
	ReloadAdmins()
}

func ReloadAdmins() {
	adminUsers, _ := DAO.getAdminUsers()
	adminsLock.Lock()
	defer adminsLock.Unlock()
	admins := make([]int64, len(adminUsers))
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
	return DAO.storeGameSesion(session)
}

func GetTopRating(gameId string) ([]RatingEntry, error) {
	return DAO.getGameTop(gameId, 100)
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
