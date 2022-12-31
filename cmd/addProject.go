/*
Copyright © 2022 orion0616 earth.nobu.light@gmail.com

*/
package cmd

import (
	"fmt"

	"github.com/orion0616/sealion/todoist"
	"github.com/spf13/cobra"
)

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Add a project",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("You need to give one argument")
			return
		}
		client, err := todoist.NewClient()
		if err != nil {
			fmt.Println(err)
			return
		}
		err = client.AddProject(args[0])
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("added a project successfully")
		}
	},
}

func init() {
	addCmd.AddCommand(projectCmd)
}
