// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

// Test_AcceptQuery_Middleware_ValidRequest tests that valid QUERY requests pass validation.
func Test_AcceptQuery_Middleware_ValidRequest(t *testing.T) {
	s := g.Server(guid.S())
	// Apply Accept-Query middleware with default formats
	s.Use(ghttp.MiddlewareAcceptQuery())

	s.BindHandler("QUERY:/query-valid", func(r *ghttp.Request) {
		r.Response.Write("query-success")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		url := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(url)

		// Send QUERY request with JSON content type (default supported)
		client.SetContentType("application/json")
		result, err := client.DoRequest(context.TODO(), "QUERY", "/query-valid", g.Map{"key": "value"})
		t.AssertNil(err)
		t.Assert(result.StatusCode, http.StatusOK)
		t.Assert(result.ReadAll(), "query-success")

		// Verify Accept-Query header in response
		t.Assert(result.Header.Get(ghttp.AcceptQueryHeader), ghttp.DefaultAcceptQueryFormats)
		result.Close()
	})
}

// Test_AcceptQuery_Middleware_MissingContentType tests that QUERY requests without Content-Type are rejected.
func Test_AcceptQuery_Middleware_MissingContentType(t *testing.T) {
	s := g.Server(guid.S())
	// Apply Accept-Query middleware
	s.Use(ghttp.MiddlewareAcceptQuery())

	s.BindHandler("QUERY:/query-no-content-type", func(r *ghttp.Request) {
		r.Response.Write("should-not-reach")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		url := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(url)

		// Send QUERY request without Content-Type header
		result, err := client.DoRequest(context.TODO(), "QUERY", "/query-no-content-type", "test-data")
		t.AssertNil(err)
		// Should return 415 Unsupported Media Type
		t.Assert(result.StatusCode, http.StatusUnsupportedMediaType)
		result.Close()
	})
}

// Test_AcceptQuery_Middleware_UnsupportedContentType tests that unsupported Content-Type is rejected.
func Test_AcceptQuery_Middleware_UnsupportedContentType(t *testing.T) {
	s := g.Server(guid.S())
	// Apply Accept-Query middleware with only JSON format
	s.Use(ghttp.MiddlewareAcceptQuery("application/json"))

	s.BindHandler("QUERY:/query-unsupported", func(r *ghttp.Request) {
		r.Response.Write("should-not-reach")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		url := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(url)

		// Send QUERY request with unsupported content type
		client.SetContentType("application/xml")
		result, err := client.DoRequest(context.TODO(), "QUERY", "/query-unsupported", "<test>data</test>")
		t.AssertNil(err)
		// Should return 415 Unsupported Media Type
		t.Assert(result.StatusCode, http.StatusUnsupportedMediaType)
		result.Close()
	})
}

// Test_AcceptQuery_Middleware_CustomFormats tests custom Accept-Query formats.
func Test_AcceptQuery_Middleware_CustomFormats(t *testing.T) {
	s := g.Server(guid.S())
	// Apply Accept-Query middleware with custom formats
	s.Use(ghttp.MiddlewareAcceptQuery("application/xml", "text/plain"))

	s.BindHandler("QUERY:/query-custom", func(r *ghttp.Request) {
		r.Response.Write("custom-format-success")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		url := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(url)

		// Send QUERY request with supported XML content type
		client.SetContentType("application/xml")
		result, err := client.DoRequest(context.TODO(), "QUERY", "/query-custom", "<data>test</data>")
		t.AssertNil(err)
		t.Assert(result.StatusCode, http.StatusOK)
		t.Assert(result.ReadAll(), "custom-format-success")

		// Verify Accept-Query header in response
		t.Assert(result.Header.Get(ghttp.AcceptQueryHeader), "application/xml, text/plain")
		result.Close()
	})
}

// Test_AcceptQuery_Middleware_WithCharset tests that Content-Type with charset parameter is handled correctly.
func Test_AcceptQuery_Middleware_WithCharset(t *testing.T) {
	s := g.Server(guid.S())
	s.Use(ghttp.MiddlewareAcceptQuery("application/json"))

	s.BindHandler("QUERY:/query-charset", func(r *ghttp.Request) {
		r.Response.Write("charset-success")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		url := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(url)

		// Send QUERY request with Content-Type containing charset parameter
		client.SetContentType("application/json; charset=utf-8")
		result, err := client.DoRequest(context.TODO(), "QUERY", "/query-charset", g.Map{"key": "value"})
		t.AssertNil(err)
		// Should be accepted because charset parameter should be ignored
		t.Assert(result.StatusCode, http.StatusOK)
		t.Assert(result.ReadAll(), "charset-success")
		result.Close()
	})
}

