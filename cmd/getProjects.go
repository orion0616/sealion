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
	"fmt"

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
		projects, err := client.GetProjects()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf(createResult(projects))
	},
}

func createResult(projects []todoist.Project) string {
	var bytes []byte
	bytes = append(bytes, "ID         NAME\n"...)
	for _, project := range projects {
		str := fmt.Sprintf("%-10d %s\n", project.ID, project.Name)
		bytes = append(bytes, str...)
	}
	return string(bytes)
}

func init() {
	getCmd.AddCommand(getProjectsCmd)
}
