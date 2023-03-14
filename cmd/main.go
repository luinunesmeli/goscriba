package main

import (
	"fmt"
	"os"

	"github.com/luinunesmeli/goscriba/cmd/app"
	"github.com/luinunesmeli/goscriba/cmd/install"
	"github.com/luinunesmeli/goscriba/cmd/version"
	"github.com/luinunesmeli/goscriba/scriba"
)

const actualVersion = "0.1.0"

func main() {
	cfg, err := scriba.LoadConfig()
	if err != nil {
		handleErr(err)
	}

	switch {
	case cfg.Autoinstall:
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
