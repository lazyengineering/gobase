// Copyright 2013 Jesse Allen. All rights reserved
// Released under the MIT license found in the LICENSE file.

// Provide some useful functions for use in templates
package filters

import (
	"github.com/russross/blackfriday"
	"html/template"
	"strings"
)

// Wraps github.com/russross/blackfriday.MarkdownCommon() for use in templates
func MarkdownCommon(raw string) template.HTML {
	return template.HTML(blackfriday.MarkdownCommon([]byte(raw)))
}

// Wraps github.com/russross/blackfriday.MarkdownBasic() for use in templates
func MarkdownBasic(raw string) template.HTML {
	return template.HTML(blackfriday.MarkdownBasic([]byte(raw)))
}

// does nothing but replace "@" with `at`
func CloakEmail(email, at string) template.HTML {
	return template.HTML(strings.Replace(email, "@", at, -1))
}

// All available Filters
var All = template.FuncMap{
	"markdownCommon": MarkdownCommon,
	"markdownBasic":  MarkdownBasic,
	"cloakEmail":     CloakEmail,
}
