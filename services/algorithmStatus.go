package services

import (
	"vortex/repository"
)

type AlgorithmStatusService struct {
}

var algorithmStatusRepository repository.AlgorithmStatusRepository

func (a *AlgorithmStatusService) Get() ([](*repository.AlgorithmStatusRepository), error) {
	return algorithmStatusRepository.GetAll()
}

func (a *AlgorithmStatusService) GetById(id int) (*repository.AlgorithmStatusRepository, error) {
	res, err := algorithmStatusRepository.GetById(id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *AlgorithmStatusService) Insert(algorithmStatus *repository.AlgorithmStatusRepository) error {
	return algorithmStatus.Insert(algorithmStatus)
}

func (a *AlgorithmStatusService) Update(algorithmStatus *repository.AlgorithmStatusRepository) error {
	return algorithmStatus.Update(algorithmStatus)
}

func (a *AlgorithmStatusService) Delete(algorithmStatus *repository.AlgorithmStatusRepository) error {
	return algorithmStatus.Delete(algorithmStatus)
}
