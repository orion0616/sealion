package todoist

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/orion0616/sealion/util"
)

// Project express a todoist project
type Project struct {
	ID   int64
	Name string
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
	data, err := util.DecodeResponse(resp)
	if err != nil {
		return nil, err
	}
	projects, ok := data["projects"]
	if !ok {
		return nil, errors.New("no key named 'projects' in data")
	}
	castedProjects, isCasted := projects.([]interface{})
	if !isCasted {
		return nil, errors.New("failed to convert projects")
	}

	var p []Project
	for _, project := range castedProjects {
		var castedProject map[string]interface{}
		castedProject, isCasted = project.(map[string]interface{})
		if !isCasted {
			return nil, errors.New("failed to convert a project")
		}
		var pro Project
		for k, v := range castedProject {
			if k == "name" {
				var name string
				name, ok := v.(string)
				if !ok {
					return nil, errors.New("failed to get project name")
				}
				pro.Name = name
			}
			if k == "id" {
				var id json.Number
				id, ok := v.(json.Number)
				if !ok {
					return nil, errors.New("failed to get project id")
				}
				pro.ID, _ = id.Int64()
			}
		}
		p = append(p, pro)
	}
	return p, nil
}
