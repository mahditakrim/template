package rpc

import (
	"context"

	"github.com/mahditakrim/template/entity"
	"github.com/mahditakrim/template/transport/rpc/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (rpc *rpc) GetBook(ctx context.Context, book *pb.BookID) (*pb.Book, error) {

	b, err := rpc.service.GetBook(ctx, book.GetId())
	if err != nil {
		return nil, err
	}

	return &pb.Book{
		Id:      b.ID,
		Name:    b.Name,
		Writer:  b.Writer,
		PageNum: uint32(b.PageNum),
	}, nil
}

func (rpc *rpc) GetBooks(ctx context.Context, _ *emptypb.Empty) (*pb.Books, error) {

	bs, err := rpc.service.GetBooks(ctx)
	if err != nil {
		return nil, err
	}

	var books pb.Books
	for _, b := range bs {
		books.Books = append(books.Books, &pb.Book{
			Id:      b.ID,
			Name:    b.Name,
			Writer:  b.Writer,
			PageNum: uint32(b.PageNum),
		})
	}

	return &books, nil
}

func (rpc *rpc) EditeBook(ctx context.Context, book *pb.Book) (*emptypb.Empty, error) {

	b := &entity.Book{
		ID:      book.GetId(),
		Name:    book.GetName(),
		Writer:  book.GetWriter(),
		PageNum: uint(book.GetPageNum()),
	}

	err := rpc.service.Validate(b)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, rpc.service.EditeBook(ctx, b)
}

func (rpc *rpc) CreateBook(ctx context.Context, book *pb.Book) (*pb.BookID, error) {

	b := &entity.Book{
		Name:    book.GetName(),
		Writer:  book.GetWriter(),
		PageNum: uint(book.GetPageNum()),
	}

	err := rpc.service.Validate(b)
	if err != nil {
		return nil, err
	}

	err = rpc.service.AddBook(ctx, b)
	if err != nil {
		return nil, err
	}

	return &pb.BookID{Id: b.ID}, nil
}

func (rpc *rpc) RemoveBook(ctx context.Context, book *pb.BookID) (*emptypb.Empty, error) {

	return &emptypb.Empty{}, rpc.service.RemoveBook(ctx, book.GetId())
}
