package service

import (
	"context"

	"github.com/mahditakrim/template/entity"
)

type Service interface {
	Validate(*entity.Book) error
	AddBook(ctx context.Context, book *entity.Book) error
	RemoveBook(ctx context.Context, bookID int64) error
	EditeBook(ctx context.Context, book *entity.Book) error
	GetBook(ctx context.Context, bookID int64) (*entity.Book, error)
	GetBooks(ctx context.Context) ([]entity.Book, error)
}
