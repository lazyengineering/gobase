// Copyright 2013 Jesse Allen. All rights reserved
// Released under the MIT license found in the LICENSE file.

package layouts

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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

func NilNilAction() Action {
	return Action(func(r *http.Request) (map[string]interface{}, error) {
		return nil, nil
	})
}

func CountNilAction(t *testing.T) Action {
	counter := 0
	return Action(func(r *http.Request) (map[string]interface{}, error) {
		counter++
		data := map[string]interface{}{
			"Count": counter,
		}
		t.Log(data)
		return data, nil
	})
}

func ErrorAction() Action {
	return Action(func(r *http.Request) (map[string]interface{}, error) {
		return nil, errors.New("Stock Error")
	})
}

func DefaultError(t *testing.T) ErrorHandler {
	errorCount := 0
	return ErrorHandler(func(w http.ResponseWriter, r *http.Request, e error) {
		errorCount++
		http.Error(w, fmt.Sprint(errorCount), 500)
	})
}

func TestAct(t *testing.T) {
	l, err := New(nil, "base", ".test/helper")
	if err != nil {
		t.Fatal(err)
	}

	type testResponse struct {
		Status int
		Body   string
	}
	type testCase struct {
		A       Action
		E       ErrorHandler
		V       Volatility
		T       []string
		R1      testResponse
		R2      testResponse
		RDelay1 testResponse
		RDelay2 testResponse
	}
	testCases := []testCase{
		{
			A: NilNilAction(),
			E: nil,
			V: NoVolatility,
			T: nil,
			R1: testResponse{
				Status: 200,
				Body:   "",
			},
			R2: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "",
			},
		},
		{
			A: NilNilAction(),
			E: nil,
			V: NoVolatility,
			T: []string{".test/base"},
			R1: testResponse{
				Status: 200,
				Body:   "",
			},
			R2: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "",
			},
		},
		{
			A: NilNilAction(),
			E: nil,
			V: LowVolatility,
			T: nil,
			R1: testResponse{
				Status: 200,
				Body:   "",
			},
			R2: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "",
			},
		},
		{
			A: NilNilAction(),
			E: nil,
			V: LowVolatility,
			T: []string{".test/base"},
			R1: testResponse{
				Status: 200,
				Body:   "",
			},
			R2: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "",
			},
		},
		{
			A: NilNilAction(),
			E: nil,
			V: MediumVolatility,
			T: nil,
			R1: testResponse{
				Status: 200,
				Body:   "",
			},
			R2: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "",
			},
		},
		{
			A: NilNilAction(),
			E: nil,
			V: MediumVolatility,
			T: []string{".test/base"},
			R1: testResponse{
				Status: 200,
				Body:   "",
			},
			R2: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "",
			},
		},
		{
			A: NilNilAction(),
			E: nil,
			V: HighVolatility,
			T: nil,
			R1: testResponse{
				Status: 200,
				Body:   "",
			},
			R2: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "",
			},
		},
		{
			A: NilNilAction(),
			E: nil,
			V: HighVolatility,
			T: []string{".test/base"},
			R1: testResponse{
				Status: 200,
				Body:   "",
			},
			R2: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "",
			},
		},
		{
			A: NilNilAction(),
			E: nil,
			V: ExtremeVolatility,
			T: nil,
			R1: testResponse{
				Status: 200,
				Body:   "",
			},
			R2: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "",
			},
		},
		{
			A: NilNilAction(),
			E: nil,
			V: ExtremeVolatility,
			T: []string{".test/base"},
			R1: testResponse{
				Status: 200,
				Body:   "",
			},
			R2: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "",
			},
		},

		{
			A: NilNilAction(),
			E: DefaultError(t),
			V: NoVolatility,
			T: nil, // should cause error
			R1: testResponse{
				Status: 500,
				Body:   "1",
			},
			R2: testResponse{
				Status: 500,
				Body:   "2",
			},
			RDelay1: testResponse{
				Status: 500,
				Body:   "3",
			},
			RDelay2: testResponse{
				Status: 500,
				Body:   "4",
			},
		},
		{
			A: NilNilAction(),
			E: DefaultError(t),
			V: NoVolatility,
			T: []string{".test/base"},
			R1: testResponse{
				Status: 200,
				Body:   "",
			},
			R2: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "",
			},
		},
		{
			A: NilNilAction(),
			E: DefaultError(t),
			V: LowVolatility,
			T: nil, // should cause error
			R1: testResponse{
				Status: 500,
				Body:   "1",
			},
			R2: testResponse{
				Status: 500,
				Body:   "2",
			},
			RDelay1: testResponse{
				Status: 500,
				Body:   "3",
			},
			RDelay2: testResponse{
				Status: 500,
				Body:   "4",
			},
		},
		{
			A: NilNilAction(),
			E: DefaultError(t),
			V: LowVolatility,
			T: []string{".test/base"},
			R1: testResponse{
				Status: 200,
				Body:   "",
			},
			R2: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "",
			},
		},
		{
			A: NilNilAction(),
			E: DefaultError(t),
			V: MediumVolatility,
			T: nil, // error condition
			R1: testResponse{
				Status: 500,
				Body:   "1",
			},
			R2: testResponse{
				Status: 500,
				Body:   "2",
			},
			RDelay1: testResponse{
				Status: 500,
				Body:   "3",
			},
			RDelay2: testResponse{
				Status: 500,
				Body:   "4",
			},
		},
		{
			A: NilNilAction(),
			E: DefaultError(t),
			V: MediumVolatility,
			T: []string{".test/base"},
			R1: testResponse{
				Status: 200,
				Body:   "",
			},
			R2: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "",
			},
		},
		{
			A: NilNilAction(),
			E: DefaultError(t),
			V: HighVolatility,
			T: nil, // error
			R1: testResponse{
				Status: 500,
				Body:   "1",
			},
			R2: testResponse{
				Status: 500,
				Body:   "2",
			},
			RDelay1: testResponse{
				Status: 500,
				Body:   "3",
			},
			RDelay2: testResponse{
				Status: 500,
				Body:   "4",
			},
		},
		{
			A: NilNilAction(),
			E: DefaultError(t),
			V: HighVolatility,
			T: []string{".test/base"},
			R1: testResponse{
				Status: 200,
				Body:   "",
			},
			R2: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "",
			},
		},
		{
			A: NilNilAction(),
			E: DefaultError(t),
			V: ExtremeVolatility,
			T: nil, // error
			R1: testResponse{
				Status: 500,
				Body:   "1",
			},
			R2: testResponse{
				Status: 500,
				Body:   "2",
			},
			RDelay1: testResponse{
				Status: 500,
				Body:   "3",
			},
			RDelay2: testResponse{
				Status: 500,
				Body:   "4",
			},
		},
		{
			A: NilNilAction(),
			E: DefaultError(t),
			V: ExtremeVolatility,
			T: []string{".test/base"},
			R1: testResponse{
				Status: 200,
				Body:   "",
			},
			R2: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "",
			},
		},

		{
			A: CountNilAction(t),
			E: nil,
			V: NoVolatility,
			T: nil,
			R1: testResponse{
				Status: 200,
				Body:   "",
			},
			R2: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "",
			},
		},
		{
			A: CountNilAction(t),
			E: nil,
			V: NoVolatility,
			T: []string{".test/base"},
			R1: testResponse{
				Status: 200,
				Body:   "1",
			},
			R2: testResponse{
				Status: 200,
				Body:   "1",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "1",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "1",
			},
		},
		{
			A: CountNilAction(t),
			E: nil,
			V: LowVolatility,
			T: nil,
			R1: testResponse{
				Status: 200,
				Body:   "",
			},
			R2: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "",
			},
		},
		{
			A: CountNilAction(t),
			E: nil,
			V: LowVolatility,
			T: []string{".test/base"},
			R1: testResponse{
				Status: 200,
				Body:   "1",
			},
			R2: testResponse{
				Status: 200,
				Body:   "1",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "1",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "1",
			},
		},
		{
			A: CountNilAction(t),
			E: nil,
			V: MediumVolatility,
			T: nil,
			R1: testResponse{
				Status: 200,
				Body:   "",
			},
			R2: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "",
			},
		},
		{
			A: CountNilAction(t),
			E: nil,
			V: MediumVolatility,
			T: []string{".test/base"},
			R1: testResponse{
				Status: 200,
				Body:   "1",
			},
			R2: testResponse{
				Status: 200,
				Body:   "1",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "1",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "1",
			},
		},
		{
			A: CountNilAction(t),
			E: nil,
			V: HighVolatility,
			T: nil,
			R1: testResponse{
				Status: 200,
				Body:   "",
			},
			R2: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "",
			},
		},
		{
			A: CountNilAction(t),
			E: nil,
			V: HighVolatility,
			T: []string{".test/base"},
			R1: testResponse{
				Status: 200,
				Body:   "1",
			},
			R2: testResponse{
				Status: 200,
				Body:   "1",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "2",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "2",
			},
		},
		{
			A: CountNilAction(t),
			E: nil,
			V: ExtremeVolatility,
			T: nil,
			R1: testResponse{
				Status: 200,
				Body:   "",
			},
			R2: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "",
			},
		},
		{
			A: CountNilAction(t),
			E: nil,
			V: ExtremeVolatility,
			T: []string{".test/base"},
			R1: testResponse{
				Status: 200,
				Body:   "1",
			},
			R2: testResponse{
				Status: 200,
				Body:   "2",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "3",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "4",
			},
		},

		{
			A: CountNilAction(t),
			E: DefaultError(t),
			V: NoVolatility,
			T: nil, // error
			R1: testResponse{
				Status: 500,
				Body:   "1",
			},
			R2: testResponse{
				Status: 500,
				Body:   "2",
			},
			RDelay1: testResponse{
				Status: 500,
				Body:   "3",
			},
			RDelay2: testResponse{
				Status: 500,
				Body:   "4",
			},
		},
		{
			A: CountNilAction(t),
			E: DefaultError(t),
			V: NoVolatility,
			T: []string{".test/base"},
			R1: testResponse{
				Status: 200,
				Body:   "1",
			},
			R2: testResponse{
				Status: 200,
				Body:   "1",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "1",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "1",
			},
		},
		{
			A: CountNilAction(t),
			E: DefaultError(t),
			V: LowVolatility,
			T: nil, // error
			R1: testResponse{
				Status: 500,
				Body:   "1",
			},
			R2: testResponse{
				Status: 500,
				Body:   "2",
			},
			RDelay1: testResponse{
				Status: 500,
				Body:   "3",
			},
			RDelay2: testResponse{
				Status: 500,
				Body:   "4",
			},
		},
		{
			A: CountNilAction(t),
			E: DefaultError(t),
			V: LowVolatility,
			T: []string{".test/base"},
			R1: testResponse{
				Status: 200,
				Body:   "1",
			},
			R2: testResponse{
				Status: 200,
				Body:   "1",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "1",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "1",
			},
		},
		{
			A: CountNilAction(t),
			E: DefaultError(t),
			V: MediumVolatility,
			T: nil, // error
			R1: testResponse{
				Status: 500,
				Body:   "1",
			},
			R2: testResponse{
				Status: 500,
				Body:   "2",
			},
			RDelay1: testResponse{
				Status: 500,
				Body:   "3",
			},
			RDelay2: testResponse{
				Status: 500,
				Body:   "4",
			},
		},
		{
			A: CountNilAction(t),
			E: DefaultError(t),
			V: MediumVolatility,
			T: []string{".test/base"},
			R1: testResponse{
				Status: 200,
				Body:   "1",
			},
			R2: testResponse{
				Status: 200,
				Body:   "1",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "1",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "1",
			},
		},
		{
			A: CountNilAction(t),
			E: DefaultError(t),
			V: HighVolatility,
			T: nil, // error
			R1: testResponse{
				Status: 500,
				Body:   "1",
			},
			R2: testResponse{
				Status: 500,
				Body:   "2",
			},
			RDelay1: testResponse{
				Status: 500,
				Body:   "3",
			},
			RDelay2: testResponse{
				Status: 500,
				Body:   "4",
			},
		},
		{
			A: CountNilAction(t),
			E: DefaultError(t),
			V: HighVolatility,
			T: []string{".test/base"},
			R1: testResponse{
				Status: 200,
				Body:   "1",
			},
			R2: testResponse{
				Status: 200,
				Body:   "1",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "2",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "2",
			},
		},
		{
			A: CountNilAction(t),
			E: DefaultError(t),
			V: ExtremeVolatility,
			T: nil, // error
			R1: testResponse{
				Status: 500,
				Body:   "1",
			},
			R2: testResponse{
				Status: 500,
				Body:   "2",
			},
			RDelay1: testResponse{
				Status: 500,
				Body:   "3",
			},
			RDelay2: testResponse{
				Status: 500,
				Body:   "4",
			},
		},
		{
			A: CountNilAction(t),
			E: DefaultError(t),
			V: ExtremeVolatility,
			T: []string{".test/base"},
			R1: testResponse{
				Status: 200,
				Body:   "1",
			},
			R2: testResponse{
				Status: 200,
				Body:   "2",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "3",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "4",
			},
		},

		{
			A: ErrorAction(),
			E: nil, // error
			V: NoVolatility,
			T: nil,
			R1: testResponse{
				Status: 200,
				Body:   "",
			},
			R2: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "",
			},
		},
		{
			A: ErrorAction(),
			E: nil,
			V: NoVolatility,
			T: []string{".test/base"},
			R1: testResponse{
				Status: 200,
				Body:   "",
			},
			R2: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "",
			},
		},
		{
			A: ErrorAction(),
			E: nil,
			V: LowVolatility,
			T: nil,
			R1: testResponse{
				Status: 200,
				Body:   "",
			},
			R2: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "",
			},
		},
		{
			A: ErrorAction(),
			E: nil,
			V: LowVolatility,
			T: []string{".test/base"},
			R1: testResponse{
				Status: 200,
				Body:   "",
			},
			R2: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "",
			},
		},
		{
			A: ErrorAction(),
			E: nil, // error
			V: MediumVolatility,
			T: nil,
			R1: testResponse{
				Status: 200,
				Body:   "",
			},
			R2: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "",
			},
		},
		{
			A: ErrorAction(),
			E: nil,
			V: MediumVolatility,
			T: []string{".test/base"},
			R1: testResponse{
				Status: 200,
				Body:   "",
			},
			R2: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "",
			},
		},
		{
			A: ErrorAction(),
			E: nil,
			V: HighVolatility,
			T: nil, // error
			R1: testResponse{
				Status: 200,
				Body:   "",
			},
			R2: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "",
			},
		},
		{
			A: ErrorAction(),
			E: nil,
			V: HighVolatility,
			T: []string{".test/base"},
			R1: testResponse{
				Status: 200,
				Body:   "",
			},
			R2: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "",
			},
		},
		{
			A: ErrorAction(),
			E: nil,
			V: ExtremeVolatility,
			T: nil, // error
			R1: testResponse{
				Status: 200,
				Body:   "",
			},
			R2: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "",
			},
		},
		{
			A: ErrorAction(),
			E: nil,
			V: ExtremeVolatility,
			T: []string{".test/base"},
			R1: testResponse{
				Status: 200,
				Body:   "",
			},
			R2: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay1: testResponse{
				Status: 200,
				Body:   "",
			},
			RDelay2: testResponse{
				Status: 200,
				Body:   "",
			},
		},

		{
			A: ErrorAction(),
			E: DefaultError(t),
			V: NoVolatility,
			T: nil, // error
			R1: testResponse{
				Status: 500,
				Body:   "1",
			},
			R2: testResponse{
				Status: 500,
				Body:   "2",
			},
			RDelay1: testResponse{
				Status: 500,
				Body:   "3",
			},
			RDelay2: testResponse{
				Status: 500,
				Body:   "4",
			},
		},
		{
			A: ErrorAction(),
			E: DefaultError(t),
			V: NoVolatility,
			T: []string{".test/base"},
			R1: testResponse{
				Status: 500,
				Body:   "1",
			},
			R2: testResponse{
				Status: 500,
				Body:   "2",
			},
			RDelay1: testResponse{
				Status: 500,
				Body:   "3",
			},
			RDelay2: testResponse{
				Status: 500,
				Body:   "4",
			},
		},
		{
			A: ErrorAction(),
			E: DefaultError(t),
			V: LowVolatility,
			T: nil,
			R1: testResponse{
				Status: 500,
				Body:   "1",
			},
			R2: testResponse{
				Status: 500,
				Body:   "2",
			},
			RDelay1: testResponse{
				Status: 500,
				Body:   "3",
			},
			RDelay2: testResponse{
				Status: 500,
				Body:   "4",
			},
		},
		{
			A: ErrorAction(),
			E: DefaultError(t),
			V: LowVolatility,
			T: []string{".test/base"},
			R1: testResponse{
				Status: 500,
				Body:   "1",
			},
			R2: testResponse{
				Status: 500,
				Body:   "2",
			},
			RDelay1: testResponse{
				Status: 500,
				Body:   "3",
			},
			RDelay2: testResponse{
				Status: 500,
				Body:   "4",
			},
		},
		{
			A: ErrorAction(),
			E: DefaultError(t),
			V: MediumVolatility,
			T: nil,
			R1: testResponse{
				Status: 500,
				Body:   "1",
			},
			R2: testResponse{
				Status: 500,
				Body:   "2",
			},
			RDelay1: testResponse{
				Status: 500,
				Body:   "3",
			},
			RDelay2: testResponse{
				Status: 500,
				Body:   "4",
			},
		},
		{
			A: ErrorAction(),
			E: DefaultError(t),
			V: MediumVolatility,
			T: []string{".test/base"},
			R1: testResponse{
				Status: 500,
				Body:   "1",
			},
			R2: testResponse{
				Status: 500,
				Body:   "2",
			},
			RDelay1: testResponse{
				Status: 500,
				Body:   "3",
			},
			RDelay2: testResponse{
				Status: 500,
				Body:   "4",
			},
		},
		{
			A: ErrorAction(),
			E: DefaultError(t),
			V: HighVolatility,
			T: nil,
			R1: testResponse{
				Status: 500,
				Body:   "1",
			},
			R2: testResponse{
				Status: 500,
				Body:   "2",
			},
			RDelay1: testResponse{
				Status: 500,
				Body:   "3",
			},
			RDelay2: testResponse{
				Status: 500,
				Body:   "4",
			},
		},
		{
			A: ErrorAction(),
			E: DefaultError(t),
			V: HighVolatility,
			T: []string{".test/base"},
			R1: testResponse{
				Status: 500,
				Body:   "1",
			},
			R2: testResponse{
				Status: 500,
				Body:   "2",
			},
			RDelay1: testResponse{
				Status: 500,
				Body:   "3",
			},
			RDelay2: testResponse{
				Status: 500,
				Body:   "4",
			},
		},
		{
			A: ErrorAction(),
			E: DefaultError(t),
			V: ExtremeVolatility,
			T: nil,
			R1: testResponse{
				Status: 500,
				Body:   "1",
			},
			R2: testResponse{
				Status: 500,
				Body:   "2",
			},
			RDelay1: testResponse{
				Status: 500,
				Body:   "3",
			},
			RDelay2: testResponse{
				Status: 500,
				Body:   "4",
			},
		},
		{
			A: ErrorAction(),
			E: DefaultError(t),
			V: ExtremeVolatility,
			T: []string{".test/base"},
			R1: testResponse{
				Status: 500,
				Body:   "1",
			},
			R2: testResponse{
				Status: 500,
				Body:   "2",
			},
			RDelay1: testResponse{
				Status: 500,
				Body:   "3",
			},
			RDelay2: testResponse{
				Status: 500,
				Body:   "4",
			},
		},
	}
	queue := []func(){}

	for idx, tc := range testCases {
		service := httptest.NewServer(l.Act(tc.A, tc.E, tc.V, tc.T...))

		if r1, err := http.Get(service.URL); err != nil {
			t.Error(err)
		} else if r1.StatusCode != tc.R1.Status {
			t.Error("test\t", idx, "\t1st call: expected:\tstatus ", tc.R1.Status, "\tactual:\tstatus ", r1.StatusCode)
		} else {
			body, errr := ioutil.ReadAll(r1.Body)
			if errr != nil {
				t.Error(errr)
			}
			if strings.TrimSpace(string(body)) != tc.R1.Body {
				t.Error("test\t", idx, "\t1st call: expected:\t", tc.R1.Body, "\tactual:\t", string(body))
			}
		}
		if r2, err := http.Get(service.URL); err != nil {
			t.Error(err)
		} else if r2.StatusCode != tc.R2.Status {
			t.Error("test\t", idx, "\t2nd call: expected:\tstatus ", tc.R2.Status, "\tactual:\tstatus ", r2.StatusCode)
		} else {
			body, errr := ioutil.ReadAll(r2.Body)
			if errr != nil {
				t.Error(errr)
			}
			if strings.TrimSpace(string(body)) != tc.R2.Body {
				t.Error("test\t", idx, "\t2nd call: expected:\t", tc.R2.Body, "\tactual:\t", string(body))
			}
		}

		if !testing.Short() {
			queue = append(queue, (func(s *httptest.Server, tCase testCase, i int) func() {
				return func() {
					if rDelay1, err := http.Get(s.URL); err != nil {
						t.Error(err)
					} else if rDelay1.StatusCode != tCase.R1.Status {
						t.Error("test\t", i, "\t1st Delay call: expected:\tstatus ", tCase.RDelay1.Status, "\tactual:\tstatus ", rDelay1.StatusCode)
					} else {
						body, errr := ioutil.ReadAll(rDelay1.Body)
						if errr != nil {
							t.Error(errr)
						}
						if strings.TrimSpace(string(body)) != tCase.RDelay1.Body {
							t.Error("test\t", i, "\t1st Delay call: expected:\t", tCase.RDelay1.Body, "\tactual:\t", string(body))
						}
					}
					if rDelay2, err := http.Get(s.URL); err != nil {
						t.Error(err)
					} else if rDelay2.StatusCode != tCase.RDelay2.Status {
						t.Error("test\t", i, "\t2nd Delay call: expected:\tstatus ", tCase.RDelay2.Status, "\tactual:\tstatus ", rDelay2.StatusCode)
					} else {
						body, errr := ioutil.ReadAll(rDelay2.Body)
						if errr != nil {
							t.Error(errr)
						}
						if strings.TrimSpace(string(body)) != tCase.RDelay2.Body {
							t.Error("test\t", i, "\t2nd Delay call: expected:\t", tCase.RDelay2.Body, "\tactual:\t", string(body))
						}
					}
				}
			})(service, tc, idx))
		}
	}

	if !testing.Short() {
		start := time.Now()
		for now := range time.Tick(10 * time.Second) {
			log.Print(now)
			if time.Since(start) > 6*time.Minute {
				break
			}
		}
		for _, fn := range queue {
			fn()
		}
	}
}
