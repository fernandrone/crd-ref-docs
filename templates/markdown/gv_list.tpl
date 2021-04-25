{{- define "gvList" -}}
{{- $groupVersions := . -}}

<!-- Generated documentation. Please do not edit. -->
# API Reference

## Packages
{{ range $groupVersions }}
- {{ mdRenderGVLink . }}
{{- end -}}

{{- range $groupVersions }}
{{ template "gvDetails" . }}
{{ end }}

{{- end -}}
