package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-git/go-billy/v5/osfs"

	"github.com/luinunesmeli/goscriba/cmd/app"
	"github.com/luinunesmeli/goscriba/cmd/install"
	"github.com/luinunesmeli/goscriba/cmd/version"
	"github.com/luinunesmeli/goscriba/pkg/config"
)

const actualVersion = "1.0.0"

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		handleErr(err)
	}

	cfg, err := config.LoadConfig(homeDir)
	if err != nil {
		handleErr(err)
	}

	logFile := initLog(cfg.LogPath)
	defer logFile.Close()

	switch {
	case cfg.Install:
		path, err := os.Executable()
		if err != nil {
			handleErr(err)
		}
		err = install.Run(cfg.HomeDir, path, osfs.New("/"))
	case cfg.Version:
		err = version.Run(actualVersion)
	default:
		err = app.Run(cfg)
	}

	if err != nil {
		handleErr(err)
	}

	os.Exit(0)
}

func handleErr(err error) {
	fmt.Println(err)
	log.Println(err)
	os.Exit(1)
}

func initLog(filename string) *os.File {
	logFile, err := os.OpenFile(filename, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	return logFile
}
