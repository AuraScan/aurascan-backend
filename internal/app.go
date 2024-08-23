package internal

import (
	"aurascan-backend/internal/config"
	"aurascan-backend/internal/router"
	"aurascan-backend/internal/websocket/example"
	"aurascan-backend/internal/websocket/pubsub"
	"aurascan-backend/util"
	"ch-common-package/cache"
	"ch-common-package/exit"
	"ch-common-package/logger"
	"ch-common-package/middleware"
	"ch-common-package/mongodb"
	"ch-common-package/ssdb"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(ctx context.Context, configPath string) error {
	config.MustLoad(configPath)
	logger.SetLog(config.Global.Log)
	if err := mongodb.SetupMongoDB(config.Global.MongoDB); err != nil {
		panic(err)
	}
	if err := ssdb.Setup(config.Global.SSDB); err != nil {
		panic(err)
	}
	cache.SetupRedis(config.Global.Redis)

	//go cron.NewCron()

	ctx, cancel := context.WithCancel(ctx)
	handler, err := initHttpServer()
	if err != nil {
		panic(err)
	}

	addr := fmt.Sprintf(":%d", config.Global.HTTP.Port)

	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  time.Duration(config.Global.HTTP.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.Global.HTTP.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(config.Global.HTTP.IdleTimeout) * time.Second,
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		fmt.Printf("HTTP server is running at %s, %v\n", addr, time.Now().Local())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("start http server err: ", err.Error())
			panic(err)
		}
	}()

	state := 1

EXIT:
	for {
		sig := <-sc
		logger.Infof("receive signal[%s]", sig.String())
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			state = 0
			break EXIT
		case syscall.SIGHUP:
		default:
			break EXIT
		}
	}

	cancel()

	exit.Exit()

	logger.Info("service exit")
	os.Exit(state)

	return nil
}

func initHttpServer() (*gin.Engine, error) {
	app := gin.New()
	app.Use(gin.Recovery())
	app.Use(middleware.Logger())
	app.Use(middleware.CrossDomain())
	app.Use(util.SaveLanguageToContext())

	app.NoMethod(middleware.NoMethodHandler())
	app.NoRoute(middleware.NoRouteHandler())

	pubsub.Start()
	pubsub.RegisterPublisher(&example.Block{})

	router.Register(app)

	return app, nil
}
