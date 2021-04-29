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
		GameComplete   string
		Error          string
		TopicComplete  string
		IncorrectFinal string
		IncorrectRetry string
		Correct        string
		WrongBranch    string
		UnknownTopic   string
	}
	Post  GamePost
	Rules struct {
		InstantWin bool
		NumTries   int
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
	Id          int64 `bson:"_id"`
	ApiKey      string
	ConfirmCode string
	Name        string
	Secret      string
	Active      bool
}

type User struct {
	Id       int64 `bson:"_id"`
	Img      string
	Name     string
	LastName string
}

type TopicResult struct {
	PostId   int64
	Complete bool
	Result   bool
	Attempt  int
	Question int
}

type Answer struct {
	Id           *string `bson:"_id"`
	Complete     bool
	CurrentTopic int
	Score        int
	GameId       string
	UserId       int64
	Topics       []*TopicResult
}

type RatingEntry struct {
	UserId   int64
	Score    int
	Name     string
	Lastname string
	Img      string
}

type AdminUser struct {
	Id    int64 `bson:"_id"`
	Name  string
	Image string
}
