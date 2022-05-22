package domain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	apiVersion = "5.103"
	apiURL     = "https://api.vk.com/method/%s"
)

type vkAPIResponse struct {
	Response      json.RawMessage `json:"response"`
	ResponseError vkErrorResponse `json:"error"`
}

type vkErrorResponse struct {
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

func (e *vkErrorResponse) Error() string {
	return fmt.Sprintf("[%d] %s", e.ErrorCode, e.ErrorMsg)
}

var Replier *VkReplier

var sendInterval = 1100 * time.Millisecond

type ReplyMsg struct {
	PostOwnerId string `json:"owner_id"`
	PostId      string `json:"post_id"`
	CommentId   string `json:"reply_to_comment"`
	Message     string `json:"message"`
	AccessToken string `json:"access_token"`
}

type VkReplier struct {
	queue  chan ReplyMsg
	client *http.Client
}

func (r *VkReplier) Send(msg ReplyMsg) {
	r.queue <- msg
}

func (r *VkReplier) worker() {
	for {
		select {
		case m, more := <-r.queue:
			if !more {
				return
			}
			now := time.Now()
			r.replyWithRetry(&m)
			wait := sendInterval - time.Now().Sub(now)
			<-time.After(wait)
		}
	}
}

func (r VkReplier) replyWithRetry(msg *ReplyMsg) {
	for {
		err := r.vkReply(msg)
		if err != nil {
			if vkErr, ok := err.(*vkErrorResponse); ok && vkErr.ErrorCode == 9 {
				<-time.After(15 * time.Second)
				continue
			} else {
				log.Printf("Failed to response to %s:%s: %v", msg.PostOwnerId, msg.CommentId, err)
			}
		}
		break
	}
}

func (r *VkReplier) vkReply(msg *ReplyMsg) error {
	log.Printf("Replying %s to comment %s", msg.Message, msg.CommentId)
	params := url.Values{}
	params.Set("v", apiVersion)
	params.Set("access_token", msg.AccessToken)
	params.Set("owner_id", msg.PostOwnerId)
	params.Set("post_id", msg.PostId)
	params.Set("reply_to_comment", msg.CommentId)
	params.Set("message", msg.Message)
	methodUrl := fmt.Sprintf(apiURL, "wall.createComment")
	response, err := r.client.PostForm(methodUrl, params)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	var jsonResponse vkAPIResponse
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return err
	}
	if jsonResponse.ResponseError.ErrorCode != 0 {
		return &jsonResponse.ResponseError
	}
	return nil
}

func InitReplier() {
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	var netClient = &http.Client{
		Timeout:   10 * time.Second,
		Transport: netTransport,
	}
	Replier = &VkReplier{
		client: netClient,
		queue:  make(chan ReplyMsg, 3000),
	}
	go Replier.worker()
}
