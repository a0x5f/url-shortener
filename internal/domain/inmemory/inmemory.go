package inmemory

import (
	"context"
	"errors"
	"url-shortener/internal/domain/interfaces"
	"url-shortener/internal/domain/model"
)

var (
	ErrorItemNotFound = errors.New("item not found")
)

type InmemoryDb map[uint64]string

type inmemoryService struct {
	db *InmemoryDb
}

func (s *inmemoryService) Select(ctx context.Context, l *model.Link) (*model.Link, error) {
	if url, in := (*s.db)[l.Id]; in {
		return &model.Link{Id: l.Id, Url: url}, nil
	}

	if l.Url == "" {
		return nil, ErrorItemNotFound
	}

	for id, url := range *s.db {
		if url == l.Url {
			return &model.Link{Id: id, Url: url}, nil
		}
	}

	return nil, ErrorItemNotFound
}

func (s *inmemoryService) Insert(ctx context.Context, l *model.Link) error {
	for _, url := range *s.db {
		if l.Url == url {
			return nil
		}
	}

	id := 1 + uint64(len(*s.db))
	(*s.db)[id] = l.Url

	return nil
}

func New(db *InmemoryDb) interfaces.DataStorageService {
	return &inmemoryService{
		db: db,
	}
}
