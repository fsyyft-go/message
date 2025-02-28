// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package http

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/middleware"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/mux"
)

type (
	matcher interface {
		Use(ms ...middleware.Middleware)
		Add(selector string, ms ...middleware.Middleware)
		Match(operation string) []middleware.Middleware
	}

	// serverAccessor 是一个用于访问 kratos http.Server 内部字段的结构体。
	// 通过 unsafe.Pointer 转换实现对私有字段的访问。
	// 字段说明：
	// - Server: 底层的 http.Server 实例
	// - router: gorilla/mux 路由器实例
	// 其他字段使用空标识符 _ 占位，保持内存布局与原始结构体一致。
	serverAccessor struct {
		*http.Server
		_      net.Listener
		_      *tls.Config
		_      *url.URL
		_      error
		_      string
		_      string
		_      time.Duration
		_      []kratoshttp.FilterFunc
		_      matcher
		_      kratoshttp.DecodeRequestFunc
		_      kratoshttp.DecodeRequestFunc
		_      kratoshttp.DecodeRequestFunc
		_      kratoshttp.EncodeResponseFunc
		_      kratoshttp.EncodeErrorFunc
		_      bool
		router *mux.Router
	}

	RouteInfo struct {
		method string
		path   string
	}
)

// getRouter 从 kratos http.Server 中获取 mux.Router 实例。
// 参数 s 为 kratos http.Server 指针。
// 返回 gorilla/mux 路由器指针。
func getRouter(s *kratoshttp.Server) *mux.Router {
	sa := (*serverAccessor)(unsafe.Pointer(s))
	return sa.router
}

// GetPaths 获取 http server 中注册的所有路由路径。
// 参数 s 为 kratos http.Server 指针。
// 返回包含所有注册路由路径的字符串切片。
func GetPaths(s *kratoshttp.Server) []RouteInfo {
	routeInfos := make([]RouteInfo, 0)

	router := getRouter(s)
	_ = router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err != nil {
			return err
		}

		method, err := route.GetMethods()
		if err != nil {
			return err
		}

		for _, m := range method {
			routeInfos = append(routeInfos, RouteInfo{
				method: m,
				path:   path,
			})
		}

		return nil
	})
	return routeInfos
}

// Parse 将 kratos http.Server 中的路由注册到 gin.Engine 中。
// 参数 s 为 kratos http.Server 指针。
// 参数 e 为 gin.Engine 指针。
// 返回错误信息。
func Parse(s *kratoshttp.Server, e *gin.Engine) {
	routeInfos := GetPaths(s)

	for _, routeInfo := range routeInfos {
		e.Handle(routeInfo.method, routeInfo.path, func(c *gin.Context) {
			s.ServeHTTP(c.Writer, c.Request)
		})
	}
}
