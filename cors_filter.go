package restful

// Copyright 2013 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// Type aliases for types defined elsewhere in the package
type Request struct {
	Request *http.Request
}

type Response struct {
	http.ResponseWriter
}

type FilterChain struct{}

func (f *FilterChain) ProcessFilter(req *Request, resp *Response) {}

// Container type
type Container struct{}

var DefaultContainer = &Container{}

func (c *Container) computeAllowedMethods(req *Request) []string {
	return []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
}

// Header constants
const (
	HEADER_Origin                        = "Origin"
	HEADER_AccessControlRequestMethod    = "Access-Control-Request-Method"
	HEADER_AccessControlRequestHeaders   = "Access-Control-Request-Headers"
	HEADER_AccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HEADER_AccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HEADER_AccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HEADER_AccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HEADER_AccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HEADER_AccessControlMaxAge           = "Access-Control-Max-Age"
)

// Tracing variables
var (
	trace       = false
	traceLogger interface {
		Print(v ...interface{})
		Printf(format string, v ...interface{})
	}
)

// CrossOriginResourceSharing is used to create a Container Filter that implements CORS.
// Cross-origin resource sharing (CORS) is a mechanism that allows JavaScript on a web page
// to make XMLHttpRequests to another domain, not the domain the JavaScript originated from.
//
// http://en.wikipedia.org/wiki/Cross-origin_resource_sharing
// http://enable-cors.org/server.html
// http://www.html5rocks.com/en/tutorials/cors/#toc-handling-a-not-so-simple-request
type CrossOriginResourceSharing struct {
	ExposeHeaders []string // list of Header names

	// AllowedHeaders is a list of Header names. Checking is case-insensitive.
	// The list may contain the special wildcard string ".*" ; all is allowed
	AllowedHeaders []string

	// AllowedDomains is a list of allowed values for Http Origin.
	// The list may contain the special wildcard string ".*" ; all is allowed
	// If empty all are allowed.
	AllowedDomains []string

	// AllowedDomainFunc is optional and is a function that will do the check
	// when the origin is not part of the AllowedDomains and it does not contain the wildcard ".*".
	AllowedDomainFunc func(origin string) bool

	// AllowedMethods is either empty or has a list of http methods names. Checking is case-insensitive.
	AllowedMethods []string
	MaxAge         int // number of seconds before requiring new Options request
	CookiesAllowed bool

	allowedOriginPatterns []*regexp.Regexp // internal field for origin regexp check.
}

// Filter is a filter function that implements the CORS flow as documented on http://enable-cors.org/server.html
// and http://www.html5rocks.com/static/images/cors_server_flowchart.png
func (c CrossOriginResourceSharing) Filter(req *Request, resp *Response, chain *FilterChain) {
	origin := req.Request.Header.Get(HEADER_Origin)
	if len(origin) == 0 {
		if trace {
			traceLogger.Print("no Http header Origin set")
		}
		chain.ProcessFilter(req, resp)
		return
	}
	if !c.isOriginAllowed(origin) { // check whether this origin is allowed
		if trace {
			traceLogger.Printf("HTTP Origin:%s is not part of %v, neither matches any part of %v", origin, c.AllowedDomains, c.allowedOriginPatterns)
		}
		chain.ProcessFilter(req, resp)
		return
	}
	if req.Request.Method != "OPTIONS" {
		c.doActualRequest(req, resp)
		chain.ProcessFilter(req, resp)
		return
	}
	if acrm := req.Request.Header.Get(HEADER_AccessControlRequestMethod); acrm != "" {
		c.doPreflightRequest(req, resp)
	} else {
		c.doActualRequest(req, resp)
		chain.ProcessFilter(req, resp)
		return
	}
}

func (c CrossOriginResourceSharing) doActualRequest(req *Request, resp *Response) {
	c.setOptionsHeaders(req, resp)
	// continue processing the response
}

func (c *CrossOriginResourceSharing) doPreflightRequest(req *Request, resp *Response) {
	if len(c.AllowedMethods) == 0 {
		c.AllowedMethods = DefaultContainer.computeAllowedMethods(req)
	}

	acrm := req.Request.Header.Get(HEADER_AccessControlRequestMethod)
	if !c.isValidAccessControlRequestMethod(acrm, c.AllowedMethods) {
		if trace {
			traceLogger.Printf("Http header %s:%s is not in %v",
				HEADER_AccessControlRequestMethod,
				acrm,
				c.AllowedMethods)
		}
		return
	}
	acrhs := req.Request.Header.Get(HEADER_AccessControlRequestHeaders)
	if len(acrhs) > 0 {
		for _, each := range strings.Split(acrhs, ",") {
			if !c.isValidAccessControlRequestHeader(strings.Trim(each, " ")) {
				if trace {
					traceLogger.Printf("Http header %s:%s is not in %v",
						HEADER_AccessControlRequestHeaders,
						acrhs,
						c.AllowedHeaders)
				}
				return
			}
		}
	}
	resp.Header().Add(HEADER_AccessControlAllowMethods, strings.Join(c.AllowedMethods, ","))
	resp.Header().Add(HEADER_AccessControlAllowHeaders, acrhs)
	c.setOptionsHeaders(req, resp)

	// return http 200 response, no body
}

func (c CrossOriginResourceSharing) setOptionsHeaders(req *Request, resp *Response) {
	c.checkAndSetExposeHeaders(resp)
	c.setAllowOriginHeader(req, resp)
	c.checkAndSetAllowCredentials(resp)
	if c.MaxAge > 0 {
		resp.Header().Add(HEADER_AccessControlMaxAge, strconv.Itoa(c.MaxAge))
	}
}

func (c CrossOriginResourceSharing) isOriginAllowed(origin string) bool {
	if len(origin) == 0 {
		return false
	}
	lowerOrigin := strings.ToLower(origin)
	if len(c.AllowedDomains) == 0 {
		if c.AllowedDomainFunc != nil {
			return c.AllowedDomainFunc(lowerOrigin)
		}
		return true
	}

	// exact match on each allowed domain
	for _, domain := range c.AllowedDomains {
		if domain == ".*" || strings.ToLower(domain) == lowerOrigin {
			return true
		}
	}
	if c.AllowedDomainFunc != nil {
		return c.AllowedDomainFunc(origin)
	}
	return false
}

func (c CrossOriginResourceSharing) setAllowOriginHeader(req *Request, resp *Response) {
	origin := req.Request.Header.Get(HEADER_Origin)
	if c.isOriginAllowed(origin) {
		resp.Header().Add(HEADER_AccessControlAllowOrigin, origin)
	}
}

func (c CrossOriginResourceSharing) checkAndSetExposeHeaders(resp *Response) {
	if len(c.ExposeHeaders) > 0 {
		resp.Header().Add(HEADER_AccessControlExposeHeaders, strings.Join(c.ExposeHeaders, ","))
	}
}

func (c CrossOriginResourceSharing) checkAndSetAllowCredentials(resp *Response) {
	if c.CookiesAllowed {
		resp.Header().Add(HEADER_AccessControlAllowCredentials, "true")
	}
}

func (c CrossOriginResourceSharing) isValidAccessControlRequestMethod(method string, allowedMethods []string) bool {
	for _, each := range allowedMethods {
		if each == method {
			return true
		}
	}
	return false
}

func (c CrossOriginResourceSharing) isValidAccessControlRequestHeader(header string) bool {
	for _, each := range c.AllowedHeaders {
		if strings.ToLower(each) == strings.ToLower(header) {
			return true
		}
		if each == "*" {
			return true
		}
	}
	return false
}
