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

	"uot-exam/internal/adapters/inbound/http/handler"
	"uot-exam/internal/adapters/inbound/http/helpers"
	"uot-exam/internal/adapters/inbound/http/middleware"
	"uot-exam/internal/adapters/inbound/http/router"
	"uot-exam/internal/adapters/outbound/config"
	"uot-exam/internal/adapters/outbound/logger"

	"go.uber.org/zap"
)

const shutdownTimeout int = 30 // seconds

type Server struct {
	router *router.Router
	server *http.Server
	logger *logger.Logger
}

// NewServer creates and returns a new Server instance.
func NewServer(ctrl *handler.Handler, authMiddleware *middleware.AuthMiddleware, logger *logger.Logger) *Server {
	// initialize the router
	router := router.NewRouter(ctrl, authMiddleware)
	router.SetupRoutes()
	return &Server{
		router: router,
		logger: logger,
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
				s.logger.LogInfo(logger.UseExistedTLSCert.Type, logger.UseExistedTLSCert.Msg)
				return certFile, keyFile
			}
		}
	}

	generatedCert, generatedKey, err := helpers.GenerateSelfSignedTLSCert(cfg.HTTPServer.Addr, cfg.HTTPServer.CertDir)
	if err != nil {
		s.logger.LogErrorWithLevel("fatal", logger.FailedToGenerateTLSCert.Type, logger.FailedToGenerateTLSCert.Code, logger.FailedToGenerateTLSCert.Msg, err)
	}

	s.logger.LogInfo(logger.GenerateTLSCertOK.Type, logger.GenerateTLSCertOK.Msg)
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
			s.logger.LogErrorWithLevel("warn", logger.FailedToStartWithTLS.Type, logger.FailedToStartWithTLS.Code, logger.FailedToStartWithTLS.Msg, errors.New("unable to get TLS cert"))
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
		s.logger.LogErrorWithLevel("fatal", logger.FailedToStartServer.Type, logger.FailedToStartServer.Code, logger.FailedToStartServer.Msg, err)

	case exitSignal := <-exitChan:
		s.logger.LogInfo(logger.ServerShutdownSignalOK.Type, logger.ServerShutdownSignalOK.Msg, zap.String("signal_type", exitSignal.String()))
		s.shutdownServer()
	}
}

func (s *Server) shutdownServer() {
	// create a context with timeout for the shutdown process
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(shutdownTimeout)*time.Second)
	defer cancel()

	s.logger.LogInfo(logger.ServerShutdown.Type, logger.ServerShutdown.Msg)

	if s.server != nil {
		if err := s.server.Shutdown(ctx); err != nil {
			s.logger.LogErrorWithLevel("fatal", logger.FailedToShutdownServer.Type, logger.FailedToShutdownServer.Code, logger.FailedToShutdownServer.Msg, err)
		}
		s.logger.LogInfo(logger.ServerShutdownOK.Type, logger.ServerShutdownOK.Msg)
	}
}
