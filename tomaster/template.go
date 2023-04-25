package tomaster

import (
	"time"
)

type templateData struct {
	Version      string
	Now          string
	Author       Author
	Features     PRs
	Fixes        PRs
	Enhancements PRs
	Bugfixes     PRs
}

func newTemplateData(session Session, author Author, prs PRs) templateData {
	data := time.Now()
	return templateData{
		Version:      session.ChosenVersion,
		Now:          data.Format("2006-01-02"),
		Author:       author,
		Features:     prs.Filter(Feature),
		Fixes:        prs.Filter(Fix),
		Enhancements: prs.Filter(Enhancement),
		Bugfixes:     prs.Filter(Bugfix),
	}
}

const ChangelogTemplate = `## Version {{ .Version }}
{{- println ""}}
{{- if .Features }}
### Features
	{{- range $pr := .Features }}
* [{{ $pr.Title }}]({{ $pr.Link }}) by {{ $pr.Author.Name }}
	{{- end }}
{{- end }}
{{- if .Enhancements }}
### Enhancements
	{{- range $pr := .Enhancements }}
* [{{ $pr.Title }}]({{ $pr.Link }}) by {{ $pr.Author.Name }}
	{{- end }}
{{- end }}
{{- if .Fixes }}
### Fixes
	{{- range $pr := .Fixes }}
* [{{ $pr.Title }}]({{ $pr.Link }}) by {{ $pr.Author.Name }}
	{{- end }}
{{- end }}
{{- if .Bugfixes }}
### Fixes
	{{- range $pr := .Bugfixes }}
* [{{ $pr.Title }}]({{ $pr.Link }}) by {{ $pr.Author.Name }}
	{{- end }}
{{- end }}
{{- println ""}}`
