package gitlab

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// TriggerPipelineResponse 通过API触发管道响应结构
type TriggerPipelineResponse struct {
	ID         int    `json:"id"`
	Sha        string `json:"sha"`
	Ref        string `json:"master"`
	Status     string `json:"status"`
	WebURL     string `json:"web_url"`
	BeforeSha  string `json:"before_sha"`
	Tag        bool   `json:"tag"`
	YamlErrors string `json:"yaml_errors"`
	User       struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Username  string `json:"username"`
		State     string `json:"state"`
		AvatarURL string `json:"avatar_url"`
		WebURL    string `json:"web_url"`
	} `json:"user"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	StartedAt   string `json:"started_at"`
	FinishedAt  string `json:"finished_at"`
	CommittedAt string `json:"committed_at"`
	Duration    int    `json:"duration"`
	Coverage    bool   `json:"coverage"`
}

// TriggerPipeline 通过API触发管道, projectID为触发项目ID
func (c *Client) TriggerPipeline(projectID int, triggerToken string, variables map[string]string) error {
	client := http.Client{}

	// 使用form提交数据
	data := url.Values{}
	data.Set("token", triggerToken)
	data.Set("ref", "master") // 此处为部署仓库的ref

	for k, v := range variables {
		data.Set(fmt.Sprintf("variables[%s]", k), v)
	}

	url := fmt.Sprintf("%s/api/v4/projects/%v/trigger/pipeline", strings.TrimSuffix(c.BaseURL, "/"), projectID)
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("request error response status code %v", res.StatusCode)
	}

	return nil
}
