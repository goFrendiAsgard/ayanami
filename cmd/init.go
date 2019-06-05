package cmd

import (
	"github.com/spf13/cobra"
	"github.com/state-alchemists/ayanami/generator"
	"log"
	"os"
)

var dirName, projectName, repoName string

func init() {
	cwd, _ := os.Getwd()
	initCmd.Flags().StringVarP(&dirName, "dir", "d", cwd, "project's parent directory")
	initCmd.Flags().StringVarP(&projectName, "project", "p", "", "name of the project")
	initCmd.Flags().StringVarP(&repoName, "repo", "r", "", "project's package repository")
	initCmd.MarkFlagRequired("name")
	initCmd.MarkFlagRequired("repo")
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a project",
	Long:  `Create a project in current working directory`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("[INFO] Project's parent directory : %s", dirName)
		log.Printf("[INFO] Project name               : %s", projectName)
		log.Printf("[INFO] Project repository         : %s", repoName)
		generator.GenerateProject(dirName, projectName, repoName)
	},
}
