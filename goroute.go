// Copyright 2013 John Lee <john@0xlab.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// goroute is a very simple URL router based on named submatches of
// regular expression that works well with http.Handler
package goroute

import (
	"container/list"
	"log"
	"net/http"
	"regexp"
)

// Handler differs from http.Handler in that it requires an extra
// argument to pass in path parameters parsed from the named sub
// matches of path pattern.
type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request, map[string]string)
}

type patternHandler struct {
	Regexp  *regexp.Regexp
	Handler Handler
}

// RouteHandler stores patterns and matching handlers of a path.
type RouteHandler struct {
	path            string
	patternHandlers *list.List
}

// ServeHTTP parses the path parameters, calls SetPathParameters of
// the corresponding hander, then directs traffic to it.
func (r *RouteHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	pathstr := req.URL.String()
	log.Printf("%s: '%s'", req.Method, pathstr)
	pathstr = pathstr[len(r.path):]
	var match []string
	var pattern *regexp.Regexp
	var handler Handler
	for e := r.patternHandlers.Front(); e != nil; e = e.Next() {
		h := e.Value.(*patternHandler)
		pattern = h.Regexp
		handler = h.Handler
		match = pattern.FindStringSubmatch(pathstr)
		if match != nil {
			break
		}
	}
	if match == nil {
		log.Printf("Cannot find a matching pattern for Path `%s'",
			pathstr)
		http.NotFound(w, req)
		return
	}
	kvpairs := make(map[string]string)
	for i, name := range pattern.SubexpNames() {
		// ignore full match and unnamed submatch
		if i == 0 || name == "" {
			continue
		}
		kvpairs[name] = match[i]
	}
	log.Printf("Parsed path parameters: %s", kvpairs)
	handler.ServeHTTP(w, req, kvpairs)
}

// AddPatternHandler adds an additional pair of pattern and handler into
// RouteHandler. Last added will be matched first.
func (r *RouteHandler) AddPatternHandler(pattern string, handler Handler) {
	reg, err := regexp.Compile(pattern)
	if err != nil {
		panic(err)
	}
	r.patternHandlers.PushFront(&patternHandler{reg, handler})
}

// Handle acts like http.Handle except it requires one more argument:
// pattern. `pattern' must be a regular expression with named sub
// matches, while path acts just like the `pattern' argument of
// http.Handle .
func Handle(path string, pattern string, handler Handler) (r *RouteHandler) {
	r = &RouteHandler{path, list.New()}
	r.AddPatternHandler(pattern, handler)
	http.Handle(path, r)
	return
}

type wrapHandler struct {
	handle func(http.ResponseWriter, *http.Request, map[string]string)
}

func (wh *wrapHandler) ServeHTTP(w http.ResponseWriter, r *http.Request,
	kvpairs map[string]string) {
	wh.handle(w, r, kvpairs)
}

// HandleFunc acts like like http.HandleFunc except one more argument
// of handle func is required to get the parsed path parameters.
func HandleFunc(path string, pattern string, handle func(
	http.ResponseWriter, *http.Request, map[string]string)) (
	r *RouteHandler) {
	handler := &wrapHandler{handle: handle}
	r = Handle(path, pattern, handler)
	return
}

// AddPatternHandlerFunc adds an additional pair of pattern and
// handler function into RouteHandler. Last added will be matched
// first.
func (r *RouteHandler) AddPatternHandlerFunc(pattern string, handle func(
	http.ResponseWriter, *http.Request, map[string]string)) {
	handler := &wrapHandler{handle: handle}
	r.AddPatternHandler(pattern, handler)
}

// ServeMux acts just like http.ServeMux except it accepts `pattern'
// as an extra argument just like Handle and HandleFunc
type ServeMux struct {
	HTTPServeMux *http.ServeMux
}

func NewServeMux() *ServeMux {
	return &ServeMux{http.NewServeMux()}
}

func (mux *ServeMux) Handle(path string, pattern string, handler Handler) (
	r *RouteHandler) {
	r = &RouteHandler{path, list.New()}
	r.AddPatternHandler(pattern, handler)
	mux.HTTPServeMux.Handle(path, r)
	return
}

func (mux *ServeMux) HandleFunc(path string, pattern string, handle func(
	http.ResponseWriter, *http.Request, map[string]string)) (
	r *RouteHandler) {
	handler := &wrapHandler{handle: handle}
	r = mux.Handle(path, pattern, handler)
	return
}

func (mux *ServeMux) Handler(r *http.Request) (h http.Handler, pattern string) {
	h, pattern = mux.HTTPServeMux.Handler(r)
	return
}

func (mux *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux.HTTPServeMux.ServeHTTP(w, r)
}
