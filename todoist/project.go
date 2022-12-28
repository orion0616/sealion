package todoist

import (
	"encoding/json"
	"fmt"
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
	// resp, err := c.HTTPClient.PostForm("https://api.todoist.com/sync/v9/sync", values)

	req, err := http.NewRequest("POST", "https://api.todoist.com/sync/v9/sync", nil)

	req.Header = map[string][]string{
		"Authorization": {"Bearer " + c.Token},
	}

	req.Header.Add("Authorization", "Bearer "+c.Token)
	values := url.Values{}
	values.Add("sync_token", "*")
	values.Add("resource_types", "[\"projects\"]")
	req.Form = values
	resp, err := c.HTTPClient.Do(req)
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
