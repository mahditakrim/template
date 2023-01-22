package rest

import (
	"context"
	"errors"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/mahditakrim/template/luncher"
	"github.com/mahditakrim/template/service"
)

type rest struct {
	addr   string
	router *iris.Application
}

func NewHttp(service service.Service, addr string) (luncher.Runnable, error) {

	if service == nil {
		return nil, errors.New("nil service reference")
	}

	router := iris.New()
	router.PartyFunc("/library", func(library iris.Party) {
		library.Post("/book", addBookHandler(service))
		library.Delete("/book", removeBookHandler(service))
		library.Put("/book", editeBookHandler(service))
		library.Get("/book", getBookHandler(service))
		library.Get("/books", getBooksHandler(service))

	})

	return &rest{addr, router}, nil
}

func (http *rest) Run() error {

	return http.router.Listen(http.addr, iris.WithoutServerError(iris.ErrServerClosed))
}

func (http *rest) Shutdown() error {

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	return http.router.Shutdown(ctx)
}
