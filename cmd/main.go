package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/luinunesmeli/goscriba/scriba"
	"github.com/luinunesmeli/goscriba/view"
)

func main() {
	_, err := scriba.LoadConfig()
	if err != nil {
		handleErr(err)
	}

	gitRepo, err := scriba.NewGitRepo()
	if err != nil {
		handleErr(err)
	}

	p := tea.NewProgram(view.NewView(gitRepo))
	if _, err = p.Run(); err != nil {
		handleErr(err)
	}
}

func handleErr(err error) {
	fmt.Println(err)
	os.Exit(1)
}
