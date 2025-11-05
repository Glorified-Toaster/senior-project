// Package server provides server-related functionalities.
package server

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Glorified-Toaster/senior-project/internal/config"
	"github.com/Glorified-Toaster/senior-project/internal/routers"
)

var shutdownTimeout int = 30 // seconds

type Server struct {
	router *routers.Router
	server *http.Server
}

// NewServer creates and returns a new Server instance.
func NewServer() *Server {
	// initialize the router
	router := routers.NewRouter()
	router.SetupRoutes()
	return &Server{
		router: router,
	}
}

// StartOverTLS defines and starts the HTTP server with TLS configuration.
func (s *Server) StartOverTLS(certFile, keyFile string) {
	// adding TLS configuration
	tlsConfig := &tls.Config{
		NextProtos: []string{"h2", "http/1.1"},
		MinVersion: tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
		},
		PreferServerCipherSuites: true,
	}

	// initialize the HTTP server
	s.server = &http.Server{
		Addr:      net.JoinHostPort(config.GlobalConfig.HTTPServer.Addr, config.GlobalConfig.HTTPServer.Port),
		Handler:   s.router.GetHandler(),
		TLSConfig: tlsConfig,
	}
	s.startWithGracefulShutdown(certFile, keyFile)
}

// startWithGracefulShutdown starts the server and sets up graceful shutdown handling.
func (s *Server) startWithGracefulShutdown(certFile, keyFile string) {
	go func() {
		var err error

		// check if certFile and keyFile are provided
		if certFile != "" && keyFile != "" {
			err = s.server.ListenAndServeTLS(certFile, keyFile)
		} else {
			log.Println("starting without TLS...")
			err = s.server.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start HTTP server: %v", err)
		}
	}()
	s.waitForShutdownSignal()
}

// waitForShutdownSignal waits for a shutdown signal and gracefully shuts down the server.
func (s *Server) waitForShutdownSignal() {
	exitChan := make(chan os.Signal, 1)

	// listen for interrupt signals In this case, os.Interrupt, syscall.SIGTERM, syscall.SIGINT
	signal.Notify(exitChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	exitSignal := <-exitChan // wait for a signal
	log.Printf("Received shutdown signal: %v", exitSignal)

	// create a context with timeout for the shutdown process
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(shutdownTimeout)*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP server Shutdown failed: %v", err)
	}

	log.Println("HTTP server exited gracefully...")
}

// GetRouter returns the server's router.
func (s *Server) GetRouter() *routers.Router {
	return s.router
}
