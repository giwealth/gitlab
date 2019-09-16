package gitlab

import (
	// "net/url"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Group 组信息
type Group struct {
	ID                    int    `json:"id"`
	Name                  string `json:"name"`
	Path                  string `json:"path"`
	Description           string `json:"description"`
	Visibility            string `json:"visibility"`
	LFSEnabled            bool   `json:"lfs_enabled"`
	AvatarURL             string `json:"avatar_url"`
	WebURL                string `json:"web_url"`
	RequestAccessEnabled  bool   `json:"request_access_enabled"`
	FullName              string `json:"full_name"`
	FullPath              string `json:"full_path"`
	FileTemplateProjectID int    `json:"file_template_project_id"`
	ParentID              int    `json:"parent_id"`
}

// ListGroups 获取所有组
func (c *Client) ListGroups() ([]Group, error) {
	var groups []Group
	err := getAPI("/groups", c.AccessToken, &groups)
	if err != nil {
		return nil, err
	}
	// page := 1
	// for {
	// 	var g []Group
	// 	url := fmt.Sprintf("%s/api/v4/groups?page=%v", c.BaseURL, page)
	// 	client := &http.Client{}

	// 	resp, err := httpGetRequest(url, c.AccessToken, client)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	defer resp.Body.Close()

	// 	body, err := ioutil.ReadAll(resp.Body)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	// var groups []Group
	// 	err = json.Unmarshal(body, &g)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	groups = append(groups, g...)

	// 	totalPages, err := strconv.Atoi(resp.Header.Get("X-Total-Pages"))
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	if page == totalPages {
	// 		break
	// 	}

	// 	page++
	// }

	return groups, nil
}

// ListSubGroups 获取指定组下的子组
func (c *Client) ListSubGroups(groupID int) ([]Group, error) {
	var subGroups []Group
	page := 1
	for {
		url := fmt.Sprintf("%s/api/v4/groups/%v/subgroups?page=%v", c.BaseURL, groupID, page)
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

		var g []Group
		err = json.Unmarshal(body, &g)
		if err != nil {
			return nil, err
		}
		subGroups = append(subGroups, g...)

		totalPages, err := strconv.Atoi(res.Header.Get("X-Total-Pages"))
		if err != nil {
			return nil, err
		}

		if page == totalPages {
			break
		}

		page++
	}

	return subGroups, nil
}

// ListGroupsProjects 获取指定组下面的仓库
func (c *Client) ListGroupsProjects(groupID int) ([]Project, error) {
	var projects []Project
	page := 1
	for {
		url := fmt.Sprintf("%s/api/v4/groups/%v/projects?simple=yes&page=%v", c.BaseURL, groupID, page)
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

		var p []Project
		err = json.Unmarshal(body, &p)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p...)

		totalPages, err := strconv.Atoi(res.Header.Get("X-Total-Pages"))
		if err != nil {
			return nil, err
		}

		if page == totalPages {
			break
		}

		page++
	}

	return projects, nil
}

// CreateSubGroup 增加子组
func (c *Client) CreateSubGroup(newGroupName string, parentID int) (Group, error) {
	var subGroup Group
	client := http.Client{}

	url := fmt.Sprintf("%s/api/v4/groups?name=%s&path=%s&parent_id=%v&visibility=private", c.BaseURL, newGroupName, newGroupName, parentID)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return subGroup, err
	}
	req.Header.Set("Private-Token", c.AccessToken)

	res, err := client.Do(req)
	if err != nil {
		return subGroup, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return subGroup, fmt.Errorf("request error response status code %v", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return subGroup, err
	}

	if err = json.Unmarshal(body, &subGroup); err != nil {
		return subGroup, err
	}

	return subGroup, nil
}