// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gclient_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

// Test_Query_Method tests the HTTP QUERY method (RFC 10008) in gclient.
func Test_Query_Method(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/query-test", func(r *ghttp.Request) {
		// QUERY method should carry body data like POST
		r.Response.Writef("method=%s, name=%s, age=%s", r.Method, r.Get("name"), r.Get("age"))
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		url := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(url)

		// Test QUERY method with body data
		result, err := client.Query(ctx, "/query-test", g.Map{
			"name": "test",
			"age":  "25",
		})
		t.AssertNil(err)
		t.Assert(result.ReadAll(), "method=QUERY, name=test, age=25")
	})
}

// Test_Query_Method_ContentTypes tests QUERY method with different content types.
func Test_Query_Method_ContentTypes(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/query-json", func(r *ghttp.Request) {
		r.Response.Writef("method=%s, content-type=%s, body=%s",
			r.Method,
			r.Header.Get("Content-Type"),
			r.Get("name"))
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		url := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(url)

		// Test QUERY with JSON content type
		client.SetContentType("application/json")
		result, err := client.Query(ctx, "/query-json", g.Map{"name": "json-test"})
		t.AssertNil(err)
		t.Assert(result.ReadAll(), "method=QUERY, content-type=application/json, body=json-test")
	})
}

// Test_Query_Method_With_AcceptQuery tests QUERY method with Accept-Query middleware.
func Test_Query_Method_With_AcceptQuery(t *testing.T) {
	s := g.Server(guid.S())
	// Apply Accept-Query middleware
	s.Use(ghttp.MiddlewareAcceptQuery("application/json"))

	s.BindHandler("/query-accept", func(r *ghttp.Request) {
		// Verify Accept-Query header is set in response
		r.Response.Writef("accept-query=%s", r.Response.Header().Get("Accept-Query"))
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		url := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(url)

		// Send QUERY request
		result, err := client.Query(ctx, "/query-accept", g.Map{"test": "data"})
		t.AssertNil(err)
		// Verify Accept-Query header in response
		t.Assert(result.ReadAll(), "accept-query=application/json")
	})
}

// Test_Query_Method_FromClient tests the Query method from gclient directly.
func Test_Query_Method_FromClient(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/client-query", func(r *ghttp.Request) {
		r.Response.Writef("method=%s, data=%s", r.Method, r.Get("key"))
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		url := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := gclient.New()
		client.SetPrefix(url)

		// Use the Query method directly
		result, err := client.Query(context.Background(), "/client-query", map[string]string{
			"key": "value",
		})
		t.AssertNil(err)
		t.Assert(result.ReadAll(), "method=QUERY, data=value")
	})
}