// Test_AcceptQuery_Middleware_Wildcard tests wildcard pattern matching in Accept-Query formats.
func Test_AcceptQuery_Middleware_Wildcard(t *testing.T) {
	s := g.Server(guid.S())
	// Apply Accept-Query middleware with wildcard
	s.Use(ghttp.MiddlewareAcceptQuery("application/*"))

	s.BindHandler("QUERY:/query-wildcard", func(r *ghttp.Request) {
		r.Response.Write("wildcard-success")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		url := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(url)

		// Send QUERY request with any application/* content type
		client.SetContentType("application/json")
		result, err := client.DoRequest(context.TODO(), "QUERY", "/query-wildcard", g.Map{"key": "value"})
		t.AssertNil(err)
		t.Assert(result.StatusCode, http.StatusOK)
		t.Assert(result.ReadAll(), "wildcard-success")

		// Also test with application/xml
		client.SetContentType("application/xml")
		result2, err := client.DoRequest(context.TODO(), "QUERY", "/query-wildcard", "<data>test</data>")
		t.AssertNil(err)
		t.Assert(result2.StatusCode, http.StatusOK)
		t.Assert(result2.ReadAll(), "wildcard-success")
		result2.Close()

		result.Close()
	})
}

// Test_AcceptQuery_Middleware_Wildcard_Double tests the "*/*" wildcard pattern
// which matches all media types.
func Test_AcceptQuery_Middleware_Wildcard_Double(t *testing.T) {
	s := g.Server(guid.S())
	s.Use(ghttp.MiddlewareAcceptQuery("*/*"))

	s.BindHandler("QUERY:/query-double-wildcard", func(r *ghttp.Request) {
		r.Response.Write("double-wildcard-success")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		url := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(url)

		// "*/*" should match application/json
		client.SetContentType("application/json")
		result, err := client.DoRequest(context.TODO(), "QUERY", "/query-double-wildcard", g.Map{"key": "value"})
		t.AssertNil(err)
		t.Assert(result.StatusCode, http.StatusOK)
		t.Assert(result.ReadAll(), "double-wildcard-success")
		result.Close()

		// "*/*" should match text/plain
		client.SetContentType("text/plain")
		result, err = client.DoRequest(context.TODO(), "QUERY", "/query-double-wildcard", "plain-text")
		t.AssertNil(err)
		t.Assert(result.StatusCode, http.StatusOK)
		t.Assert(result.ReadAll(), "double-wildcard-success")
		result.Close()

		// "*/*" should match application/xml
		client.SetContentType("application/xml")
		result, err = client.DoRequest(context.TODO(), "QUERY", "/query-double-wildcard", "<data>test</data>")
		t.AssertNil(err)
		t.Assert(result.StatusCode, http.StatusOK)
		t.Assert(result.ReadAll(), "double-wildcard-success")
		result.Close()
	})
}

// Test_AcceptQuery_Middleware_Wildcard_Reject tests that a wildcard pattern
// rejects media types that do not match the pattern.
func Test_AcceptQuery_Middleware_Wildcard_Reject(t *testing.T) {
	s := g.Server(guid.S())
	// Only application/* is supported
	s.Use(ghttp.MiddlewareAcceptQuery("application/*"))

	s.BindHandler("QUERY:/query-wildcard-reject", func(r *ghttp.Request) {
		r.Response.Write("should-not-reach")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		url := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(url)

		// text/plain should NOT match application/*
		client.SetContentType("text/plain")
		result, err := client.DoRequest(context.TODO(), "QUERY", "/query-wildcard-reject", "plain-text")
		t.AssertNil(err)
		t.Assert(result.StatusCode, http.StatusUnsupportedMediaType)
		result.Close()

		// text/html should NOT match application/*
		client.SetContentType("text/html")
		result, err = client.DoRequest(context.TODO(), "QUERY", "/query-wildcard-reject", "<html></html>")
		t.AssertNil(err)
		t.Assert(result.StatusCode, http.StatusUnsupportedMediaType)
		result.Close()
	})
}

