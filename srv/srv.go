package srv

import (
	"net/http"

	"github.com/bugfan/de"
	"github.com/bugfan/rest"
	"github.com/bugfan/trojan-auth/env"
	"github.com/bugfan/trojan-auth/models"
	"github.com/bugfan/trojan-auth/srv/apis"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func init() {
	de.SetKey(env.Get("des_key"))
}

type route struct {
	httpMethod   string
	relativePath string
	handlers     []gin.HandlerFunc
}

var globalRoutes = make([]*route, 0)
var routes = make([]*route, 0)

func RegisterGlobal(method, path string, handlers ...gin.HandlerFunc) {
	for _, r := range globalRoutes {
		if r.httpMethod == method && r.relativePath == path {
			return
		}
	}

	r := &route{
		httpMethod:   method,
		relativePath: path,
		handlers:     handlers,
	}
	globalRoutes = append(globalRoutes, r)
}
func Register(method, path string, handlers ...gin.HandlerFunc) {
	for _, r := range routes {
		if r.httpMethod == method && r.relativePath == path {
			return
		}
	}

	r := &route{
		httpMethod:   method,
		relativePath: path,
		handlers:     handlers,
	}
	routes = append(routes, r)
}

type Router struct {
	G *gin.Engine
}

func NewAPIBackend() *Router {
	g := gin.Default()
	g.Use(gin.Recovery())
	g.Use(gin.ErrorLogger())

	api := g.Group("/")
	// global routes
	for _, r := range globalRoutes {
		api.Handle(r.httpMethod, r.relativePath, r.handlers...)
	}
	api.Use(apis.AuthMiddleware)

	// common routes
	for _, r := range routes {
		api.Handle(r.httpMethod, r.relativePath, r.handlers...)
	}

	rest.NewAPIBackend(api, models.GetEngine(), "")
	return &Router{
		G: g,
	}
}

func (s *Router) Run(addr string) {
	logrus.Infof("HTTP Server Listen On %v\n", addr)
	logrus.Fatal(s.G.Run(addr))
}
func (s *Router) RunTLS(addr, certFile, keyFile string) {
	logrus.Infof("HTTPS Server Listen On %v\n", addr)
	logrus.Fatal(s.G.RunTLS(addr, certFile, keyFile))
}
func (s *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.G.ServeHTTP(w, r)
}

func init() {
	Register(http.MethodGet, "/auth", apis.VerifyVPNRequest)
	Register(http.MethodPost, "/auth", apis.GenerateCredential)
}
