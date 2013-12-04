// Copyright 2013 Jesse Allen. All rights reserved
// Released under the MIT license found in the LICENSE file.

package layouts

type layoutError string

func (e layoutError) Error() string {
  return string(e)
}

// Error Messages
const (
	errNoBaseTemplate layoutError = "baseTemplate required but not provided"
)

