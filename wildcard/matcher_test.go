/*-
 * Copyright 2018 Square Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package wildcard

import (
	"fmt"
	"testing"
)

func ExampleCompile_simple() {
	matcher, err := Compile("spiffe://some/*/pattern")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%t\n", matcher.Matches("spiffe://some/test/pattern"))
	fmt.Printf("%t\n", matcher.Matches("spiffe://some/test/example"))
	// Output:
	// true
	// false
}

func ExampleCompile_doubleWildcard() {
	matcher, err := Compile("spiffe://some/*/pattern/**")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%t\n", matcher.Matches("spiffe://some/test/pattern"))
	fmt.Printf("%t\n", matcher.Matches("spiffe://some/test/pattern/that/continues"))
	// Output:
	// true
	// true
}

func testMatches(t *testing.T, pattern string, matches []string, invalids []string) {
	matcher, err := Compile(pattern)
	if err != nil {
		t.Fatalf("bad pattern: '%s' (%s)", pattern, err)
	}

	t.Logf("testing pattern: '%s' => '%s'", pattern, matcher.(regexpMatcher).pattern.String())

	for _, candidate := range matches {
		if !matcher.Matches(candidate) {
			t.Errorf("missed: pattern '%s' didn't match string '%s', but should have", pattern, candidate)
		}
	}

	for _, candidate := range invalids {
		if matcher.Matches(candidate) {
			t.Errorf("bad match: pattern '%s' matched string '%s', but shouldn't have", pattern, candidate)
		}
	}
}

func TestMatchingSimple(t *testing.T) {
	testMatches(t,
		"spiffe://foo/*/bar",
		[]string{
			"spiffe://foo/baz/bar",
			"spiffe://foo/baz/bar/",
			"spiffe://foo/spam/bar",
			"spiffe://foo/spam/bar/",
			"spiffe://foo/asdf/bar",
			"spiffe://foo/asdf/bar/",
		},
		[]string{
			"invalid",
			"spiffe://foo/bar",
			"spiffe://foo//bar",
			"spiffe://foo/baz/bar/spam",
			"spiff://foo/asdf/bar",
			"spiffe://foox/baz/bar",
			"spiffe://foox/baz/bar/",
			"spiffe://foox/spam/bar",
			"spiffe://foox/spam/bar/",
			"spiffe://foox/asdf/bar",
			"spiffe://foox/asdf/bar/",
			"spiffe://foox/baz/barx",
			"spiffe://foox/baz/barx/",
			"spiffe://foox/spam/barx",
			"spiffe://foox/spam/barx/",
			"spiffe://foox/asdf/barx",
			"spiffe://foox/asdf/barx/",
		})
	testMatches(t,
		"spiffe://*/*/bar",
		[]string{
			"spiffe://foo/baz/bar",
			"spiffe://foo/baz/bar/",
			"spiffe://foo/spam/bar",
			"spiffe://foo/spam/bar/",
			"spiffe://foo/asdf/bar",
			"spiffe://foo/asdf/bar/",
			"spiffe://foox/baz/bar",
			"spiffe://foox/baz/bar/",
			"spiffe://foox/spam/bar",
			"spiffe://foox/spam/bar/",
			"spiffe://foox/asdf/bar",
			"spiffe://foox/asdf/bar/",
		},
		[]string{
			"invalid",
			"spiffe://foo/bar",
			"spiffe://foo//bar",
			"spiffe://foo/baz/bar/spam",
			"spiff://foo/asdf/bar",
			"spiffe://foox/baz/barx",
			"spiffe://foox/baz/barx/",
			"spiffe://foox/spam/barx",
			"spiffe://foox/spam/barx/",
			"spiffe://foox/asdf/barx",
			"spiffe://foox/asdf/barx/",
		})
	testMatches(t,
		"spiffe://foo/*/*",
		[]string{
			"spiffe://foo/baz/bar",
			"spiffe://foo/baz/bar/",
			"spiffe://foo/spam/bar",
			"spiffe://foo/spam/bar/",
			"spiffe://foo/asdf/bar",
			"spiffe://foo/asdf/bar/",
			"spiffe://foo/baz/barx",
			"spiffe://foo/baz/barx/",
			"spiffe://foo/spam/barx",
			"spiffe://foo/spam/barx/",
			"spiffe://foo/asdf/barx",
			"spiffe://foo/asdf/barx/",
		},
		[]string{
			"invalid",
			"spiffe://foo/bar",
			"spiffe://foo//bar",
			"spiffe://foo/baz/bar/spam",
			"spiff://foo/asdf/bar",
		})
	testMatches(t,
		"spiffe://*/*/*",
		[]string{
			"spiffe://foo/baz/bar",
			"spiffe://foo/baz/bar/",
			"spiffe://foo/spam/bar",
			"spiffe://foo/spam/bar/",
			"spiffe://foo/asdf/bar",
			"spiffe://foo/asdf/bar/",
			"spiffe://foo/baz/barx",
			"spiffe://foo/baz/barx/",
			"spiffe://foo/spam/barx",
			"spiffe://foo/spam/barx/",
			"spiffe://foo/asdf/barx",
			"spiffe://foo/asdf/barx/",
			"spiffe://foox/baz/barx",
			"spiffe://foox/baz/barx/",
			"spiffe://foox/spam/barx",
			"spiffe://foox/spam/barx/",
			"spiffe://foox/asdf/barx",
			"spiffe://foox/asdf/barx/",
		},
		[]string{
			"invalid",
			"spiffe://foo/bar",
			"spiffe://foo//bar",
			"spiffe://foo/baz/bar/spam",
			"spiff://foo/asdf/bar",
		})
}

