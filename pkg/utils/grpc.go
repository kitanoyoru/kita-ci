package utils

import (
	"net/http"

	ptypes "github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/status"

	grpcTypes "github.com/kitanoyoru/kita-proto/go"

	"github.com/kitanoyoru/kita-ci/pkg/structs"
)

func ParseBuildToProto(build *structs.Build) (*grpcTypes.Build, error) {
	timestamp, err := ptypes.TimestampProto(build.CreatedAt)
	if err != nil {
		return &grpcTypes.Build{}, status.New(http.StatusInternalServerError, "failed to parse model timestamp field to proto").Err()
	}

	return &grpcTypes.Build{
		ID:            build.ID,
		GithubRepoID:  build.GithubRepoID,
		Branch:        build.Branch,
		IsSuccessfull: build.IsSuccessfull,
		CreatedAt:     timestamp,
		Stdout:        build.Stdout,
	}, nil
}
