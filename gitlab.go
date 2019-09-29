package gitlab

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	apiVersionPath = "/api/v4"
)

// Client gitlab api client
type Client struct {
	BaseURL     string
	AccessToken string
}

// GetResource get gitlab resource detail
func (c *Client) GetResource(api string, v interface{}) error {
	httpClient := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s%s", c.BaseURL, apiVersionPath, api), nil)
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
		return fmt.Errorf("request error response status code %v, response:%v", res.StatusCode, string(body))
	}

	err = json.Unmarshal(body, &v)
	if err != nil {
		return err
	}

	return nil
}

// GetResourceList get gitlab resource list
func (c *Client) GetResourceList(api string, v interface{}) error {
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

// CreateResource 创建
func (c *Client) CreateResource(api string, v interface{}) error {
	client := http.Client{}

	u, err := url.Parse(fmt.Sprintf("%s%s%s", c.BaseURL, apiVersionPath, api))
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Private-Token", c.AccessToken)
	fmt.Println(u.String())

	res, err := client.Do(req)
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

	if err = json.Unmarshal(body, &v); err != nil {
		return err
	}

	return nil
}
