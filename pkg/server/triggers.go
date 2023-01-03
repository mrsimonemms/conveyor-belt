package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrsimonemms/conveyor-belt/pkg/config"
	"github.com/mrsimonemms/conveyor-belt/pkg/pipeline"
	"github.com/rs/zerolog/log"
)

func webhookHandler(p *pipeline.Pipeline) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Run as async go routine
		go func() {
			_ = p.Run()
		}()

		ctx.Writer.WriteHeader(http.StatusAccepted)
	}
}

func (s *Server) Triggers(p *pipeline.Pipeline) error {
	if len(s.cfg.Spec.Triggers) == 0 {
		return ErrNoTriggers
	}

	for _, t := range s.cfg.Spec.Triggers {
		switch t.Type {
		case config.PipelineTriggerTypeWebHook:
			relativePath := fmt.Sprintf("/webhook/%s", s.cfg.Metadata.Name)

			log.Debug().Str("path", relativePath).Msg("Registering webhook trigger")
			s.gin.POST(relativePath, webhookHandler(p))
		default:
			return fmt.Errorf("%w: %s", ErrUnknownTrigger, t.Type)
		}
	}

	return nil
}
