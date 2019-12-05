/*
Copyright © 2019 orion0616 earth.nobu.light@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/orion0616/sealion/todoist"
	"github.com/spf13/cobra"
)

// getProjectsCmd represents the getProjects command
var getProjectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "list of todoist project",

	Run: func(cmd *cobra.Command, args []string) {
		client, err := todoist.NewClient()
		if err != nil {
			fmt.Println(err)
		}
		values := url.Values{}
		values.Add("token", client.Token)
		values.Add("sync_token", "*")
		values.Add("resource_types", "[\"projects\"]")
		resp, err := client.HTTPClient.PostForm("https://todoist.com/api/v8/sync", values)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		projects, err := extractProjects(resp)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("ID         NAME\n")
		for _, project := range projects {
			fmt.Printf("%-10d %s\n", project.ID, project.Name)
		}
	},
}

func extractProjects(resp *http.Response) ([]todoist.Project, error) {
	data, err := decodeResponse(resp)
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

	var p []todoist.Project
	for _, project := range castedProjects {
		var castedProject map[string]interface{}
		castedProject, isCasted = project.(map[string]interface{})
		if !isCasted {
			return nil, errors.New("failed to convert a project")
		}
		var pro todoist.Project
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

func decodeResponse(resp *http.Response) (map[string]interface{}, error) {
	var data map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	err := decoder.Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func init() {
	getCmd.AddCommand(getProjectsCmd)
}
