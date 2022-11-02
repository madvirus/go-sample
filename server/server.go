package server

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go-sample/appctx"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HandlerRegisterer interface {
	RegisterHandler(router *gin.Engine)
}

func CreateApiServer(config appctx.Config, ctx appctx.AppCtx) (*ApiServer, error) {
	registers := appctx.GetAllByType[HandlerRegisterer](ctx)

	router := gin.Default()
	for _, r := range registers {
		r.RegisterHandler(router)
	}

	srv := &http.Server{
		Addr:    config.ServerAddr,
		Handler: router,
	}

	return &ApiServer{
		server: srv,
	}, nil
}

type ApiServer struct {
	server *http.Server
}

func (s *ApiServer) Run() {
	go func() {
		if err := s.server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Infof("listen: %s", err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
}
