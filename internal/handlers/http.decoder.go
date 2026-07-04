package handlers

import "github.com/gorilla/schema"

// muxDecoder decodes URL query parameters into the request structs of the
// public read endpoints, following their `schema` tags.
var muxDecoder = schema.NewDecoder()
