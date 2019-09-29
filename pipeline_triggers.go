package gitlab

import (
	"fmt"
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
func (c *Client) CreateTrigger(projectID int, description string) (Trigger, error) {
	var trigger Trigger
	err := c.CreateResource(fmt.Sprintf("/projects/%v/triggers?description=%s", projectID, description), &trigger)
	if err != nil {
		return trigger, err
	}

	return trigger, nil
}

// GetTrigger 获取仓库触发器
func (c *Client) GetTrigger(projectID int) (triggerToken string, err error) {
	var triggers []Trigger
	err = c.GetResourceList(fmt.Sprintf("/projects/%v/triggers", projectID), &triggers)
	if err != nil {
		return triggerToken, err
	}

	for _, trigger := range triggers {
		if trigger.Token != "" {
			return trigger.Token, nil
		}
	}

	return
}
