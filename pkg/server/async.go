package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrsimonemms/conveyor-belt/pkg/pipeline"
)

func async(r *gin.Engine) {
	r.POST("/async/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")

		// @todo(sje): do some logging

		fn, err := pipeline.GetCallback(id)
		if err != nil {
			// Check the output
			ctx.JSON(http.StatusBadRequest, err)
			return
		}

		if err := fn(); err != nil {
			// Check the output
			ctx.JSON(http.StatusBadRequest, err)
			return
		}

		// All finished - nothing to report
		ctx.Writer.WriteHeader(http.StatusNoContent)
	})
}
