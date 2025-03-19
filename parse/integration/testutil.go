package integration

import (
	"strings"
	"testing"

	"github.com/theplant/testingutils"
	"github.com/zhangshanwen/html2go/parse"
)

// HTMLGoTestCase defines a single test case for HTML to Go code conversion
type HTMLGoTestCase struct {
	Name         string
	Pkg          string
	ChildrenMode bool
	HTML         string
	GoCode       string
}

// RunHTMLGoTestCase runs a single test case and reports any differences between expected and actual output
func RunHTMLGoTestCase(t *testing.T, tc HTMLGoTestCase) {
	t.Helper()

	// Convert HTML to Go code
	gocode := parse.GenerateHTMLGo(tc.Pkg, tc.ChildrenMode, strings.NewReader(
		strings.ReplaceAll(tc.HTML, "|backquote|", "`"),
	))

	// Compare expected and actual output
	diff := testingutils.PrettyJsonDiff(strings.ReplaceAll(tc.GoCode, "|backquote|", "`"), gocode)

	if len(diff) > 0 {
		t.Error(diff)
	}
}

// RunHTMLGoTestCases runs all test cases in the provided slice
func RunHTMLGoTestCases(t *testing.T, testCases []HTMLGoTestCase) {
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			RunHTMLGoTestCase(t, tc)
		})
	}
}
