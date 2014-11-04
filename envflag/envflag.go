// Wrap the standard flag package parser to look for environment variables
// after command flags
package envflag

import (
	"flag"
	"os"
	"strings"
)

// Map the command flag to the environment variable
// This is only needed where a natural mapping of
//     flag-name -> FLAG_NAME
// is not there or a filter is needed to adapt the
// value from the environment variable
type FlagMap map[string]Flag

// Define a Flag's behavior where Name is the Environment variable
// and the Filter is function to be used on the value in order to
// transform into the expected value for a command-line flag
type Flag struct {
	Name   string
	Filter func(string) string
}

// Parse flags where command > environment > default
func Parse(m FlagMap) {
	flag.Parse()
	flag.VisitAll(func(f *flag.Flag) {
		mapping := Flag{}
		if s, ok := m[f.Name]; ok {
			if len(s.Name) > 0 {
				mapping.Name = s.Name
			}
			if s.Filter != nil {
				mapping.Filter = s.Filter
			}
		}
		if len(mapping.Name) == 0 {
			mapping.Name = strings.ToUpper(strings.Replace(f.Name, "-", "_", -1))
		}
		if mapping.Filter == nil {
			mapping.Filter = func(s string) string { return s }
		}
		if v := os.Getenv(mapping.Name); len(v) > 0 {
			f.Value.Set(mapping.Filter(v))
		}
	})
}
