// Package server provides server-related functionalities.
package server

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Glorified-Toaster/senior-project/internal/config"
	"github.com/Glorified-Toaster/senior-project/internal/helpers"
	"github.com/Glorified-Toaster/senior-project/internal/routers"
	"github.com/Glorified-Toaster/senior-project/internal/utils"
	"go.uber.org/zap"
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
func (s *Server) StartOverTLS(cfg *config.Config) {
	certFile, keyFile := s.assignCertFile(cfg.HTTPServer.CertFile, cfg.HTTPServer.KeyFile, cfg)

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
				utils.LogInfo(utils.UseExistedTLSCert.Type, utils.UseExistedTLSCert.Msg)
				return certFile, keyFile
			}
		}
	}

	generatedCert, generatedKey, err := helpers.GenerateSelfSignedTLSCert(cfg.HTTPServer.Addr, cfg.HTTPServer.CertDir)
	if err != nil {
		utils.LogErrorWithLevel("fatal", utils.FailedToGenerateTLSCert.Type, utils.FailedToGenerateTLSCert.Code, utils.FailedToGenerateTLSCert.Msg, err)
	}

	utils.LogInfo(utils.GenerateTLSCertOK.Type, utils.GenerateTLSCertOK.Msg)
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
			utils.LogErrorWithLevel("warn", utils.FailedToStartWithTLS.Type, utils.FailedToStartWithTLS.Code, utils.FailedToStartWithTLS.Msg, errors.New("unable to get TLS cert"))
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
		utils.LogErrorWithLevel("fatal", utils.FailedToStartServer.Type, utils.FailedToStartServer.Code, utils.FailedToStartServer.Msg, err)

	case exitSignal := <-exitChan:
		utils.LogInfo(utils.ServerShutdownSignalOK.Type, utils.ServerShutdownSignalOK.Msg, zap.String("signal_type", exitSignal.String()))
		s.shutdownServer()
	}
}

func (s *Server) shutdownServer() {
	// create a context with timeout for the shutdown process
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(shutdownTimeout)*time.Second)
	defer cancel()

	utils.LogInfo(utils.ServerShutdown.Type, utils.ServerShutdown.Msg)

	if s.server != nil {
		if err := s.server.Shutdown(ctx); err != nil {
			utils.LogErrorWithLevel("fatal", utils.FailedToShutdownServer.Type, utils.FailedToShutdownServer.Code, utils.FailedToShutdownServer.Msg, err)
		}
		utils.LogInfo(utils.ServerShutdownOK.Type, utils.ServerShutdownOK.Msg)
	}
}

// GetRouter returns the server's router.
func (s *Server) GetRouter() *routers.Router {
	return s.router
}
