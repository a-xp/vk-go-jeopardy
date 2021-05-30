package domain

import (
	"encoding/json"
	vkapi "github.com/himidori/golang-vk-api"
	"log"
	"net/url"
	"strconv"
)

type VKExt struct {
	client *vkapi.VKClient
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
		log.Print("Failed to create VK client", err)
		return nil, err
	}
	return &VKExt{client: client}, nil
}

func (client *VKExt) GetConfirmCode(id int64) (*string, error) {
	v := url.Values{}
	v.Add("group_id", strconv.FormatInt(id, 10))

	resp, err := client.client.MakeRequest("groups.getCallbackConfirmationCode", v)
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

	resp, err := client.client.MakeRequest("groups.addCallbackServer", v)

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

	resp, err := client.client.MakeRequest("groups.getCallbackServers", v)

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

	_, err := client.client.MakeRequest("groups.deleteCallbackServer", v)

	return err
}

func (client *VKExt) SetCallbackSettings(groupId int64, serverId int64) error {
	v := url.Values{}
	v.Add("group_id", strconv.FormatInt(groupId, 10))
	v.Add("server_id", strconv.FormatInt(serverId, 10))
	v.Add("wall_reply_new", "1")

	_, err := client.client.MakeRequest("groups.setCallbackSettings", v)

	return err
}

func (client *VKExt) GetGroup() ([]*GroupInfo, error) {
	v := url.Values{}
	resp, err := client.client.MakeRequest("groups.getById", v)
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

	resp, err := client.client.MakeRequest("users.get", v)
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
