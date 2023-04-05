package tomaster

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"text/template"

	"github.com/go-git/go-git/v5"
)

type (
	Changelog struct {
		filename string
		content  []string
	}
)

const (
	tempPath = "temp-file.txt"
)

func NewChangelog(filename string) Changelog {
	return Changelog{
		filename: filename,
	}
}

func (c *Changelog) LoadChangelog() Task {
	return Task{
		Desc: "Load actual changelog",
		Help: fmt.Sprintf("Not found! Changelog should exist at %s.", c.filename),
		Func: func(session Session) (error, string, Session) {
			file, err := os.Open(c.filename)
			if err != nil {
				return err, "", session
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				c.content = append(c.content, scanner.Text())
			}
			return scanner.Err(), "", session
		},
	}
}

func (c *Changelog) UpdateChangelog(session Session, author string, tree *git.Worktree) (error, string) {
	t, err := template.New("changelog").Parse(changelogTemplate)
	if err != nil {
		return err, ""
	}

	buf := bytes.NewBufferString("")
	err = t.Execute(buf, newTemplateData(session.ChosenVersion, author, session.PRs))

	return writeChangelogContent(c.filename, buf.String(), tree), buf.String()
}

func writeChangelogContent(path string, content string, tree *git.Worktree) error {
	// temporary file
	temp, err := tree.Filesystem.Create(tempPath)
	if err != nil {
		return err
	}
	defer temp.Close()

	// existing changelog file
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = temp.Write([]byte(content)); err != nil {
		return err
	}
	if _, err = temp.Write([]byte("\n")); err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if match, _ := regexp.Match("#.Changelog", []byte(text)); match {
			continue
		}
		if _, err = temp.Write([]byte(text)); err != nil {
			return err
		}
		if _, err = temp.Write([]byte("\n")); err != nil {
			return err
		}
	}

	if err = scanner.Err(); err != nil {
		return err
	}

	return tree.Filesystem.Rename(tempPath, path)
}
