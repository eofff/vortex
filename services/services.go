package services

import "vortex/repository"

type Deployer interface {
	CreatePod(name string) error
	DeletePod(name string) error
	GetPodList() ([]string, error)
}

type IClientService interface {
	Init(algorithmStatusService IAlgorithmStatusService, deployerService Deployer)
	Add(client Client) error
	Update(client Client) error
	Delete(client Client) error
	UpdateAlgorithmStatus() error
}

type ICheckScheduleService interface {
	Init(clientService IClientService)
	StartWatcher()
}

type IAlgorithmStatusService interface {
	Init(algorithmStatusRepository repository.IAlgorithmStatusRepository)
	GetAll() ([](*repository.AlgorithmStatus), error)
}
