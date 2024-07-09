package services

import "time"

type CheckScheduleService struct {
	clientService *ClientService
}

func (c *CheckScheduleService) InitService(clientService *ClientService) {
	c.clientService = clientService
}

func (c *CheckScheduleService) StartWatcher() {
	for {
		c.clientService.UpdateAlgorithmStatus()
		time.Sleep(300 * time.Second)
	}
}
