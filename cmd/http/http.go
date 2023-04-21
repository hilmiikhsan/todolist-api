package http

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"todolist-api/cmd/http/handlers/activity"
	"todolist-api/cmd/http/handlers/todo"
	"todolist-api/cmd/http/routers"
	"todolist-api/config"
	activityRepository "todolist-api/data/repositories/activity"
	todoRepository "todolist-api/data/repositories/todo"
	"todolist-api/infra/db"

	"todolist-api/infra/context/repository"
	"todolist-api/infra/context/service"

	activityService "todolist-api/cmd/services/activity"
	todoService "todolist-api/cmd/services/todo"

	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	routerCMD = &cobra.Command{
		Use:   "serve-http",
		Short: "Run http server",
		Long:  "API Todolist",
		RunE:  runHTTP,
	}
)

// initRepoCtx for context repository
func initRepoCtx(db *db.DB) *repository.RepoCtx {
	activityRepository := activityRepository.NewActivityRepository(db)
	todoRepository := todoRepository.NewTodoRepository(db)

	return &repository.RepoCtx{
		DB:                 db,
		ActivityRepository: activityRepository,
		TodoRepository:     todoRepository,
	}
}

// initServiceCtx for contextService
func initServiceCtx(ctx *repository.RepoCtx) *service.Ctx {
	activityService := activityService.NewActivityService(ctx)
	todoService := todoService.NewTodoService(ctx)

	return &service.Ctx{
		ActivityService: activityService,
		TodoService:     todoService,
	}
}

func runHTTP(cmd *cobra.Command, args []string) error {
	// initial config
	ctx := context.Background()
	cfg := config.InitConfig()

	// this Pings the database trying to connect, panics on error
	// use sqlx.Open() for sql.Open() semantics
	db, err := db.Open(&cfg.DB)
	if err != nil {
		log.Fatalln(err)
	}

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*time.Duration(cfg.Server.GraceFulTimeout), "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	// init repo ctx
	repoCtx := initRepoCtx(db)

	// init service ctx
	serviceCtx := initServiceCtx(repoCtx)

	// init handler
	activityHandler := activity.NewActivityHandler(serviceCtx)
	todoHandler := todo.NewTodoHandler(serviceCtx)

	// initial router
	r := routers.InitialRouter(
		activityHandler,
		todoHandler,
	)

	corsHandler := cors.New(cors.Options{
		AllowedHeaders: []string{"Origin", "Authorization", "Content-Type", "Access-Control-Allow-Origin", "API-KEY"},
		AllowedMethods: []string{"HEAD", "PUT", "PATCH", "GET", "POST", "DELETE", "OPTIONS"},
		AllowedOrigins: []string{
			"http://localhost:3030",
		},
		OptionsPassthrough: false,
		AllowCredentials:   true,
	})

	// server conf
	srv := &http.Server{
		Handler: corsHandler.Handler(r),
		Addr:    cfg.Server.Addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
	}

	fmt.Printf("API Listening on %s", cfg.Server.Addr)
	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	err = srv.Shutdown(ctx)
	if err != nil {
		logrus.Error(err)
	}
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)

	return nil
}

// ServeHTTP return instance of serve HTTP command object
func ServeHTTP() *cobra.Command {
	return routerCMD
}
