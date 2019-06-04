package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

const longDescription = `
    _                                     _ 
   / \  _   _  ____ _ __   ____ _ ______ (_)
  / _ \| | | |/ _  | '_ \ / _  | '_   _ \| |
 / ___ \ |_| | (_| | | | | (_| | | | | | | |
/_/   \_\__, |\__,_|_| |_|\__,_|_| |_| |_|_|
        |___/       

Ayanami is a FaaS-like framework for your own infrastructure.
To make a project, please use --help for seeing available commands`

var rootCmd = &cobra.Command{
	Use:   "ayanami",
	Short: "FaaS-like framework for your own infrastructure.",
	Long:  longDescription,
}

// Execute ...
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
