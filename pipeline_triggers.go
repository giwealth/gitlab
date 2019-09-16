package gitlab

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Trigger 触发器信息
type Trigger struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	LastUsed    string `json:"last_used"`
	Token       string `json:"token"`
	UpdatedAt   string `json:"updated_at"`
	Owner       struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Username  string `json:"username"`
		State     string `json:"state"`
		AvatarURL string `json:"avatar_url"`
		WebURL    string `json:"web_url"`
	}
}

// CreateTrigger 创建触发器
func (c *Client) CreateTrigger(projectID int) (tiggerToken string, err error) {
	client := http.Client{}

	url := fmt.Sprintf("%s/api/v4/projects/%v/triggers?description=deploy", strings.TrimSuffix(c.BaseURL, "/"), projectID)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("Private-Token", c.AccessToken)

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return "", fmt.Errorf("request error response status code %v", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	var tigger Trigger
	if err = json.Unmarshal(body, &tigger); err != nil {
		return
	}

	return tigger.Token, nil
}

// GetTrigger 获取仓库触发器
func (c *Client) GetTrigger(projectID int) (triggerToken string, err error) {
	url := fmt.Sprintf("%s/api/v4/projects/%v/triggers", c.BaseURL, projectID)
	client := &http.Client{}
	var triggers []Trigger

	res, err := httpGetRequest(url, c.AccessToken, client)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &triggers)
	if err != nil {
		return
	}

	for _, trigger := range triggers {
		if trigger.Token != "" {
			return trigger.Token, nil
		}
	}

	return
}
