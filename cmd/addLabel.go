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

// addLabelCmd represents the addLabel command
var addLabelCmd = &cobra.Command{
	Use:   "labels",
	Short: "Add labels to a task",

	Run: func(cmd *cobra.Command, args []string) {
		project, err := cmd.Flags().GetString("project")
		if err != nil {
			fmt.Println(err)
			return
		}
		client, err := todoist.NewClient()
		if err != nil {
			fmt.Println(err)
			return
		}
		err = client.AddLabels(args, project)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("added label successfully")
		}
	},
}

func init() {
	addCmd.AddCommand(addLabelCmd)
	addLabelCmd.Flags().StringP("project", "p", "", "select a project")
}
