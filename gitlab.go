package gitlab

import (
	// "strings"
	"fmt"
	"net/http"
)

// Client gitlab api client
type Client struct {
	BaseURL     string
	AccessToken string
}

// httpRequest http请求，返回response
func httpGetRequest(url, token string, c *http.Client) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Private-Token", token)

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("request error response status code %v", resp.StatusCode)
	}

	return resp, err
}
