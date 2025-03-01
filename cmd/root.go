/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

const (
	defaultProfile       = "default"
	defaultProfileNumber = 0
	profileRegexp        = `(\[profile .*\])`
)

var (
	profiles = []string{defaultProfile}
	re       = regexp.MustCompile(profileRegexp)
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "awsp",
	Run: func(cmd *cobra.Command, args []string) {
		homeDir := os.Getenv("HOME")
		currentProfile := os.Getenv("AWS_PROFILE")
		chooseProfile := defaultProfileNumber

		file, err := os.Open(fmt.Sprintf("%s/.aws/config", homeDir))
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		defer file.Close()

		buf, err := io.ReadAll(file)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		matches := re.FindAllString(string(buf), -1)
		for _, match := range matches {
			profile := match[9 : len(match)-1]
			profiles = append(profiles, profile)
			if profile == currentProfile {
				chooseProfile = len(profiles) - 1
			}
		}
		fmt.Println("choosed profile:", profiles[chooseProfile])

		prompt := promptui.Select{
			Label: "Select profile",
			Items: profiles,
		}
		_, result, err := prompt.Run()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		if result == defaultProfile {
			result = ""
		}
		shellCommand := fmt.Sprintf("./cmd/run.sh %s", result)
		if err := exec.Command(shellCommand).Run(); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.awsp.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
