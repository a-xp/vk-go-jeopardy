//go:build integration
// +build integration

package domain

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func getTestKeys() []string {
	buf, err := os.ReadFile("../test_resources/key.txt")
	if err != nil {
		log.Panic(err)
	}
	keys := strings.Split(string(buf), "\n")
	return keys
}

func TestThatSendMessageThrowsNoErrors(t *testing.T) {
	InitReplier()
	tokens := getTestKeys()
	for i := 0; i < 1000; i++ {
		Replier.Send(ReplyMsg{
			PostOwnerId: "-196065343",
			PostId:      "501",
			CommentId:   "677",
			Message:     fmt.Sprintf("Hello %d", i),
			AccessToken: tokens[i%len(tokens)],
		})
	}
	<-time.After(time.Second * 360)
}
