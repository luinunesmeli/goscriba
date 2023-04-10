package tomaster

import (
	"time"
)

type templateData struct {
	Version      string
	Now          string
	Author       string
	PRNumber     int
	PRURL        string
	Features     PRs
	Fix          PRs
	Enhancements PRs
}

func newTemplateData(session Session, author string, prs PRs) templateData {
	data := time.Now()
	return templateData{
		Version:      session.ChosenVersion,
		Now:          data.Format("2006-01-02"),
		Author:       author,
		Features:     prs.Filter(Feature),
		Fix:          prs.Filter(Fix),
		Enhancements: prs.Filter(Enhancement),
		PRNumber:     session.PRNumber,
		PRURL:        session.PRUrl,
	}
}

const changelogTemplate = `## Version {{ .Version }}
**Created at {{ .Now }} by @{{ .Author }}**
{{- println ""}}
{{- if .Features }}
### Features
	{{- range $pr := .Features }}
* [{{ $pr.Title }}]({{ $pr.PRLink }}) by @{{ $pr.Author }}
	{{- end }}
{{- end }}
{{- if .Enhancements }}
### Enhancements
	{{- range $pr := .Enhancements }}
* [{{ $pr.Title }}]({{ $pr.PRLink }}) by @{{ $pr.Author }}
	{{- end }}
{{- end }}
{{- if .Fix }}
### Fixes
	{{- range $pr := .Fix }}
* [{{ $pr.Title }}]({{ $pr.PRLink }}) by @{{ $pr.Author }}
	{{- end }}
{{- end }}
{{- println ""}}`
