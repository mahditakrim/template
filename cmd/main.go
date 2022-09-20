package main

import (
	"os"

	"github.com/mahditakrim/template/cmd/config"
	"github.com/mahditakrim/template/internal/shutdown"
	"github.com/mahditakrim/template/repository"
	"github.com/mahditakrim/template/service"
	"github.com/mahditakrim/template/transport/rest"
	"github.com/mahditakrim/template/transport/rpc"
)

func main() {

	// init config
	err := config.Init()
	if err != nil {
		panic(err)
	}

	// make dependencies in order
	repository, err := repository.NewPostgres(
		config.Get().DB.Postgres.Addr,
		config.Get().DB.Postgres.Username,
		config.Get().DB.Postgres.Password,
		config.Get().DB.Postgres.Db,
	)
	if err != nil {
		panic(err)
	}

	library, err := service.NewLibrary(repository)
	if err != nil {
		panic(err)
	}

	httpServer, err := rest.NewHttp(library)
	if err != nil {
		panic(err)
	}

	rpcServer, err := rpc.NewRPC(library)
	if err != nil {
		panic(err)
	}

	// run app
	go func() {
		err := httpServer.Run(config.Get().Transport.HttpAddr)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		err := rpcServer.Run(config.Get().Transport.RpcAddr)
		if err != nil {
			panic(err)
		}
	}()

	// gracefull shutdown
	err = shutdown.Graceful(
		func() error {
			err := httpServer.Shutdown()
			if err != nil {
				return err
			}
			_ = rpcServer.Shutdown()
			return repository.Close()
		},
		make(chan os.Signal, 1),
	)
	if err != nil {
		panic(err)
	}
}
