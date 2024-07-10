package services

import "time"

type CheckScheduleService struct {
	clientService IClientService
}

func (c *CheckScheduleService) Init(clientService IClientService) {
	c.clientService = clientService
}

func (c *CheckScheduleService) StartWatcher() {
	go func() {
		for {
			c.clientService.UpdateAlgorithmStatus()
			time.Sleep(300 * time.Second)
		}
	}()
}
