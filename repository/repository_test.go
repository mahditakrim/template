package repository

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
	"github.com/mahditakrim/template/entity"
)

var book = &entity.Book{
	ID:      42224,
	Name:    "test",
	Writer:  "tester",
	PageNum: 113,
}

func newMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("sqlmock got error = %v", err)
	}

	return db, mock
}

func Test_postgres_InsertBook(t *testing.T) {

	db, mock := newMock(t)
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO books (name, writer, page_num) VALUES ($1, $2, $3) RETURNING id")).
		WithArgs(book.Name, book.Writer, book.PageNum).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(13))

	postgres := &postgres{db}
	err := postgres.InsertBook(context.Background(), book)
	if err != nil {
		t.Errorf("postgres.InsertBook() error = %v", err)
		return
	}

	if book.ID != 13 {
		t.Errorf("invalid id returning in book struct")
		return
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("mock.ExpectationsWereMet() error = %v", err)
	}
}

func Test_postgres_DeleteBook(t *testing.T) {

	db, mock := newMock(t)
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM books WHERE id = $1")).WithArgs(13).
		WillReturnResult(sqlmock.NewResult(0, 1))

	postgres := &postgres{db}
	err := postgres.DeleteBook(context.Background(), 13)
	if err != nil {
		t.Errorf("postgres.DeleteBook() error = %v", err)
		return
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("mock.ExpectationsWereMet() error = %v", err)
	}
}

func Test_postgres_UpdateBook(t *testing.T) {

	db, mock := newMock(t)
	mock.ExpectExec(regexp.QuoteMeta("UPDATE books SET name = $1, writer = $2, page_num = $3 WHERE id = $4")).
		WithArgs(book.Name, book.Writer, book.PageNum, book.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	postgres := &postgres{db}
	err := postgres.UpdateBook(context.Background(), book)
	if err != nil {
		t.Errorf("postgres.UpdateBook() error = %v", err)
		return
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("mock.ExpectationsWereMet() error = %v", err)
	}
}

func Test_postgres_FindBook(t *testing.T) {

	dbProper, mockProper := newMock(t)
	mockProper.ExpectQuery(regexp.QuoteMeta("SELECT * FROM books WHERE id = $1")).WithArgs(book.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "writer", "page_num"}).
			AddRow(book.ID, book.Name, book.Writer, book.PageNum))

	dbNotFound, mockNotfound := newMock(t)
	mockNotfound.ExpectQuery(regexp.QuoteMeta("SELECT * FROM books WHERE id = $1")).WithArgs(book.ID).
		WillReturnError(sql.ErrNoRows)

	dbErr, mockErr := newMock(t)
	mockErr.ExpectQuery(regexp.QuoteMeta("SELECT * FROM books WHERE id = $1")).WithArgs(book.ID).
		WillReturnError(nil)

	type fields struct {
		conn *sql.DB
		mock sqlmock.Sqlmock
	}
	tests := []struct {
		name    string
		fields  fields
		want    *entity.Book
		wantErr bool
	}{
		{
			name:    "proper",
			fields:  fields{dbProper, mockProper},
			want:    book,
			wantErr: false,
		},
		{
			name:    "not found",
			fields:  fields{dbNotFound, mockNotfound},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error",
			fields:  fields{dbErr, mockErr},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			postgres := &postgres{tt.fields.conn}
			got, err := postgres.FindBook(context.Background(), book.ID)
			if (err != nil) != tt.wantErr {
				t.Errorf("postgres.FindBook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.name == "not found" && err.Error() != "book not found" {
				t.Errorf("sql.ErrNoRows missing translation")
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("postgres.FindBook() = %v, want %v", got, tt.want)
				return
			}

			err = tt.fields.mock.ExpectationsWereMet()
			if err != nil {
				t.Errorf("mock.ExpectationsWereMet() error = %v", err)
			}
		})
	}
}

func Test_postgres_GetBooks(t *testing.T) {

	dbProper, mockProper := newMock(t)
	mockProper.ExpectQuery(regexp.QuoteMeta("SELECT * FROM books")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "writer", "page_num"}).
			AddRow(book.ID, book.Name, book.Writer, book.PageNum))

	dbNoBooks, mockNoBooks := newMock(t)
	mockNoBooks.ExpectQuery(regexp.QuoteMeta("SELECT * FROM books")).WillReturnRows(sqlmock.NewRows(nil))

	dbRowsErr, mockRowsErr := newMock(t)
	mockRowsErr.ExpectQuery(regexp.QuoteMeta("SELECT * FROM books")).
		WillReturnRows(sqlmock.NewRows([]string{"t"}).AddRow("t").RowError(0, sqlmock.ErrCancelled))

	dbScanErr, mockScanErr := newMock(t)
	mockScanErr.ExpectQuery(regexp.QuoteMeta("SELECT * FROM books")).
		WillReturnRows(sqlmock.NewRows([]string{"t"}).AddRow("t").RowError(0, nil))

	dbErr, mockErr := newMock(t)
	mockErr.ExpectQuery(regexp.QuoteMeta("SELECT * FROM books")).WillReturnError(nil)

	type fields struct {
		conn *sql.DB
		mock sqlmock.Sqlmock
	}
	tests := []struct {
		name    string
		fields  fields
		want    []entity.Book
		wantErr bool
	}{
		{
			name:    "proper",
			fields:  fields{dbProper, mockProper},
			want:    []entity.Book{*book},
			wantErr: false,
		},
		{
			name:    "not found",
			fields:  fields{dbNoBooks, mockNoBooks},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "rows error",
			fields:  fields{dbRowsErr, mockRowsErr},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "scan error",
			fields:  fields{dbScanErr, mockScanErr},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error",
			fields:  fields{dbErr, mockErr},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			postgres := &postgres{tt.fields.conn}
			got, err := postgres.GetBooks(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("postgres.GetBooks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.name == "not found" && err.Error() != "no books" {
				t.Errorf("missing no books error translation")
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("postgres.GetBooks() = %v, want %v", got, tt.want)
				return
			}

			err = tt.fields.mock.ExpectationsWereMet()
			if err != nil {
				t.Errorf("mock.ExpectationsWereMet() error = %v", err)
			}
		})
	}
}

func Test_postgres_Close(t *testing.T) {

	dbProper, mockProper := newMock(t)
	mockProper.ExpectClose()

	dbErr, mockErr := newMock(t)
	mockErr.ExpectClose().WillReturnError(errors.New(""))

	type fields struct {
		conn *sql.DB
		mock sqlmock.Sqlmock
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "proper",
			fields:  fields{dbProper, mockProper},
			wantErr: false,
		},
		{
			name:    "error",
			fields:  fields{dbErr, mockErr},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			postgres := &postgres{tt.fields.conn}
			err := postgres.Close()
			if (err != nil) != tt.wantErr {
				t.Errorf("postgres.Close() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			err = tt.fields.mock.ExpectationsWereMet()
			if err != nil {
				t.Errorf("mock.ExpectationsWereMet() error = %v", err)
			}
		})
	}
}
