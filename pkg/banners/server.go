package banners

import (
	"context"
	"errors"
	"sync"
)

var startId int64

//Service представляет собой сервис по управления баннерами
type Service struct {
	mu    sync.RWMutex
	items []*Banner
}

//NewService создаёт сервис
func NewService() *Service {
	return &Service{items: make([]*Banner, 0)}
}

//Banner предсталяет собой баннер
type Banner struct {
	ID      int64
	Title   string
	Content string
	Button  string
	Link    string
}

//All ...
func (s *Service) All(ctx context.Context) ([]*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.items, nil
}

//ByID ...
func (s *Service) ByID(ctx context.Context, id int64) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, banner := range s.items {
		if banner.ID == id {
			return banner, nil
		}
	}

	return nil, errors.New("item not found")
}

//Save ...
func (s *Service) Save(ctx context.Context, item *Banner) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if item.ID == 0 {
		startId++
		item.ID = startId
		s.items = append(s.items, item)
		return item, nil
	}

	for key, banner := range s.items {
		if banner.ID == item.ID {
			s.items[key] = item
			return item, nil
		}
	}

	return nil, errors.New("item not found")
}

//RemoveByID ... Метод для удаления
func (s *Service) RemoveByID(ctx context.Context, id int64) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for k, banner := range s.items {
		if banner.ID == id {
			s.items = append(s.items[:k], s.items[k+1:]...)
			return banner, nil
		}
	}

	return nil, errors.New("item not found")
}
