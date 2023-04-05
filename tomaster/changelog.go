package tomaster

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"text/template"

	"github.com/go-git/go-git/v5"

	"github.com/luinunesmeli/goscriba/pkg/config"
)

type (
	Changelog struct {
		tree *git.Worktree
		cfg  config.Config
	}
)

const (
	tempPath = "temp-file.txt"
)

func NewChangelog(cfg config.Config, tree *git.Worktree) Changelog {
	return Changelog{
		tree: tree,
		cfg:  cfg,
	}
}

func (c *Changelog) LoadChangelog() Task {
	return Task{
		Desc: "Verify actual changelog",
		Help: fmt.Sprintf("Not found! Changelog should exist at %s.", c.cfg.Changelog),
		Func: func(session Session) (error, string, Session) {
			file, err := os.Open(c.cfg.Changelog)
			defer file.Close()
			return err, "", session
		},
	}
}

func (c *Changelog) WriteChangelog() Task {
	return Task{
		Desc: "Update changelog",
		Help: fmt.Sprintf("Not found! Changelog should exist at %s.", c.cfg.Changelog),
		Func: func(session Session) (error, string, Session) {
			t, err := template.New("changelog").Parse(changelogTemplate)
			if err != nil {
				return err, "", Session{}
			}

			buf := bytes.NewBufferString("")
			err = t.Execute(buf, newTemplateData(session, c.cfg.Repo.Author, session.PRs))

			session.Changelog = buf.String()

			return writeChangelogContent(c.cfg.Changelog, buf.String(), c.tree), "", session
		},
	}
}

func writeChangelogContent(path string, content string, tree *git.Worktree) error {
	temp, err := tree.Filesystem.Create(tempPath)
	if err != nil {
		return err
	}
	defer temp.Close()

	file, err := tree.Filesystem.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = temp.Write([]byte(fmt.Sprintf("# Changelog\n\n%s\n", content))); err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if match, _ := regexp.Match("#.Changelog", []byte(text)); match {
			continue
		}
		if _, err = temp.Write([]byte(text + "\n")); err != nil {
			return err
		}
	}

	if err = scanner.Err(); err != nil {
		return err
	}

	return tree.Filesystem.Rename(tempPath, path)
}
