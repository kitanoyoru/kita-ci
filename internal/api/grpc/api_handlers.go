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
		s.log.Error("Filed to find build by ID", err)
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

func (s *GrpcAPIServer) GetAllBuilds(ctx context.Context, q grpcTypes.BuildQuery) (*grpcTypes.BuildList, error) {
	builds, err := s.dbClient.AllBuilds(q.GithubRepoID, q.Branch)
	if err != nil {
		s.log.Error("Failed to find all builds by GithubRepoID and branch", err)
		return &grpcTypes.BuildList{}, status.New(http.StatusInternalServerError, "").Err()
	}
	if builds == nil {
		return &grpcTypes.BuildList{}, status.New(http.StatusNotFound, "cannot find builds specified by query").Err()
	}

	size, err := s.dbClient.CountBuilds(q.GithubRepoID, q.Branch)
	if err != nil {
		s.log.Error("Failed to find all builds by GithubRepoID and branch", err)
		return &grpcTypes.BuildList{}, status.New(http.StatusInternalServerError, "").Err()
	}
	if builds == nil {
		return &grpcTypes.BuildList{}, status.New(http.StatusNotFound, "cannot find builds specified by query").Err()
	}

	buildList := &grpcTypes.BuildList{}
	buildList.Size = int64(size)
	for _, build := range builds {
		parsedBuild, err := utils.ParseBuildToProto(build)
		if err != nil {
			return &grpcTypes.BuildList{}, err
		}
		buildList.Builds = append(buildList.Builds, parsedBuild)
	}

	return buildList, nil
}
