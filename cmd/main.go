package main

import (
	"io"
	"log"
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

	// init logger
	// it's way better to implement logger interface and inject it to domains.
	f, err := os.OpenFile(config.Get().LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(io.MultiWriter(f, os.Stdout))

	// init repository
	repository, err := repository.NewPostgres(
		config.Get().DB.Postgres.Addr,
		config.Get().DB.Postgres.Username,
		config.Get().DB.Postgres.Password,
		config.Get().DB.Postgres.Db,
	)
	if err != nil {
		log.Panic(err)
	}

	// init service usecase
	library, err := service.NewLibrary(repository)
	if err != nil {
		log.Panic(err)
	}

	// init restful api server
	httpServer, err := rest.NewHttp(library, config.Get().Transport.HttpAddr)
	if err != nil {
		log.Panic(err)
	}

	// init rpc api server
	rpcServer, err := rpc.NewRPC(library, config.Get().Transport.RpcAddr)
	if err != nil {
		log.Panic(err)
	}

	// lunch rest server
	go func() {
		log.Println("start http server on port", config.Get().Transport.HttpAddr)
		if err := httpServer.Run(); err != nil {
			log.Panic(err)
		}
	}()

	// lunch rpc server
	go func() {
		log.Println("start rpc server on port", config.Get().Transport.RpcAddr)
		if err := rpcServer.Run(); err != nil {
			log.Panic(err)
		}
	}()

	// gracefull shutdown
	shutdown.Graceful(func() {
		log.Println("shutting down")
		err = httpServer.Shutdown()
		if err != nil {
			log.Println(err)
		}
		err = rpcServer.Shutdown()
		if err != nil {
			log.Println(err)
		}
		err = repository.Close()
		if err != nil {
			log.Println(err)
		}
	}, make(chan os.Signal))
}
