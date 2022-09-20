package rest

import (
	"context"
	"errors"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/mahditakrim/template/service"
	"github.com/mahditakrim/template/transport"
)

type rest struct {
	router *iris.Application
}

func NewHttp(service service.Library) (transport.Transport, error) {

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

	return &rest{router}, nil
}

func (http *rest) Run(addr string) error {

	return http.router.Listen(addr, iris.WithoutServerError(iris.ErrServerClosed))
}

func (http *rest) Shutdown() error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return http.router.Shutdown(ctx)
}
