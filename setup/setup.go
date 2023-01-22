package setup

import (
	"log"

	"github.com/mahditakrim/template/config"
	"github.com/mahditakrim/template/internal/rest"
	"github.com/mahditakrim/template/internal/rpc"
	"github.com/mahditakrim/template/luncher"
	"github.com/mahditakrim/template/repository"
	"github.com/mahditakrim/template/service"
)

func Init(conf *config.Config) ([]luncher.Runnable, error) {

	repository, err := repository.NewPostgres(
		conf.DB.Postgres.Addr,
		conf.DB.Postgres.Username,
		conf.DB.Postgres.Password,
		conf.DB.Postgres.Db,
	)
	if err != nil {
		log.Panic(err)
	}

	library, err := service.NewLibrary(repository)
	if err != nil {
		log.Panic(err)
	}

	httpServer, err := rest.NewHttp(library, conf.Web.HttpAddr)
	if err != nil {
		log.Panic(err)
	}

	rpcServer, err := rpc.NewRPC(library, conf.Web.RpcAddr)
	if err != nil {
		log.Panic(err)
	}

	return []luncher.Runnable{httpServer, rpcServer}, nil
}
