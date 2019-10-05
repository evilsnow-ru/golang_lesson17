package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"golang_lesson17/internal/service"
	"os"
	"strconv"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts server at custom port. Default port is 9000.",
	Run: func(cmd *cobra.Command, args []string) {
		var port uint32

		if len(args) == 0 {
			port = defaultPort
		} else {
			parsedValue, err := strconv.Atoi(args[0])

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			port = uint32(parsedValue)
		}
		fmt.Println(service.StartServer(port))
	},
}
