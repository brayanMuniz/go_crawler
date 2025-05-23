package main

import (
	"reflect"
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected string
	}{
		{
			name:     "remove scheme",
			inputURL: "https://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove scheme and trailing slash",
			inputURL: "https://blog.boot.dev/path/",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "lowercase capital letters",
			inputURL: "https://BLOG.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove scheme, trailing slash, lowercase letters",
			inputURL: "https://BLOG.boot.dev/path/",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "root path with slash",
			inputURL: "https://blog.boot.dev/",
			expected: "blog.boot.dev/",
		},
		{
			name:     "root path without slash",
			inputURL: "https://blog.boot.dev",
			expected: "blog.boot.dev/",
		},
		{
			name:     "remove fragment",
			inputURL: "https://blog.boot.dev/path#section",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove query string",
			inputURL: "https://blog.boot.dev/path?foo=bar",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove query and fragment",
			inputURL: "https://blog.boot.dev/path/?foo=bar#section",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "keep port number",
			inputURL: "https://blog.boot.dev:8080/path",
			expected: "blog.boot.dev:8080/path",
		},
		{
			name:     "empty path",
			inputURL: "https://blog.boot.dev",
			expected: "blog.boot.dev/",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestGetUrlsFromHTML(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		inputBody string
		expected  []string
	}{
		{
			name:     "absolute and relative URLs",
			inputURL: "https://blog.boot.dev",
			inputBody: `
<html>
	<body>
		<a href="/path/one">
			<span>Boot.dev</span>
		</a>
		<a href="https://other.com/path/one">
			<span>Boot.dev</span>
		</a>
	</body>
</html>
`,
			expected: []string{"https://blog.boot.dev/path/one",
				"https://other.com/path/one"},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getURLsFromHTML(tc.inputBody, tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}

}
