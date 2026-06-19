package game

import (
	"goj/domain"
	"testing"
)

func Test_isCorrectAnswer(t *testing.T) {
	type args struct {
		gameAnswers []string
		userAnswer  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Ответ одно число",
			args{[]string{"1"}, "1"},
			true,
		},
		{
			"Ответ одно слово кириллицей",
			args{[]string{"ответ"}, "ответ"},
			true,
		},
		{
			"Ответ несколько слов кириллицей",
			args{[]string{"правильный ответ"}, "правильный ответ"},
			true,
		},
		{
			"Ответ несколько слов кириллицей с запятыми",
			args{[]string{"правильный, ответ"}, "правильный, ответ"},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isCorrectAnswer(tt.args.gameAnswers, tt.args.userAnswer); got != tt.want {
				t.Errorf("isCorrectAnswer() = %v, want %v", got, tt.want)
			}
		})
	}
}

var testGameMessages = struct {
	GameComplete   string `bson:"gameComplete"`
	Error          string `bson:"error"`
	TopicComplete  string `bson:"topicComplete"`
	IncorrectFinal string `bson:"incorrectFinal"`
	IncorrectRetry string `bson:"incorrectRetry"`
	Correct        string `bson:"correct"`
	WrongBranch    string `bson:"wrongBranch"`
	UnknownTopic   string `bson:"unknownTopic"`
}{
	GameComplete:   "",
	Error:          "",
	TopicComplete:  "",
	IncorrectFinal: "",
	IncorrectRetry: "",
	Correct:        "",
	WrongBranch:    "",
	UnknownTopic:   "",
}

var testGameRules = struct {
	InstantWin bool `bson:"instantWin" json:"instantWin"`
	NumTries   int  `bson:"numTries" json:"numTries"`
}{}

func Test_searchTopicByText(t *testing.T) {
	game := domain.Game{
		GameHeader: domain.GameHeader{},
		Messages:   testGameMessages,
		Post:       domain.GamePost{},
		Rules:      testGameRules,
		Topics: []domain.Topic{
			{
				Name:   "Тема",
				Points: 1,
				Q:      nil,
			},
			{
				Name:   "загАдки-рАзгадки",
				Points: 1,
				Q:      nil,
			},
			{
				Name:   "ЕЛКИ",
				Points: 1,
				Q:      nil,
			},
			{
				Name:   "ёжики и елки",
				Points: 1,
				Q:      nil,
			},
			{
				Name:   "pines 45!",
				Points: 1,
				Q:      nil,
			},
		},
	}
	ctx := processingContext{
		text:    "",
		event:   nil,
		game:    &game,
		user:    nil,
		group:   nil,
		session: nil,
		client:  nil,
	}

	tests := []struct {
		name     string
		answer   string
		topicNum int
		hasMatch bool
	}{
		{"Ответ кирилицей", "тема", 0, true},
		{"Ответ с тире", "загадки разгадки", 1, true},
		{"Ответ капсом", "елки", 2, true},
		{"Ответ с ё", "ежики и елки", 3, true},
		{"Ответ на латинице со спецсимволами", "pines 45", 4, true},
		{"Неправильный ответ", "не знаю", 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx.text = tt.answer
			topicNum, hasMatch := searchTopicByText(&ctx)
			if topicNum != tt.topicNum {
				t.Errorf("searchTopicByText() got = %v, want %v", topicNum, tt.topicNum)
			}
			if hasMatch != tt.hasMatch {
				t.Errorf("searchTopicByText() got1 = %v, want %v", hasMatch, tt.hasMatch)
			}
		})
	}
}
