package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrsimonemms/conveyor-belt/pkg/config"
	"github.com/rs/zerolog/log"
)

func webhookHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusAccepted, "hello world2")
	}
}

func (s *Server) Triggers() error {
	if len(s.cfg.Spec.Triggers) == 0 {
		return ErrNoTriggers
	}

	for _, t := range s.cfg.Spec.Triggers {
		switch t.Type {
		case config.PipelineTriggerTypeWebHook:
			relativePath := fmt.Sprintf("/webhook/%s", s.cfg.Metadata.Name)

			log.Debug().Str("path", relativePath).Msg("Registering webhook trigger")
			s.gin.POST(relativePath, webhookHandler())
		default:
			return fmt.Errorf("%w: %s", ErrUnknownTrigger, t.Type)
		}
	}

	return nil
}
