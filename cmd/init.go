package cmd

import (
	"github.com/spf13/cobra"
	"github.com/state-alchemists/ayanami/projectgenerator"
	"log"
	"os"
	"path/filepath"
)

var dirPath, templatePath, genPath, projectName, repoName string

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	execFile, err := os.Executable()
	if err != nil {
		panic(err)
	}
	// define default values
	execDir := filepath.Dir(execFile)
	templatePath = filepath.Join(execDir, "templates")
	genPath = filepath.Join(execDir, "gen")
	// define flags
	initCmd.Flags().StringVarP(&genPath, "gen", "g", genPath, "Gen directory")
	initCmd.Flags().StringVarP(&templatePath, "template", "t", templatePath, "project generator's template directory")
	initCmd.Flags().StringVarP(&dirPath, "dir", "d", cwd, "project's parent directory")
	initCmd.Flags().StringVarP(&projectName, "project", "p", "", "name of the project, e.g: myProject")
	initCmd.Flags().StringVarP(&repoName, "repo", "r", "", "project's package repository, e.g: github.com/myUser/myProject")
	// mark name and repo as required
	initCmd.MarkFlagRequired("name")
	initCmd.MarkFlagRequired("repo")
	// register command
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a project",
	Long:  `Create a project in current working directory`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("[INFO] Gen directory                  : %s", genPath)
		log.Printf("[INFO] Generator's template directory : %s", templatePath)
		log.Printf("[INFO] Project's parent directory     : %s", dirPath)
		log.Printf("[INFO] Project name                   : %s", projectName)
		log.Printf("[INFO] Project repository             : %s", repoName)
		generator, err := projectgenerator.NewProjectGenerator(dirPath, projectName, repoName, templatePath, genPath)
		if err != nil {
			log.Printf("[ERROR] Cannot init generator : %s", err)
			return
		}
		err = generator.Generate()
		if err != nil {
			log.Printf("[ERROR] %s", err)
			return
		}
		log.Printf("[INFO] Done")
	},
}
