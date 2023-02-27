package scriba

import (
	"fmt"
	"os"
)

const githubAccessToken = "GH_ACCESS_TOKEN"
const errMsg = "`%s` enviroment variable not found! Please refer to README for help."

type Config struct {
	GithubTokenAPI string
}

func LoadConfig() (Config, error) {
	token := os.Getenv(githubAccessToken)
	if token == "" {
		return Config{},
			fmt.Errorf(errMsg, githubAccessToken)
	}
	return Config{
		GithubTokenAPI: token,
	}, nil
}
