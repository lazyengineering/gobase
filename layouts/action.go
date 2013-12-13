// Copyright 2013 Jesse Allen. All rights reserved
// Released under the MIT license found in the LICENSE file.

package layouts

import (
	"net/http"
	"sync"
	"time"
)

// An Action does the unique work for an http response where the result should be
// a page rendered from a template executed with unique data.
type Action func(*http.Request) (map[string]interface{}, error)

// Returns an Action that runs the original Action when there is no cached value.
// The cached value is unset after the given ttl (time to live) duration.
// A negative ttl will permanently cache
func (a Action) cache(ttl time.Duration) Action {
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
			if ttl > 0 {
				time.AfterFunc(ttl, func() {
					lock.Lock()
					data = nil
					lock.Unlock()
				})
			}
		}
		return data, err
	}
}
