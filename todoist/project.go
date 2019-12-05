package todoist

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
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
	IsFavorite   int         `json:"is_favorite"`
	Color        int         `json:"color"`
	Collapsed    int         `json:"collapsed"`
	InboxProject bool        `json:"inbox_project,omitempty"`
	ChildOrder   int         `json:"child_order"`
	ID           int         `json:"id"`
	Name         string      `json:"name"`
	IsDeleted    int         `json:"is_deleted"`
	ParentID     interface{} `json:"parent_id"`
	LegacyID     int         `json:"legacy_id,omitempty"`
	Shared       bool        `json:"shared"`
	IsArchived   int         `json:"is_archived"`
}

// GetProjects returns a list of todoist projects
func (c *Client) GetProjects() ([]Project, error) {
	values := url.Values{}
	values.Add("token", c.Token)
	values.Add("sync_token", "*")
	values.Add("resource_types", "[\"projects\"]")

	resp, err := c.HTTPClient.PostForm("https://todoist.com/api/v8/sync", values)
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
		return nil, err
	}
	return getProjectsResult.Projects, nil
}
