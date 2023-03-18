package tomaster

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"text/template"
)

type (
	Changelog struct {
		filename  string
		content   []string
		Generated string
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
		Help: fmt.Sprintf("Changelog should exist at %s", c.filename),
		Func: func(session Session) (error, string) {
			file, err := os.Open(c.filename)
			if err != nil {
				return err, ""
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				c.content = append(c.content, scanner.Text())
			}
			return scanner.Err(), ""
		},
	}
}

func (c *Changelog) Update() Task {
	return Task{
		Desc: "Load actual changelog",
		Help: "Changelog should exist at ",
		Func: func(session Session) (error, string) {
			t, err := template.New("changelog").Parse(changelogTemplate)
			if err != nil {
				return err, ""
			}

			s := ""
			buf := bytes.NewBufferString(s)
			err = t.Execute(buf, newTemplateData(session.ChosenVersion, session.PRs))

			c.Generated = buf.String()

			//content := c.content[1:]
			//
			//file, err := os.Create(c.filename)
			//if err != nil {
			//	return err
			//}
			//defer file.Close()
			//
			//w := bufio.NewWriter(file)
			//for _, line := range lines {
			//	fmt.Fprintln(w, line)
			//}
			//return w.Flush()

			return err, ""
		},
	}
}

func writeChangelogContent(path string, content string) error {
	// temporary file
	temp, err := os.Create(tempPath)
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

	if _, err = temp.WriteString(content); err != nil {
		return err
	}
	if _, err = temp.WriteString("\n"); err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if match, _ := regexp.Match("#.Changelog", []byte(text)); match {
			continue
		}
		temp.WriteString(text)
		temp.WriteString("\n")
	}

	if err = scanner.Err(); err != nil {
		return err
	}

	temp.Sync()
	if err := os.Rename(tempPath, path); err != nil {
		return err
	}

	return nil
}
