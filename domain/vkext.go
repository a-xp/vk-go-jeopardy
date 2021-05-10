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
	json.Unmarshal(resp.Response, &code)
	return &code.Code, nil
}

func (client *VKExt) CreateListener(groupId int64, cbUrl string, name string, secret string) error {
	v := url.Values{}
	v.Add("group_id", strconv.FormatInt(groupId, 10))
	v.Add("url", cbUrl)
	v.Add("title", name)
	v.Add("secret_key", secret)

	_, err := client.client.MakeRequest("groups.addCallbackServer", v)
	return err
}

func (client *VKExt) GetGroup() ([]*GroupInfo, error) {
	v := url.Values{}
	resp, err := client.client.MakeRequest("groups.getById", v)
	if err != nil {
		return nil, err
	}

	var items []*GroupInfo
	json.Unmarshal(resp.Response, &items)

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
	json.Unmarshal(resp.Response, &userList)

	return userList, nil
}
