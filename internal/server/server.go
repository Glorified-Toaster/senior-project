// Package server provides server-related functionalities.
package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Glorified-Toaster/senior-project/internal/config"
	"github.com/Glorified-Toaster/senior-project/internal/helpers"
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
	// get the config
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("failed to get config: %v", err)
	}

	certFile, keyFile = s.assignCertFile(certFile, keyFile, cfg)

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
		Addr:      net.JoinHostPort(cfg.HTTPServer.Addr, cfg.HTTPServer.Port),
		Handler:   s.router.GetHandler(),
		TLSConfig: tlsConfig,
	}
	s.startWithGracefulShutdown(certFile, keyFile)
}

// assignCertFile checks for existing cert and key files or generates new ones if not found.
func (s *Server) assignCertFile(certFile, keyFile string, cfg *config.Config) (string, string) {
	// check if certFile and keyFile are provided, else check for existing files
	if certFile != "" && keyFile != "" {
		if _, err := os.Stat(certFile); err == nil {
			if _, err := os.Stat(keyFile); err == nil {
				log.Println("Using an existing TLS certificate...")
				return certFile, keyFile
			}
		}
	}
	log.Println("Generating self-signed TLS certificate...")

	generatedCert, generatedKey, err := helpers.GenerateSelfSignedTLSCert(cfg.HTTPServer.Addr, cfg.HTTPServer.CertDir)
	if err != nil {
		log.Fatalf("failed to generate self-signed TLS certificate: %v", err)
	}

	log.Println("Self-signed TLS certificate generated successfully")
	return generatedCert, generatedKey
}

// startWithGracefulShutdown starts the server and sets up graceful shutdown handling.
func (s *Server) startWithGracefulShutdown(certFile, keyFile string) {
	errChan := make(chan error, 1)
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
			errChan <- fmt.Errorf("failed to start HTTP server: %v", err)
		}
	}()
	s.waitForShutdownSignal(errChan)
}

// waitForShutdownSignal : waits for a shutdown signal and gracefully shuts down the server.
func (s *Server) waitForShutdownSignal(errChan chan error) {
	exitChan := make(chan os.Signal, 1)

	// listen for interrupt signals In this case, os.Interrupt, syscall.SIGTERM, syscall.SIGINT
	signal.Notify(exitChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-errChan:
		log.Fatalf("failed to start the HTTP server: %v", err)
	case exitSignal := <-exitChan:
		log.Printf("Received shutdown signal: %v", exitSignal)
		s.shutdownServer()
	}
}

func (s *Server) shutdownServer() {
	// create a context with timeout for the shutdown process
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(shutdownTimeout)*time.Second)
	defer cancel()

	log.Println("Shutting down server...")
	if s.server != nil {
		if err := s.server.Shutdown(ctx); err != nil {
			log.Fatalf("failed to shutdown the HTTP server: %v", err)
		}
		log.Println("HTTP server shut down gracefully...")
	}
}

// GetRouter returns the server's router.
func (s *Server) GetRouter() *routers.Router {
	return s.router
}
