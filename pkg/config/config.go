package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// https://github.com/settings/tokens
const classicToken = "GH_PERSONAL_ACCESS_TOKEN_CLASSIC"

// https://github.com/settings/tokens?type=beta
const finegrainedToken = "GH_PERSONAL_ACCESS_TOKEN_FINEGRAINED"
const errMsg = "`%s` or `%s` enviroment variable not found! Please refer to README for help."

type Config struct {
	ClassicToken  string
	GranularToken string
	Path          string
	Base          string
	Changelog     string
	Version       bool
	AutoPR        bool
	Autoinstall   bool
}

func LoadConfig() (Config, error) {
	classic, granular, err := getGHTokenEnv()
	if err != nil {
		return Config{}, err
	}

	path, baseBranch, changelog, pr, auto, version := loadCliParams()
	return Config{
		ClassicToken:  classic,
		GranularToken: granular,
		Path:          path,
		Base:          baseBranch,
		Changelog:     changelog,
		AutoPR:        pr,
		Autoinstall:   auto,
		Version:       version,
	}, nil
}

func (c Config) AuthStrategy() transport.AuthMethod {
	if c.GranularToken != "" {
		return &http.TokenAuth{Token: c.GranularToken}
	}
	return &http.BasicAuth{
		Username: "token_user", // yes, this can be anything except an empty string
		Password: c.ClassicToken,
	}
}

func (c Config) GetPersonalAccessToken() string {
	if c.GranularToken != "" {
		return c.GranularToken
	}
	return c.ClassicToken
}

func loadCliParams() (path, base, changelog string, pr, auto, version bool) {
	flag.BoolVar(&pr, "autopr", false, "automatically generate Pull Request (optional)")
	flag.BoolVar(&auto, "install", false, "automatically install ToMaster on environment")
	flag.BoolVar(&version, "version", false, "show actual version")
	flag.StringVar(&path, "path", "./", "project path you want to generate a release")
	flag.StringVar(&base, "base", "master", "provide the base: master or main")
	flag.StringVar(&changelog, "changelog", path+"docs/guide/pages/changelog.md", "provide the changelog filename")

	flag.Parse()

	return path, base, changelog, pr, auto, version
}

func getGHTokenEnv() (string, string, error) {
	classic := os.Getenv(classicToken)
	if classic != "" {
		return classic, "", nil
	}

	finegrained := os.Getenv(finegrainedToken)
	if finegrained != "" {
		return "", finegrained, nil
	}

	return "", "", fmt.Errorf(errMsg, classicToken, finegrainedToken)
}
