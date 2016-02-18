package httprouter

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const (
	CONNECT = "CONNECT"
	DELETE  = "DELETE"
	GET     = "GET"
	HEAD    = "HEAD"
	OPTIONS = "OPTIONS"
	PATCH   = "PATCH"
	POST    = "POST"
	PUT     = "PUT"
	TRACE   = "TRACE"
)

type route struct {
	patternRegexp *regexp.Regexp
	methodHandler map[string]http.HandlerFunc
	params        map[int]string
}

type RouteMux struct {
	routes []*route
}

func New() *RouteMux {
	return &RouteMux{}
}

// Get adds a new Route for GET requests.
func (m *RouteMux) Get(pattern string, handler http.HandlerFunc) {
	m.AddRoute(GET, pattern, handler)
}

// Put adds a new Route for PUT requests.
func (m *RouteMux) Put(pattern string, handler http.HandlerFunc) {
	m.AddRoute(PUT, pattern, handler)
}

// Del adds a new Route for DELETE requests.
func (m *RouteMux) Del(pattern string, handler http.HandlerFunc) {
	m.AddRoute(DELETE, pattern, handler)
}

// Patch adds a new Route for PATCH requests.
func (m *RouteMux) Patch(pattern string, handler http.HandlerFunc) {
	m.AddRoute(PATCH, pattern, handler)
}

// Post adds a new Route for POST requests.
func (m *RouteMux) Post(pattern string, handler http.HandlerFunc) {
	m.AddRoute(POST, pattern, handler)
}

// Adds a new Route to the Handler
func (m *RouteMux) AddRoute(method string, pattern string, handler http.HandlerFunc) {

	//split the url into sections
	parts := strings.Split(pattern, "/")
	//find params that start with ":"
	//replace with regular expressions
	params := make(map[int]string)
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			expr := "([^/]+)"
			//a user may choose to override the defult expression
			// similar to expressjs: ‘/user/:id([0-9]+)’
			if index := strings.Index(part, "("); index != -1 {
				expr = part[index:]
				part = part[:index]
			}
			index := i - 1
			params[index] = part
			parts[i] = expr
		}
	}
	//recreate the url pattern, with parameters replaced
	//by regular expressions. then compile the regex
	pattern = strings.Join(parts, "/")
	regex, regexErr := regexp.Compile(pattern)
	if regexErr != nil {
		//TODO add error handling here to avoid panic
		panic(regexErr)
		return
	}

	//now create the Route
	route := &route{methodHandler: make(map[string]http.HandlerFunc)}
	route.patternRegexp = regex
	route.methodHandler[method] = handler
	route.params = params

	//and finally append to the list of Routes
	m.routes = append(m.routes, route)
}

// Required by http.Handler interface. This method is invoked by the
// http server and will handle all page routing
func (m *RouteMux) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	requestPath := r.URL.Path

	//find a matching Route
	for _, route := range m.routes {

		//check if Route pattern matches url
		if !route.patternRegexp.MatchString(requestPath) {
			continue
		}

		//get submatches (params)
		matches := route.patternRegexp.FindStringSubmatch(requestPath)

		//double check that the Route matches the URL pattern.
		if len(matches[0]) != len(requestPath) {
			continue
		}

		if handler, err := route.methodHandler[r.Method]; !err {
			continue
		} else {
			if len(route.params) > 0 {
				//add url parameters to the query param map
				values := r.URL.Query()
				for i, match := range matches[1:] {
					fmt.Println(i, route.params[i], match)
					values.Add(route.params[i], match)
				}
				r.URL.RawQuery = url.Values(values).Encode()
			}
			//Invoke the request handler
			handler(rw, r)
			break
		}
	}

}
