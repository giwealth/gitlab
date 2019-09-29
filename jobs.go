package gitlab

import (
	"fmt"
)

// Job gitlab job
type Job struct {
	ID           int     `json:"id"`
	Status       string  `json:"status"`
	Stage        string  `json:"stage"`
	Name         string  `json:"name"`
	Ref          string  `json:"ref"`
	Tag          bool    `json:"tag"`
	Coverage     string  `json:"coverage"`
	AllowFailure bool    `json:"allow_failure"`
	CreatedAt    string  `json:"created_at"`
	StartedAt    string  `json:"started_at"`
	FinishedAt   string  `json:"finished_at"`
	Duration     float32 `json:"duration"`
	User         struct {
		ID           int    `json:"id"`
		Name         string `json:"name"`
		Username     string `json:"username"`
		State        string `json:"state"`
		AvatarURL    string `json:"avatar_url"`
		WebURL       string `json:"web_url"`
		CreatedAT    string `json:"created_at"`
		Bio          string `json:"bio"`
		Location     string `json:"location"`
		PublicEmail  string `json:"public_email"`
		Skype        string `json:"skype"`
		Linkedin     string `json:"linkedin"`
		Twitter      string `json:"twitter"`
		WebsiteURL   string `json:"website_url"`
		Organization string `json:"organization"`
	} `json:"user"`
	Commit struct {
		ID             string   `json:"id"`
		ShortID        string   `json:"short_id"`
		CreatedAt      string   `json:"created_at"`
		ParentIds      []string `json:"parent_ids"`
		Title          string   `json:"title"`
		Message        string   `json:"message"`
		AuthorName     string   `json:"author_name"`
		AuthorEmail    string   `json:"author_email"`
		AuthoredDate   string   `json:"authored_date"`
		CommitterName  string   `json:"committer_name"`
		CommitterEmail string   `json:"committer_email"`
		CommittedDate  string   `json:"committed_date"`
	} `json:"commit"`
	Pipeline struct {
		ID     int    `json:"id"`
		Sha    string `json:"sha"`
		Ref    string `json:"ref"`
		Status string `json:"status"`
		WebURL string `json:"web_url"`
	} `json:"pipeline"`
	WebURL    string     `json:"web_url"`
	Artifacts []Artifact `json:"artifacts"`
	Runner    struct {
		ID          int    `json:"id"`
		Description string `json:"description"`
		IPAddress   string `json:"ip_address"`
		Active      bool   `json:"active"`
		IsShared    bool   `json:"is_shared"`
		Name        string `json:"name"`
		Online      bool   `json:"online"`
		Status      string `json:"status"`
	} `json:"runner"`
	ArtifactsExpireAt string `json:"artifacts_expire_at"`
}

// Artifact job artifact
type Artifact struct {
	FileType   string `json:"file_type"`
	Size       int64  `json:"size"`
	Filename   string `json:"filename"`
	FileFormat string `json:"file_format"`
}

// ListPipelineJobs get a list of jobs for a pipeline
func (c *Client) ListPipelineJobs(projectID, pipelineID int) ([]Job, error) {
	var jobs []Job
	err := c.GetResourceList(fmt.Sprintf("/projects/%v/pipelines/%v/jobs", projectID, pipelineID), &jobs)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

// ListProjectJobs get a list of jobs in a project
func (c *Client) ListProjectJobs(projectID int) ([]Job, error) {
	var jobs []Job
	err := c.GetResourceList(fmt.Sprintf("/projects/%v/jobs", projectID), &jobs)
	if err != nil {
		return jobs, err
	}

	return jobs, nil
}

// ActionJob play or retry a job
func (c *Client) ActionJob(projectID, jobID int, action string) (Job, error) {
	var job Job
	err := c.CreateResource(fmt.Sprintf("/projects/%v/jobs/%v/%s", projectID, jobID, action), &job)
	if err != nil {
		return job, err
	}

	return job, nil
}

// GetJob get a single job
func (c *Client) GetJob(projectID, jobID int) (Job, error) {
	var job Job
	err := c.GetResource(fmt.Sprintf("/projects/%v/jobs/%v", projectID, jobID), &job)
	if err != nil {
		return job, err
	}

	return job, nil
}