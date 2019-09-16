package gitlab

import (
	"strconv"
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

// File 仓库文件列表信息
type File struct {
	ID   string `json:"id,omitempty"`   // commit ID
	Name string `json:"name,omitempty"` // 文件或目录名
	Type string `json:"type,omitempty"` // 类型: tree, blob
	Path string `json:"path,omitempty"` // 路径
	Mode string `json:"mode,omitempty"`
}

// CreateFileOptions 创建文件选项
type CreateFileOptions struct {
	Branch        string   `json:"branch,omitempty"`         // 分支名称
	AuthorEmail   string   `json:"author_email,omitempty"`   // 提交者Email
	AuthorName    string   `json:"author_name,omitempty"`    // 提交者
	Actions       []Action `json:"actions,omitempty"`        // 动作
	CommitMessage string   `json:"commit_message,omitempty"` // 提交消息
}

// Action 动作
type Action struct {
	Action   string `json:"action,omitempty"`    // 动作，包含create,delete,move,update,chmod
	FilePath string `json:"file_path,omitempty"` // 提交文件的完整路径. Ex app/main.go
	Content  string `json:"content,omitempty"`   // 文件内容
	Encoding string `json:"encoding,omitempty"`  // text or base64,默认text
}

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

// Trigger 触发器信息
type Trigger struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	LastUsed    string `json:"last_used"`
	Token       string `json:"token"`
	UpdatedAt   string `json:"updated_at"`
	Owner       struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Username  string `json:"username"`
		State     string `json:"state"`
		AvatarURL string `json:"avatar_url"`
		WebURL    string `json:"web_url"`
	}
}

// GetRepRootList 获取仓库根目录文件和目录列表
func (c *Client) GetRepRootList(projectID int, branch string) ([]File, error) {
	url := fmt.Sprintf("%s/api/v4/projects/%v/repository/tree?per_page=100&ref=%s", c.BaseURL, projectID, branch)
	client := &http.Client{}

	resp, err := httpGetRequest(url, c.AccessToken, client)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var files []File
	err = json.Unmarshal(body, &files)
	if err != nil {
		return nil, err
	}

	return files, nil
}

// CheckCIFile 检查gitlab仓库根目录文件是否存在
func (c *Client) CheckCIFile(projectID int, branch string) (bool, error) {
	files, err := c.GetRepRootList(projectID, branch)
	if err != nil {
		return false, err
	}

	for _, file := range files {
		if file.Type == "blob" && file.Name == ".gitlab-ci.yml" {
			return true, nil
		}
	}

	return false, nil

}

// httpRequest http请求，返回body
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
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("request error response status code %v", resp.StatusCode)
	}

	return resp, err
}

// AnalysisRepLanguage 分析存储库语言
func (c *Client) AnalysisRepLanguage(projectID int, branch string) (string, error) {
	var language string
	files, err := c.GetRepRootList(projectID, branch)
	if err != nil {
		return language, err
	}

	for _, file := range files {
		if file.Type == "blob" {
			if file.Name == "package.json" {
				language = "html"
				break
			}
			if file.Name == "go.mod" || file.Name == "Gopkg.toml" {
				language = "go"
				break
			}
			if file.Name == "composer.json" {
				language = "php"
				break
			}
		}
	}

	return language, nil
}

// CreateFile 仓库创建文件,其中files参数为需要创建的文件信息,key:文件路径; value:文件内容
func (c *Client) CreateFile(projectID int, branch, commitMsg string, files map[string]string) (response string, err error) {
	client := http.Client{}
	var actions []Action
	for k, v := range files {
		actions = append(actions, Action{Action: "create", FilePath: k, Content: v})
	}

	url := fmt.Sprintf("%v/api/v4/projects/%v/repository/commits", c.BaseURL, projectID)
	cf := &CreateFileOptions{
		Branch:        branch,
		Actions:       actions,
		CommitMessage: commitMsg,
	}
	reqBody, err := json.Marshal(cf)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(reqBody)))
	if err != nil {
		return
	}

	req.Header.Add("Private-Token", c.AccessToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	return string(body), nil
}

// GetTrigger 获取仓库触发器
func (c *Client) GetTrigger(projectID int) (triggerToken string, err error) {
	url := fmt.Sprintf("%s/api/v4/projects/%v/triggers", c.BaseURL, projectID)
	client := &http.Client{}
	var triggers []Trigger

	resp, err := httpGetRequest(url, c.AccessToken, client)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &triggers)
	if err != nil {
		return
	}

	for _, trigger := range triggers {
		if trigger.Token != "" {
			return trigger.Token, nil
		}
	}

	return
}

// ListGroups 获取所有组
func (c *Client) ListGroups() ([]Group, error) {
	var groups []Group
	page := 1
	for {
		var g []Group
		url := fmt.Sprintf("%s/api/v4/groups?page=%v", c.BaseURL, page)
		client := &http.Client{}
	
		resp, err := httpGetRequest(url, c.AccessToken, client)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
	
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
	
		// var groups []Group
		err = json.Unmarshal(body, &g)
		if err != nil {
			return nil, err
		}

		groups = append(groups, g...)

		totalPages, err := strconv.Atoi(resp.Header.Get("X-Total-Pages"))
		if err != nil {
			return nil, err
		}

		if page == totalPages {
			break
		}

		page++
	}

	return groups, nil
}

// ListSubGroups 获取指定组下的子组
func (c *Client) ListSubGroups(groupID int) ([]Group, error) {
	url := fmt.Sprintf("%s/api/v4/groups/%v/subgroups", c.BaseURL, groupID)
	client := &http.Client{}

	resp, err := httpGetRequest(url, c.AccessToken, client)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var subGroups []Group
	err = json.Unmarshal(body, &subGroups)
	if err != nil {
		return nil, err
	}

	return subGroups, nil
}

// ListGroupsProjects 获取指定组下面的仓库
func (c *Client) ListGroupsProjects(groupID int) ([]Project, error) {
	url := fmt.Sprintf("%s/api/v4/groups/%v/projects?simple=yes", c.BaseURL, groupID)
	client := &http.Client{}

	resp, err := httpGetRequest(url, c.AccessToken, client)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var projects []Project
	err = json.Unmarshal(body, &projects)
	if err != nil {
		return nil, err
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

	resp, err := client.Do(req)
	if err != nil {
		return subGroup, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return subGroup, err
	}

	if err = json.Unmarshal(body, &subGroup); err != nil {
		return subGroup, err
	}

	return subGroup, nil
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

	resp, err := client.Do(req)
	if err != nil {
		return project, err
	}
	defer resp.Body.Close()

	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return project, err
	}

	if err = json.Unmarshal(respbody, &project); err != nil {
		return project, err
	}

	return project, nil
}

// CreateTrigger 创建触发器
func (c *Client) CreateTrigger(projectID int) (tiggerToken string, err error) {
	client := http.Client{}

	url := fmt.Sprintf("%s/api/v4/projects/%v/triggers?description=deploy", strings.TrimSuffix(c.BaseURL, "/"), projectID)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("Private-Token", c.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var tigger Trigger
	if err = json.Unmarshal(body, &tigger); err != nil {
		return
	}

	return tigger.Token, nil
}
