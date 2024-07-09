package services

import (
	"fmt"
	"log"
	"sync"
	"time"
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
	clients  map[int64]Client
	mapMutex sync.Mutex
}

func (c *ClientService) InitService() {
	c.clients = make(map[int64]Client)
	log.Println("Client service initialized")
}

func (c *ClientService) Add(client Client) error {
	_, ok := c.clients[client.ID]
	if ok {
		return fmt.Errorf("create error, client with id: %d alredy exists", client.ID)
	}

	c.mapMutex.Lock()
	defer c.mapMutex.Unlock()

	if client.ID < 0 {
		return fmt.Errorf("create error, client id must be greater then 0, got: %d", client.ID)
	}

	c.clients[client.ID] = client
	log.Printf("Client added with id: %d", client.ID)

	return nil
}

func (c *ClientService) Update(client Client) error {
	_, ok := c.clients[client.ID]
	if !ok {
		return fmt.Errorf("update error, client with id: %d not exists", client.ID)
	}

	c.mapMutex.Lock()
	c.clients[client.ID] = client
	c.mapMutex.Unlock()

	log.Printf("Client updated with id: %d", client.ID)

	return nil
}

func (c *ClientService) Delete(client Client) error {
	c.mapMutex.Lock()
	delete(c.clients, client.ID)
	c.mapMutex.Unlock()

	log.Printf("Client deleted with id: %d", client.ID)

	return nil
}

func (c *ClientService) UpdateAlgorithmStatus() error {
	log.Println("Update alogrithm status called")

	return nil
}
