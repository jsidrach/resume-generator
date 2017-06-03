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
{{- with .Webpage.Url}}
- Webpage: {{.}}
{{- end -}}
{{- with .Github.Url}}
- GitHub: {{.}}
{{- end -}}
{{- with .Linkedin.Url}}
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
{{- with .Url}} - {{.}}{{end}}
{{- with .Description}}

{{.}}{{end}}
{{- range .Details}}
  - {{.}}{{end}}
{{end}}
{{end -}}
