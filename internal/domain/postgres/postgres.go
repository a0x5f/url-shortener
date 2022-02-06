package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"url-shortener/internal/domain/interfaces"
	"url-shortener/internal/domain/model"
)

const (
	postgresUniqueConstraintError = "23505"
)

type postgresService struct {
	db *sql.DB
}

func (s *postgresService) Select(ctx context.Context, l *model.Link) (*model.Link, error) {
	link := &model.Link{}
	sqlUrl := fmt.Sprintf("'%s'", l.Url)
	query := fmt.Sprintf("SELECT * FROM links WHERE id = %d OR url = %s", l.Id, sqlUrl)

	err := s.db.QueryRowContext(ctx, query).Scan(&link.Id, &link.Url)

	if err != nil {
		return nil, err
	}

	return link, nil

}

func (s *postgresService) Insert(ctx context.Context, l *model.Link) error {
	query := fmt.Sprintf("INSERT INTO links (url) VALUES ('%s')", l.Url)
	_, err := s.db.ExecContext(ctx, query)

	if err != nil {
		switch err.(type) {
		case *pq.Error:
			if err.(*pq.Error).Code != postgresUniqueConstraintError {
				return err
			}
		default:
			return err
		}
	}

	return nil
}

func New(db *sql.DB) interfaces.DataStorageService {
	return &postgresService{
		db: db,
	}
}
