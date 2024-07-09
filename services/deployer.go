package services

import (
	"fmt"
	"log"
	"sync"
)

type Deployer interface {
	CreatePod(name string) error
	DeletePod(name string) error
	GetPodList() ([]string, error)
}

type DeployerService struct {
	pods      map[string]bool
	podsMutex sync.Mutex
}

func (d *DeployerService) InitService() {
	d.pods = make(map[string]bool)
}

func (d *DeployerService) CreatePod(name string) error {
	_, ok := d.pods[name]
	if ok {
		return fmt.Errorf("create pod error: pod with name %s exists", name)
	}

	d.podsMutex.Lock()
	d.pods[name] = true
	d.podsMutex.Unlock()

	log.Printf("Pod %s created\n", name)

	return nil
}

func (d *DeployerService) DeletePod(name string) error {
	_, ok := d.pods[name]
	if !ok {
		return fmt.Errorf("delete pod error: pod with name %s not exists", name)
	}

	d.podsMutex.Lock()
	delete(d.pods, name)
	d.podsMutex.Unlock()

	log.Printf("Pod %s deleted\n", name)
	return nil
}

func (d *DeployerService) GetPodList() ([]string, error) {
	log.Println("GetPodList called")

	result := make([]string, len(d.pods))
	cnt := 0
	for k := range d.pods {
		result[cnt] = k
		cnt++
	}

	return result, nil
}
