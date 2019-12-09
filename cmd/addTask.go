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

// addTaskCmd represents the addTask command
var addTaskCmd = &cobra.Command{
	Use:   "task",
	Short: "add tasks written in a file to a project",
	Run: func(cmd *cobra.Command, args []string) {
		fileName, err := cmd.Flags().GetString("file")
		if err != nil {
			fmt.Println(err)
			return
		}
		client, err := todoist.NewClient()
		if err != nil {
			fmt.Println(err)
			return
		}
		err = client.AddTasks(fileName)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	addCmd.AddCommand(addTaskCmd)
	addTaskCmd.Flags().StringP("file", "f", "", "select a file")
}
