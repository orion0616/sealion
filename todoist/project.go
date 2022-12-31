package todoist

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/orion0616/sealion/util"
)

// GetProjectsResult express a result of getting all projects
type GetProjectsResult struct {
	SyncToken     string `json:"sync_token"`
	TempIDMapping struct {
	} `json:"temp_id_mapping"`
	FullSync bool `json:"full_sync"`
	Projects []Project
}

// Project express a todoist project
type Project struct {
	IsFavorite bool        `json:"is_favorite"`
	Color      string      `json:"color"`
	Collapsed  bool        `json:"collapsed"`
	ChildOrder int         `json:"child_order"`
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	IsDeleted  bool        `json:"is_deleted"`
	ParentID   interface{} `json:"parent_id"`
	Shared     bool        `json:"shared"`
	IsArchived bool        `json:"is_archived"`
	SyncID     interface{} `json:"sync_id"`
	ViewStyle  string      `json:"view_style"`
}

// GetProjects returns a list of todoist projects
func (c *Client) GetProjects() ([]Project, error) {
	resp, err := c.do("[\"projects\"]", "", "")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	projects, err := ExtractProjects(resp)
	if err != nil {
		return nil, err
	}
	return projects, err
}

// ExtractProjects extracts projects from http.Response
func ExtractProjects(resp *http.Response) ([]Project, error) {
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var getProjectsResult GetProjectsResult
	if err := json.Unmarshal(data, &getProjectsResult); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal in ExtractProjects. data = " + string(data))
	}
	return getProjectsResult.Projects, nil
}

func (c *Client) AddProject(project string) error {
	command, err := createProjectAddCommand(project)
	if err != nil {
		return err
	}
	resp, err := c.do("", command, "")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.Status != "200 OK" {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to add a project.\nStatus -> %s.\nBody -> %s\n", resp.Status, string(b))
	}
	return nil
}

func createProjectAddCommand(project string) (string, error) {
	uuid1, err := util.CreateUUID()
	if err != nil {
		return "", err
	}
	uuid2, err := util.CreateUUID()
	if err != nil {
		return "", err
	}
	command := fmt.Sprintf("[{\"type\": \"project_add\", \"temp_id\": \"%s\", \"uuid\": \"%s\",\"args\": {\"name\": \"%s\"}}]",
		uuid1, uuid2, project)
	return command, nil
}
