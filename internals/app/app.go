package app

import (
	"context"
	"net/http"
	"slurm/go-on-practice-2/http_06/api"
	"slurm/go-on-practice-2/http_06/api/middleware"
	db3 "slurm/go-on-practice-2/http_06/internals/app/db"
	"slurm/go-on-practice-2/http_06/internals/app/handlers"
	"slurm/go-on-practice-2/http_06/internals/app/processors"
	"slurm/go-on-practice-2/http_06/internals/cfg"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type Server struct {
	config cfg.Cfg
	ctx    context.Context
	srv    *http.Server
	db     *pgxpool.Pool
}

func NewServer(config cfg.Cfg, ctx context.Context) *Server {
	server := new(Server)
	server.ctx = ctx
	server.config = config
	return server
}

func (server *Server) Serve() {
	logrus.Println("Starting server")
	var err error
	server.db, err = pgxpool.Connect(server.ctx, server.config.GetDBString())
	if err != nil {
		logrus.Fatalln(err)
	}

	carsStorage := db3.NewCarsStorage(server.db)
	usersStorage := db3.NewUsersStorage(server.db)

	carsProcessor := processors.NewCarsProcessor(carsStorage)
	usersProcessor := processors.NewUsersProcessor(usersStorage)

	carsHandler := handlers.NewCarsHandler(carsProcessor)
	usersHandler := handlers.NewUsersHandler(usersProcessor)

	routes := api.CreateRoutes(usersHandler, carsHandler)
	routes.Use(middleware.RequestLog)

	server.srv = &http.Server{
		Addr:    ":" + server.config.Port,
		Handler: routes,
	}

	logrus.Println("server started")

	err = server.srv.ListenAndServe()

	if err != nil {
		logrus.Fatalln(err)
	}
}

func (server *Server) Shutdown() {
	logrus.Println("server stopped")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	server.db.Close()

	defer func() {
		cancel()
	}()

	var err error
	err = server.srv.Shutdown(ctxShutdown)
	if err != nil {
		logrus.Fatalf("server Shutdown failed: %s\n", err.Error())
	}

	logrus.Println("server exited properly")

	if err == http.ErrServerClosed {
		err = nil
	}
}
