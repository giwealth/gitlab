package gitlab

import (
	"fmt"
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

	err := c.GetResourceList("/groups", &groups)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

// ListSubGroups 获取指定组下的子组
func (c *Client) ListSubGroups(groupID int) ([]Group, error) {
	var subGroups []Group

	err := c.GetResourceList(fmt.Sprintf("/groups/%v/subgroups", groupID), &subGroups)
	if err != nil {
		return nil, err
	}

	return subGroups, nil
}

// ListGroupsProjects 获取指定组下面的仓库
func (c *Client) ListGroupsProjects(groupID int) ([]Project, error) {
	var projects []Project

	err := c.GetResourceList(fmt.Sprintf("/groups/%v/projects", groupID), &projects)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

// CreateSubGroup 增加子组
func (c *Client) CreateSubGroup(newGroupName string, parentID int) (Group, error) {
	var subGroup Group

	err := c.CreateResource(fmt.Sprintf("/groups?name=%s&path=%s&parent_id=%v&visibility=private", newGroupName, newGroupName, parentID), &subGroup)
	if err != nil {
		return subGroup, err
	}

	return subGroup, nil
}

// GetGroup details of a group
func (c *Client) GetGroup(groupID int) (Group, error) {
	var group Group
	err := c.GetResource(fmt.Sprintf("/groups/%v", groupID), &group)
	if err != nil {
		return group, err
	}

	return group, nil
}
