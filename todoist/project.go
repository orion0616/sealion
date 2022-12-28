package todoist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	resp, err := c.do("[\"projects\"]", "")
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
		fmt.Println("Failed to unmarshal in ExtractProjects. data = " + string(data))
		return nil, err
	}
	return getProjectsResult.Projects, nil
}
