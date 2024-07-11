package repository_mock

import (
	"database/sql"
	"vortex/repository"
)

type AlgorithmStatusRepositoryMock struct {
	AlgorithmStatuses []*repository.AlgorithmStatus
}

func (a *AlgorithmStatusRepositoryMock) Init(db *sql.DB) {
	a.AlgorithmStatuses = make([]*repository.AlgorithmStatus, 0)
}

func (a *AlgorithmStatusRepositoryMock) GetAll() ([](*repository.AlgorithmStatus), error) {
	return a.AlgorithmStatuses, nil
}
