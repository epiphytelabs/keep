package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/epiphytelabs/keep/pkg/cert"
	"github.com/miekg/dns"
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
	if !certificateExists() {
		if err := generateCertificate(); err != nil {
			return err
		}
	}

	dns.HandleFunc("keep.", resolve)

	ds := &dns.Server{Addr: ":53", Net: "udp"}
	go ds.ListenAndServe()

	return http.ListenAndServeTLS(addr, "/etc/keep/ssl/cert.pem", "/etc/keep/ssl/cert.key", s.router())
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
			h.ServeHTTP(w, r)
		} else {
			fmt.Println("404")
			w.WriteHeader(http.StatusNotFound)
		}
	})
}

func certificateExists() bool {
	if _, err := os.Stat("/etc/keep/ssl/cert.pem"); err != nil {
		return false
	}

	if _, err := os.Stat("/etc/keep/ssl/cert.key"); err != nil {
		return false
	}

	return true
}

func generateCertificate() error {
	c, err := cert.SelfSigned("*.app.keep")
	if err != nil {
		return err
	}

	pub, key, err := cert.Parts(c)
	if err != nil {
		return err
	}

	if err := os.MkdirAll("/etc/keep/ssl", 0755); err != nil {
		return err
	}

	if err := ioutil.WriteFile("/etc/keep/ssl/cert.pem", []byte(pub), 0644); err != nil {
		return err
	}

	if err := ioutil.WriteFile("/etc/keep/ssl/cert.key", []byte(key), 0600); err != nil {
		return err
	}

	return nil
}
