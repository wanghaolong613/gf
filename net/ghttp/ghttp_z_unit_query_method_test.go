// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

// import (
// 	"context"
// 	"fmt"
// 	"testing"
// 	"time"

// 	"github.com/gogf/gf/v2/frame/g"
// 	"github.com/gogf/gf/v2/net/ghttp"
// 	"github.com/gogf/gf/v2/test/gtest"
// 	"github.com/gogf/gf/v2/util/guid"
// )

// // Test_Server_Query_Route tests that QUERY method is supported in server routing.
// func Test_Server_Query_Route(t *testing.T) {
// 	s := g.Server(guid.S())
// 	// Register handler for QUERY method
// 	s.BindHandler("QUERY:/query-route", func(r *ghttp.Request) {
// 		r.Response.Writef("method=%s, name=%s", r.Method, r.Get("name"))
// 	})
// 	s.SetDumpRouterMap(false)
// 	s.Start()
// 	defer s.Shutdown()

// 	time.Sleep(100 * time.Millisecond)
// 	gtest.C(t, func(t *gtest.T) {
// 		url := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
// 		client := g.Client()
// 		client.SetPrefix(url)

// 		// Test QUERY method route
// 		result, err := client.DoRequest(context.TODO(), "QUERY", "/query-route", g.Map{"name": "test"})
// 		t.AssertNil(err)
// 		t.Assert(result.ReadAll(), "method=QUERY, name=test")
// 	})
// }

// // Test_Server_Query_WithQueryParams tests QUERY method with query string parameters.
// func Test_Server_Query_WithQueryParams(t *testing.T) {
// 	s := g.Server(guid.S())
// 	// Register handler that reads query params
// 	s.BindHandler("QUERY:/query-params", func(r *ghttp.Request) {
// 		// GetQuery should work for QUERY method too
// 		name := r.GetQuery("name")
// 		age := r.GetQuery("age")
// 		r.Response.Writef("name=%s, age=%s", name, age)
// 	})
// 	s.SetDumpRouterMap(false)
// 	s.Start()
// 	defer s.Shutdown()

// 	time.Sleep(100 * time.Millisecond)
// 	gtest.C(t, func(t *gtest.T) {
// 		url := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
// 		client := g.Client()
// 		client.SetPrefix(url)

// 		// Test QUERY with body params
// 		result, err := client.DoRequest(context.TODO(), "QUERY", "/query-params", g.Map{
// 			"name": "query-test",
// 			"age":  "30",
// 		})
// 		t.AssertNil(err)
// 		t.Assert(result.ReadAll(), "name=query-test, age=30")
// 	})
// }

// // Test_Server_Query_AcceptQuery tests the Accept-Query middleware functionality.
// func Test_Server_Query_AcceptQuery(t *testing.T) {
// 	s := g.Server(guid.S())
// 	// Apply Accept-Query middleware
// 	s.Use(ghttp.MiddlewareAcceptQuery("application/json", "application/xml"))

// 	s.BindHandler("QUERY:/accept-query", func(r *ghttp.Request) {
// 		// Verify Accept-Query header is set in response
// 		acceptQuery := r.Response.Header().Get(ghttp.AcceptQueryHeader)
// 		r.Response.Writef("accept-query=%s", acceptQuery)
// 	})
// 	s.SetDumpRouterMap(false)
// 	s.Start()
// 	defer s.Shutdown()

// 	time.Sleep(100 * time.Millisecond)
// 	gtest.C(t, func(t *gtest.T) {
// 		url := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
// 		client := g.Client()
// 		client.SetPrefix(url)

// 		// Send QUERY request
// 		result, err := client.DoRequest(context.TODO(), "QUERY", "/accept-query", g.Map{"test": "data"})
// 		t.AssertNil(err)
// 		t.Assert(result.ReadAll(), "accept-query=application/json, application/xml")
// 	})
// }

// // Test_Server_Query_SetAcceptQuery tests the SetAcceptQuery helper method.
// func Test_Server_Query_SetAcceptQuery(t *testing.T) {
// 	s := g.Server(guid.S())
// 	s.BindHandler("QUERY:/set-accept-query", func(r *ghttp.Request) {
// 		// Use SetAcceptQuery helper to set custom formats
// 		r.SetAcceptQuery("application/custom-format")
// 		r.Response.Write("ok")
// 	})
// 	s.SetDumpRouterMap(false)
// 	s.Start()
// 	defer s.Shutdown()

// 	time.Sleep(100 * time.Millisecond)
// 	gtest.C(t, func(t *gtest.T) {
// 		url := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
// 		client := g.Client()
// 		client.SetPrefix(url)

// 		// Send QUERY request
// 		result, err := client.DoRequest(context.TODO(), "QUERY", "/set-accept-query", nil)
// 		t.AssertNil(err)

// 		// Verify response header
// 		resp := result.Response
// 		t.Assert(resp.Header.Get(ghttp.AcceptQueryHeader), "application/custom-format")
// 	})
// }

// // Test_Server_Query_UnsupportedMethod tests that QUERY method returns 404 if not registered.
// func Test_Server_Query_UnsupportedMethod(t *testing.T) {
// 	s := g.Server(guid.S())
// 	// Only register GET handler, not QUERY (using explicit GET method)
// 	s.BindHandler("GET:/only-get", func(r *ghttp.Request) {
// 		r.Response.Write("get-only")
// 	})
// 	s.SetDumpRouterMap(false)
// 	s.Start()
// 	defer s.Shutdown()

// 	time.Sleep(100 * time.Millisecond)
// 	gtest.C(t, func(t *gtest.T) {
// 		url := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
// 		client := g.Client()
// 		client.SetPrefix(url)

// 		// QUERY request to a GET-only route should not match
// 		result, err := client.DoRequest(context.TODO(), "QUERY", "/only-get", nil)
// 		t.AssertNil(err)
// 		// QUERY method should not match GET route, so it returns 404
// 		t.Assert(result.StatusCode, 404)
// 	})
// }

// // Test_Server_Query_MultipleFormats tests QUERY with multiple Accept-Query formats.
// func Test_Server_Query_MultipleFormats(t *testing.T) {
// 	s := g.Server(guid.S())
// 	// Apply Accept-Query middleware with multiple formats
// 	s.Use(ghttp.MiddlewareAcceptQuery("application/json", "application/x-www-form-urlencoded"))

// 	s.BindHandler("QUERY:/multi-format", func(r *ghttp.Request) {
// 		r.Response.Write("query-success")
// 	})
// 	s.SetDumpRouterMap(false)
// 	s.Start()
// 	defer s.Shutdown()

// 	time.Sleep(100 * time.Millisecond)
// 	gtest.C(t, func(t *gtest.T) {
// 		url := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
// 		client := g.Client()
// 		client.SetPrefix(url)

// 		// Send QUERY request with JSON content type
// 		client.SetContentType("application/json")
// 		result, err := client.DoRequest(context.TODO(), "QUERY", "/multi-format", g.Map{"key": "value"})
// 		t.AssertNil(err)
// 		t.Assert(result.ReadAll(), "query-success")
// 	})
// }
