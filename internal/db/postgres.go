package db

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"

	"github.com/kitanoyoru/kita-ci/pkg/config"
	"github.com/kitanoyoru/kita-ci/pkg/structs"
)

const (
	ConnectAttempts = 5
)

type PostgresClient struct {
	pg *pg.DB
}

func NewPostgresClient(config config.PostgresConfig) *PostgresClient {
	return &PostgresClient{
		pg: pg.Connect(&pg.Options{
			Database: config.DB,
			User:     config.DBUser,
			Password: config.DBPassword,
			Addr:     config.DBAddr,

			MaxRetries: ConnectAttempts,
		}),
	}
}

func (p *PostgresClient) CreateSchema() error {
	models := []interface{}{
		(*structs.Build)(nil),
		(*structs.Artifact)(nil),
	}

	for _, model := range models {
		err := p.pg.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PostgresClient) InsertBuild(build *structs.Build) error {
	return p.pg.Insert(build)
}

func (p *PostgresClient) AllBuilds(repoID int64, branch string) ([]*structs.Build, error) {
	var builds []*structs.Build

	err := p.pg.Model(&builds).Where("github_repo_id = ?", repoID).Where("branch = ?", branch).Select()

	return builds, err
}

func (p *PostgresClient) CountBuildsInRepoWithBranch(repoID int64, branch string) (int, error) {
	count, err := p.pg.Model(&structs.Build{}).Where("github_repo_id = ?", repoID).Where("branch = ?", branch).Count()

	return count, err
}

func (p *PostgresClient) FindBuildByID(buildID int64) (*structs.Build, error) {
	build := &structs.Build{
		ID: buildID,
	}

	err := p.pg.Select(build)

	if err != nil {
		if err == pg.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return build, nil
}

func (p *PostgresClient) InsertArtifact(artifact *structs.Artifact) error {
	return p.pg.Insert(artifact)
}

func (p *PostgresClient) AllArtifacts(repoID int64, branch string) ([]*structs.Artifact, error) {
	var builds []int

	err := p.pg.Model((*structs.Build)(nil)).Where("github_repo_id = ?", repoID).Where("branch = ?", branch).ColumnExpr("array_agg(id)").Select(pg.Array(&builds))
	if err != nil {
		return nil, err
	}

	var artifacts []*structs.Artifact

	err = p.pg.Model(&artifacts).Where("build_id in (?)", pg.In(builds)).Select()
	if err != nil {
		return nil, err
	}

	return artifacts, nil
}

func (p *PostgresClient) FindArtifactByBuildID(buildID int64) (*structs.Artifact, error) {
	artifact := &structs.Artifact{}

	err := p.pg.Model(&artifact).Where("build_id = ?", buildID).Select()
	if err != nil {
		return nil, err
	}

	return artifact, nil
}

func (p *PostgresClient) FindArtifactByID(id int64) (*structs.Artifact, error) {
	artifact := &structs.Artifact{
		ID: id,
	}

	err := p.pg.Select(artifact)
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return artifact, nil
}

func (p *PostgresClient) Disconnect() error {
	return p.pg.Close()
}
