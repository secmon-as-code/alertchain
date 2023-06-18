{{$countMarkdown := 0}}
- ID: {{ .ID }}
- Created At: {{ .CreatedAt }}
- Schema: {{ .Schema }}
- Detected by: {{ .Source }}

## Description
{{.Description}}

## Attributes

| Name | Value | Type |
|------|-------|------|
{{range .Attrs}}{{ if ne .Type "markdown" }} | {{ .Name }} | `{{ .Value }}` | {{ .Type }} |
{{else}}{{ $countMarkdown = add $countMarkdown 1 }}{{end}}{{end}}

{{ if gt $countMarkdown 0 }}
## Comments

{{range .Attrs}}{{ if eq .Type "markdown" }}
### {{ .Name }}

{{ .Value }}

{{end}}{{end}}
{{end}}

## Alert

<details>
<summary>Raw data</summary>

```json

{{ .Raw }}

```

</details>
