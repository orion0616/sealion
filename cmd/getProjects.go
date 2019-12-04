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
	// "io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

// getProjectsCmd represents the getProjects command
var getProjectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "list of todoist project",
	Run: func(cmd *cobra.Command, args []string) {
		token, err := getToken()
		if err != nil {
			fmt.Println(err)
			return
		}
		values := url.Values{}
		values.Add("token", token)
		values.Add("sync_token", "*")
		values.Add("resource_types", "[\"projects\"]")
		resp, err := http.PostForm("https://todoist.com/api/v8/sync", values)
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
		fmt.Println(projects)
	},
}

func getToken() (string, error) {
	token := os.Getenv("TODOIST_TOKEN")
	if token == "" {
		return "", errors.New("token is not set")
	}
	return token, nil
}

// Project express a todoist project
type Project struct {
	Id   int64
	Name string
}

// extractProjects have to return ([]Project, error) TODO
func extractProjects(resp *http.Response) (Project, error) {
	// TODO id is converted wrongly.
	var data map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return Project{}, err
	}
	projects, ok := data["projects"]
	if ok == false {
		return Project{}, err
	}
	castedProjects, isCasted := projects.([]interface{})
	if !isCasted {
		return Project{}, errors.New("failed to convert projects")
	}
	for _, project := range castedProjects {
		var castedProject map[string]interface{}
		castedProject, isCasted = project.(map[string]interface{})
		if !isCasted {
			return Project{}, errors.New("failed to convert a project")
		}
		fmt.Println(castedProject)
	}
	return Project{}, nil
}

func init() {
	getCmd.AddCommand(getProjectsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getProjectsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getProjectsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
