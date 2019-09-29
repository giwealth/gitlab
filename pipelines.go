package gitlab

import (
	"encoding/json"
	"strconv"
	"strings"
	"io/ioutil"
	"net/http"
	"net/url"
	"fmt"
)

// Pipeline gitlab pipeline struct
type Pipeline struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
	Ref    string `json:"ref"`
	Sha    string `json:"sha"`
	WebURL string `json:"web_url"`
}

// Variable gitlab pipeline vriable
type Variable struct {
	Key          string `json:"key"`
	Value        string `json:"value"`
	VariableType string `json:"variable_type"`
}

// ListPipelines list project pipelines
func (c *Client) ListPipelines(projectID int) ([]Pipeline, error) {
	var pipelines []Pipeline
	err := c.GetResourceList(fmt.Sprintf("/projects/%v/pipelines", projectID), &pipelines)
	if err != nil {
		return nil, err
	}

	return pipelines, nil
}

// ListPipelineVar get variables of a pipeline
func (c *Client) ListPipelineVar(projectID, pipelineID int) ([]Variable, error) {
	var variables []Variable
	err := c.GetResourceList(fmt.Sprintf("/projects/%v/pipelines/%v/variables", projectID, pipelineID), &variables)
	if err != nil {
		return nil, err
	}

	return variables, nil
}

// GetPipeline get a single pipeline
func (c *Client) GetPipeline(projectID, pipelineID int) (Pipeline, error) {
	var pipeline Pipeline
	err := c.GetResource(fmt.Sprintf("/projects/%v/pipelines/%v", projectID, pipelineID), &pipeline)
	if err != nil {
		return pipeline, err
	}

	return pipeline, nil
}

// GetResourceList get gitlab resource list
func (c *Client) requestPiplines(api string, v interface{}) error {
	var response string
	page := 1
	for {
		u, err := url.Parse(fmt.Sprintf("%s%s%s", c.BaseURL, apiVersionPath, api))
		if err != nil {
			return err
		}
		q := u.Query()
		q.Set("page", fmt.Sprint(page))
		u.RawQuery = q.Encode()

		httpClient := &http.Client{}

		req, err := http.NewRequest("GET", u.String(), nil)
		if err != nil {
			return err
		}

		req.Header.Set("Private-Token", c.AccessToken)

		res, err := httpClient.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		if res.StatusCode < 200 || res.StatusCode > 299 {
			return fmt.Errorf("request error response status code %v, response:%s", res.StatusCode, string(body))
		}

		response += string(body)

		s := res.Header.Get("X-Total-Pages")
		if s == "" {
			break
		}
		totalPages, err := strconv.Atoi(s)
		if err != nil {
			return err
		}

		if page == totalPages {
			break
		}

		page++
	}

	response = strings.ReplaceAll(response, "][", ",")

	err := json.Unmarshal([]byte(response), &v)
	if err != nil {
		return err
	}

	return nil
}