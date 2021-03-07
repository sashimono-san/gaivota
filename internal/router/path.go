package router

import (
	"fmt"
	"path"
	"strings"
)

func NewPath(p string) Path {
	return Path(cleanPath(p))
}

type Path string

func (path Path) Match(candidatePathFields []string) bool {
	pathFields := path.Fields()

	if len(pathFields) != len(candidatePathFields) {
		return false
	}

	for i, pathField := range pathFields {
		if isParam(pathField) {
			if candidatePathFields[i] == "" {
				return false
			}

			continue
		}

		if pathField != candidatePathFields[i] {
			return false
		}
	}

	return true
}

// Concatenates the passed elements with the current path string
func (path Path) Join(elem ...string) Path {
	elem = append([]string{string(path)}, elem...)
	// Join elements with `/` just in case passed strings are not cleaned
	return NewPath(strings.Join(elem, "/"))
}

func (path Path) Fields() []string {
	return strings.FieldsFunc(string(path), isPathSeparator)
}

// Points param names to their position in the path
type ParamsPos map[string]int

func (path Path) extractParamsPos() ParamsPos {
	paramsPos := ParamsPos{}
	for i, subpath := range path.Fields() {
		if isParam(subpath) {
			param := subpath[1:]
			if _, ok := paramsPos[param]; ok {
				panic(fmt.Sprintf("Duplicated param for path: '%s'", path))
			}
			paramsPos[param] = i
		}
	}

	return paramsPos
}

// ----------------------------------------------------------------------------
// Path helpers
// ----------------------------------------------------------------------------

type Params map[string]string

func isParam(subpath string) bool {
	return len(subpath) > 0 && subpath[0] == ':'
}

func isPathSeparator(r rune) bool {
	return r == '/'
}

// CleanPath returns the canonical path for p, eliminating . and .. elements and trailing slashes.
// Copyright 2009 The Go Authors. All rights reserved.
// Source: https://github.com/golang/go/blob/b8fd3cab3944d5dd5f2a50f3cc131b1048897ee1/src/net/http/http.go
func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	// path.Clean removes trailing slash except for root
	return path.Clean(p)
}
