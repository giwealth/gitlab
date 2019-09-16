package gitlab

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Webhook 增加webhooks请求数据结构
type Webhook struct {
	ID                    int    `json:"id"`  // project_id
	URL                   string `json:"url"` // webhook url
	PushEvents            bool   `json:"push_events"`
	PipelineEvents        bool   `json:"pipeline_events"`
	EnableSSLVerification bool   `json:"enable_ssl_verification"`
	Token                 string `json:"token"` // access webhook token
}

// AddWebhooks 增加项目webhooks pushEventsURL, pipelineEventsURL,
func (c *Client) AddWebhooks(projectID int, webhooks []Webhook) error {
	api := fmt.Sprintf("%s/api/v4/projects/%v/hooks", c.BaseURL, projectID)

	for _, webhook := range webhooks {
		reqBody, err := json.Marshal(webhook)
		if err != nil {
			return err
		}

		client := &http.Client{}
		req, err := http.NewRequest("POST", api, strings.NewReader(string(reqBody)))
		req.Header.Set("Private-Token", c.AccessToken)
		req.Header.Set("Content-Type", "application/json")

		res, err := client.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode < 200 || res.StatusCode > 299 {
			return fmt.Errorf("request error response status code %v", res.StatusCode)
		}
	}

	return nil
}
