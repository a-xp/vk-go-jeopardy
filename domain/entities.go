package domain

type GamePost struct {
	PostId      int64 `json:"postId"`
	PostOwnerId int64 `json:"postOwnerId"`
	GroupId     int64 `json:"groupId"`
}

type GameHeader struct {
	Id        *string `bson:"_id,omitempty" json:"id"`
	Name      string  `bson:"name" json:"name"`
	Active    bool    `bson:"active" json:"active"`
	New       bool    `bson:"new" json:"new"`
	RatingUrl *string
}

type Game struct {
	GameHeader `bson:",inline"`
	Messages   struct {
		GameComplete   string `bson:"gameComplete"`
		Error          string `bson:"error"`
		TopicComplete  string `bson:"topicComplete"`
		IncorrectFinal string `bson:"incorrectFinal"`
		IncorrectRetry string `bson:"incorrectRetry"`
		Correct        string `bson:"correct"`
		WrongBranch    string `bson:"wrongBranch"`
		UnknownTopic   string `bson:"unknownTopic"`
	} `json:"messages"`
	Post  GamePost `json:"post"`
	Rules struct {
		InstantWin bool `bson:"instantWin"`
		NumTries   int  `bson:"numTries"`
	} `json:"rules"`
	Topics []struct {
		Name   string `json:"name"`
		Points int    `json:"points"`
		Q      []struct {
			Ans  []string `json:"ans"`
			Text string   `json:"text"`
		} `json:"q"`
	} `json:"topics"`
}

type Group struct {
	Id          int64  `bson:"_id" json:"id"`
	ApiKey      string `bson:"apiKey" json:"apiKey"`
	ConfirmCode string `bson:"confirmCode" json:"confirmCode"`
	Name        string `json:"name"`
	Secret      string `json:"secret"`
	Active      bool   `json:"active"`
	Image       string `json:"image"`
}

type User struct {
	Id       int64 `bson:"_id"`
	Img      string
	Name     string
	Lastname string `bson:"lastname"`
}

type TopicResult struct {
	PostId   int64 `bson:"postId"`
	Complete bool
	Result   bool
	Attempt  int
	Question int
}

type Answer struct {
	Id           *string `bson:"_id"`
	Complete     bool
	CompleteTime int64 `bson:"completeTime"`
	CurrentTopic int   `bson:"currentTopic"`
	Score        int
	GameId       *string `bson:"gameId"`
	UserId       int64   `bson:"userId"`
	Topics       []*TopicResult
}

type RatingEntry struct {
	Pos      int    `json:"pos"`
	UserId   int64  `json:"userId"`
	Score    int    `json:"score"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Img      string `json:"img"`
}

type AdminUser struct {
	Id    int64  `bson:"_id" json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}
