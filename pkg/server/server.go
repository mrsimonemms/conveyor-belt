package server

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/mrsimonemms/conveyor-belt/pkg/config"
	logger "github.com/mrsimonemms/gin-structured-logger"
	"github.com/rs/zerolog/log"
)

type Server struct {
	cfg  *config.Config
	gin  *gin.Engine
	port int
}

func (s *Server) Start() error {
	log.Info().Int("port", s.port).Msg("Starting server...")

	return s.gin.Run(fmt.Sprintf(":%d", s.port))
}

func New(cfg *config.Config) Server {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(
		requestid.New(),
		logger.New(),
		gin.Recovery(),
		func(ctx *gin.Context) {
			logger.Get(ctx).Debug().Str("path", ctx.Request.URL.Path).Msg("New HTTP call")
		},
	)
	r.NoRoute(func(ctx *gin.Context) {
		statusCode := http.StatusNotFound
		message := fmt.Sprintf("Route %s:%s not found", ctx.Request.Method, ctx.Request.URL)

		ctx.JSON(statusCode, map[string]interface{}{
			"message":    message,
			"error":      http.StatusText(statusCode),
			"statusCode": statusCode,
		})
	})

	port := cfg.Spec.Port
	if port == 0 {
		log.Debug().Msg("Port not set")
		port = 3000
	}

	return Server{
		cfg:  cfg,
		gin:  r,
		port: port,
	}
}
