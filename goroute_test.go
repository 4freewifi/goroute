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

package goroute

import (
	"fmt"
	"net/http"
	"testing"
)

type MySrv struct {
	kvpairs map[string]string
}

func (m *MySrv) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!\n", m.kvpairs["userid"])
}

func (m *MySrv) SetPathParameters(kvpairs map[string]string) {
	m.kvpairs = kvpairs
}

func TestRouteHandler(t *testing.T) {
	srv := MySrv{nil}
	Handle("/", `users/(?P<userid>[^/]+)/?`, &srv)
	fmt.Println("Try visit http://localhost:8080/users/john")
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		panic(err)
	}
	return
}
