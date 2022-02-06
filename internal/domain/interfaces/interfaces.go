package interfaces

import (
	"context"
	"url-shortener/internal/domain/model"
)

type DataStorageService interface {
	Select(ctx context.Context, l *model.Link) (*model.Link, error)
	Insert(ctx context.Context, l *model.Link) error
}
