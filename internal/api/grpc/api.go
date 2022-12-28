package api

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	grpcTypes "github.com/kitanoyoru/kita-proto/go"

	"github.com/kitanoyoru/kita-ci/internal/db"
	"github.com/kitanoyoru/kita-ci/pkg/config"
	"github.com/kitanoyoru/kita-ci/pkg/log"
)

// TODO: Change Files Types to Interfaces
type GrpcAPIServer struct {
	config   config.GrpcAPIServerConfig
	log      log.Logger
	dbClient db.DatabaseClient
}

func NewGrpcAPIServer(cfg config.GrpcAPIServerConfig, logger log.Logger, dbClient db.DatabaseClient) *GrpcAPIServer {
	return &GrpcAPIServer{
		config:   cfg,
		log:      logger,
		dbClient: dbClient,
	}
}

func (s *GrpcAPIServer) Start() {
	defer s.dbClient.Disconnect()
	err := s.dbClient.CreateSchema()
	if err != nil {
		s.log.Fatal("Failed to create db schema", err)
	}

	s.log.Info(fmt.Sprintf("Starting GRPC API Server at port %d", s.config.Port))
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.Port))
	if err != nil {
		s.log.Fatal("Failed to start GRPC API Server", err)
	}

	grpcServer := grpc.NewServer()
	grpcTypes.RegisterCIServer(grpcServer, s)
	if err := grpcServer.Serve(lis); err != nil {
		s.log.Fatal("Failed to serve GRPC API Server", err)
	}

	s.log.Info("GRPC API Server started successfully")
}
