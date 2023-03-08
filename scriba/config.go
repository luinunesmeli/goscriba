package scriba

import (
	"flag"
	"fmt"
	"os"
)

const githubAccessToken = "GH_ACCESS_TOKEN"
const errMsg = "`%s` enviroment variable not found! Please refer to README for help."

type Config struct {
	GithubTokenAPI string
	Path           string
	Base           string
	Changelog      string
}

func LoadConfig() (Config, error) {
	token := os.Getenv(githubAccessToken)
	if token == "" {
		return Config{},
			fmt.Errorf(errMsg, githubAccessToken)
	}

	path, baseBranch, changelog := loadCliParams()
	return Config{
		GithubTokenAPI: token,
		Path:           path,
		Base:           baseBranch,
		Changelog:      changelog,
	}, nil
}

func loadCliParams() (path, base, changelog string) {
	flag.StringVar(&path, "path", "./", "project path you want to generate a release")
	flag.StringVar(&base, "base", "master", "provide the base: master or main")
	flag.StringVar(&changelog, "changelog", path+"docs/guide/pages/changelog.md", "provide the changelog filename")

	flag.Parse()

	return path, base, changelog
}
