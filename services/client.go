package services

import (
	"fmt"
	"log"
	"sync"
	"time"
	"vortex/repository"
)

type Client struct {
	ID          int64     `json:"id"`
	ClientName  string    `json:"clientName"`
	Version     int       `json:"version"`
	Image       string    `json:"image"`
	CPU         string    `json:"cpu"`
	Memory      string    `json:"memory"`
	Priority    float64   `json:"priority"`
	NeedRestart bool      `json:"needRestart"`
	SpawnedAt   time.Time `json:"spawnedAt"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ClientService struct {
	algorithmStatusService IAlgorithmStatusService
	deployerService        Deployer
	Clients                map[int64]Client
	clientMapMutex         sync.Mutex
}

func (c *ClientService) Init(algorithmStatusService IAlgorithmStatusService, deployerService Deployer) {
	c.Clients = make(map[int64]Client)
	c.deployerService = deployerService
	c.algorithmStatusService = algorithmStatusService
	log.Println("Client service initialized")
}

func (c *ClientService) Add(client Client) error {
	_, ok := c.Clients[client.ID]
	if ok {
		return fmt.Errorf("create error, client with id: %d alredy exists", client.ID)
	}

	c.clientMapMutex.Lock()
	defer c.clientMapMutex.Unlock()

	if client.ID < 0 {
		return fmt.Errorf("create error, client id must be greater then 0, got: %d", client.ID)
	}

	c.Clients[client.ID] = client
	log.Printf("Client added with id: %d", client.ID)

	return nil
}

func (c *ClientService) Update(client Client) error {
	_, ok := c.Clients[client.ID]
	if !ok {
		return fmt.Errorf("update error, client with id: %d not exists", client.ID)
	}

	c.clientMapMutex.Lock()
	c.Clients[client.ID] = client
	c.clientMapMutex.Unlock()

	log.Printf("Client updated with id: %d", client.ID)

	return nil
}

func (c *ClientService) Delete(client Client) error {
	c.clientMapMutex.Lock()
	delete(c.Clients, client.ID)
	c.clientMapMutex.Unlock()

	log.Printf("Client deleted with id: %d", client.ID)

	return nil
}

func (c *ClientService) UpdateAlgorithmStatus() error {
	log.Println("Update alogrithm status called")

	statuses, err := c.algorithmStatusService.GetAll()
	if err != nil {
		log.Fatal(err)
	}

	algorithmStatuses := make(map[int64][]repository.AlgorithmStatus)
	for _, status := range statuses {
		_, ok := algorithmStatuses[status.ClientId]
		if !ok {
			algorithmStatuses[status.ClientId] = make([]repository.AlgorithmStatus, 0)
		}
		algorithmStatuses[status.ClientId] = append(algorithmStatuses[status.ClientId], *status)
	}

	pods, err := c.deployerService.GetPodList()
	if err != nil {
		log.Fatal("error while get pods")
	}

	podMap := make(map[string]bool)
	for _, pod := range pods {
		podMap[pod] = false
	}

	for cid := range c.Clients {
		clientStatuses, ok := algorithmStatuses[cid]
		if !ok {
			log.Printf("Client with id %d have no statuses in db", cid)
			continue
		}

		for _, clientStatus := range clientStatuses {
			c.UpdateAlgorithm(podMap, &clientStatus, "HFT")
			c.UpdateAlgorithm(podMap, &clientStatus, "VWAP")
			c.UpdateAlgorithm(podMap, &clientStatus, "TWAP")
		}
	}

	for pod, val := range podMap {
		if !val {
			c.deployerService.DeletePod(pod)
		}
	}

	return nil
}

func (c *ClientService) UpdateAlgorithm(
	podMap map[string]bool,
	clientStatus *repository.AlgorithmStatus,
	algorithmName string,
) {
	podName := fmt.Sprintf("pod_%d_%s", clientStatus.ClientId, algorithmName)
	_, podExists := podMap[podName]
	algorithmEnabled := false

	switch algorithmName {
	case "HFT":
		algorithmEnabled = clientStatus.HFT
	case "VWAP":
		algorithmEnabled = clientStatus.VWAP
	case "TWAP":
		algorithmEnabled = clientStatus.TWAP
	default:
		return
	}

	if algorithmEnabled {
		if !podExists {
			c.deployerService.CreatePod(podName)
		}
		podMap[podName] = true
	} else {
		if podExists {
			c.deployerService.DeletePod(podName)
			delete(podMap, podName)
		}
	}
}
