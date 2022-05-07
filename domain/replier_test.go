//go:build integration
// +build integration

package domain

import (
	"fmt"
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
	for i := 0; i < 10; i++ {
		Replier.Input() <- ReplyMsg{
			PostOwnerId: "-196065343",
			PostId:      "501",
			CommentId:   "503",
			Message:     fmt.Sprintf("Hello %d", i),
			AccessToken: token,
		}
	}
	<-time.After(time.Second * 5)
}
