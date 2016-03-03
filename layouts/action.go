// Copyright 2013-2016 Jesse Allen. All rights reserved
// Released under the MIT license found in the LICENSE file.

package layouts

import (
	"errors"
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

// Calls several actions together in order to return a superset of data.
// Duplicate keys in the returned map are resolved with the first action
// in the arguments taking priority, except where the value is a slice of
// strings, which will be merged.
func MergeActions(actions ...Action) Action {
	return func(req *http.Request) (map[string]interface{}, error) {
		done := make(chan struct{})
		defer close(done)

		c, errc := mergeActions(done, actions, req)

		unmergedData := make([]map[string]interface{}, len(actions), len(actions))
		for r := range c {
			if r.Err != nil {
				return nil, r.Err
			}
			unmergedData[r.Index] = r.Data
		}
		if err := <-errc; err != nil {
			return nil, err
		}
		mergedData := make(map[string]interface{})
		for _, d := range unmergedData {
			for key, val := range d {
				if mergedData[key] == nil {
					mergedData[key] = val
				} else if v1, isStringSlice := mergedData[key].([]string); isStringSlice {
					if v2, bothStringSlices := val.([]string); bothStringSlices {
						mergedData[key] = append(v1, v2...)
					}
				}
			}
		}
		return mergedData, nil
	}
}

type actionResult struct {
	Index int
	Data  map[string]interface{} // result data
	Err   error
}

func mergeActions(done <-chan struct{}, actions []Action, req *http.Request) (<-chan actionResult, <-chan error) {
	c := make(chan actionResult)
	errc := make(chan error, 1)
	go func() {
		var wg sync.WaitGroup
		var canceled error
		// loop over seed data
		for i, a := range actions {
			wg.Add(1)
			go func(idx int, action Action) {
				// do the work
				data, err := action(req)
				select {
				case c <- actionResult{idx, data, err}:
				case <-done:
				}
				wg.Done()
			}(i, a)
			select {
			case <-done:
				canceled = errors.New("actions canceled")
			default:
				canceled = nil
			}
		}
		go func() {
			wg.Wait()
			close(c)
		}()
		errc <- canceled
	}()
	return c, errc
}
