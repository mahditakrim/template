package rest

import (
	"github.com/kataras/iris/v12"
	"github.com/mahditakrim/template/entity"
	"github.com/mahditakrim/template/service"
)

func addBookHandler(library service.Service) func(iris.Context) {
	return func(ctx iris.Context) {

		book := &entity.Book{}
		err := ctx.ReadJSON(book)
		if err != nil {
			ctx.StatusCode(iris.StatusUnprocessableEntity)
			ctx.JSON(struct{ Err string }{err.Error()})
			return
		}

		err = library.Validate(book)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(struct{ Err string }{err.Error()})
			return
		}

		err = library.AddBook(ctx.Request().Context(), book)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(struct{ Err string }{err.Error()})
			return
		}

		ctx.JSON(struct{ ID int64 }{book.ID})
	}
}

func removeBookHandler(library service.Service) func(iris.Context) {
	return func(ctx iris.Context) {

		id, err := ctx.URLParamInt64("book_id")
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(struct{ Err string }{err.Error()})
			return
		}

		err = library.RemoveBook(ctx.Request().Context(), id)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(struct{ Err string }{err.Error()})
			return
		}

		ctx.WriteString("ok")
	}
}

func editeBookHandler(library service.Service) func(iris.Context) {
	return func(ctx iris.Context) {

		book := &entity.Book{}
		err := ctx.ReadJSON(book)
		if err != nil {
			ctx.StatusCode(iris.StatusUnprocessableEntity)
			ctx.JSON(struct{ Err string }{err.Error()})
			return
		}

		err = library.Validate(book)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(struct{ Err string }{err.Error()})
			return
		}

		err = library.EditeBook(ctx.Request().Context(), book)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(struct{ Err string }{err.Error()})
			return
		}

		ctx.WriteString("ok")
	}
}

func getBookHandler(library service.Service) func(iris.Context) {
	return func(ctx iris.Context) {

		id, err := ctx.URLParamInt64("book_id")
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(struct{ Err string }{err.Error()})
			return
		}

		book, err := library.GetBook(ctx.Request().Context(), id)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(struct{ Err string }{err.Error()})
			return
		}

		ctx.JSON(book)
	}
}

func getBooksHandler(library service.Service) func(iris.Context) {
	return func(ctx iris.Context) {

		book, err := library.GetBooks(ctx.Request().Context())
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(struct{ Err string }{err.Error()})
			return
		}

		ctx.JSON(book)
	}
}
