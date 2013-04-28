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
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// Handler differs from http.Handler that it requires func
// SetPathParameters, which is used to pass in path parameters parsed
// from the named sub matches of path pattern.
type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	SetPathParameters(map[string]string)
}

type routeHandler struct {
	path    string
	pattern *regexp.Regexp
	handler Handler
}

func (route *routeHandler) parsePathParameters(url *url.URL) (
	kvpairs map[string]string) {
	pathstr := url.String()
	log.Printf("URL: '%s'", pathstr)
	pathstr = strings.TrimLeft(pathstr, route.path)
	kvpairs = make(map[string]string)
	match := route.pattern.FindStringSubmatch(pathstr)
	if match == nil {
		return
	}
	for i, name := range route.pattern.SubexpNames() {
		// ignore full match and unnamed submatch
		if i == 0 || name == "" {
			continue
		}
		kvpairs[name] = match[i]
	}
	return
}

func (route *routeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	kvpairs := route.parsePathParameters(r.URL)
	route.handler.SetPathParameters(kvpairs)
	route.handler.ServeHTTP(w, r)
}

// Handle acts like http.Handle except pattern must be a regular
// expression with named sub matches, while path acts just like the
// `pattern` argument of http.Handle .
func Handle(path string, pattern string, handler Handler) {
	r, err := regexp.Compile(pattern)
	if err != nil {
		panic(err)
	}
	route := routeHandler{path, r, handler}
	http.Handle(path, &route)
}

type wrapHandler struct {
	handle  func(http.ResponseWriter, *http.Request, map[string]string)
	kvpairs map[string]string
}

func (wh *wrapHandler) SetPathParameters(kvpairs map[string]string) {
	wh.kvpairs = kvpairs
}

func (wh *wrapHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wh.handle(w, r, wh.kvpairs)
}

// HandleFunc acts like like http.HandleFunc except one more argument
// of handle func is required to get the parsed path parameters.
func HandleFunc(path string, pattern string, handle func(
	http.ResponseWriter, *http.Request, map[string]string),) {
	handler := wrapHandler{handle, nil}
	Handle(path, pattern, &handler)
}
