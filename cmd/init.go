package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/state-alchemists/ayanami/projectgenerator"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
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
	initCmd.Flags().StringVarP(&projectName, "project", "p", "", "[REQUIRED] name of the project, e.g: myProject")
	initCmd.Flags().StringVarP(&repoName, "repo", "r", "", "[REQUIRED] project's package repository, e.g: github.com/myUser/myProject")
	// register command
	rootCmd.AddCommand(initCmd)
}

func getRandomUser() string {
	rand.Seed(time.Now().UTC().UnixNano())
	users := []string{"rei", "shinji", "asuka", "toji", "kaworu", "mari"}
	user := users[rand.Intn(len(users))]
	return user
}

func getRandomProject() string {
	rand.Seed(time.Now().UTC().UnixNano())
	colors := []string{"red", "green", "blue", "yellow", "black", "white"}
	angels := []string{"Adam", "Lilith", "Ramiel", "Leliel", "Bardiel", "Zeruel", "Tabris"}
	color := colors[rand.Intn(len(colors))]
	angel := angels[rand.Intn(len(angels))]
	number := 1000 + rand.Intn(8999)
	return fmt.Sprintf("%s%s%d", color, angel, number)
}

func askProjectIfNotExists() {
	if projectName == "" {
		defaultProjectName := getRandomProject()
		fmt.Printf("Enter your project name (default: %s): ", defaultProjectName)
		_, err := fmt.Scanln(&projectName)
		if err != nil {
			fmt.Printf("Failed to read input: %s", err)
		}
		if projectName == "" {
			projectName = defaultProjectName
		}
	}
}

func askRepoIfNotExists() {
	if repoName == "" {
		defaultUserName := getRandomUser()
		defaultRepoName := fmt.Sprintf("github.com/%s/%s", defaultUserName, projectName)
		fmt.Printf("Enter your repo name (default: %s): ", defaultRepoName)
		_, err := fmt.Scanln(&repoName)
		if err != nil {
			fmt.Printf("Failed to read input: %s", err)
		}
		if repoName == "" {
			repoName = defaultRepoName
		}
	}
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a project",
	Long:  `Create a project in current working directory`,
	Run: func(cmd *cobra.Command, args []string) {
		askProjectIfNotExists()
		askRepoIfNotExists()
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
