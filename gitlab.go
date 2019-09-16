package gitlab

import (
	"strconv"
	"io/ioutil"
	"encoding/json"
	// "strings"
	"fmt"
	"net/http"
	"net/url"
)

// Client gitlab api client
type Client struct {
	BaseURL     string
	AccessToken string
	APIVersionPath string
}

// New create gitlab client
func New(baseURL, accessToken string) Client {
	return Client{
		BaseURL: baseURL,
		AccessToken: accessToken,
		APIVersionPath: "/api/v4",
	}	
}

// httpRequest http请求，返回response
func httpGetRequest(urlAddr, token string, c *http.Client) (*http.Response, error) {
	req, err := http.NewRequest("GET", urlAddr, nil)
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

// GetAPI get gitlab api
func (c *Client) GetAPI(addr string, v interface{}) error {
	var list []interface{}
	page := 1
	for {
		u, err := url.Parse(addr)
		if err != nil {
			return err
		}
		q := u.Query()
		q.Set("page", fmt.Sprint(page))
		u.RawQuery = q.Encode()

		httpClient := &http.Client{}

		req, err := http.NewRequest("GET", u.String(), nil)
		if err != nil {
			return nil
		}

		req.Header.Set("Private-Token", c.AccessToken)

		res, err := httpClient.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode < 200 || res.StatusCode > 299 {
			return fmt.Errorf("request error response status code %v", res.StatusCode)
		}

		var l []interface{}
		f, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		err = json.Unmarshal(f, &l)
		if err != nil {
			return err
		}

		list = append(list, l...)

		totalPages, err := strconv.Atoi(res.Header.Get("X-Total-Pages"))
		if err != nil {
			return err
		}

		if page == totalPages {
			break
		}

		page++
	}

	v = list

	return nil
}