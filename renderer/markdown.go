// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
package renderer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/elastic/crd-ref-docs/config"
	"github.com/elastic/crd-ref-docs/types"
)

const (
	markdownDefaultOutputFile = "out.md"
)

type MarkdownRenderer struct {
	conf *config.Config
	*Functions
}

func NewMarkdownRenderer(conf *config.Config) (*MarkdownRenderer, error) {
	baseFuncs, err := NewFunctions(conf)
	if err != nil {
		return nil, err
	}
	return &MarkdownRenderer{conf: conf, Functions: baseFuncs}, nil
}

func (md *MarkdownRenderer) Render(gvd []types.GroupVersionDetails) error {
	funcMap := combinedFuncMap(funcMap{prefix: "md", funcs: md.ToFuncMap()}, funcMap{funcs: sprig.TxtFuncMap()})
	tmpl, err := loadTemplate(md.conf.TemplatesDir, funcMap)
	if err != nil {
		return err
	}

	outputFile := md.conf.OutputPath

	if outputFile == "" {
		outputFile = markdownDefaultOutputFile
	}

	finfo, err := os.Stat(outputFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if finfo != nil && finfo.IsDir() {
		outputFile = filepath.Join(outputFile, markdownDefaultOutputFile)
	}

	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.ExecuteTemplate(f, mainTemplate, gvd)
}

func (md *MarkdownRenderer) ToFuncMap() template.FuncMap {
	return template.FuncMap{
		"GroupVersionID":     md.GroupVersionID,
		"RenderAnchorID":     md.RenderAnchorID,
		"RenderExternalLink": md.RenderExternalLink,
		"RenderGVLink":       md.RenderGVLink,
		"RenderLocalLink":    md.RenderLocalLink,
		"RenderType":         md.RenderType,
		"RenderTypeLink":     md.RenderTypeLink,
		"SafeID":             md.SafeID,
		"ShouldRenderType":   md.ShouldRenderType,
		"TypeID":             md.TypeID,
	}
}

func (md *MarkdownRenderer) ShouldRenderType(t *types.Type) bool {
	return t != nil && (t.GVK != nil || len(t.References) > 0)
}

func (md *MarkdownRenderer) RenderType(t *types.Type) string {
	var sb strings.Builder
	switch t.Kind {
	case types.MapKind:
		sb.WriteString("object (")
		sb.WriteString("keys:")
		sb.WriteString(md.RenderTypeLink(t.KeyType))
		sb.WriteString(", values:")
		sb.WriteString(md.RenderTypeLink(t.ValueType))
		sb.WriteString(")")
	case types.ArrayKind, types.SliceKind:
		sb.WriteString(md.RenderTypeLink(t.UnderlyingType))
		sb.WriteString(" array")
	default:
		sb.WriteString(md.RenderTypeLink(t))
	}

	return sb.String()
}

func (md *MarkdownRenderer) RenderTypeLink(t *types.Type) string {
	text := md.SimplifiedTypeName(t)

	link, local := md.LinkForType(t)
	if link == "" {
		return text
	}

	if local {
		return md.RenderLocalLink(link, text)
	} else {
		return md.RenderExternalLink(link, text)
	}
}

func (md *MarkdownRenderer) RenderLocalLink(link, text string) string {
	return fmt.Sprintf("[%s](#%s)", text, link)
}

func (md *MarkdownRenderer) RenderExternalLink(link, text string) string {
	return fmt.Sprintf("[%s](%s)", text, link)
}

func (md *MarkdownRenderer) RenderGVLink(gv types.GroupVersionDetails) string {
	return md.RenderLocalLink(md.GroupVersionID(gv), gv.GroupVersionString())
}

func (md *MarkdownRenderer) RenderAnchorID(id string) string {
	return fmt.Sprintf("%s", md.SafeID(id))
}
