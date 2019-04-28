package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

const longDescription = `Ayanami is a code deployment generator.
With Ayanami, you only need to define 3 things:
	* Functions
	* Templates
	* Function Compositions
Ayanami will do everything else for you`

func init() {
	rootCmd.PersistentFlags().StringP("template", "t", "templates", "template directory")
	rootCmd.PersistentFlags().StringP("service", "s", "services", "service directory")
	rootCmd.PersistentFlags().StringP("composition", "c", "compositions", "composition directory")
}

var rootCmd = &cobra.Command{
	Use:   "ayanami",
	Short: "Ayanami is a code deployment generator.",
	Long:  longDescription,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		fmt.Println(os.Getwd())
		fmt.Println(cmd.Flags().GetString("template"))
		fmt.Println(cmd.Flags().GetString("service"))
		fmt.Println(cmd.Flags().GetString("composition"))
		fmt.Println(args)
	},
}

// Execute ...
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
