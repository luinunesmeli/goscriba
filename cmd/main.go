package main

import (
	"fmt"
	"os"

	"github.com/luinunesmeli/goscriba/cmd/app"
	"github.com/luinunesmeli/goscriba/cmd/install"
	"github.com/luinunesmeli/goscriba/cmd/version"
	"github.com/luinunesmeli/goscriba/pkg/config"
)

const actualVersion = "1.0.0"

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		handleErr(err)
	}

	switch {
	case cfg.Install:
		err = install.Run()
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
	os.Exit(1)
}