func TestMatchingWithDouble(t *testing.T) {
	testMatches(t,
		"spiffe://foo/*/bar/**",
		[]string{
			"spiffe://foo/baz/bar",
			"spiffe://foo/baz/bar/",
			"spiffe://foo/baz/bar/spam",
			"spiffe://foo/spam/bar",
			"spiffe://foo/spam/bar/",
			"spiffe://foo/spam/bar/asdf",
			"spiffe://foo/spam/bar/asdf/qwer",
		},
		[]string{
			"invalid",
			"spiffe://foo/bar",
			"spiffe://foo//bar",
			"spiffe://foo//bar/asdf",
			"spiff://foo/asdf/bar",
			"spiffe://foo/baz/barf",
		})
	testMatches(t,
		"spiffe://foo/*/bar/**",
		[]string{
			"spiffe://foo/baz/bar",
			"spiffe://foo/baz/bar/",
			"spiffe://foo/baz/bar/spam",
			"spiffe://foo/spam/bar",
			"spiffe://foo/spam/bar/",
			"spiffe://foo/spam/bar/asdf",
			"spiffe://foo/spam/bar/asdf/qwer",
		},
		[]string{
			"invalid",
			"spiffe://foo/bar",
			"spiffe://foo//bar",
			"spiffe://foo//bar/asdf",
			"spiff://foo/asdf/bar",
			"spiffe://foo/baz/barf",
		})
}

func TestMatchingWithMetaChars(t *testing.T) {
	testMatches(t,
		// The '.' should not be interpreted as a regex char
		"spiffe://foo/./bar",
		[]string{
			"spiffe://foo/./bar",
			"spiffe://foo/./bar/",
		},
		[]string{
			"spiffe://foo/x/bar",
			"spiffe://foo/x/bar/",
		})
	testMatches(t,
		// The '.' should not be interpreted as a regex char
		"spiffe://././.",
		[]string{
			"spiffe://././.",
			"spiffe://./././",
		},
		[]string{
			"spiffe://a/b/c",
			"spiffe://a/b/c/",
		})
	testMatches(t,
		// The meta chars should not be interpreted as a regex
		".+",
		[]string{
			".+",
		},
		[]string{
			"invalid",
		})
}

func TestInvalidPatterns(t *testing.T) {
	for _, pattern := range []string{
		"",
		"test://foo*/asdf",
		"test://*foo/asdf",
		"test://**/asdf",
		"**://foo/asdf",
	} {
		_, err := Compile(pattern)
		if err == nil {
			t.Errorf("should reject invalid pattern '%s'", pattern)
		}
	}
}

func TestMustCompile(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("call to MustCompile did not panic with invalid pattern")
		}
	}()

	// Compile with valid pattern
	p := MustCompile("test/**")
	if p == nil {
		t.Error("MustCompile returned nil with valid pattern?")
	}

	// Compile with invalid pattern (should panic)
	MustCompile("**/test")
}

func TestCompileList(t *testing.T) {
	// Compile with valid patterns
	ms, err := CompileList([]string{
		"test",
		"test/**",
	})
	if err != nil {
		t.Errorf("CompileList failed with valid patterns: %s", err)
	}
	if len(ms) != 2 {
		t.Errorf("CompileList returned bad number of matchers (%d, wanted 2)", len(ms))
	}

	// Compile with invalid patterns
	ms, err = CompileList([]string{
		"test",
		"**/test",
	})
	if err == nil {
		t.Error("CompileList failed with invalid pattern in input")
	}
	if len(ms) != 0 {
		t.Errorf("CompileList returned bad number of matchers (%d, wanted 0)", len(ms))
	}
}
