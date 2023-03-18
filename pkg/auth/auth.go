package auth

import (
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"

	"github.com/luinunesmeli/goscriba/pkg/config"
)

func AuthMethod(cfg config.Config) transport.AuthMethod {
	if cfg.FinegrainedToken != "" {
		return &http.TokenAuth{Token: cfg.FinegrainedToken}
	}
	return &http.BasicAuth{
		Username: "token_user", // yes, this can be anything except an empty string
		Password: cfg.ClassicToken,
	}
}
