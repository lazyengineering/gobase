// Copyright 2013-2016 Jesse Allen. All rights reserved
// Released under the MIT license found in the LICENSE file.

package layouts

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Volatility uint

const (
	NoVolatility      Volatility = iota // A request is not expected to have a different result for the lifetime of the application
	LowVolatility                       // A request should have a different result within one day of changes to source data
	MediumVolatility                    // A request should have a different result within an hour of changes to source data
	HighVolatility                      // A request should have a different result with five minutes of changes to source data
	ExtremeVolatility                   // A request should immediately reflect changes to source data
)

// Layout defines a collection of templates we can use throughout a site,
// including the "default" template that we execute.
type Layout struct {
	patterns     []string
	functions    template.FuncMap
	baseTemplate string
}

// The signature for a function that will be used when an error occurs with an Action
type ErrorHandler func(http.ResponseWriter, *http.Request, error)

// Create a new Layout from the provided function map, base template name, and set of
// Globs where template files can be located.
func New(functions template.FuncMap, baseTemplate string, patterns ...string) (*Layout, error) {
	l := new(Layout)
	err := l.Init(functions, baseTemplate, patterns...)
	if err != nil {
		return nil, err
	}
	return l, nil
}

// Initialize a layout with the provided function map, base template name, and set of
// Globs where template files can be located.
func (l *Layout) Init(functions template.FuncMap, baseTemplate string, patterns ...string) error {
	if len(baseTemplate) == 0 {
		return errNoBaseTemplate
	}
	l.functions = functions
	l.baseTemplate = baseTemplate
	l.patterns = patterns
	return nil
}

// Use Act in order to create an http.Handler that fills a template with the data from an executed Action
// or executes the ErrorHandler in case of an error.
func (l *Layout) Act(respond Action, eh ErrorHandler, volatility Volatility, templates ...string) http.Handler {
	var loadTemplates func() (*template.Template, error)
	var ttl time.Duration
	if eh == nil {
		eh = func(w http.ResponseWriter, r *http.Request, e error) {}
	}
	switch volatility {
	case NoVolatility:
		// Load templates so that we can clone instead of loading every time
		var storedTemplates *template.Template
		lock := sync.Mutex{}
		loadTemplates = func() (*template.Template, error) {
			var err error
			lock.Lock()
			defer lock.Unlock()
			if storedTemplates == nil {
				storedTemplates, err = l.load(templates...)
				if err != nil {
					return nil, err
				}
			}
			return storedTemplates.Clone()
		}
		respond = respond.cache(-1) // cache permanently
		ttl = 7 * 24 * time.Hour
	case LowVolatility:
		ttl = 24 * time.Hour
		fallthrough
	case MediumVolatility:
		if ttl == 0 {
			ttl = 1 * time.Hour
		}
		fallthrough
	case HighVolatility:
		if ttl == 0 {
			ttl = 5 * time.Minute
		}
		var storedTemplates *template.Template
		lock := sync.Mutex{}
		loadTemplates = func() (*template.Template, error) {
			var err error
			// lock to ensure we don't have multiple requests attempting to reload the
			// templates at the same time
			lock.Lock()
			defer lock.Unlock()
			if storedTemplates == nil {
				storedTemplates, err = l.load(templates...)
				if err != nil {
					return nil, err
				}
				time.AfterFunc(ttl, func() {
					lock.Lock()
					defer lock.Unlock()
					storedTemplates = nil
				})
			}
			return storedTemplates.Clone()
		}
		respond = respond.cache(ttl)
	case ExtremeVolatility:
		fallthrough // make this the default value
	default:
		loadTemplates = func() (*template.Template, error) {
			return l.load(templates...)
		}
	}
	// ensure that template loading will work
	template.Must(loadTemplates())
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		t, err := loadTemplates()
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
		// Add Client-Side caching
		if volatility < ExtremeVolatility {
			res.Header().Set("Cache-Control", "public, max-age="+strconv.FormatFloat(ttl.Seconds(), 'f', 0, 64))
			res.Header().Set("Expires", time.Now().Add(ttl).Format(time.RFC1123))
		}
		if _, err = b.WriteTo(res); err != nil {
			eh(res, req, err)
		}
	})
}

func (l *Layout) load(patterns ...string) (*template.Template, error) {
	t := time.Now()
	var err error
	// add some key helper functions to the templates
	b := template.New("base").Funcs(l.functions)
	for _, p := range append(l.patterns, patterns...) {
		_, err = b.ParseGlob(p)
		if err != nil {
			return nil, err
		}
	}
	log.Printf("\x1b[1;35mTemplates:\x1b[0m \x1b[34m%6d\x1b[0mµs \x1b[33m%v\x1b[0m", time.Since(t).Nanoseconds()/1000, append(l.patterns, patterns...))
	return b, nil
}
