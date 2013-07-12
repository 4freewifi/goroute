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
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MrFriendly struct {
}

func (m *MrFriendly) ServeHTTP(w http.ResponseWriter, r *http.Request,
	kvpairs map[string]string) {
	for k, v := range kvpairs {
		switch k {
		case "userid":
			fmt.Fprintf(w, "Hello, %s!", v)
		case "sitename":
			fmt.Fprintf(w, "Welcome to %s!", v)
		}
	}
}

func expect(t *testing.T, url string, expect string) {
	resp, err := http.Get(url)
	if err != nil {
		t.Error(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	bodytext := string(body)
	if bodytext != expect {
		t.Errorf("Expect: `%s', got: `%s'", expect, bodytext)
		return
	}
	log.Printf("Got expected response: `%s'", expect)
}

func TestRouteHandler(t *testing.T) {
	mr := MrFriendly{}
	r := Handle("/", `users/(?P<userid>[^/]+)/?`, &mr)
	r.AddPatternHandler(`sites/(?P<sitename>[^/]+)/?`, &mr)
	s := httptest.NewServer(r)
	defer s.Close()
	expect(t, s.URL+"/users/John", "Hello, John!")
	expect(t, s.URL+"/sites/Taipei", "Welcome to Taipei!")
}
