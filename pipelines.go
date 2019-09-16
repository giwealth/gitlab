package gitlab

import (
	"strconv"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"fmt"
)

// Pipeline gitlab pipeline struct
type Pipeline struct {
	ID int `json:"id"`
	Status string `json:"status"`
	Ref string `json:"ref"`
	Sha string `json:"sha"`
	WebURL string `json:"web_url"`
}

// ListPipelines list project pipelines
func (c *Client) ListPipelines(projectID int) ([]Pipeline, error) {
	var pipelines []Pipeline
	page := 1
	for {
		url := fmt.Sprintf("%s/api/v4/projects/%v/pipelines?page=%v", c.BaseURL, projectID, page)
		client := &http.Client{}

		res, err := httpGetRequest(url, c.AccessToken, client)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		var p []Pipeline
		err = json.Unmarshal(body, &p)
		if err != nil {
			return nil, err
		}
		pipelines = append(pipelines, p...)

		totalPages, err := strconv.Atoi(res.Header.Get("X-Total-Pages"))
		if err != nil {
			return nil, err
		}

		if page == totalPages {
			break
		}

		page++
	}

	return pipelines, nil
}