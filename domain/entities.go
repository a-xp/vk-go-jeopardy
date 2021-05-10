package domain

type GamePost struct {
	PostId      int64
	PostOwnerId int64
}

type Game struct {
	Id       string `bson:"_id"`
	Name     string
	Active   bool
	Messages struct {
		GameComplete   string `bson:"gameComplete"`
		Error          string `bson:"error"`
		TopicComplete  string `bson:"topicComplete"`
		IncorrectFinal string `bson:"incorrectFinal"`
		IncorrectRetry string `bson:"incorrectRetry"`
		Correct        string `bson:"correct"`
		WrongBranch    string `bson:"wrongBranch"`
		UnknownTopic   string `bson:"unknownTopic"`
	}
	Post  GamePost
	Rules struct {
		InstantWin bool `bson:"instantWin"`
		NumTries   int  `bson:"numTries"`
	}
	Topics []struct {
		Name   string
		Points int
		Q      []struct {
			Ans  []string
			Text string
		}
	}
}

type Group struct {
	Id          int64  `bson:"_id"`
	ApiKey      string `bson:"apiKey"`
	ConfirmCode string `bson:"confirmCode"`
	Name        string
	Secret      string
	Active      bool
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
	GameId       string `bson:"gameId"`
	UserId       int64  `bson:"userId"`
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
	Id    int64 `bson:"_id"`
	Name  string
	Image string
}
