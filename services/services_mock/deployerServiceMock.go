package services_mock

import (
	"fmt"
	"sync"
)

type DeployerServiceMock struct {
	Pods      map[string]bool
	podsMutex sync.Mutex
}

func (d *DeployerServiceMock) InitService() {
	d.Pods = make(map[string]bool)
}

func (d *DeployerServiceMock) CreatePod(name string) error {
	_, ok := d.Pods[name]
	if ok {
		return fmt.Errorf("create pod error: pod with name %s exists", name)
	}

	d.podsMutex.Lock()
	d.Pods[name] = true
	d.podsMutex.Unlock()

	return nil
}

func (d *DeployerServiceMock) DeletePod(name string) error {
	_, ok := d.Pods[name]
	if !ok {
		return fmt.Errorf("delete pod error: pod with name %s not exists", name)
	}

	d.podsMutex.Lock()
	delete(d.Pods, name)
	d.podsMutex.Unlock()

	return nil
}

func (d *DeployerServiceMock) GetPodList() ([]string, error) {
	result := make([]string, len(d.Pods))
	cnt := 0
	for k := range d.Pods {
		result[cnt] = k
		cnt++
	}

	return result, nil
}
