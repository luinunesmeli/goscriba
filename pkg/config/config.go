package config

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
)

// https://github.com/settings/tokens
const classicToken = "GH_PERSONAL_ACCESS_TOKEN_CLASSIC"

// https://github.com/settings/tokens?type=beta
const finegrainedToken = "GH_PERSONAL_ACCESS_TOKEN_FINEGRAINED"
const errMsg = "`%s` or `%s` enviroment variable not found! Please refer to README for help"

const logPath = "%s/.tomaster/debug.log"

type (
	Config struct {
		ClassicToken     string
		FinegrainedToken string
		Path             string
		Base             string
		Changelog        string
		HomeDir          string
		LogPath          string
		Version          bool
		Install          bool
		Repo             Repo
	}

	Repo struct {
		URL    string
		Author string
		Owner  string
		Name   string
	}
)

func LoadConfig(homeDir string) (Config, error) {
	path, baseBranch, changelog, install, version := loadCliParams()

	cfg := Config{
		Path:      path,
		Base:      baseBranch,
		Changelog: changelog,
		Install:   install,
		Version:   version,
		HomeDir:   homeDir,
		LogPath:   fmt.Sprintf(logPath, homeDir),
	}

	cfg, err := getRepoConfig(cfg)
	if err != nil {
		return Config{}, err
	}

	cfg, err = getGHTokenEnv(cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
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

	flag.BoolVar(&install, "install", false, "automatically install ToMaster on environment")
	flag.BoolVar(&version, "version", false, "show actual version")
	flag.StringVar(&path, "path", basePath, "project path you want to generate a release")
	flag.StringVar(&base, "base", "master", "provide the base: master or main")
	flag.StringVar(&changelog, "changelog", "docs/guide/pages/changelog.md", "provide the changelog filename")
	flag.Parse()

	changelogPath := basePath + changelog

	return path, base, changelogPath, install, version
}

func getGHTokenEnv(cfg Config) (Config, error) {
	cfg.ClassicToken = os.Getenv(classicToken)
	if cfg.ClassicToken != "" {
		return cfg, nil
	}

	cfg.FinegrainedToken = os.Getenv(finegrainedToken)
	if cfg.FinegrainedToken != "" {
		return cfg, nil
	}

	return Config{}, fmt.Errorf(errMsg, classicToken, finegrainedToken)
}

func getRepoConfig(cfg Config) (Config, error) {
	plainRepo, err := git.PlainOpen(cfg.Path)
	if err != nil {
		return Config{}, err
	}

	repoCfg, err := plainRepo.Config()
	if err != nil {
		return Config{}, err
	}

	remotes := repoCfg.Remotes["origin"]
	cfg.Repo.URL = remotes.URLs[0]

	repoConfig, err := plainRepo.ConfigScoped(gitconfig.SystemScope)
	if err != nil {
		return Config{}, err
	}
	cfg.Repo.Author = repoConfig.User.Name

	parts := strings.Split(remotes.URLs[0], "/")
	cfg.Repo.Owner = parts[len(parts)-2]
	cfg.Repo.Name = strings.TrimSuffix(parts[len(parts)-1], ".git")

	return cfg, nil
}
