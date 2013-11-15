// Copyright 2013 Jesse Allen. All rights reserved
// Released under the MIT license found in the LICENSE file.

package layouts

import (
	"bytes"
	"github.com/russross/blackfriday"
	"html/template"
	"net/http"
)

type Layout struct {
	patterns     []string
	functions    template.FuncMap
	baseTemplate string
}

func New(functions template.FuncMap, baseTemplate string, patterns ...string) *Layout {
	l := new(Layout)
	l.Init(functions, baseTemplate, patterns...)
	return l
}

func (l *Layout) Init(functions template.FuncMap, baseTemplate string, patterns ...string) {
	l.functions = functions
	l.baseTemplate = baseTemplate
	l.patterns = patterns
}

type Action func(*http.Request) (map[string]interface{}, error)

type ErrorHandler func(http.ResponseWriter, *http.Request, error)

func (l *Layout) Act(respond Action, eh ErrorHandler, templates ...string) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		t, err := l.load(templates...)
		if err != nil {
			eh(res, req, err)
			return
		}
		var data map[string]interface{}
		data, err = respond(req)
		if err != nil {
			eh(res, req, err)
			return
		}

		b := new(bytes.Buffer)
		if err = t.ExecuteTemplate(b, l.baseTemplate, data); err != nil {
			eh(res, req, err)
			return
		}
		if _, err = b.WriteTo(res); err != nil {
			eh(res, req, err)
		}
	})
}

// TODO: if performance becomes an issue, we can start caching the base templates, and cloning
func (l *Layout) load(patterns ...string) (*template.Template, error) {
	var err error
	// add some key helper functions to the templates
	b := template.New("base").Funcs(l.functions)
	for _, p := range append(l.patterns, patterns...) {
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
