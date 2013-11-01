// Copyright 2013 Jesse Allen. All rights reserved
// Released under the MIT license found in the LICENSE file.

package main

import (
	"flag"
	"github.com/russross/blackfriday"
	"html/template"
)

var (
	LayoutTemplateGlob = flag.String("layouts", "templates/layouts/*.html", "Pattern for layout templates")
	HelperTemplateGlob = flag.String("helpers", "templates/helpers/*.html", "Pattern for helper templates")
)

/*
func init() {
	var err error

	if !flag.Parsed() {
		flag.Parse()
	}
}
*/

// Load base templates and templates from the provided pattern
// TODO: if performance becomes an issue, we can start caching the base templates, and cloning
func LoadTemplates(patterns ...string) (*template.Template, error) {
	var err error
	// add some key helper functions to the templates
	b := template.New("base").Funcs(template.FuncMap{
		"markdownCommon": func(raw string) template.HTML {
			return template.HTML(blackfriday.MarkdownCommon([]byte(raw)))
		},
		"markdownBasic": func(raw string) template.HTML {
			return template.HTML(blackfriday.MarkdownBasic([]byte(raw)))
		},
	})
	b, err = b.ParseGlob(*LayoutTemplateGlob)
	if err != nil {
		return nil, err
	}
	for _, p := range append([]string{*HelperTemplateGlob}, patterns...) {
		_, err = b.ParseGlob(p)
		if err != nil {
			return nil, err
		}
	}
	return b, nil
}
