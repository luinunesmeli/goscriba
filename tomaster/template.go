package tomaster

import (
	"time"
)

type templateData struct {
	Version      string
	Now          string
	Author       string
	Features     PRs
	Fix          PRs
	Enhancements PRs
}

func newTemplateData(version, author string, prs PRs) templateData {
	data := time.Now()
	return templateData{
		Version:      version,
		Author:       author,
		Features:     prs.Filter(Feature),
		Enhancements: prs.Filter(Enhancement),
		Fix:          prs.Filter(Fix),
		Now:          data.Format("2006-01-02"),
	}
}

const changelogTemplate = `## Version {{ .Version }}
**Created at {{ .Now }} by {{ .Author }}**
{{- println ""}}
{{- if .Features }}
### Features
	{{- range $pr := .Features }}
* [{{ $pr.Title }}]({{ $pr.PRLink }}) by {{ $pr.Author }}
	{{- end }}
{{- end }}
{{- if .Enhancements }}
### Enhancements
	{{- range $pr := .Enhancements }}
* [{{ $pr.Title }}]({{ $pr.PRLink }}) by {{ $pr.Author }}
	{{- end }}
{{- end }}
{{- if .Fix }}
### Fixes
	{{- range $pr := .Fix }}
* [{{ $pr.Title }}]({{ $pr.PRLink }}) by {{ $pr.Author }}
	{{- end }}
{{- end }}
{{- println ""}}`
