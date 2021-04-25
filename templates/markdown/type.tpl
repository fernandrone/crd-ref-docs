{{- define "type" -}}
{{- $type := . -}}
{{- if mdShouldRenderType $type -}}

### {{ $type.Name }}
{{- if $type.IsAlias }}({{- mdRenderTypeLink $type.UnderlyingType  }}) {{- end }}

{{ $type.Doc -}}

{{ if $type.References }}

*Appears In:*

***
{{ range $type.SortedReferences }}
- {{ mdRenderTypeLink . }}
{{- end }}

***
{{- end }}

{{ if $type.Members -}}
| Field | Description |
|-------|-------------|
{{ if $type.GVK -}}
| *`apiVersion`* __string__ | `{{ $type.GVK.Group }}/{{ $type.GVK.Version }}` |
| *`kind`* __string__ | `{{ $type.GVK.Kind }}` |
{{ end -}}
{{ range $type.Members -}}
| *`{{ .Name  }}`* __{{ mdRenderType .Type }}__ | {{ template "type_members" . }} |
{{ end -}}
{{ end -}}

{{- end -}}
{{- end -}}