// Test_AcceptQuery_Middleware_Wildcard_Multiple tests multiple wildcard patterns.
func Test_AcceptQuery_Middleware_Wildcard_Multiple(t *testing.T) {
	s := g.Server(guid.S())
	// Support both application/* and text/*
	s.Use(ghttp.MiddlewareAcceptQuery("application/*", "text/*"))

	s.BindHandler("QUERY:/query-multi-wildcard", func(r *ghttp.Request) {
		r.Response.Write("multi-wildcard-success")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		url := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(url)

		// application/json matches application/*
		client.SetContentType("application/json")
		result, err := client.DoRequest(context.TODO(), "QUERY", "/query-multi-wildcard", g.Map{"key": "value"})
		t.AssertNil(err)
		t.Assert(result.StatusCode, http.StatusOK)
		t.Assert(result.ReadAll(), "multi-wildcard-success")
		result.Close()

		// text/plain matches text/*
		client.SetContentType("text/plain")
		result, err = client.DoRequest(context.TODO(), "QUERY", "/query-multi-wildcard", "plain-text")
		t.AssertNil(err)
		t.Assert(result.StatusCode, http.StatusOK)
		t.Assert(result.ReadAll(), "multi-wildcard-success")
		result.Close()

		// image/png should NOT match application/* or text/*
		client.SetContentType("image/png")
		result, err = client.DoRequest(context.TODO(), "QUERY", "/query-multi-wildcard", "binary-data")
		t.AssertNil(err)
		t.Assert(result.StatusCode, http.StatusUnsupportedMediaType)
		result.Close()
	})
}

// Test_AcceptQuery_Middleware_CaseInsensitive tests that media type matching
// is case-insensitive per RFC 7231.
func Test_AcceptQuery_Middleware_CaseInsensitive(t *testing.T) {
	s := g.Server(guid.S())
	// Configure with uppercase format
	s.Use(ghttp.MiddlewareAcceptQuery("Application/JSON"))

	s.BindHandler("QUERY:/query-case-insensitive", func(r *ghttp.Request) {
		r.Response.Write("case-insensitive-success")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		url := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(url)

		// Lowercase content type should match uppercase supported format
		client.SetContentType("application/json")
		result, err := client.DoRequest(context.TODO(), "QUERY", "/query-case-insensitive", g.Map{"key": "value"})
		t.AssertNil(err)
		t.Assert(result.StatusCode, http.StatusOK)
		t.Assert(result.ReadAll(), "case-insensitive-success")
		result.Close()

		// Mixed-case content type should also match
		client.SetContentType("Application/JSON")
		result, err = client.DoRequest(context.TODO(), "QUERY", "/query-case-insensitive", g.Map{"key": "value"})
		t.AssertNil(err)
		t.Assert(result.StatusCode, http.StatusOK)
		t.Assert(result.ReadAll(), "case-insensitive-success")
		result.Close()
	})
}

// Test_AcceptQuery_Middleware_Wildcard_CaseInsensitive tests that wildcard patterns
// work with case-insensitive matching.
func Test_AcceptQuery_Middleware_Wildcard_CaseInsensitive(t *testing.T) {
	s := g.Server(guid.S())
	// Configure with uppercase wildcard
	s.Use(ghttp.MiddlewareAcceptQuery("Application/*"))

	s.BindHandler("QUERY:/query-wildcard-case", func(r *ghttp.Request) {
		r.Response.Write("wildcard-case-success")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		url := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(url)

		// Lowercase content type should match uppercase wildcard
		client.SetContentType("application/json")
		result, err := client.DoRequest(context.TODO(), "QUERY", "/query-wildcard-case", g.Map{"key": "value"})
		t.AssertNil(err)
		t.Assert(result.StatusCode, http.StatusOK)
		t.Assert(result.ReadAll(), "wildcard-case-success")
		result.Close()

		// Mixed-case content type should also match
		client.SetContentType("Application/XML")
		result, err = client.DoRequest(context.TODO(), "QUERY", "/query-wildcard-case", "<data>test</data>")
		t.AssertNil(err)
		t.Assert(result.StatusCode, http.StatusOK)
		t.Assert(result.ReadAll(), "wildcard-case-success")
		result.Close()
	})
}

