package api

import (
	"context"
	"net/http"

	grpcTypes "github.com/kitanoyoru/kita-proto/go"
	"google.golang.org/grpc/status"

	"github.com/kitanoyoru/kita-ci/pkg/utils"
)

func (s *GrpcAPIServer) GetBuildByID(ctx context.Context, buildID grpcTypes.BuildID) (*grpcTypes.Build, error) {
	build, err := s.dbClient.FindBuildByID(buildID.ID)
	if err != nil {
		s.log.Error("Filed to find build bu ID", err)
		return &grpcTypes.Build{}, status.New(http.StatusInternalServerError, "").Err()
	}
	if build == nil {
		return &grpcTypes.Build{}, status.New(http.StatusNotFound, "build not found").Err()
	}

	parsedBuild, err := utils.ParseBuildToProto(build)
	if err != nil {
		return &grpcTypes.Build{}, err
	}

	return parsedBuild, nil
}
