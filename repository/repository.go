package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/mahditakrim/template/entity"
)

type Repository interface {
	Close() error
	InsertBook(context.Context, *entity.Book) error
	DeleteBook(context.Context, int64) error
	UpdateBook(context.Context, *entity.Book) error
	FindBook(context.Context, int64) (*entity.Book, error)
	GetBooks(context.Context) ([]entity.Book, error)
}

type postgres struct {
	conn *sql.DB
}

func NewPostgres(addr, username, password, database string) (Repository, error) {

	conn, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		username,
		password,
		addr,
		database))
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	return &postgres{conn}, nil
}

func (db *postgres) InsertBook(ctx context.Context, book *entity.Book) error {

	return db.conn.QueryRowContext(ctx, "INSERT INTO books (name, writer, page_num) VALUES ($1, $2, $3) RETURNING id",
		book.Name, book.Writer, book.PageNum).Scan(&book.ID)
}

func (db *postgres) DeleteBook(ctx context.Context, bookID int64) error {

	_, err := db.conn.ExecContext(ctx, "DELETE FROM books WHERE id = $1", bookID)

	return err
}

func (db *postgres) UpdateBook(ctx context.Context, book *entity.Book) error {

	_, err := db.conn.ExecContext(ctx, "UPDATE books SET name = $1, writer = $2, page_num = $3 WHERE id = $4",
		book.Name, book.Writer, book.PageNum, book.ID)

	return err
}

func (db *postgres) FindBook(ctx context.Context, bookID int64) (*entity.Book, error) {

	var book entity.Book
	if err := db.conn.QueryRowContext(ctx, "SELECT * FROM books WHERE id = $1", bookID).
		Scan(&book.ID, &book.Name, &book.Writer, &book.PageNum); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("book not found")
		}
		return nil, err
	}

	return &book, nil
}

func (db *postgres) GetBooks(ctx context.Context) ([]entity.Book, error) {

	rows, err := db.conn.QueryContext(ctx, "SELECT * FROM books")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("no books")
		}
		return nil, err
	}
	defer rows.Close()

	var books []entity.Book
	for rows.Next() {
		var book entity.Book
		err = rows.Scan(&book.ID, &book.Name, &book.Writer, &book.PageNum)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (db *postgres) Close() error {

	return db.conn.Close()
}
