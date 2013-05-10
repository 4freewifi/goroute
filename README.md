# goroute

A very simple URL router for net/http package based on regular
expression.

## Why another

The idea of `http.Handler` is great because it is just an interface so
the developer can keep all the necessary info in its own
context. However, most of the web framework or routing package out
there only support `http.HandleFunc`, which means the context must be
kept elsewhere. For example, consider the following program:

```go
package main

import (
	"fmt"
	"net/http"
)

type MySrv struct {
	hello string
}

func (m *MySrv) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", m.hello)
}

func main() {
	srv := MySrv{"World"}
	http.Handle("/", &srv)
	http.ListenAndServe("localhost:8080", nil)
}
```

In case of [web.go](https://github.com/hoisie/web), if keeping the
original context `MySrv` is required, this program will become
something like:

```go
package main

import (
	"fmt"
	"net/http"
	"github.com/hoisie/web"
)

type MySrv struct {
	hello string
}

func (m *MySrv) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!\n", m.hello)
}

func main() {
	srv := MySrv{"World"}
	hello := func(ctx *web.Context, val string) {
		for k,v := range ctx.Params {
			fmt.Fprintln(ctx, k, v)
		}
		srv.ServeHTTP(ctx, ctx.Request)
	}
	web.Get("/(.*)", hello)
	web.Run("localhost:8080")
}
```

Because of function closure, `srv MySrv` is implicitly included in
`hello`, and modifications made to `srv` will implicitly affect
`hello`. *Explicit is better than implicit*. Many other projects
suffer from similiar issue. However, the requirement of a simple URL
path router still exists, e.g. in case of RESTful like `GET
/users/<userid>`. A very thin layer that handles just routing and
doesn't get in the way of the original structure of net/http will be
handy.

## Example

```go
package main

import (
	"fmt"
	"net/http"
	"github.com/johncylee/goroute"
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

func main() {
	srv := MySrv{nil}
	goroute.Handle("/", `users/(?P<userid>[^/]+)/?`, &srv)
	fmt.Println("try visit http://localhost:8080/users/john")
	http.ListenAndServe("localhost:8080", nil)
}
```

Comparing to the original http.Handler, a new function
`SetPathParameters(map[string]string)` is required to pass in the
key-value pairs parsed from the named submatches of the regular
expression `users/(?P<userid>[^/]+)/?`. Besides that, `goroute.Handle`
acts the same as `http.Handle`.

## API

Visit <http://godoc.org/github.com/johncylee/goroute>
