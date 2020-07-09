package notificators

import (
	"errors"
	"sync"
)

// TODO: Add persistent storage
type NotificatorsStorage interface {
	Add(services ...NotificatorService) error
	GetByChannel(channel string) (NotificatorService, error)
	Channels() ([]string, error)
}

func NewMemoryNotificationsStorage() NotificatorsStorage {
	return &memoryNotificatorsStorage{}
}

type memoryNotificatorsStorage struct {
	storage sync.Map
}

func (s *memoryNotificatorsStorage) Add(services ...NotificatorService) error {
	for _, service := range services {
		for _, channel := range service.Channels {
			s.storage.Store(channel, service)
		}
	}
	return nil
}

func (s *memoryNotificatorsStorage) GetByChannel(channel string) (NotificatorService, error) {
	service, ok := s.storage.Load(channel)
	if !ok {
		return NotificatorService{}, errors.New("service for this channel is not registered")
	}
	return service.(NotificatorService), nil
}

func (s *memoryNotificatorsStorage) Channels() ([]string, error) {
	keys := make([]string, 0)
	s.storage.Range(func(key, value interface{}) bool {
		keys = append(keys, key.(string))
		return true
	})
	return keys, nil
}
