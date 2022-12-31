package db

import "github.com/kitanoyoru/kita-ci/pkg/structs"

type DatabaseClient interface {
	CreateSchema() error

	InsertBuild(build *structs.Build) error
	AllBuilds(repoID int64, branch string) ([]*structs.Build, error)
  CountBuildsInRepoWithBranch(repoID int64, branch string) (int, error)
	FindBuildByID(buildID int64) (*structs.Build, error)
  
  InsertArtifact(artifact *structs.Artifact) error
  AllArtifacts(repoID int64, branch string) ([]*structs.Artifact, error)
  FindArtifactByBuildID(buildID int64) (*structs.Artifact, error)
  FindArtifactByID(id int64) (*structs.Artifact, error)

	Disconnect() error
}
