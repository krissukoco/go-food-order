package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
)

type Server struct {
	restServer    *fiber.App
	restPort      uint
	grpcServer    *grpc.Server
	grpcPort      uint
	exitTimeout   time.Duration
	shutdownFuncs map[string]ShutdownFunc
}

type ShutdownFunc = func(ctx context.Context) error

type ServerOption = func(s *Server)

func NewServer(
	restPort uint,
	grpcPort uint,
	exitTimeout time.Duration,
	shutdownFuncs map[string]ShutdownFunc,
	opts ...ServerOption,
) *Server {
	fiberApp := fiber.New()
	grpcApp := grpc.NewServer()
	s := &Server{
		restServer:    fiberApp,
		restPort:      restPort,
		grpcServer:    grpcApp,
		grpcPort:      grpcPort,
		exitTimeout:   exitTimeout,
		shutdownFuncs: shutdownFuncs,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func WithFiberConfig(cfg fiber.Config) ServerOption {
	return func(s *Server) {
		s.restServer = fiber.New(cfg)
	}
}

func (s *Server) RestGroup(prefix string) fiber.Router {
	return s.restServer.Group(prefix)
}

func (s *Server) GrpcServer() *grpc.Server {
	return s.grpcServer
}

func (s *Server) AddShutdownFunc(key string, fn ShutdownFunc) {
	s.shutdownFuncs[key] = fn
}

func (s *Server) ListenAndStartGracefully(ctx context.Context) error {
	errChan := make(chan error, 1)
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.grpcPort))
		if err != nil {
			errChan <- err
			return
		}
		errChan <- err
		log.Printf("Starting gRPC server on port %d\n", s.grpcPort)
		if err = s.grpcServer.Serve(lis); err != nil {
			errChan <- err
			return
		}
	}()

	go func() {
		log.Printf("Starting REST server on port %d\n", s.restPort)
		if err := s.restServer.Listen(fmt.Sprintf(":%d", s.restPort)); err != nil {
			errChan <- err
			return
		}
	}()

	// Gracefully shutdown
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	for {
		select {
		case err := <-errChan:
			if err != nil {
				log.Printf("ERROR starting server: %v\n", err)
				os.Exit(1)
			}
		case <-exit:
			log.Println("Shutdown signal received. Shutting down server...")

			exitFunc := time.AfterFunc(s.exitTimeout, func() {
				log.Printf("Graceful shutdown taking too long after %v. Forcefully shutdown!", s.exitTimeout)
				os.Exit(1)
			})
			defer exitFunc.Stop()

			// Shutdown servers
			log.Println("Shutting down gRPC server...")
			s.grpcServer.Stop()
			log.Println("Shutting down REST server...")
			s.restServer.Shutdown()

			// Cleanup operations (DB, connections, etc.)
			var wg sync.WaitGroup

			for k, op := range s.shutdownFuncs {
				wg.Add(1)
				go func(key string, operation ShutdownFunc) {
					defer wg.Done()
					log.Printf("Performing shutdown for '%s' ... \n", key)
					err := operation(ctx)
					if err != nil {
						log.Printf("ERROR on shutting down '%s': %v", key, err)
					}
				}(k, op)
			}

			wg.Wait()
			log.Println("Shutting down gracefully completed")
			return nil
		}
	}
}
