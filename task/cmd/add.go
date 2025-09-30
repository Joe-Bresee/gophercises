/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Joe-Bresee/gophercises/task/db"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add a task to your task manager table",
	Run: func(cmd *cobra.Command, args []string) {
		task := strings.Join(args, " ")
		_, err := db.CreateTask(task)
		if err != nil {
			fmt.Println("Something went wrong:", err.Error())
			os.Exit(1)
		}
		fmt.Printf("Added \"%s\" to your task list.\n", task)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
