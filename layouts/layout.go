// Copyright 2013 Jesse Allen. All rights reserved
// Released under the MIT license found in the LICENSE file.

package layouts

import (
	"bytes"
	"github.com/russross/blackfriday"
	"html/template"
	"net/http"
	"sync"
	"time"
)

// Layout defines a collection of templates we can use throughout a site,
// including the "default" template that we execute.
type Layout struct {
	patterns     []string
	functions    template.FuncMap
	baseTemplate string
}

// Create a new Layout from the provided function map, base template name, and set of
// Globs where template files can be located.
func New(functions template.FuncMap, baseTemplate string, patterns ...string) *Layout {
	l := new(Layout)
	l.Init(functions, baseTemplate, patterns...)
	return l
}

// Initialize a layout with the provided function map, base template name, and set of
// Globs where template files can be located.
func (l *Layout) Init(functions template.FuncMap, baseTemplate string, patterns ...string) {
	l.functions = functions
	l.baseTemplate = baseTemplate
	l.patterns = patterns
}

// An Action does the unique work for an http response where the result should be
// a page rendered from a template executed with unique data.
type Action func(*http.Request) (map[string]interface{}, error)

// Returns an Action that runs the original Action when there is no cached value.
// The cached value is unset after the given ttl (time to live) duration.
func (a Action) Cache(ttl time.Duration) Action {
	var data map[string]interface{}
	lock := sync.RWMutex{}
	return func(r *http.Request) (map[string]interface{}, error) {
		lock.RLock()
		if data != nil {
			lock.RUnlock()
			return data, nil
		}
		lock.RUnlock()

		lock.Lock()
		defer lock.Unlock()
		var err error
		data, err = a(r)
		if data != nil {
			time.AfterFunc(ttl, func() {
				lock.Lock()
				data = nil
				lock.Unlock()
			})
		}
		return data, err
	}
}

// The signature for a function that will be used when an error occurs with an Action
type ErrorHandler func(http.ResponseWriter, *http.Request, error)

// Use Act in order to create an http.Handler that fills a template with the data from an executed Action
// or executes the ErrorHandler in case of an error.
func (l *Layout) Act(respond Action, eh ErrorHandler, templates ...string) http.Handler {
	// Load templates so that we can clone instead of loading every time
	permanentTemplates := template.Must(l.load(templates...))
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		t, err := permanentTemplates.Clone()
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

// BasicFunctionMap returns a simple template.FuncMap containing the following helper functions
//
//     markdownCommon(raw string) template.HTML // process markdown formatted string using github.com/russross/blackfriday MarkdownCommon
//     markdownBasic(raw string)  template.HTML // process markdown formatted string using github.com/russross/blackfriday MarkdownBasic
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
