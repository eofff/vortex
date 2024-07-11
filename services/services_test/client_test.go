package services_test

import (
	"fmt"
	"testing"
	"vortex/repository"
	"vortex/repository/repository_mock"
	"vortex/services"
	"vortex/services/services_mock"
)

func TestAdd(t *testing.T) {
	var clientService services.IClientService = &services.ClientService{}
	clientService.Init(nil, nil)

	var client services.Client
	client.ID = 1

	err := clientService.Add(client)
	if err != nil {
		t.Errorf("Error when add client in client service")
	}

	err = clientService.Add(client)
	if err == nil {
		t.Errorf("Error successful try to add client with existing id")
	}

	client.ID = -42
	err = clientService.Add(client)
	if err == nil {
		t.Errorf("Error successful try to add client with id lesser than 0")
	}
}

func TestUpdate(t *testing.T) {
	var clientService services.IClientService = &services.ClientService{}
	clientService.Init(nil, nil)

	var client services.Client
	client.ID = 1

	err := clientService.Add(client)
	if err != nil {
		t.Errorf("Error when add client in client service")
	}

	client.ClientName = "name123"
	err = clientService.Update(client)
	if err != nil {
		t.Errorf("Error when update client")
	}
	if client.ClientName != "name123" {
		t.Errorf("Error update client")
	}

	client.ID = -42
	err = clientService.Update(client)
	if err == nil {
		t.Errorf("Updated not existing entity")
	}
}

func TestDelete(t *testing.T) {
	var clientServiceStruct services.ClientService
	var clientService services.IClientService = &clientServiceStruct
	clientService.Init(nil, nil)

	var client services.Client
	client.ID = 1

	err := clientService.Add(client)
	if err != nil {
		t.Errorf("Error when add client in client service")
	}

	clientService.Delete(client)
	_, ok := clientServiceStruct.Clients[client.ID]
	if ok {
		t.Errorf("Client not deleted")
	}
}

func TestUpdateAlgorithmStatus(t *testing.T) {
	var algorithmStatusRepositoryMockStruct repository_mock.AlgorithmStatusRepositoryMock
	var algorithmStatusRepositoryMock repository.IAlgorithmStatusRepository = &algorithmStatusRepositoryMockStruct
	algorithmStatusRepositoryMock.Init(nil)

	var deployerMockStruct services_mock.DeployerService
	deployerMockStruct.InitService()
	var deployerMock services.Deployer = &deployerMockStruct

	var algorithmStatusService services.IAlgorithmStatusService = &services.AlgorithmStatusService{}
	algorithmStatusService.Init(algorithmStatusRepositoryMock)

	var clientServiceStruct services.ClientService
	var clientService services.IClientService = &clientServiceStruct
	clientService.Init(algorithmStatusService, deployerMock)

	var client services.Client
	client.ID = 1
	clientService.Add(client)

	var algorithmStatus repository.AlgorithmStatus
	algorithmStatus.ClientId = 1
	algorithmStatus.HFT = true
	algorithmStatus.TWAP = false
	algorithmStatus.VWAP = false
	algorithmStatusRepositoryMockStruct.AlgorithmStatuses = append(algorithmStatusRepositoryMockStruct.AlgorithmStatuses, &algorithmStatus)

	clientService.UpdateAlgorithmStatus()

	podNameHft := fmt.Sprintf("pod_%d_%s", algorithmStatus.ClientId, "HFT")
	_, ok := deployerMockStruct.Pods[podNameHft]
	if !ok {
		t.Errorf("needed pod does not exist")
	}

	algorithmStatus.HFT = false
	algorithmStatus.TWAP = false
	algorithmStatus.VWAP = true

	clientService.UpdateAlgorithmStatus()

	podNameVWAP := fmt.Sprintf("pod_%d_%s", algorithmStatus.ClientId, "VWAP")
	_, ok = deployerMockStruct.Pods[podNameVWAP]
	if !ok {
		t.Errorf("needed pod does not exist")
	}
	_, ok = deployerMockStruct.Pods[podNameHft]
	if ok {
		t.Errorf("pod was not deleted")
	}

	algorithmStatus.HFT = false
	algorithmStatus.TWAP = true
	algorithmStatus.VWAP = false

	clientService.UpdateAlgorithmStatus()

	podNameTWAP := fmt.Sprintf("pod_%d_%s", algorithmStatus.ClientId, "TWAP")
	_, ok = deployerMockStruct.Pods[podNameTWAP]
	if !ok {
		t.Errorf("needed pod does not exist")
	}
	_, ok = deployerMockStruct.Pods[podNameVWAP]
	if ok {
		t.Errorf("pod was not deleted")
	}

	clientService.Delete(client)

	clientService.UpdateAlgorithmStatus()

	_, ok = deployerMockStruct.Pods[podNameHft]
	if ok {
		t.Errorf("pod was not deleted")
	}
	_, ok = deployerMockStruct.Pods[podNameTWAP]
	if ok {
		t.Errorf("pod was not deleted")
	}
	_, ok = deployerMockStruct.Pods[podNameVWAP]
	if ok {
		t.Errorf("pod was not deleted")
	}
}
