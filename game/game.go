package game

import (
	"goj/domain"
	"log"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
)

type processingContext struct {
	text    string
	event   *Event
	game    *domain.Game
	user    *domain.User
	group   *domain.Group
	session *domain.Answer
	client  *domain.VKExt
}

func playSession(ctx *processingContext) {
	if !validateReplyBranch(ctx) {
		sendReply(ctx, ctx.game.Messages.WrongBranch)
		return
	}
	if ctx.session.CurrentTopic == -1 {
		chooseTopic(ctx)
	} else {
		answerQuestion(ctx)
	}
}

func chooseTopic(ctx *processingContext) {
	if ctx.session.Complete {
		sendReply(ctx, ctx.game.Messages.GameComplete, strconv.Itoa(ctx.session.Score))
		return
	}
	topicNum, found := searchTopicByText(ctx)
	if !found {
		sendReply(ctx, ctx.game.Messages.UnknownTopic)
		return
	}
	if ctx.session.Topics != nil &&
		ctx.session.Topics[topicNum] != nil &&
		ctx.session.Topics[topicNum].Complete {
		sendReply(ctx, ctx.game.Messages.TopicComplete)
	} else {
		playTopic(ctx, topicNum)
	}
}

func searchTopicByText(ctx *processingContext) (int, bool) {
	topicNum, err := strconv.Atoi(ctx.text)
	found := false
	if err != nil {
		for num, topic := range ctx.game.Topics {
			if strings.ToLower(topic.Name) == ctx.text {
				topicNum = num
				found = true
				break
			}
		}
	} else {
		if topicNum >= 1 && topicNum <= len(ctx.game.Topics) {
			found = true
			topicNum--
		}
	}
	return topicNum, found
}

func playTopic(ctx *processingContext, topicNum int) {
	ctx.session.CurrentTopic = topicNum
	questionNum := rand.Intn(len(ctx.game.Topics[topicNum].Q))
	postId := ctx.event.Details.Id
	if len(ctx.event.Details.ParentsStack) > 0 {
		postId = ctx.event.Details.ParentsStack[0]
	}
	if ctx.session.Topics == nil {
		numTopics := len(ctx.game.Topics)
		ctx.session.Topics = make([]*domain.TopicResult, numTopics)
	}
	ctx.session.Topics[topicNum] = &domain.TopicResult{
		PostId:   postId,
		Question: questionNum,
	}
	if err := domain.StoreGameSession(ctx.session); err != nil {
		sendReply(ctx, ctx.game.Messages.Error)
	} else {
		sendReply(ctx, ctx.game.Topics[topicNum].Q[questionNum].Text)
	}
}

func answerQuestion(ctx *processingContext) {
	game := ctx.game
	topicId := ctx.session.CurrentTopic
	answer := ctx.session.Topics[topicId]
	question := game.Topics[topicId].Q[answer.Question]
	if isCorrectAnswer(question.Ans, ctx.text) {
		answer.Complete = true
		answer.Result = true
		ctx.session.CurrentTopic = -1
	} else {
		answer.Attempt++
		if game.Rules.NumTries > 0 && answer.Attempt >= game.Rules.NumTries {
			answer.Complete = true
			answer.Result = false
			ctx.session.CurrentTopic = -1
		}
	}
	ctx.session.Score = calcGameScore(ctx.session, game)
	if isGameComplete(ctx.session) {
		ctx.session.Complete = true
	}
	if err := domain.StoreGameSession(ctx.session); err != nil {
		sendReply(ctx, game.Messages.Error)
	} else {
		if answer.Result {
			sendReply(ctx, game.Messages.Correct, strconv.Itoa(game.Topics[topicId].Points))
		} else {
			if answer.Complete {
				sendReply(ctx, game.Messages.IncorrectFinal)
			} else {
				sendReply(ctx, game.Messages.IncorrectRetry)
			}
		}
		if ctx.session.Complete {
			sendReply(ctx, game.Messages.GameComplete, strconv.Itoa(ctx.session.Score))
		}
	}
}

func calcGameScore(session *domain.Answer, game *domain.Game) int {
	score := 0
	for id, topic := range session.Topics {
		if topic != nil && topic.Complete && topic.Result {
			score += game.Topics[id].Points
		}
	}
	return score
}

func isGameComplete(session *domain.Answer) bool {
	for _, topic := range session.Topics {
		if topic == nil || !topic.Complete {
			return false
		}
	}
	return true
}

func isCorrectAnswer(gameAnswers []string, userAnswer string) bool {
	for _, s := range gameAnswers {
		if strings.EqualFold(s, userAnswer) {
			return true
		}
	}
	return false
}

func validateReplyBranch(ctx *processingContext) bool {
	if len(ctx.event.Details.ParentsStack) > 0 {
		rootId := ctx.event.Details.ParentsStack[0]
		for _, topic := range ctx.session.Topics {
			if topic != nil && topic.PostId == rootId {
				return true
			}
		}
		return false
	}
	return true
}

func sendReply(ctx *processingContext, template string, param ...string) {
	message := template
	if len(param) > 0 {
		message = strings.ReplaceAll(template, "#X", param[0])
	}
	if !domain.MockResponse {
		vkParams := url.Values{}
		vkParams.Add("reply_to_comment", strconv.FormatInt(ctx.event.Details.Id, 10))
		_, err := ctx.client.WallPostComment(int(ctx.game.Post.PostOwnerId), int(ctx.game.Post.PostId), message, vkParams)
		if err != nil {
			log.Print("Failed to respond: ", err)
		}
	} else {
		log.Printf("Sending response to %d with message %s", ctx.game.Post.PostId, message)
	}
}
