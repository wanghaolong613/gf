// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"mime"
	"net/http"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gstr"
)

const (
	// AcceptQueryHeader is the HTTP header name for advertising supported QUERY media types (RFC 10008).
	AcceptQueryHeader = "Accept-Query"
)

// DefaultAcceptQueryFormats is the default list of supported QUERY formats.
// By default, ghttp supports common content types for QUERY requests.
const DefaultAcceptQueryFormats = "application/json, application/x-www-form-urlencoded, multipart/form-data"

// MiddlewareAcceptQuery returns a middleware that validates QUERY requests and
// adds the Accept-Query header to responses. This implements the server-side
// part of RFC 10008, allowing servers to:
// 1. Validate QUERY request Content-Type against supported formats
// 2. Advertise which QUERY media types they support via Accept-Query header
//
// Parameters:
//   - formats: One or more supported media types. If empty, DefaultAcceptQueryFormats is used.
//     Multiple formats should be passed as separate arguments, e.g.:
//     MiddlewareAcceptQuery("application/json", "application/xml")
//     Or as a single comma-separated string:
//     MiddlewareAcceptQuery("application/json, application/xml")
func MiddlewareAcceptQuery(formats ...string) HandlerFunc {
	var acceptQuery string
	acceptQuery = gstr.JoinIgnoreEmpty(formats, ", ")
	if acceptQuery == "" {
		acceptQuery = DefaultAcceptQueryFormats
	}

	return func(r *Request) {
		// Advertise supported QUERY media types via Accept-Query header.
		// Per RFC 10008, this header can be included in any response so that
		// clients can discover supported formats before making a QUERY request.
		r.Response.Header().Set(AcceptQueryHeader, acceptQuery)

		// Validate QUERY method requests
		if r.Method == httpMethodQuery {
			if err := r.validateQueryRequest(acceptQuery); err != nil {
				r.Response.WriteStatus(http.StatusUnsupportedMediaType, err.Error())
				return
			}
		}
		r.Middleware.Next()
	}
}

// SetAcceptQuery sets the Accept-Query header for the current response.
// This is a helper function that can be called from handlers to set supported QUERY formats.
// Parameters:
//   - formats: One or more supported media types. If empty, DefaultAcceptQueryFormats is used.
func (r *Request) SetAcceptQuery(formats ...string) {
	acceptQuery := gstr.JoinIgnoreEmpty(formats, ", ")
	if acceptQuery == "" {
		acceptQuery = DefaultAcceptQueryFormats
	}
	r.Response.Header().Set(AcceptQueryHeader, acceptQuery)
}

// validateQueryRequest validates a QUERY request according to RFC 10008.
// It checks:
// 1. Content-Type header is required for QUERY requests
// 2. Content-Type must be in the server's supported formats list
func (r *Request) validateQueryRequest(supportedFormats string) error {
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		// RFC 10008 requires Content-Type for QUERY requests
		return gerror.NewCode(gcode.CodeInvalidParameter, "QUERY method requires Content-Type header")
	}

	// Parse the Content-Type to get the media type (ignore parameters like charset)
	mediaType := parseContentType(contentType)

	// Validate Content-Type against server's supported formats
	supported := false
	for _, format := range gstr.SplitAndTrim(supportedFormats, ",") {
		if matchMediaType(mediaType, format) {
			supported = true
			break
		}
	}
	if !supported {
		return gerror.NewCodef(
			gcode.CodeInvalidParameter,
			"Unsupported Media Type for QUERY request: %s. Supported formats: %s",
			mediaType, supportedFormats,
		)
	}

	return nil
}

// parseContentType extracts the media type from a Content-Type header value,
// ignoring parameters like charset.
func parseContentType(contentType string) string {
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return contentType
	}
	return mediaType
}

// matchMediaType checks if the Content-Type matches the supported media type.
// Supports wildcard patterns like "application/*" or "*/*".
// Media types are case-insensitive per RFC 7231.
func matchMediaType(contentType, supported string) bool {
	// Normalize to lowercase for case-insensitive comparison
	contentType = gstr.ToLower(contentType)
	supported = gstr.ToLower(supported)

	// Extract main type and sub type
	ctMain, ctSub := parseMediaType(contentType)
	supMain, supSub := parseMediaType(supported)

	// Check main type
	if supMain != "*" && supMain != ctMain {
		return false
	}

	// Check sub type
	if supSub != "*" && supSub != ctSub {
		return false
	}

	return true
}

// parseMediaType parses a media type string into main type and sub type.
func parseMediaType(mediaType string) (mainType, subType string) {
	for i := 0; i < len(mediaType); i++ {
		if mediaType[i] == '/' {
			return mediaType[:i], mediaType[i+1:]
		}
	}
	return mediaType, "*"
}
