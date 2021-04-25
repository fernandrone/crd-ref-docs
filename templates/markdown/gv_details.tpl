{{- define "gvDetails" -}}
{{- $gv := . }}
## {{ $gv.GroupVersionString }}

{{ $gv.Doc }}

{{- if $gv.Kinds  }}
### Resource Types
{{ range $gv.SortedKinds }}
- {{ $gv.TypeForKind . | mdRenderTypeLink }}
{{- end }}
{{ end }}

{{ range $gv.SortedTypes }}
{{ template "type" . }}
{{ end }}

{{- end -}}
