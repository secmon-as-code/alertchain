{{$countMarkdown := 0}}
- ID: {{ .ID }}
- Created At: {{ .CreatedAt }}
- Schema: {{ .Schema }}
- Detected by: {{ .Source }}

## Description
{{.Description}}

## Parameters

| Key | Value | Type |
|-----|-------|------|
{{range .Params}}{{ if ne .Type "markdown" }} | {{ .Key }} | `{{ .Value }}` | {{ .Type }} |
{{else}}{{ $countMarkdown = add $countMarkdown 1 }}{{end}}{{end}}

{{ if gt $countMarkdown 0 }}
## Comments

{{range .Params}}{{ if eq .Type "markdown" }}
### {{ .Key }}

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
