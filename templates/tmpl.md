# {{.Name}}
{{- with .Title}} | {{.}}{{end}}

{{- with .Contact}}

## Contact
{{- with .Email}}
- Email: {{.}}
{{- end -}}
{{- with .Phone}}
- Phone: {{.}}
{{- end -}}
{{- with .Address}}
- Address: {{.}}
{{- end -}}
{{- with .Webpage.URL}}
- Webpage: {{.}}
{{- end -}}
{{- with .Github.URL}}
- GitHub: {{.}}
{{- end -}}
{{- with .Linkedin.URL}}
- LinkedIn: {{.}}
{{- end}}
{{- end}}

{{with .Summary -}}
{{.}}
{{end}}
{{range .Sections -}}
## {{.Name}}
{{- range .Entries}}
{{- with .Where}}
### {{.}}{{end}}{{with .Location}} | {{.}}{{end}}
{{- with .What}}
**{{.}}**{{end}}
{{- with .When}} | *{{.}}*{{end}}
{{- with .URL}} - {{.}}{{end}}
{{- with .Description}}

{{.}}{{end}}
{{- range .Details}}
  - {{.}}{{end}}
{{end}}
{{end -}}
