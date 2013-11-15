// Copyright 2013 Jesse Allen. All rights reserved
// Released under the MIT license found in the LICENSE file.

package main

import (
	"github.com/russross/blackfriday"
	"html/template"
)

type TemplateSet struct {
	LayoutGlob, HelperGlob string
	Functions              template.FuncMap
}

// TODO: if performance becomes an issue, we can start caching the base templates, and cloning
func (s *TemplateSet) Load(patterns ...string) (*template.Template, error) {
	var err error
	// add some key helper functions to the templates
	b := template.New("base").Funcs(s.Functions)
	b, err = b.ParseGlob(s.LayoutGlob)
	if err != nil {
		return nil, err
	}
	for _, p := range append([]string{s.HelperGlob}, patterns...) {
		_, err = b.ParseGlob(p)
		if err != nil {
			return nil, err
		}
	}
	return b, nil
}

func BasicFunctionMap() template.FuncMap {
	return template.FuncMap{
		"markdownCommon": func(raw string) template.HTML {
			return template.HTML(blackfriday.MarkdownCommon([]byte(raw)))
		},
		"markdownBasic": func(raw string) template.HTML {
			return template.HTML(blackfriday.MarkdownBasic([]byte(raw)))
		},
	}
}
