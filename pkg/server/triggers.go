package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mrsimonemms/conveyor-belt/pkg/config"
	"github.com/mrsimonemms/conveyor-belt/pkg/pipeline"
	logger "github.com/mrsimonemms/gin-structured-logger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TriggerCallback(url string, p *pipeline.Pipeline, l zerolog.Logger) error {
	lg := l.With().Str("url", url).Logger()

	lg.Debug().Msg("Triggering pipeline completion callback")
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	data, err := json.Marshal(p)
	if err != nil {
		lg.Error().Err(err).Msg("Failed to decode pipeline to JSON")
		return err
	}

	if _, err := client.Post(url, "application/json", bytes.NewBuffer(data)); err != nil {
		lg.Error().Err(err).Msg("Failed to send completion callback")
		return err
	}

	lg.Debug().Msg("Completion callback triggered successfully")
	return nil
}

func webhookHandler(p *pipeline.Pipeline) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Run as async go routine
		go func() {
			completionAction := make([]pipeline.CompletionAction, 0)
			// If x-callback-url exists, trigger a POST to the callback with the collated data
			if callbackUrl := ctx.GetHeader("x-callback-url"); callbackUrl != "" {
				completionAction = append(completionAction, func(p *pipeline.Pipeline, l zerolog.Logger) {
					// For HTTP, don't use the received logger, use the HTTP context one
					_ = TriggerCallback(callbackUrl, p, logger.Get(ctx).With().Logger())
				})
			}

			completionAction = append(completionAction, func(p *pipeline.Pipeline, l zerolog.Logger) {
				logger.Get(ctx).Debug().Msg("Pipeline run completed")
			})

			_ = p.Run(completionAction...)
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
