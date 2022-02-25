package main

import (
	"github.com/gorilla/mux"
	httpController "github.com/norby7/questions-rest-api/interfaceAdapters/http"
	httpServer "github.com/norby7/questions-rest-api/server/http"
	"github.com/norby7/questions-rest-api/usecases/repository"
	ucService "github.com/norby7/questions-rest-api/usecases/service"
	"log"
	"os"
	"strconv"
)

func main() {
	dbPath := "./database/questions.db"
	l := log.New(os.Stdout, "question-api", log.LstdFlags)

	err := repository.CreateDatabase(dbPath)
	if err != nil {
		l.Fatalln(err.Error())
	}

	repo, err := repository.NewSqliteRepository(dbPath)
	if err != nil {
		l.Fatalln("unable to create new repository: " + err.Error())
	}

	defer repo.Handler.Close()

	err = repository.ValidateSchema(repo.Handler)
	if err != nil {
		l.Fatalln(err.Error())
	}

	service := ucService.NewService(repo)
	controller := httpController.NewController(service, l)

	muxRouter := mux.NewRouter()
	httpServer.RegisterRoutes(muxRouter, *controller)

	port := os.Getenv("PORT")

	portAdr, err := strconv.Atoi(port)
	if err != nil {
		portAdr = 3000
	}
	httpServer.StartServer(muxRouter, portAdr)
}
