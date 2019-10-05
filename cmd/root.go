package cmd

import "github.com/spf13/cobra"

const defaultPort = 9000

var root = &cobra.Command{}

func init() {
	root.AddCommand(startCmd)
	root.AddCommand(eventAddCmd)
	root.AddCommand(getEventCmd)
}

func Execute() error {
	return root.Execute()
}
