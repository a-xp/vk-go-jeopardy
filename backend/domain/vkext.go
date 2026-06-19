package domain

import (
	"encoding/json"
	vkapi "github.com/himidori/golang-vk-api"
	"log"
	"net/url"
	"strconv"
)

type VKExt struct {
	Client *vkapi.VKClient
}

type GroupCodeResponse struct {
	Code string `json:"code"`
}

type GroupInfo struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Photo string `json:"photo_50"`
}

type CallbackServerInfo struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
}

type CallbackServersResponse struct {
	Count int                   `json:"count"`
	Items []*CallbackServerInfo `json:"items"`
}

type AddCallbackServerResponse struct {
	ServerId int64 `json:"server_id"`
}

func CreateClient(apiKey string) (*VKExt, error) {
	client, err := vkapi.NewVKClientWithToken(apiKey, nil, false)
	if err != nil {
		log.Print("Failed to create VK Client", err)
		return nil, err
	}
	return &VKExt{Client: client}, nil
}

func (client *VKExt) GetConfirmCode(id int64) (*string, error) {
	v := url.Values{}
	v.Add("group_id", strconv.FormatInt(id, 10))

	resp, err := client.Client.MakeRequest("groups.getCallbackConfirmationCode", v)
	if err != nil {
		return nil, err
	}
	var code GroupCodeResponse
	err = json.Unmarshal(resp.Response, &code)
	if err != nil {
		return nil, err
	}
	return &code.Code, nil
}

func (client *VKExt) AddCallbackServer(groupId int64, cbUrl string, name string, secret string) (int64, error) {
	v := url.Values{}
	v.Add("group_id", strconv.FormatInt(groupId, 10))
	v.Add("url", cbUrl)
	v.Add("title", name)
	v.Add("secret_key", secret)

	resp, err := client.Client.MakeRequest("groups.addCallbackServer", v)

	if err != nil {
		return 0, err
	}
	var result AddCallbackServerResponse
	err = json.Unmarshal(resp.Response, &result)
	if err != nil {
		return 0, err
	}
	return result.ServerId, nil
}

func (client *VKExt) GetCallbackServers(groupId int64) ([]*CallbackServerInfo, error) {
	v := url.Values{}
	v.Add("group_id", strconv.FormatInt(groupId, 10))

	resp, err := client.Client.MakeRequest("groups.getCallbackServers", v)

	if err != nil {
		return nil, err
	}

	var result CallbackServersResponse
	err = json.Unmarshal(resp.Response, &result)

	if err != nil {
		return nil, err
	}

	return result.Items, nil
}

func (client *VKExt) DeleteCallbackServer(groupId int64, serverId int64) error {
	v := url.Values{}
	v.Add("group_id", strconv.FormatInt(groupId, 10))
	v.Add("server_id", strconv.FormatInt(serverId, 10))

	_, err := client.Client.MakeRequest("groups.deleteCallbackServer", v)

	return err
}

func (client *VKExt) SetCallbackSettings(groupId int64, serverId int64) error {
	v := url.Values{}
	v.Add("group_id", strconv.FormatInt(groupId, 10))
	v.Add("server_id", strconv.FormatInt(serverId, 10))
	v.Add("wall_reply_new", "1")

	_, err := client.Client.MakeRequest("groups.setCallbackSettings", v)

	return err
}

func (client *VKExt) GetGroup() ([]*GroupInfo, error) {
	v := url.Values{}
	resp, err := client.Client.MakeRequest("groups.getById", v)
	if err != nil {
		return nil, err
	}

	var items []*GroupInfo
	err = json.Unmarshal(resp.Response, &items)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (client *VKExt) GetUser(id string) ([]*vkapi.User, error) {
	v := url.Values{}
	v.Add("user_ids", id)
	v.Add("fields", "photo,nickname,screen_name")
	v.Add("lang", "ru")

	resp, err := client.Client.MakeRequest("users.get", v)
	if err != nil {
		return nil, err
	}

	var userList []*vkapi.User
	err = json.Unmarshal(resp.Response, &userList)
	if err != nil {
		return nil, err
	}
	return userList, nil
}

type WallPostCommentResponse struct {
	CommentId int `json:"comment_id"`
}

func (client *VKExt) WallPostComment(ownerID int, postID int, message string, params url.Values) (int, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Set("owner_id", strconv.Itoa(ownerID))
	params.Set("post_id", strconv.Itoa(postID))
	params.Set("message", message)

	resp, err := client.Client.MakeRequest("wall.createComment", params)
	if err != nil {
		return 0, err
	}
	m := WallPostCommentResponse{}
	if err = json.Unmarshal(resp.Response, &m); err != nil {
		return 0, err
	}
	return m.CommentId, nil
}

func (client *VKExt) Execute(code string) error {
	params := url.Values{}
	params.Set("code", code)
	_, err := client.Client.MakeRequest("execute", params)
	return err
}
