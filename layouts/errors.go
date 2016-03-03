// Copyright 2013-2016 Jesse Allen. All rights reserved
// Released under the MIT license found in the LICENSE file.

package layouts

// simple error to eliminate the need for the errors package
type layoutError string

func (e layoutError) Error() string {
	return string(e)
}

// Error Messages used in this package
const (
	errNoBaseTemplate layoutError = "layouts: baseTemplate required but not provided"
)
