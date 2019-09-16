package gitlab

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

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

// GetRepRootList 获取仓库根目录文件和目录列表
func (c *Client) GetRepRootList(projectID int, branch string) ([]File, error) {
	var files []File
	page := 1
	for {
		url := fmt.Sprintf("%s/api/v4/projects/%v/repository/tree?per_page=100&ref=%s&page=%v", c.BaseURL, projectID, branch, page)
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

		var f []File
		err = json.Unmarshal(body, &f)
		if err != nil {
			return nil, err
		}
		files = append(files, f...)

		totalPages, err := strconv.Atoi(res.Header.Get("X-Total-Pages"))
		if err != nil {
			return nil, err
		}

		if page == totalPages {
			break
		}

		page++
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
func (c *Client) CreateFile(projectID int, branch, commitMsg string, files map[string]string) error {
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
		return err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(reqBody)))
	if err != nil {
		return err
	}

	req.Header.Add("Private-Token", c.AccessToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("request error response status code %v", res.StatusCode)
	}

	return nil
}
