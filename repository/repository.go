package repository

import (
	"context"

	_ "github.com/lib/pq"
	"github.com/mahditakrim/template/entity"
)

type Repository interface {
	Close() error
	InsertBook(ctx context.Context, book *entity.Book) error
	DeleteBook(ctx context.Context, bookID int64) error
	UpdateBook(ctx context.Context, book *entity.Book) error
	FindBook(ctx context.Context, bookID int64) (*entity.Book, error)
	GetBooks(ctx context.Context) ([]entity.Book, error)
}
