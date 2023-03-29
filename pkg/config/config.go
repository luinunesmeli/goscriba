package config

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

// https://github.com/settings/tokens
const classicToken = "GH_PERSONAL_ACCESS_TOKEN_CLASSIC"

// https://github.com/settings/tokens?type=beta
const finegrainedToken = "GH_PERSONAL_ACCESS_TOKEN_FINEGRAINED"
const errMsg = "`%s` or `%s` enviroment variable not found! Please refer to README for help."

type Config struct {
	ClassicToken     string
	FinegrainedToken string
	Path             string
	Base             string
	Changelog        string
	Version          bool
	Install          bool
}

func LoadConfig() (Config, error) {
	classic, finegrained, err := getGHTokenEnv()
	if err != nil {
		return Config{}, err
	}

	path, baseBranch, changelog, install, version := loadCliParams()

	return Config{
		ClassicToken:     classic,
		FinegrainedToken: finegrained,
		Path:             path,
		Base:             baseBranch,
		Changelog:        changelog,
		Install:          install,
		Version:          version,
	}, nil
}

func (c Config) GetPersonalAccessToken() string {
	if c.FinegrainedToken != "" {
		return c.FinegrainedToken
	}
	return c.ClassicToken
}

func (c Config) ReadGitignore() ([]string, error) {
	file, err := os.Open(c.Path + "/.gitignore")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	var filtered []string
	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		filtered = append(filtered, line)
	}

	return filtered, scanner.Err()
}

func loadCliParams() (path, base, changelog string, install, version bool) {
	dir, _ := os.Getwd()
	basePath := dir + "/"
	changelogPath := basePath + "docs/guide/pages/changelog.md"

	flag.BoolVar(&install, "install", false, "automatically install ToMaster on environment")
	flag.BoolVar(&version, "version", false, "show actual version")
	flag.StringVar(&path, "path", basePath, "project path you want to generate a release")
	flag.StringVar(&base, "base", "master", "provide the base: master or main")
	flag.StringVar(&changelog, "changelog", changelogPath, "provide the changelog filename")

	flag.Parse()

	return path, base, changelog, install, version
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
