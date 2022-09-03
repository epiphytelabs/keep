package server

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

type Routes map[string]http.Handler

type Server struct {
	lock   sync.Mutex
	routes Routes
}

func New(routes Routes) *Server {
	s := &Server{}
	s.UpdateRoutes(routes)
	return s
}

func (s *Server) Listen(addr string) error {
	return http.ListenAndServe(addr, s.router())
}

func (s *Server) UpdateRoutes(routes Routes) {
	s.lock.Lock()
	defer s.lock.Unlock()

	fmt.Println("updating routes")

	for k := range routes {
		fmt.Printf("  %s\n", k)
	}

	s.routes = routes
}

func (s *Server) route(host string) http.Handler {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.routes[host]
}

func (s *Server) router() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host := strings.Split(r.Host, ":")[0]

		if h := s.route(host); h != nil {
			fmt.Printf("host: %+v\n", host)
			h.ServeHTTP(w, r)
		} else {
			fmt.Println("404")
			w.WriteHeader(http.StatusNotFound)
		}
	})
}
