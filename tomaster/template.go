package tomaster

type templateData struct {
	Version      string
	Features     PRs
	Fix          PRs
	Enhancements PRs
}

func newTemplateData(version string, prs PRs) templateData {
	return templateData{
		Version:      version,
		Features:     prs.Filter(Feature),
		Enhancements: prs.Filter(Enhancement),
		Fix:          prs.Filter(Fix),
	}
}

const changelogTemplate = `
## Version {{ .Version }}
{{ if gt (len .Features) 0 }}
### Features
	{{ range $pr := .Features }}
* [{{ $pr.Title }}]({{ $pr.PRLink }}) by {{ $pr.Author }}
	{{ end }}
{{ end }}
{{ if gt (len .Enhancements) 0 }}
### Enhancements
	{{ range $pr := .Enhancements }}
* [{{ $pr.Title }}]({{ $pr.PRLink }}) by {{ $pr.Author }}
	{{ end }}
{{ end }}
{{ if gt (len .Enhancements) 0 }}
### Enhancements
	{{ range $pr := .Enhancements }}
* [{{ $pr.Title }}]({{ $pr.PRLink }}) by {{ $pr.Author }}
	{{ end }}
{{ end }}
`
