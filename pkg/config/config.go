package config

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
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
		HomeDir          string
		LogPath          string
		GenerateTemplate bool
		Version          bool
		Install          bool
		Uninstall        bool
		Repo             Repo
		Template         Template
	}

	Repo struct {
		URL    string
		Author string
		Owner  string
		Name   string
	}

	Template struct {
		Path           string `yaml:"path"`
		CustomTemplate string
	}
)

func LoadConfig(homeDir string) (Config, error) {
	path, baseBranch, changelog, install, uninstall, version, generate := loadCliParams()

	cfg := Config{
		Path:             path,
		Base:             baseBranch,
		Install:          install,
		Uninstall:        uninstall,
		Version:          version,
		HomeDir:          homeDir,
		LogPath:          fmt.Sprintf(logPath, homeDir),
		GenerateTemplate: generate,
		Template: Template{
			Path: changelog,
		},
	}

	if !install && !uninstall && !version {
		var err error
		cfg, err = getRepoConfig(cfg)
		if err != nil {
			return Config{}, err
		}
	}

	cfg, err := getGHTokenEnv(cfg)
	if err != nil {
		return Config{}, err
	}

	cfg, err = getChangelogTemplate(cfg)
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

func loadCliParams() (path, base, changelog string, install, uninstall, version, generate bool) {
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

	return path, base, changelog, install, uninstall, version, generate
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

	if strings.HasPrefix(cfg.Repo.URL, "git@") {
		parts := strings.Split(cfg.Repo.URL, "/")
		ownerParts := strings.Split(parts[0], ":")

		cfg.Repo.Owner = ownerParts[len(ownerParts)-1]
		cfg.Repo.Name = strings.TrimSuffix(parts[1], ".git")
		cfg.Repo.URL = fmt.Sprintf("https://github.com/%s/%s.git", cfg.Repo.Owner, cfg.Repo.Name)
	} else {
		parts := strings.Split(cfg.Repo.URL, "/")
		cfg.Repo.Owner = parts[len(parts)-2]
		cfg.Repo.Name = strings.TrimSuffix(parts[len(parts)-1], ".git")
	}

	return cfg, nil
}

func getChangelogTemplate(cfg Config) (Config, error) {
	fs := osfs.New("./")

	f, err := fs.Open(".tomaster")
	if err != nil {
		if strings.HasSuffix(err.Error(), "no such file or directory") {
			return cfg, nil
		}
		return Config{}, err
	}
	defer f.Close()

	rest, err := frontmatter.Parse(f, &cfg.Template)
	if len(rest) > 0 {
		cfg.Template.CustomTemplate = string(rest)
	}
	return cfg, nil
}
