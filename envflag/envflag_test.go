package envflag

import (
	"flag"
	"os"
	"testing"
)

func buildStringFlags() map[string]*string {
	return map[string]*string{
		"passthrough": flag.String("test-first-flag", "pass", "First Flag"),   // Passthrough
		"auto":        flag.String("test-second-flag", "fail", "Second Flag"), // Auto
		"manual no name no filter":   flag.String("test-third-flag", "fail", "Third Flag"),   // Manual Empty Value
		"manual yes name no filter":  flag.String("test-fourth-flag", "fail", "Fourth Flag"), // Manual Name No Filter
		"manual no name yes filter":  flag.String("test-fifth-flag", "fail", "Fifth Flag"),   // Manual No Name Filter
		"manual yes name yes filter": flag.String("test-sixth-flag", "fail", "Sixth Flag"),   // Manual Name Filter
	}
}

// No Need to test if the command-line flags work. So long as the flags package works, we only need to test that the ENV variables overwrite defaults
func TestParse(t *testing.T) {
	stringFlags := buildStringFlags()
	flagMap := FlagMap{
		"test-third-flag": Flag{},
		"test-fourth-flag": Flag{
			Name: "FOURTH",
		},
		"test-fifth-flag": Flag{
			Filter: func(s string) string {
				if s == "bar" {
					return "pass"
				} else {
					return "fail"
				}
			},
		},
		"test-sixth-flag": Flag{
			Name: "SIXTH",
			Filter: func(s string) string {
				if s == "bar" {
					return "pass"
				} else {
					return "fail"
				}
			},
		},
	}
	if err := os.Setenv("TEST_SECOND_FLAG", "pass"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("TEST_THIRD_FLAG", "pass"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("TEST_FOURTH_FLAG", "fail"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("FOURTH", "pass"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("TEST_FIFTH_FLAG", "bar"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("TEST_SIXTH_FLAG", "fail"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("SIXTH", "bar"); err != nil {
		t.Fatal(err)
	}
	Parse(flagMap)
	for testName, flg := range stringFlags {
		if *flg != "pass" {
			t.Error(testName)
		}
	}
}
