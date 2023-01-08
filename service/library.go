package service

import (
	"context"
	"errors"

	"github.com/mahditakrim/template/entity"
	"github.com/mahditakrim/template/repository"
)

type Library interface {
	Validate(*entity.Book) error
	AddBook(ctx context.Context, book *entity.Book) error
	RemoveBook(ctx context.Context, bookID int64) error
	EditeBook(ctx context.Context, book *entity.Book) error
	GetBook(ctx context.Context, bookID int64) (*entity.Book, error)
	GetBooks(ctx context.Context) ([]entity.Book, error)
}

type library struct {
	repo repository.Repository
}

func NewLibrary(repo repository.Repository) (Library, error) {

	if repo == nil {
		return nil, errors.New("nil respository reference")
	}

	return &library{repo}, nil
}

func (*library) Validate(book *entity.Book) error {

	if book == nil {
		return errors.New("nil slice of book")
	}

	if book.Name == "" {
		return errors.New("empty book name")
	}

	if book.Writer == "" {
		return errors.New("empty book writer")
	}

	if book.PageNum <= 0 {
		return errors.New("page num should be greater than 0")
	}

	return nil
}

func (l *library) AddBook(ctx context.Context, book *entity.Book) error {

	return l.repo.InsertBook(ctx, book)
}

func (l *library) RemoveBook(ctx context.Context, bookID int64) error {

	return l.repo.DeleteBook(ctx, bookID)
}

func (l *library) EditeBook(ctx context.Context, book *entity.Book) error {

	return l.repo.UpdateBook(ctx, book)
}

func (l *library) GetBook(ctx context.Context, bookID int64) (*entity.Book, error) {

	return l.repo.FindBook(ctx, bookID)
}

func (l *library) GetBooks(ctx context.Context) ([]entity.Book, error) {

	return l.repo.GetBooks(ctx)
}
