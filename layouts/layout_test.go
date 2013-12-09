// Copyright 2013 Jesse Allen. All rights reserved
// Released under the MIT license found in the LICENSE file.

package layouts

import (
	"testing"
)

func TestNew(t *testing.T) {
	if l, err := New(nil, ""); err != errNoBaseTemplate {
		t.Error(errNoBaseTemplate)
	} else if l != nil {
		t.Error("Layout should be nil on error")
	}
	if l, err := New(nil, "base"); err != nil {
		t.Error("New Layout with nil function map, defined baseTemplate, and no patterns")
	} else if l == nil {
		t.Error("Layout should be non-nil when no error")
	}
}

// Even though this covers exactly what TestNew covered, it's still part of the contract
func TestInit(t *testing.T) {
	l := new(Layout)
	if err := l.Init(nil, ""); err != errNoBaseTemplate {
		t.Error(errNoBaseTemplate)
	}
	if err := l.Init(nil, "base"); err != nil {
		t.Error("Init Layout with nil function map, defined baseTemplate, and no patterns")
	}
}
