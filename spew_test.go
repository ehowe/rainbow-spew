/*
 * Copyright (c) 2013-2016 Dave Collins <dave@davec.name>
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package spew_test

import (
	"bytes"
	"fmt"
	"os"

	spew "github.com/ehowe/rainbow-spew"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// spewFunc is used to identify which public function of the spew package or
// ConfigState a test applies to.
type spewFunc int

const (
	fCSFdump spewFunc = iota
	fCSFprint
	fCSFprintf
	fCSFprintln
	fCSPrint
	fCSPrintln
	fCSSdump
	fCSSprint
	fCSSprintf
	fCSSprintln
	fCSErrorf
	fCSNewFormatter
	fErrorf
	fFprint
	fFprintln
	fPrint
	fPrintln
	fSdump
	fSprint
	fSprintf
	fSprintln
)

// Map of spewFunc values to names for pretty printing.
var spewFuncStrings = map[spewFunc]string{
	fCSFdump:        "ConfigState.Fdump",
	fCSFprint:       "ConfigState.Fprint",
	fCSFprintf:      "ConfigState.Fprintf",
	fCSFprintln:     "ConfigState.Fprintln",
	fCSSdump:        "ConfigState.Sdump",
	fCSPrint:        "ConfigState.Print",
	fCSPrintln:      "ConfigState.Println",
	fCSSprint:       "ConfigState.Sprint",
	fCSSprintf:      "ConfigState.Sprintf",
	fCSSprintln:     "ConfigState.Sprintln",
	fCSErrorf:       "ConfigState.Errorf",
	fCSNewFormatter: "ConfigState.NewFormatter",
	fErrorf:         "spew.Errorf",
	fFprint:         "spew.Fprint",
	fFprintln:       "spew.Fprintln",
	fPrint:          "spew.Print",
	fPrintln:        "spew.Println",
	fSdump:          "spew.Sdump",
	fSprint:         "spew.Sprint",
	fSprintf:        "spew.Sprintf",
	fSprintln:       "spew.Sprintln",
}

func (f spewFunc) String() string {
	if s, ok := spewFuncStrings[f]; ok {
		return s
	}
	return fmt.Sprintf("Unknown spewFunc (%d)", int(f))
}

// spewTest is used to describe a test to be performed against the public
// functions of the spew package or ConfigState.
type spewTest struct {
	cs     *spew.ConfigState
	f      spewFunc
	format string
	in     interface{}
	want   string
}

type ptrTester struct {
	s *struct{}
}

type depthTester struct {
	ic    indirCir1
	arr   [1]string
	slice []string
	m     map[string]int
}

// spewTests houses the tests to be performed against the public functions of
// the spew package and ConfigState.
//
// These tests are only intended to ensure the public functions are exercised
// and are intentionally not exhaustive of types.  The exhaustive type
// tests are handled in the dump and format tests.
var spewTests []spewTest

// redirStdout is a helper function to return the standard output from f as a
// byte slice.
func redirStdout(f func()) ([]byte, error) {
	tempFile, err := os.CreateTemp("", "ss-test")
	if err != nil {
		return nil, err
	}
	fileName := tempFile.Name()
	defer os.Remove(fileName) // Ignore error

	origStdout := os.Stdout
	os.Stdout = tempFile
	f()
	os.Stdout = origStdout
	tempFile.Close()

	return os.ReadFile(fileName)
}

var _ = Describe("Spew Tests", func() {
	var scsDefault *spew.ConfigState
	var scsNoMethods *spew.ConfigState
	var scsNoPmethods *spew.ConfigState
	var scsMaxDepth *spew.ConfigState
	var scsContinue *spew.ConfigState
	var scsNoPtrAddr *spew.ConfigState
	var scsNoCap *spew.ConfigState

	var ts stringer
	var tps pstringer
	var tptr *ptrTester
	var dt depthTester
	var te customError

	BeforeEach(func() {
		scsDefault = spew.NewTestConfig()
		scsNoMethods = &spew.ConfigState{Indent: " ", DisableMethods: true}
		scsNoPmethods = &spew.ConfigState{Indent: " ", DisablePointerMethods: true}
		scsMaxDepth = &spew.ConfigState{Indent: " ", MaxDepth: 1}
		scsContinue = &spew.ConfigState{Indent: " ", ContinueOnMethod: true}
		scsNoPtrAddr = &spew.ConfigState{DisablePointerAddresses: true}
		scsNoCap = &spew.ConfigState{DisableCapacities: true}

		ts = stringer("test")
		tps = pstringer("test")
		tptr = &ptrTester{s: &struct{}{}}
		dt = depthTester{indirCir1{nil}, [1]string{"arr"}, []string{"slice"}, map[string]int{"one": 1}}
		te = customError(10)
	})

	DescribeTable(
		"all of the spew tests",
		func(scsFn func() *spew.ConfigState, spewFunc spewFunc, format string, inFunc func() interface{}, want string) {
			config := scsFn()
			in := inFunc()
			buf := new(bytes.Buffer)

			switch spewFunc {
			case fCSFdump:
				config.Fdump(buf, in)

			case fCSFprint:
				config.Fprint(buf, in)

			case fCSFprintf:
				config.Fprintf(buf, format, in)

			case fCSFprintln:
				config.Fprintln(buf, in)

			case fCSPrint:
				b, err := redirStdout(func() { config.Print(in) })
				Expect(err).To(BeNil())
				buf.Write(b)

			case fCSPrintln:
				b, err := redirStdout(func() { config.Println(in) })
				Expect(err).To(BeNil())
				buf.Write(b)

			case fCSSdump:
				str := config.Sdump(in)
				buf.WriteString(str)

			case fCSSprint:
				str := config.Sprint(in)
				buf.WriteString(str)

			case fCSSprintf:
				str := config.Sprintf(format, in)
				buf.WriteString(str)

			case fCSSprintln:
				str := config.Sprintln(in)
				buf.WriteString(str)

			case fCSErrorf:
				err := config.Errorf(format, in)
				buf.WriteString(err.Error())

			case fCSNewFormatter:
				fmt.Fprintf(buf, format, config.NewFormatter(in))

			case fErrorf:
				err := spew.Errorf(format, in)
				buf.WriteString(err.Error())

			case fFprint:
				spew.Fprint(buf, in)

			case fFprintln:
				spew.Fprintln(buf, in)

			case fPrint:
				b, err := redirStdout(func() { spew.Print(in) })
				Expect(err).To(BeNil())
				buf.Write(b)

			case fPrintln:
				b, err := redirStdout(func() { spew.Println(in) })
				Expect(err).To(BeNil())
				buf.Write(b)

			case fSdump:
				str := spew.Sdump(in)
				buf.WriteString(str)

			case fSprint:
				str := spew.Sprint(in)
				buf.WriteString(str)

			case fSprintf:
				str := spew.Sprintf(format, in)
				buf.WriteString(str)

			case fSprintln:
				str := spew.Sprintln(in)
				buf.WriteString(str)

			default:
				Fail(fmt.Sprintf("%v unrecognized function", spewFunc))
			}
			s := buf.String()
			Expect(s).To(Equal(want))
		},
		Entry("Entry 1", func() *spew.ConfigState { return scsDefault }, fCSFdump, "", func() interface{} { return int8(127) }, "(int8) 127\n"),
		Entry("Entry 2", func() *spew.ConfigState { return scsDefault }, fCSFprint, "", func() interface{} { return int16(32767) }, "32767"),
		Entry("Entry 3", func() *spew.ConfigState { return scsDefault }, fCSFprintf, "%v", func() interface{} { return int32(2147483647) }, "2147483647"),
		Entry("Entry 4", func() *spew.ConfigState { return scsDefault }, fCSFprintln, "", func() interface{} { return int(2147483647) }, "2147483647\n"),
		Entry("Entry 5", func() *spew.ConfigState { return scsDefault }, fCSPrint, "", func() interface{} { return int64(9223372036854775807) }, "9223372036854775807"),
		Entry("Entry 6", func() *spew.ConfigState { return scsDefault }, fCSPrintln, "", func() interface{} { return uint8(255) }, "255\n"),
		Entry("Entry 7", func() *spew.ConfigState { return scsDefault }, fCSSdump, "", func() interface{} { return uint8(64) }, "(uint8) 64\n"),
		Entry("Entry 8", func() *spew.ConfigState { return scsDefault }, fCSSprint, "", func() interface{} { return complex(1, 2) }, "(1+2i)"),
		Entry("Entry 9", func() *spew.ConfigState { return scsDefault }, fCSSprintf, "%v", func() interface{} { return complex(float32(3), 4) }, "(3+4i)"),
		Entry("Entry 10", func() *spew.ConfigState { return scsDefault }, fCSSprintln, "", func() interface{} { return complex(float64(5), 6) }, "(5+6i)\n"),
		Entry("Entry 11", func() *spew.ConfigState { return scsDefault }, fCSErrorf, "%#v", func() interface{} { return uint16(65535) }, "(uint16)65535"),
		Entry("Entry 12", func() *spew.ConfigState { return scsDefault }, fCSNewFormatter, "%v", func() interface{} { return uint32(4294967295) }, "4294967295"),
		Entry("Entry 13", func() *spew.ConfigState { return scsDefault }, fErrorf, "%v", func() interface{} { return uint64(18446744073709551615) }, "18446744073709551615"),
		Entry("Entry 14", func() *spew.ConfigState { return scsDefault }, fFprint, "", func() interface{} { return float32(3.14) }, "3.14"),
		Entry("Entry 15", func() *spew.ConfigState { return scsDefault }, fFprintln, "", func() interface{} { return float64(6.28) }, "6.28\n"),
		Entry("Entry 16", func() *spew.ConfigState { return scsDefault }, fPrint, "", func() interface{} { return true }, "true"),
		Entry("Entry 17", func() *spew.ConfigState { return scsDefault }, fPrintln, "", func() interface{} { return false }, "false\n"),
		Entry("Entry 18", func() *spew.ConfigState { return scsDefault }, fSdump, "", func() interface{} { return complex(-10, -20) }, "(complex128) (-10-20i)\n"),
		Entry("Entry 19", func() *spew.ConfigState { return scsDefault }, fSprint, "", func() interface{} { return complex(-1, -2) }, "(-1-2i)"),
		Entry("Entry 20", func() *spew.ConfigState { return scsDefault }, fSprintf, "%v", func() interface{} { return complex(float32(-3), -4) }, "(-3-4i)"),
		Entry("Entry 21", func() *spew.ConfigState { return scsDefault }, fSprintln, "", func() interface{} { return complex(float64(-5), -6) }, "(-5-6i)\n"),
		Entry("Entry 22", func() *spew.ConfigState { return scsNoMethods }, fCSFprint, "", func() interface{} { return ts }, "test"),
		Entry("Entry 23", func() *spew.ConfigState { return scsNoMethods }, fCSFprint, "", func() interface{} { return &ts }, "<*>test"),
		Entry("Entry 24", func() *spew.ConfigState { return scsNoMethods }, fCSFprint, "", func() interface{} { return tps }, "test"),
		Entry("Entry 25", func() *spew.ConfigState { return scsNoMethods }, fCSFprint, "", func() interface{} { return &tps }, "<*>test"),
		Entry("Entry 26", func() *spew.ConfigState { return scsNoPmethods }, fCSFprint, "", func() interface{} { return ts }, "stringer test"),
		Entry("Entry 27", func() *spew.ConfigState { return scsNoPmethods }, fCSFprint, "", func() interface{} { return &ts }, "<*>stringer test"),
		Entry("Entry 28", func() *spew.ConfigState { return scsNoPmethods }, fCSFprint, "", func() interface{} { return tps }, "test"),
		Entry("Entry 29", func() *spew.ConfigState { return scsNoPmethods }, fCSFprint, "", func() interface{} { return &tps }, "<*>stringer test"),
		Entry("Entry 30", func() *spew.ConfigState { return scsMaxDepth }, fCSFprint, "", func() interface{} { return dt }, "{{<max>} [<max>] [<max>] map[<max>]}"),
		Entry("Entry 31", func() *spew.ConfigState { return scsMaxDepth }, fCSFdump, "", func() interface{} { return dt }, "(spew_test.depthTester) {\n ic: (spew_test.indirCir1) {\n  <max depth reached>\n },\n arr: ([1]string) (len: 1 cap: 1) {\n  <max depth reached>\n },\n slice: ([]string) (len: 1 cap: 1) {\n  <max depth reached>\n },\n m: (map[string]int) (len: 1) {\n  <max depth reached>\n }\n}\n"),
		Entry("Entry 32", func() *spew.ConfigState { return scsContinue }, fCSFprint, "", func() interface{} { return ts }, "(stringer test) test"),
		Entry("Entry 33", func() *spew.ConfigState { return scsContinue }, fCSFdump, "", func() interface{} { return ts }, "(spew_test.stringer) (len: 4) (stringer test) \"test\"\n"),
		Entry("Entry 34", func() *spew.ConfigState { return scsContinue }, fCSFprint, "", func() interface{} { return te }, "(error: 10) 10"),
		Entry("Entry 35", func() *spew.ConfigState { return scsContinue }, fCSFdump, "", func() interface{} { return te }, "(spew_test.customError) (error: 10) 10\n"),
		Entry("Entry 36", func() *spew.ConfigState { return scsNoPtrAddr }, fCSFprint, "", func() interface{} { return tptr }, "<*>{<*>{}}"),
		Entry("Entry 37", func() *spew.ConfigState { return scsNoPtrAddr }, fCSSdump, "", func() interface{} { return tptr }, "(*spew_test.ptrTester)({\ns: (*struct {})({\n})\n})\n"),
		Entry("Entry 38", func() *spew.ConfigState { return scsNoCap }, fCSSdump, "", func() interface{} { return make([]string, 0, 10) }, "([]string) {\n}\n"),
		Entry("Entry 39", func() *spew.ConfigState { return scsNoCap }, fCSSdump, "", func() interface{} { return make([]string, 1, 10) }, "([]string) (len: 1) {\n(string) \"\"\n}\n"),
	)
})
