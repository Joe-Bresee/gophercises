/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"path/filepath"

	"github.com/Joe-Bresee/gophercises/task/cmd"
	"github.com/Joe-Bresee/gophercises/task/db"
	homedir "github.com/mitchellh/go-homedir"
)

func main() {
	home, _ := homedir.Dir()
	dbPath := filepath.Join(home, "tasks.db")
	db.Init(dbPath)
	cmd.Execute()
}
