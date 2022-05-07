package domain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	apiVersion = "5.103"
	apiURL     = "https://api.vk.com/method/%s"
)

type vkAPIResponse struct {
	Response      json.RawMessage `json:"response"`
	ResponseError vkError         `json:"error"`
}

type vkError struct {
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

var Replier *VkReplier

var maxDelay = 1 * time.Second
var bufWindow = 250 * time.Millisecond
var maxBatch = 20

type ReplyMsg struct {
	PostOwnerId string `json:"owner_id"`
	PostId      string `json:"post_id"`
	CommentId   string `json:"reply_to_comment"`
	Message     string `json:"message"`
	AccessToken string `json:"access_token"`
}

type VkReplier struct {
	input  chan ReplyMsg
	client *http.Client
}

func (r *VkReplier) worker() {
	buf := make([]ReplyMsg, 0)
	lastSend := time.Now()
	for {
		select {
		case m, more := <-r.input:
			if !more {
				return
			}
			buf = append(buf, m)
		case <-time.After(bufWindow):
		}
		if len(buf) > 0 && (time.Now().Sub(lastSend) > maxDelay || len(buf) >= maxBatch) {
			r.send(buf)
			buf = buf[:0]
		}
	}
}

func (r *VkReplier) send(messages []ReplyMsg) {
	groups := map[string][]ReplyMsg{}
	for _, m := range messages {
		groups[m.AccessToken] = append(groups[m.AccessToken], m)
	}
	for token, items := range groups {
		r.sendToGroup(items, token)
	}
}

func (r *VkReplier) sendToGroup(messages []ReplyMsg, token string) {
	lines := make([]string, len(messages))
	for i, m := range messages {
		b, _ := json.Marshal(m)
		lines[i] = fmt.Sprintf("r.push(API.wall.createComment(%s).comment_id);", string(b))
	}
	code := strings.Join(lines, "\n")
	code = fmt.Sprintf("var r = [];\n%s\nreturn r;\n", code)
	fmt.Printf("Sending command %s", code)
	if err, resp := r.execute(code, token); err != nil {
		log.Printf("Error responding: %v", err)
	} else {
		r.processResponse(resp, messages)
	}
}

func (r *VkReplier) processResponse(message json.RawMessage, messages []ReplyMsg) {
	var ids []*int
	if err := json.Unmarshal(message, &ids); err != nil {
		log.Printf("Failed to deserialize vk reply response %v", err)
	} else {
		for k, v := range ids {
			if v == nil {
				log.Printf("Reply to comment %s failed with unknown error", messages[k].CommentId)
			}
		}
	}
}

func (r *VkReplier) Input() chan<- ReplyMsg {
	return r.input
}

func (r *VkReplier) execute(code string, token string) (error, json.RawMessage) {
	params := url.Values{}
	params.Set("code", code)
	params.Set("v", apiVersion)
	params.Set("access_token", token)
	methodUrl := fmt.Sprintf(apiURL, "execute")
	response, err := r.client.PostForm(methodUrl, params)
	if err != nil {
		return err, nil
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err, nil
	}
	var jsonResponse vkAPIResponse
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return err, nil
	}
	if jsonResponse.ResponseError.ErrorCode != 0 {
		return fmt.Errorf("vk returened error %d %s",
			jsonResponse.ResponseError.ErrorCode,
			jsonResponse.ResponseError.ErrorMsg), nil
	}
	return nil, jsonResponse.Response
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
		input:  make(chan ReplyMsg, 100),
		client: netClient,
	}
	go Replier.worker()
}
