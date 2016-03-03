// Copyright 2013-2016 Jesse Allen. All rights reserved
// Released under the MIT license found in the LICENSE file.

package filters

import (
	"testing"
)

func TestMarkdownCommon(t *testing.T) {
	const (
		md = `Title
====

Paragraph goes here.

* list item
* list items
* listed item
`
		html = `<h1>Title</h1>

<p>Paragraph goes here.</p>

<ul>
<li>list item</li>
<li>list items</li>
<li>listed item</li>
</ul>
`
	)
	actual := MarkdownCommon(md)
	if string(actual) != html {
		t.Error("expected:\t", html, "\nactual:\t", actual)
	}
}

func TestMarkdownBasic(t *testing.T) {
	const (
		md = `Title
====

Paragraph goes here.

* list item
* list items
* listed item
`
		html = `<h1>Title</h1>

<p>Paragraph goes here.</p>

<ul>
<li>list item</li>
<li>list items</li>
<li>listed item</li>
</ul>
`
	)
	actual := MarkdownBasic(md)
	if string(actual) != html {
		t.Error("expected:\t", html, "\nactual:\t", actual)
	}
}

func TestCloakEmail(t *testing.T) {
	const (
		inEmail  = "name@email.com"
		inAt     = "[at]"
		expected = "name[at]email.com"
	)
	if actual := CloakEmail(inEmail, inAt); string(actual) != expected {
		t.Error("expected:\t", expected, "\nactual:\t", actual)
	}
}