// Test_AcceptQuery_Middleware_NonQueryMethod tests that non-QUERY methods are not
// affected by QUERY validation, but still receive the Accept-Query header per RFC 10008.
func Test_AcceptQuery_Middleware_NonQueryMethod(t *testing.T) {
	s := g.Server(guid.S())
	s.Use(ghttp.MiddlewareAcceptQuery("application/json"))

	s.BindHandler("GET:/get-test", func(r *ghttp.Request) {
		r.Response.Write("get-success")
	})
	s.BindHandler("POST:/post-test", func(r *ghttp.Request) {
		r.Response.Write("post-success")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		url := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(url)

		// GET request should work normally without Accept-Query validation
		result, err := client.Get(context.TODO(), "/get-test")
		t.AssertNil(err)
		t.Assert(result.StatusCode, http.StatusOK)
		t.Assert(result.ReadAll(), "get-success")
		// Per RFC 10008, Accept-Query header is advertised on all responses
		t.Assert(result.Header.Get(ghttp.AcceptQueryHeader), "application/json")
		result.Close()

		// POST request should work normally without Accept-Query validation
		result, err = client.Post(context.TODO(), "/post-test", g.Map{"key": "value"})
		t.AssertNil(err)
		t.Assert(result.StatusCode, http.StatusOK)
		t.Assert(result.ReadAll(), "post-success")
		// Per RFC 10008, Accept-Query header is advertised on all responses
		t.Assert(result.Header.Get(ghttp.AcceptQueryHeader), "application/json")
		result.Close()
	})
}

// Test_AcceptQuery_Middleware_SetAcceptQuery tests the SetAcceptQuery helper method.
func Test_AcceptQuery_Middleware_SetAcceptQuery(t *testing.T) {
	s := g.Server(guid.S())
	// No middleware applied, handler uses SetAcceptQuery directly
	s.BindHandler("QUERY:/set-accept-query", func(r *ghttp.Request) {
		r.SetAcceptQuery("application/custom-type")
		r.Response.Write("ok")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		url := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(url)

		// Send QUERY request with any content type
		client.SetContentType("application/json")
		result, err := client.DoRequest(context.TODO(), "QUERY", "/set-accept-query", g.Map{"test": "data"})
		t.AssertNil(err)
		t.Assert(result.StatusCode, http.StatusOK)
		t.Assert(result.ReadAll(), "ok")

		// Verify Accept-Query header is set via SetAcceptQuery
		t.Assert(result.Header.Get(ghttp.AcceptQueryHeader), "application/custom-type")
		result.Close()
	})
}

// Test_AcceptQuery_Middleware_ErrorResponseHasHeader tests that 415 error responses
// include the Accept-Query header so clients can discover supported formats.
func Test_AcceptQuery_Middleware_ErrorResponseHasHeader(t *testing.T) {
	s := g.Server(guid.S())
	s.Use(ghttp.MiddlewareAcceptQuery("application/json"))

	s.BindHandler("QUERY:/query-error-header", func(r *ghttp.Request) {
		r.Response.Write("should-not-reach")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		url := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(url)

		// Send QUERY request with unsupported content type
		client.SetContentType("application/xml")
		result, err := client.DoRequest(context.TODO(), "QUERY", "/query-error-header", "<data>test</data>")
		t.AssertNil(err)
		t.Assert(result.StatusCode, http.StatusUnsupportedMediaType)
		// Accept-Query header should be present even in error responses
		t.Assert(result.Header.Get(ghttp.AcceptQueryHeader), "application/json")
		result.Close()
	})
}

// Test_AcceptQuery_Middleware_SetAcceptQuery_MultipleFormats tests SetAcceptQuery
// with multiple format arguments.
func Test_AcceptQuery_Middleware_SetAcceptQuery_MultipleFormats(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("QUERY:/set-multi-formats", func(r *ghttp.Request) {
		r.SetAcceptQuery("application/json", "application/xml", "text/plain")
		r.Response.Write("ok")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		url := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(url)

		client.SetContentType("application/json")
		result, err := client.DoRequest(context.TODO(), "QUERY", "/set-multi-formats", g.Map{"test": "data"})
		t.AssertNil(err)
		t.Assert(result.StatusCode, http.StatusOK)
		// All formats should be joined with ", "
		t.Assert(result.Header.Get(ghttp.AcceptQueryHeader), "application/json, application/xml, text/plain")
		result.Close()
	})
}
