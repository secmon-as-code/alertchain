- ID: {{ .ID }}
- Created At: {{ .CreatedAt }}
- Schema: {{ .Schema }}
- Detected by: {{ .Source }}

## Description
{{.Description}}

## Parameters

| Key | Value | Type |
|-----|-------|------|
{{range .Params}}| {{ .Key }} | `{{ .Value }}` | {{ .Type }} |
{{end}}

## Alert

<details>
<summary>Raw data</summary>

```json

{{ .Raw }}

```

</details>
