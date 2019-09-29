package gitlab

import (
	"fmt"
)

// Project gitlab project info
type Project struct {
	ID                int      `json:"id"`
	Description       string   `json:"description"`
	Visibility		  string   `json:"visibility"`
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
	} `json:"namespace"`
	TriggerToken string `json:"trigger_token"`
	DeployProjectID int `json:"deploy_project_id"`
}

// CreateProject 新建仓库
func (c *Client) CreateProject(projectName string, namespaceID int) (Project, error) {
	var project Project
	err := c.CreateResource(fmt.Sprintf("/projects?name=%s&namespace_id=%v&visibility=private", projectName, namespaceID), &project)
	if err != nil {
		return project, err
	}

	return project, nil
}

// ListProjects list all projects
func (c *Client) ListProjects() ([]Project, error) {
	var projects []Project

	err := c.GetResourceList("/projects?per_page=100", &projects)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

// GetProject get single project
func (c *Client) GetProject(projectID int) (Project, error) {
	var project Project
	err := c.GetResource(fmt.Sprintf("/projects/%v", projectID), &project)
	if err != nil {
		return project, err
	}

	return project, nil
}
