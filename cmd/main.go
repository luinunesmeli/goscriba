package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/go-git/go-billy/v5/osfs"

	"github.com/luinunesmeli/goscriba/cmd/app"
	"github.com/luinunesmeli/goscriba/cmd/generatetemplate"
	"github.com/luinunesmeli/goscriba/cmd/install"
	"github.com/luinunesmeli/goscriba/cmd/uninstall"
	"github.com/luinunesmeli/goscriba/cmd/version"
	"github.com/luinunesmeli/goscriba/pkg/config"
)

const (
	actualVersion = "1.0.0"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		handleErr(err)
	}

	installOpt, uninstallOpt, versionOpt, generateOpt, path, base, changelog := parseCommandOptions()

	cfg, err := config.LoadConfig(homeDir, path, base, changelog)
	if err != nil {
		handleErr(err)
	}

	switch {
	case installOpt:
		path, err := os.Executable()
		if err != nil {
			handleErr(err)
		}
		err = install.Run(cfg.HomeDir, path, osfs.New("/"))
	case uninstallOpt:
		err = uninstall.Run(cfg.HomeDir, osfs.New("/"))
	case versionOpt:
		err = version.Run(actualVersion)
	case generateOpt:
		err = generatetemplate.Run(osfs.New("./"))
	default:
		logFile := initLog(cfg.LogPath)
		defer logFile.Close()

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

func parseCommandOptions() (install, uninstall, version, generate bool, path, base, changelog string) {
	dir, _ := os.Getwd()
	basePath := dir + "/"

	flag.BoolVar(&install, "install", false, "automatically install ToMaster on environment")
	flag.BoolVar(&uninstall, "uninstall", false, "uninstall ToMaster")
	flag.BoolVar(&version, "version", false, "show actual version")
	flag.BoolVar(&generate, "generate", false, "generate config template")
	flag.StringVar(&path, "path", basePath, "project path you want to generate a release")
	flag.StringVar(&base, "base", "master", "provide the base: master or main")
	flag.StringVar(&changelog, "changelog", "docs/guide/pages/changelog.md", "provide the changelog filename")

	flag.Parse()

	return install, uninstall, version, generate, path, base, changelog
}
