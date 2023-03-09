package scriba

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"text/template"
)

type (
	Changelog struct {
		filename  string
		content   []string
		PRs       PRs
		Generated string
		ChosenTag string
	}
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
		Func: func() (error, string) {
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
		Func: func() (error, string) {
			t, err := template.New("changelog").Parse(changelogTemplate)
			if err != nil {
				return err, ""
			}

			s := ""
			buf := bytes.NewBufferString(s)
			err = t.Execute(buf, newTemplateData(c.ChosenTag, c.PRs))

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
