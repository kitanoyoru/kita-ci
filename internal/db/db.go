package db

import "github.com/kitanoyoru/kita-ci/pkg/structs"

type DatabaseClient interface {
	CreateSchema() error

	InsertBuild(*structs.Build) error
	AllBuilds(repoID int64, branch string) ([]*structs.Build, error)
	CountBuilds(repoID int64, branch string) (int, error)
	FindBuildByID(buildID int64) (*structs.Build, error)

	Disconnect() error
}
