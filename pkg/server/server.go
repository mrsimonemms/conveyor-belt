package server

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	logger "github.com/mrsimonemms/gin-structured-logger"
	"github.com/rs/zerolog/log"
)

type Server struct {
	gin  *gin.Engine
	port int
}

func (s *Server) Start() error {
	log.Info().Int("port", s.port).Msg("Starting server...")

	return s.gin.Run(fmt.Sprintf(":%d", s.port))
}

func New(port int) Server {
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

	log.Info().Int("port", port).Msg("Starting server...")

	return Server{
		gin:  r,
		port: port,
	}
}
