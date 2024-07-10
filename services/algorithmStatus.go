package services

import (
	"vortex/repository"
)

type AlgorithmStatusService struct {
	algorithmStatusRepository repository.IAlgorithmStatusRepository
}

func (a *AlgorithmStatusService) Init(algorithmStatusRepository repository.IAlgorithmStatusRepository) {
	a.algorithmStatusRepository = algorithmStatusRepository
}

func (a *AlgorithmStatusService) GetAll() ([](*repository.AlgorithmStatus), error) {
	return a.algorithmStatusRepository.GetAll()
}
