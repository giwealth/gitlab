package gitlab

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Project gitlab project info
type Project struct {
	ID                int      `json:"id"`
	Description       string   `json:"description"`
	DefaultBranch     string   `json:"default_branch"`
	SSHURLToRepo      string   `json:"ssh_url_to_repo"`
	HTTPURLToRepo     string   `json:"http_url_to_repo"`
	WebURL            string   `json:"web_url"`
	ReadmeURL         string   `json:"readme_url"`
	TagList           []string `json:"tag_list"`
	Name              string   `json:"name"`
	NameWithNamespace string   `json:"name_with_namespace"`
	Path              string   `json:"path"`
	PathWithNamespace string   `json:"path_with_namespace"`
	CreatedAt         string   `json:"created_at"`
	LastActivityAt    string   `json:"last_activity_at"`
	ForksCount        int      `json:"forks_count"`
	AvatarURL         string   `json:"avatar_url"`
	StarCount         int      `json:"star_count"`
	Namespace         struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Path     string `json:"path"`
		Kind     string `json:"kind"`
		FullPath string `json:"full_path"`
		ParentID int    `json:"parent_id"`
	}
	TriggerToken string `json:"trigger_token"`
}

// CreateProject 新建仓库
func (c *Client) CreateProject(projectName string, namespaceID int) (Project, error) {
	var project Project
	client := http.Client{}

	url := fmt.Sprintf(
		"%s/api/v4/projects?name=%s&namespace_id=%v&visibility=private",
		strings.TrimSuffix(c.BaseURL, "/"),
		projectName,
		namespaceID,
	)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return project, err
	}
	req.Header.Set("Private-Token", c.AccessToken)

	res, err := client.Do(req)
	if err != nil {
		return project, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return project, fmt.Errorf("request error response status code %v", res.StatusCode)
	}

	respbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return project, err
	}

	if err = json.Unmarshal(respbody, &project); err != nil {
		return project, err
	}

	return project, nil
}
