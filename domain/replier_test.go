//go:build integration
// +build integration

package domain

import (
	"log"
	"os"
	"testing"
	"time"
)

func getTestKey() string {
	buf, err := os.ReadFile("../test_resources/key.txt")
	if err != nil {
		log.Panic(err)
	}
	return string(buf)
}

func TestThatSendMessageThrowsNoErrors(t *testing.T) {
	InitReplier()
	token := getTestKey()
	Replier.Input() <- ReplyMsg{
		PostOwnerId: "-196065343",
		PostId:      "390",
		CommentId:   "406",
		Message:     "Hello",
		AccessToken: token,
	}
	Replier.Input() <- ReplyMsg{
		PostOwnerId: "-196065343",
		PostId:      "390",
		CommentId:   "406",
		Message:     "Hello 2",
		AccessToken: token,
	}
	Replier.Input() <- ReplyMsg{
		PostOwnerId: "-196065343",
		PostId:      "390",
		CommentId:   "406",
		Message:     "Hello 3",
		AccessToken: token,
	}
	<-time.After(time.Second * 5)
}
